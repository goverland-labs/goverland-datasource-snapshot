package updates

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/goverland-labs/snapshot-sdk-go/client"
	"github.com/goverland-labs/snapshot-sdk-go/snapshot"
	"github.com/rs/zerolog/log"

	"github.com/goverland-labs/goverland-datasource-snapshot/internal/db"
	"github.com/goverland-labs/goverland-datasource-snapshot/internal/helpers"
)

const messagesPerRequest = 1000

type MessagesWorker struct {
	sdk      *snapshot.SDK
	messages *db.MessageService

	checkInterval time.Duration
}

func NewMessagesWorker(sdk *snapshot.SDK, messages *db.MessageService, checkInterval time.Duration) *MessagesWorker {
	return &MessagesWorker{
		sdk:      sdk,
		messages: messages,

		checkInterval: checkInterval,
	}
}

func (w *MessagesWorker) Start(ctx context.Context) error {
	for {
		if err := w.loop(ctx); err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			err := ctx.Err()
			if errors.Is(err, context.Canceled) {
				return nil
			}

			return err
		case <-time.After(w.checkInterval):
		}
	}
}

func (w *MessagesWorker) loop(ctx context.Context) error {
	latestMCI, err := w.messages.GetLatestMCI()
	if err != nil {
		return err
	}

	for {
		messages, err := w.fetchMessagesInternal(ctx, snapshot.ListMessageWithPagination(messagesPerRequest, 0), snapshot.ListMessageWithMCIFilter(latestMCI))
		if err != nil {
			return err
		}

		log.Info().Int("latest_mci", latestMCI).Int("count", len(messages)).Msg("fetched messages")

		if err := w.processMessages(messages); err != nil {
			return err
		}

		if len(messages) < messagesPerRequest {
			return nil
		}

		latestMCI = w.extractMaxMCI(messages)

		select {
		case <-ctx.Done():
			return nil
		default:
		}
	}
}

func (w *MessagesWorker) extractMaxMCI(messages []*client.MessageFragment) int {
	var latest int64
	for _, m := range messages {
		mci := helpers.ZeroIfNil(m.GetMci())
		if mci > latest {
			latest = mci
		}
	}

	return int(latest)
}

func (w *MessagesWorker) fetchMessagesInternal(ctx context.Context, opts ...snapshot.ListMessageOption) ([]*client.MessageFragment, error) {
	for {
		messages, err := w.sdk.ListMessage(ctx, opts...)
		if errors.Is(err, snapshot.ErrTooManyRequests) {
			log.Warn().Err(err).Msg("snapshot api limits are reached")
			<-time.After(5 * time.Second)
			continue
		}

		return messages, err
	}
}

func (w *MessagesWorker) processMessages(messages []*client.MessageFragment) error {
	for _, message := range messages {
		marshaled, err := json.Marshal(message)
		if err != nil {
			return err
		}

		s := db.Message{
			ID:        helpers.ZeroIfNil(message.GetID()),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			MCI:       int(helpers.ZeroIfNil(message.GetMci())),
			Space:     helpers.ZeroIfNil(message.GetSpace()),
			Timestamp: time.Unix(helpers.ZeroIfNil(message.GetTimestamp()), 0),
			Type:      db.MessageType(helpers.ZeroIfNil(message.GetType())),
			Snapshot:  marshaled,
		}

		if err := w.messages.Upsert(&s); err != nil {
			return err
		}
	}

	return nil
}
