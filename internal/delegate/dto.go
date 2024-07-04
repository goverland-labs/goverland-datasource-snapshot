package delegate

import (
	"encoding/json"
)

type GetDelegatesRequest struct {
	Dao       string
	Strategy  json.RawMessage
	By        string
	Addresses []string
	Limit     int
	Offset    int
}

type Delegate struct {
	Address              string  `json:"address"`
	DelegatorCount       int32   `json:"delegator_count"`
	PercentOfDelegators  int32   `json:"percent_of_delegators"`
	VotingPower          float64 `json:"voting_power"`
	PercentOfVotingPower int32   `json:"percent_of_voting_power"`
}

type GnosisTopDelegatesBodyRequest struct {
	Strategy json.RawMessage `json:"strategy"`
}

type GnosisTopDelegatesResponse struct {
	BlockNumber int64            `json:"block_number"`
	ChainID     int64            `json:"chain_id"`
	Delegates   []GnosisDelegate `json:"delegates"`
}

type GnosisDelegate struct {
	Address              string  `json:"address"`
	DelegatorCount       int32   `json:"delegator_count"`
	PercentOfDelegators  int32   `json:"percent_of_delegators"`
	VotingPower          float64 `json:"voting_power"`
	PercentOfVotingPower int32   `json:"percent_of_voting_power"`
}
