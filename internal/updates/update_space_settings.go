package updates

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"

	pevents "github.com/goverland-labs/goverland-platform-events/events/ipfs"
	client "github.com/goverland-labs/goverland-platform-events/pkg/natsclient"

	"github.com/goverland-labs/goverland-datasource-snapshot/internal/config"
	"github.com/goverland-labs/goverland-datasource-snapshot/internal/db"
	"github.com/goverland-labs/goverland-datasource-snapshot/internal/fetcher"
)

const (
	updateSpaceConsGroup = "update_space_settings"
)

type UpdateSpaceSettingsConsumer struct {
	con *client.Consumer[pevents.MessagePayload]

	spaceUpdater *SpacesUpdater
	fc           *fetcher.Client
	nc           *nats.Conn
}

func NewUpdateSpaceSettingsConsumer(
	su *SpacesUpdater,
	fc *fetcher.Client,
	nc *nats.Conn,
) *UpdateSpaceSettingsConsumer {
	return &UpdateSpaceSettingsConsumer{
		spaceUpdater: su,
		fc:           fc,
		nc:           nc,
	}
}

func (c *UpdateSpaceSettingsConsumer) Start(ctx context.Context) error {
	group := config.GenerateGroupName(updateSpaceConsGroup)
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

func (c *UpdateSpaceSettingsConsumer) stop() error {
	if c.con == nil {
		return nil
	}

	if err := c.con.Close(); err != nil {
		log.Error().Err(err).Msg("unable to close update space settings consumer")

		return err
	}

	return nil
}

func (c *UpdateSpaceSettingsConsumer) handler() pevents.MessageHandler {
	return func(payload pevents.MessagePayload) error {
		// process only update space settings
		if payload.Type != string(db.SettingsMessage) {
			return nil
		}

		space, err := c.fc.GetUpdatedSpaceByIpfsID(context.Background(), payload.IpfsID)
		if err != nil {
			log.Error().Err(err).Msg("get updated space by ipfs message")

			return err
		}

		if err = c.spaceUpdater.ProcessSpace(context.Background(), space); err != nil {
			log.Error().Err(err).Msg("process updated space message")

			return err
		}

		log.Debug().Msgf("updated space[%s] by ipfs message: %s", space, payload.IpfsID)

		return nil
	}
}
