package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/goverland-labs/goverland-platform-events/events/ipfs"
	"github.com/goverland-labs/goverland-platform-events/pkg/natsclient"
	"github.com/goverland-labs/snapshot-sdk-go/client"
	"github.com/rs/zerolog/log"

	"github.com/goverland-labs/goverland-datasource-snapshot/internal/db"
)

const (
	FetchMessagesCommandName = "fetch-messages"

	FetchTypeUnspecified    FetchType = "unspecified"
	FetchTypeDeleteProposal FetchType = "delete-proposal"
)

type FetchType string

type FetchMessages struct {
	base

	Messages  *db.MessageService
	Publisher *natsclient.Publisher
}

func (c *FetchMessages) GetName() string {
	return FetchMessagesCommandName
}

func (c *FetchMessages) GetArguments() ArgumentsDetails {
	return ArgumentsDetails{
		"type":  "import type: delete-proposal",
		"limit": "number of maximum rows for processing",
	}
}

func (c *FetchMessages) ParseArgs(args ...string) (Arguments, error) {
	return c.parseArgs(c, args...)
}

func (c *FetchMessages) Execute(args Arguments) error {
	start := time.Now()
	log.Info().Msg("import started")

	importType, err := c.getImportType(args)
	if err != nil {
		return err
	}

	limit, err := c.getLimit(args)
	if err != nil {
		return err
	}

	switch importType {
	case FetchTypeDeleteProposal:
		err = c.processDeleteProposals(int(limit))
	default:
		panic(fmt.Sprintf("import type is not implemented: %s", importType))
	}

	if err != nil {
		return fmt.Errorf("fetching messages: %w", err)
	}

	log.Info().Msgf("fetch finished. Took: %v", time.Since(start))

	return nil
}

func (c *FetchMessages) processDeleteProposals(limit int) error {
	offset := 0
	for {
		list, err := c.Messages.FindDeleteProposals(limit, offset)
		if err != nil {
			return fmt.Errorf("find delete proposals: %w", err)
		}

		for idx, snapshot := range list {
			var data client.MessageFragment
			if err := json.Unmarshal([]byte(snapshot), &data); err != nil {
				return fmt.Errorf("unmarshal delete proposal %d: %w", idx, err)
			}

			if err := c.Publisher.PublishJSON(context.Background(), ipfs.SubjectMessageCreated, ipfs.MessagePayload{
				IpfsID: *data.GetIpfs(),
				Type:   "delete-proposal",
			}); err != nil {
				return fmt.Errorf("publish message %d: %w", idx, err)
			}
		}

		if len(list) < limit {
			break
		}

		offset += limit

		log.Info().Msgf("fetched %d items with offset %d", limit, offset)
	}

	return nil
}

func (c *FetchMessages) getImportType(args Arguments) (FetchType, error) {
	switch args.Get("type") {
	case "delete-proposal":
		return FetchTypeDeleteProposal, nil
	default:
		return FetchTypeUnspecified, fmt.Errorf("type has wrong format: %s", args.Get("type"))
	}
}

func (c *FetchMessages) getLimit(args Arguments) (int64, error) {
	limit := args.Get("limit")
	if limit == "" {
		return 0, nil
	}

	return strconv.ParseInt(limit, 0, 64)
}
