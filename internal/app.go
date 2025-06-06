package internal

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Yamashou/gqlgenc/clientv2"
	"github.com/goverland-labs/goverland-ipfs-fetcher/protocol/ipfsfetcherpb"
	"github.com/goverland-labs/goverland-platform-events/pkg/natsclient"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	grcpprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/nats-io/nats.go"
	grpczerolog "github.com/pereslava/grpc_zerolog"
	"github.com/rs/zerolog/log"
	"github.com/s-larionov/process-manager"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/goverland-labs/snapshot-sdk-go/snapshot"

	"github.com/goverland-labs/goverland-datasource-snapshot/protocol/delegatepb"
	"github.com/goverland-labs/goverland-datasource-snapshot/protocol/votingpb"

	"github.com/goverland-labs/goverland-datasource-snapshot/internal/delegate"
	"github.com/goverland-labs/goverland-datasource-snapshot/internal/fetcher"
	"github.com/goverland-labs/goverland-datasource-snapshot/internal/updates"
	"github.com/goverland-labs/goverland-datasource-snapshot/pkg/gnosis"

	"github.com/goverland-labs/goverland-datasource-snapshot/internal/config"
	"github.com/goverland-labs/goverland-datasource-snapshot/internal/db"
	"github.com/goverland-labs/goverland-datasource-snapshot/internal/metrics"
	"github.com/goverland-labs/goverland-datasource-snapshot/internal/voting"
	"github.com/goverland-labs/goverland-datasource-snapshot/pkg/grpcsrv"
	"github.com/goverland-labs/goverland-datasource-snapshot/pkg/health"
	"github.com/goverland-labs/goverland-datasource-snapshot/pkg/prometheus"
)

type Application struct {
	sigChan <-chan os.Signal
	manager *process.Manager
	cfg     config.App

	proposalsRepo    *db.ProposalRepo
	proposalsService *db.ProposalService

	spacesRepo    *db.SpaceRepo
	spacesService *db.SpaceService

	votesRepo    *db.VoteRepo
	votesService *db.VoteService

	preparedVotesRepo   *db.PreparedVoteRepo
	actionVotingService *voting.ActionService

	messagesRepo    *db.MessageRepo
	messagesService *db.MessageService

	publisher *natsclient.Publisher
	natsConn  *nats.Conn

	sdk       *snapshot.SDK
	votingSDK *snapshot.SDK

	isCliMode bool
}

func NewApplication(cfg config.App) (*Application, error) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	a := &Application{
		sigChan: sigChan,
		cfg:     cfg,
		manager: process.NewManager(),
	}

	err := a.bootstrap()
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *Application) Run() {
	a.manager.StartAll()
	a.registerShutdown()
}

func (a *Application) bootstrap() error {
	initializers := []func() error{
		// Init Dependencies
		a.initDatabase,
		a.initNats,
		a.initSnapshot,
		a.initServices,
		a.initGrpc,
	}

	if !a.isCliMode {
		// Init Workers: Application
		initializers = append(initializers, a.initUpdatesWorkers)
		// Init Workers: System
		initializers = append(initializers, a.initPrometheusWorker)
		initializers = append(initializers, a.initHealthWorker)
	}

	for _, initializer := range initializers {
		if err := initializer(); err != nil {
			return err
		}
	}

	return nil
}

func (a *Application) initDatabase() error {
	conn, err := gorm.Open(postgres.Open(a.cfg.Database.DSN), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlConnection, err := conn.DB()
	if err != nil {
		return err
	}
	sqlConnection.SetMaxOpenConns(a.cfg.Database.MaxOpenConnections)
	sqlConnection.SetMaxIdleConns(a.cfg.Database.MaxIdleConnections)

	if a.cfg.Database.Debug {
		conn = conn.Debug()
	}

	// TODO: Use real migrations intead of auto migrations from gorm
	if err := conn.AutoMigrate(&db.Space{}, &db.Proposal{}, &db.Vote{}, &db.PreparedVote{}, &db.Message{}); err != nil {
		return err
	}

	a.proposalsRepo = db.NewProposalRepo(conn)
	a.spacesRepo = db.NewSpaceRepo(conn)
	a.votesRepo = db.NewVoteRepo(conn)
	a.preparedVotesRepo = db.NewPreparedVoteRepo(conn)
	a.messagesRepo = db.NewMessageRepo(conn)

	return err
}

func (a *Application) initNats() error {
	nc, err := nats.Connect(
		a.cfg.Nats.URL,
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(a.cfg.Nats.MaxReconnects),
		nats.ReconnectWait(a.cfg.Nats.ReconnectTimeout),
	)
	if err != nil {
		return err
	}

	publisher, err := natsclient.NewPublisher(nc)
	if err != nil {
		return err
	}

	a.publisher = publisher
	a.natsConn = nc

	return nil
}

func (a *Application) initSnapshot() error {
	metricsMiddleware := func(name string) clientv2.RequestInterceptor {
		return func(ctx context.Context, req *http.Request, gqlInfo *clientv2.GQLRequestInfo, res interface{}, next clientv2.RequestInterceptorFunc) (err error) {
			defer func(start time.Time) {
				metrics.CollectRequestsMetric(name, gqlInfo.Request.OperationName, err, start)
			}(time.Now())

			return next(ctx, req, gqlInfo, res)
		}
	}

	opts := []snapshot.Option{
		snapshot.WithHTTPClient(&http.Client{
			Transport: metrics.NewHeaderWatcher("general"),
		}),
		snapshot.WithInterceptors([]clientv2.RequestInterceptor{
			metricsMiddleware("general"),
		}),
	}
	if a.cfg.Snapshot.APIKey != "" {
		opts = append(opts, snapshot.WithApiKey(a.cfg.Snapshot.APIKey))
	}

	a.sdk = snapshot.NewSDK(opts...)

	votingOpts := []snapshot.Option{
		snapshot.WithHTTPClient(&http.Client{
			Transport: metrics.NewHeaderWatcher("voting"),
		}),
		snapshot.WithInterceptors([]clientv2.RequestInterceptor{
			metricsMiddleware("voting"),
		}),
	}
	if a.cfg.Snapshot.VotingAPIKey != "" {
		votingOpts = append(votingOpts, snapshot.WithApiKey(a.cfg.Snapshot.VotingAPIKey))
	}
	a.votingSDK = snapshot.NewSDK(votingOpts...)

	return nil
}

func (a *Application) initServices() error {
	a.proposalsService = db.NewProposalService(a.proposalsRepo, a.publisher)
	a.spacesService = db.NewSpaceService(a.spacesRepo, a.publisher)
	a.votesService = db.NewVoteService(a.votesRepo, a.publisher)
	a.messagesService = db.NewMessageService(a.messagesRepo, a.publisher)

	a.actionVotingService = voting.NewActionService(a.votingSDK, a.proposalsRepo, voting.NewTypedSignDataBuilder(a.cfg.Snapshot), a.preparedVotesRepo)

	return nil
}

func (a *Application) initGrpc() error {
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grcpprometheus.UnaryServerInterceptor,
			grpczerolog.NewUnaryServerInterceptor(log.Logger),
			grpcrecovery.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			grcpprometheus.StreamServerInterceptor,
			grpczerolog.NewStreamServerInterceptor(log.Logger),
			grpcrecovery.StreamServerInterceptor(),
		),
	)
	reflection.Register(grpcServer)

	votingGrpc := voting.NewGrpcServer(a.actionVotingService)
	votingpb.RegisterVotingServer(grpcServer, votingGrpc)

	delegatesGrpc := delegate.NewGrpcServer(delegate.NewService(gnosis.NewSDK()))
	delegatepb.RegisterDelegateServer(grpcServer, delegatesGrpc)

	grpcWorker := grpcsrv.NewGrpcServerWorker("snapshot", grpcServer, a.cfg.InternalAPI.Bind)
	a.manager.AddWorker(grpcWorker)

	return nil
}

func (a *Application) initUpdatesWorkers() error {
	spacesUpdater := updates.NewSpacesUpdater(a.sdk, a.spacesService)
	proposals := updates.NewProposalsWorker(a.sdk, a.proposalsService, a.cfg.Snapshot.ProposalsCheckInterval)
	activeProposals := updates.NewActiveProposalsWorker(a.sdk, a.proposalsService, a.cfg.Snapshot.ProposalsUpdatesInterval)
	spaces := updates.NewSpacesWorker(spacesUpdater, a.spacesService, a.cfg.Snapshot.UnknownSpacesCheckInterval)
	votes := updates.NewVotesWorker(a.sdk, a.votesService, a.proposalsService, a.messagesService, a.cfg.Snapshot.VotesCheckInterval)
	messages := updates.NewMessagesWorker(a.sdk, a.messagesService, a.cfg.Snapshot.MessagesCheckInterval)

	conn, err := grpc.NewClient(
		a.cfg.InternalAPI.IpfsFetcherAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("create connection with ipfs fetcher server: %v", err)
	}

	fc := ipfsfetcherpb.NewMessageClient(conn)
	fetcherWrapper := fetcher.NewClient(fc)
	deleteProposals := updates.NewDeleteProposalConsumer(a.proposalsService, fetcherWrapper, a.natsConn)

	updateSpaceConsumer := updates.NewUpdateSpaceSettingsConsumer(spacesUpdater, fetcherWrapper, a.natsConn)

	a.manager.AddWorker(process.NewCallbackWorker("snapshot proposals updates", proposals.FetchList, process.RetryOnErrorOpt{Timeout: 5 * time.Second}))
	a.manager.AddWorker(process.NewCallbackWorker("snapshot proposals mark to refetch", proposals.MarkToRefetch, process.RetryOnErrorOpt{Timeout: 5 * time.Second}))
	a.manager.AddWorker(process.NewCallbackWorker("snapshot active proposals updates", activeProposals.Start, process.RetryOnErrorOpt{Timeout: 5 * time.Second}))
	a.manager.AddWorker(process.NewCallbackWorker("snapshot unknown spaces updates", spaces.Start, process.RetryOnErrorOpt{Timeout: 5 * time.Second}))
	a.manager.AddWorker(process.NewCallbackWorker("snapshot votes load historical", votes.LoadHistorical, process.RetryOnErrorOpt{Timeout: 5 * time.Second}))
	a.manager.AddWorker(process.NewCallbackWorker("snapshot votes load active", votes.LoadActive, process.RetryOnErrorOpt{Timeout: 5 * time.Second}))
	a.manager.AddWorker(process.NewCallbackWorker("snapshot messages updates", messages.Start, process.RetryOnErrorOpt{Timeout: 5 * time.Second}))
	a.manager.AddWorker(process.NewCallbackWorker("delete-proposal-consumer", deleteProposals.Start))
	a.manager.AddWorker(process.NewCallbackWorker("update-space-settings-consumer", updateSpaceConsumer.Start))

	return nil
}

func (a *Application) initPrometheusWorker() error {
	srv := prometheus.NewServer(a.cfg.Prometheus.Listen, "/metrics")
	a.manager.AddWorker(process.NewServerWorker("prometheus", srv))

	return nil
}

func (a *Application) initHealthWorker() error {
	srv := health.NewHealthCheckServer(a.cfg.Health.Listen, "/status", health.DefaultHandler(a.manager))
	a.manager.AddWorker(process.NewServerWorker("health", srv))

	return nil
}

func (a *Application) registerShutdown() {
	go func(manager *process.Manager) {
		<-a.sigChan

		manager.StopAll()
	}(a.manager)

	a.manager.AwaitAll()
}
