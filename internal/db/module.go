package db

import "context"

type Publisher interface {
	PublishJSON(ctx context.Context, subject string, obj any) error
}
