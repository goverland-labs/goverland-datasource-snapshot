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

type VoteParams struct {
	ID  uint64
	Sig string
}

type SuccessVote struct {
	ID      string
	IPFS    string
	Relayer Relayer
}

type Relayer struct {
	Address string
	Receipt string
}

type validateVoteResult struct {
	validation
}

type validateVPResult struct {
	validation

	votingPower float64
}

func (v validateVPResult) toValidateResult() ValidateResult {
	return ValidateResult{
		OK:              v.ok,
		VotingPower:     v.votingPower,
		ValidationError: v.ValidationError,
	}
}

type validation struct {
	ok bool

	ValidationError *ValidationError
}

func (v validation) toValidateResult() ValidateResult {
	return ValidateResult{
		OK:              v.ok,
		ValidationError: v.ValidationError,
	}
}
