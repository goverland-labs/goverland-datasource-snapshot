package voting

import "fmt"

var NoVotingPowerErr = &ValidationError{
	Message: "Oops, it seems you don't have any voting power at block %s.",
	Code:    1000,
}

var NotValidErr = &ValidationError{
	Message: "Oops, you don't seem to be eligible to submit a proposal.",
	Code:    1001,
}

var BasicValidationErr = &ValidationError{
	Message: "You do not meet the minimum balance requirement to vote on this proposal.",
	Code:    1002,
}

// PassportGatedValidationErr scoreThreshold, operator, stamps params
var PassportGatedValidationErr = &ValidationError{
	Message: "You need a Gitcoin Passport with score above %s and %s of the following stamps to vote on this proposal: %s.",
	Code:    1003,
}

func (v ValidationError) WithVars(vars ...any) *ValidationError {
	v.Message = fmt.Sprintf(v.Message, vars...)
	return &v
}
