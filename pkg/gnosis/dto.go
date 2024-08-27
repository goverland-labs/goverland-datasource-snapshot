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

//	{
//	 "chainId": "1",
//	 "blockNumber": "19947439",
//	 "address": "0x37F1eE65C2F8610741cd9Dff1057F926809C4078",
//	 "votingPower": 21.41231,
//	 "incomingPower": 21.41231,
//	 "outgoingPower": 10.00000,
//	 "percentOfVotingPower": 211,
//	 "percentOfDelegators": 600,
//	 delegators: [
//	   "0x485E60C486671E932fd9C53d4110cdEab1E7F0eb"
//	 ],
//	 delegatorTree [
//	   {
//	     "delegator": "0x485E60C486671E932fd9C53d4110cdEab1E7F0eb",
//	     "weight": 10000,
//	     "delegatedPower": 21.71540506,
//	     "parents": []
//	   }
//	 ],
//	 delegates [
//	   "0xD476B79539781e499396761CE7e21ab28AeA828F"
//	 ],
//	 delegateTree [
//	   {
//	     "delegate": "0xD476B79539781e499396761CE7e21ab28AeA828F",
//	     "weight": 5000,
//	     "delegatedPower": 10.00000,
//	     "children": []
//	   }
//	 ],
//	}
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
