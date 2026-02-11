package users

import (
	"context"

	"github.com/goverland-labs/snapshot-sdk-go/client"
)

type snapshotSDK interface {
	ListUsers(
		ctx context.Context,
		addresses []string,
	) ([]*client.UserFragment, error)
}

type Info struct {
	Address string
	About   string
}
