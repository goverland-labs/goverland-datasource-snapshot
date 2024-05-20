package updates

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"

	pevents "github.com/goverland-labs/goverland-platform-events/events/ipfs"
	client "github.com/goverland-labs/goverland-platform-events/pkg/natsclient"

	"github.com/goverland-labs/goverland-datasource-snapshot/internal/config"
	"github.com/goverland-labs/goverland-datasource-snapshot/internal/db"
	"github.com/goverland-labs/goverland-datasource-snapshot/internal/fetcher"
)

const (
	maxPendingElements = 100
	rateLimit          = 500 * client.KiB
	executionTtl       = time.Minute
)

type DeleteProposalConsumer struct {
	con *client.Consumer[pevents.MessagePayload]

	proposals *db.ProposalService
	fc        *fetcher.Client
	nc        *nats.Conn
}

func NewDeleteProposalConsumer(
	s *db.ProposalService,
	fc *fetcher.Client,
	nc *nats.Conn,
) *DeleteProposalConsumer {
	return &DeleteProposalConsumer{
		proposals: s,
		fc:        fc,
		nc:        nc,
	}
}

func (c *DeleteProposalConsumer) Start(ctx context.Context) error {
	group := config.GenerateGroupName("delete_proposal")
	opts := []client.ConsumerOpt{
		client.WithRateLimit(rateLimit),
		client.WithMaxAckPending(maxPendingElements),
		client.WithAckWait(executionTtl),
	}

	var err error
	c.con, err = client.NewConsumer(ctx, c.nc, group, pevents.SubjectMessageCollected, c.handler(), opts...)
	if err != nil {
		return fmt.Errorf("create consumer: %w", err)
	}

	log.Info().Msg("sender consumers is started")

	// todo: handle correct stopping the consumer by context
	<-ctx.Done()

	return c.stop()
}

func (c *DeleteProposalConsumer) stop() error {
	if c.con == nil {
		return nil
	}

	if err := c.con.Close(); err != nil {
		log.Error().Err(err).Msg("unable to close delete proposal consumer")

		return err
	}

	return nil
}

func (c *DeleteProposalConsumer) handler() pevents.MessageHandler {
	return func(payload pevents.MessagePayload) error {
		// process only deleted proposals
		if payload.Type != string(db.DeleteProposalMessage) {
			return nil
		}

		proposalID, err := c.fc.GetDeletedProposalIDByIpfsID(context.Background(), payload.IpfsID)
		if err != nil {
			log.Error().Err(err).Msg("getting deleted proposal ipfs data")

			return err
		}

		if err = c.proposals.Delete([]string{proposalID}); err != nil {
			log.Error().Err(err).Msg("process deleted proposal message")

			return err
		}

		log.Debug().Msgf("deleted proposal[%s] by ipfs message: %s", proposalID, payload.IpfsID)

		return nil
	}
}
