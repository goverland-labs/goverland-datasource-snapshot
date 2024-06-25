package updates

import (
	"time"

	client "github.com/goverland-labs/goverland-platform-events/pkg/natsclient"
)

const (
	maxPendingElements = 100
	rateLimit          = 500 * client.KiB
	executionTtl       = time.Minute
)
