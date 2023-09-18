package config

import (
	"time"
)

type Snapshot struct {
	ProposalsCheckInterval     time.Duration `env:"SNAPSHOT_PROPOSALS_CHECK_INTERVAL" envDefault:"1m"`
	VotesCheckInterval         time.Duration `env:"SNAPSHOT_VOTES_CHECK_INTERVAL" envDefault:"5s"`
	ProposalsUpdatesInterval   time.Duration `env:"SNAPSHOT_PROPOSALS_UPDATES_INTERVAL" envDefault:"10m"`
	UnknownSpacesCheckInterval time.Duration `env:"SNAPSHOT_UNKNOWN_SPACES_CHECK_INTERVAL" envDefault:"1m"`
	APIKey                     string        `env:"SNAPSHOT_API_KEY"`
}
