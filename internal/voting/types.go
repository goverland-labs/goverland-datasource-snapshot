package voting

import "encoding/json"

type ValidateParams struct {
	Voter    string
	Proposal string
}

type ValidateResult struct {
	OK          bool
	VotingPower float64

	ValidationError *ValidationError
}

type ValidationError struct {
	Message string
	Code    uint32
}

type PrepareParams struct {
	Voter    string
	Proposal string
	Choice   json.RawMessage
	Reason   *string
}

type PrepareResult struct {
	TypedData string
}

type validateVoteResult struct {
	validation
}

type validateVPResult struct {
	validation

	votingPower float64
}

type validation struct {
	ok bool

	errorMsg  string
	errorCode uint32
}
