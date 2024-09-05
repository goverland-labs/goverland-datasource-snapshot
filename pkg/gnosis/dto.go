package gnosis

import (
	"encoding/json"
)

type TopDelegatesRequest struct {
	Dao      string
	Strategy json.RawMessage
	By       string
	Limit    int
	Offset   int
}

type DelegateProfileRequest struct {
	Dao      string
	Address  string
	Strategy json.RawMessage
}

type strategyBodyRequest struct {
	Strategy json.RawMessage `json:"strategy"`
}

type TopDelegatesResponse struct {
	BlockNumber int64      `json:"blockNumber"`
	ChainID     int64      `json:"chainId"`
	Delegates   []Delegate `json:"delegates"`
}

type DelegateProfileResponse struct {
	BlockNumber          int64               `json:"blockNumber"`
	ChainID              int64               `json:"chainId"`
	Address              string              `json:"address"`
	VotingPower          float64             `json:"votingPower"`
	IncomingPower        float64             `json:"incomingPower"`
	OutgoingPower        float64             `json:"outgoingPower"`
	PercentOfVotingPower float64             `json:"percentOfVotingPower"`
	PercentOfDelegators  float64             `json:"percentOfDelegators"`
	Delegators           []string            `json:"delegators"`
	Delegates            []string            `json:"delegates"`
	DelegateTree         []DelegateTreeItem  `json:"delegateTree"`
	DelegatorTree        []DelegatorTreeItem `json:"delegatorTree"`
}

type DelegateTreeItem struct {
	Delegate           string  `json:"delegate"`
	Weight             float64 `json:"weight"`
	DelegatedPower     float64 `json:"delegatedPower"`
	ExpirationUnixTime int     `json:"expiration"`
}

type DelegatorTreeItem struct {
	Delegator      string  `json:"delegator"`
	Weight         float64 `json:"weight"`
	DelegatedPower float64 `json:"delegatedPower"`
}

type Delegate struct {
	Address              string  `json:"address"`
	DelegatorCount       int32   `json:"delegatorCount"`
	PercentOfDelegators  float64 `json:"percentOfDelegators"`
	VotingPower          float64 `json:"votingPower"`
	PercentOfVotingPower float64 `json:"percentOfVotingPower"`
}
