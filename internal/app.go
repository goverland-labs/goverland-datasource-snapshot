package internal

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/goverland-labs/sdk-snapshot-go/snapshot"
	"github.com/nats-io/nats.go"
	"github.com/s-larionov/process-manager"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/goverland-labs/datasource-snapshot/internal/config"
	"github.com/goverland-labs/datasource-snapshot/internal/db"
	"github.com/goverland-labs/datasource-snapshot/internal/updates"
	"github.com/goverland-labs/datasource-snapshot/pkg/communicate"
	"github.com/goverland-labs/datasource-snapshot/pkg/health"
	"github.com/goverland-labs/datasource-snapshot/pkg/prometheus"
)

type Application struct {
	sigChan <-chan os.Signal
	manager *process.Manager
	cfg     config.App

	proposalsRepo    *db.ProposalRepo
	proposalsService *db.ProposalService

	spacesRepo    *db.SpaceRepo
	spacesService *db.SpaceService

	publisher *communicate.Publisher

	sdk *snapshot.SDK
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

		// Init Workers: Application
		a.initUpdatesWorkers,

		// Init Workers: System
		a.initPrometheusWorker,
		a.initHealthWorker,
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
	if err := conn.AutoMigrate(&db.Space{}, &db.Proposal{}); err != nil {
		return err
	}

	a.proposalsRepo = db.NewProposalRepo(conn)
	a.spacesRepo = db.NewSpaceRepo(conn)

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

	publisher, err := communicate.NewPublisher(nc)
	if err != nil {
		return err
	}

	a.publisher = publisher

	return nil
}

func (a *Application) initSnapshot() error {
	var opts []snapshot.Option
	if a.cfg.Snapshot.APIKey != "" {
		opts = append(opts, snapshot.WithApiKey(a.cfg.Snapshot.APIKey))
	}

	a.sdk = snapshot.NewSDK(opts...)

	return nil
}

func (a *Application) initServices() error {
	a.proposalsService = db.NewProposalService(a.proposalsRepo, a.publisher)
	a.spacesService = db.NewSpaceService(a.spacesRepo, a.publisher)

	return nil
}

func (a *Application) initUpdatesWorkers() error {
	proposals := updates.NewProposalsWorker(a.sdk, a.proposalsService, a.cfg.Snapshot.ProposalsCheckInterval)
	activeProposals := updates.NewActiveProposalsWorker(a.sdk, a.proposalsService, a.cfg.Snapshot.ProposalsUpdatesInterval)
	spaces := updates.NewSpacesWorker(a.sdk, a.spacesService, a.cfg.Snapshot.UnknownSpacesCheckInterval)

	a.manager.AddWorker(process.NewCallbackWorker("snapshot proposals updates", proposals.Start, process.RetryOnErrorOpt{Timeout: 5 * time.Second}))
	a.manager.AddWorker(process.NewCallbackWorker("snapshot active proposals updates", activeProposals.Start, process.RetryOnErrorOpt{Timeout: 5 * time.Second}))
	a.manager.AddWorker(process.NewCallbackWorker("snapshot unknown spaces updates", spaces.Start, process.RetryOnErrorOpt{Timeout: 5 * time.Second}))

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
