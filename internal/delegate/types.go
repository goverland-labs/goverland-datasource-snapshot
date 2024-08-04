package delegate

import (
	"encoding/json"
)

type GetDelegatesParams struct {
	Dao       string
	Strategy  json.RawMessage
	By        *string
	Addresses []string
	Limit     int
	Offset    int
}

type Delegate struct {
	Address              string  `json:"address"`
	DelegatorCount       int32   `json:"delegator_count"`
	PercentOfDelegators  float64 `json:"percent_of_delegators"`
	VotingPower          float64 `json:"voting_power"`
	PercentOfVotingPower float64 `json:"percent_of_voting_power"`
}

type GetDelegateProfileParams struct {
	Dao      string
	Strategy json.RawMessage
	Address  string
}

type DelegateProfile struct {
	Address              string                `json:"address"`
	VotingPower          float64               `json:"voting_power"`
	IncomingPower        float64               `json:"incoming_power"`
	OutgoingPower        float64               `json:"outgoing_power"`
	PercentOfVotingPower float64               `json:"percent_of_voting_power"`
	PercentOfDelegators  float64               `json:"percent_of_delegators"`
	Delegates            []ProfileDelegateItem `json:"delegates"`
}

type ProfileDelegateItem struct {
	Address         string  `json:"address"`
	PercentOfWeight float64 `json:"percent_of_weight"`
	DelegatedPower  float64 `json:"delegated_power"`
}
