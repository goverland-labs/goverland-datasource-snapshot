package fetcher

import (
	"context"
	"encoding/json"
	"fmt"

	fetcher "github.com/goverland-labs/goverland-ipfs-fetcher/protocol/ipfsfetcherpb"
)

type Client struct {
	ic fetcher.MessageClient
}

func NewClient(ic fetcher.MessageClient) *Client {
	return &Client{
		ic: ic,
	}
}

func (c *Client) GetDeletedProposalIDByIpfsID(ctx context.Context, ipfsID string) (string, error) {
	data, err := c.ic.GetByID(ctx, &fetcher.GetByIDRequest{IpfsId: ipfsID})
	if err != nil {
		return "", fmt.Errorf("fetching deleted proposal by ipfs id: %w", err)
	}

	var info Info
	if err = json.Unmarshal(data.GetRawMessage().GetValue(), &info); err != nil {
		return "", fmt.Errorf("unmarshalling proposal: %w", err)
	}

	return info.Data.Message.Proposal, nil
}

func (c *Client) GetUpdatedSpaceByIpfsID(ctx context.Context, ipfsID string) (string, error) {
	data, err := c.ic.GetByID(ctx, &fetcher.GetByIDRequest{IpfsId: ipfsID})
	if err != nil {
		return "", fmt.Errorf("fetching updated space by ipfs id: %w", err)
	}

	var info Info
	if err = json.Unmarshal(data.GetRawMessage().GetValue(), &info); err != nil {
		return "", fmt.Errorf("unmarshalling proposal: %w", err)
	}

	return info.Data.Message.Space, nil
}
