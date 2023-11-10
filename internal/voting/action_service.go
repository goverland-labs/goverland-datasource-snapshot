package voting

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/rs/zerolog/log"

	"github.com/goverland-labs/sdk-snapshot-go/client"
	"github.com/goverland-labs/sdk-snapshot-go/snapshot"
)

const (
	latestSnapshot = "latest"
)

type ActionService struct {
	snapshotSDK snapshotSDK

	proposalGetter       proposalGetter
	typedSignDataBuilder *TypedSignDataBuilder
}

func NewActionService(snapshotSDK snapshotSDK, proposalGetter proposalGetter, typedSignDataBuilder *TypedSignDataBuilder) *ActionService {
	return &ActionService{snapshotSDK: snapshotSDK, proposalGetter: proposalGetter, typedSignDataBuilder: typedSignDataBuilder}
}

func (a *ActionService) Validate(validateParams ValidateParams) (ValidateResult, error) {
	log.Info().Fields(map[string]any{
		"req": validateParams,
	}).Msg("got validation request")

	proposal, err := a.proposalGetter.GetByID(validateParams.Proposal)
	if err != nil {
		return ValidateResult{}, err
	}

	var pFragment client.ProposalFragment
	if err := json.Unmarshal(proposal.Snapshot, &pFragment); err != nil {
		return ValidateResult{}, err
	}

	vResult, err := a.validateVote(validateParams, &pFragment)
	if err != nil {
		return ValidateResult{}, err
	}
	if !vResult.ok {
		return buildErrorValidation(vResult.validation), nil
	}

	vpResult, err := a.validateVotingPower(validateParams, &pFragment)
	if err != nil {
		return ValidateResult{}, err
	}
	if !vpResult.ok {
		return buildErrorValidation(vpResult.validation), nil
	}

	return ValidateResult{
		OK:          true,
		VotingPower: vpResult.votingPower,
	}, nil
}

func (a *ActionService) Prepare(prepareParams PrepareParams) (PrepareResult, error) {
	proposal, err := a.proposalGetter.GetByID(prepareParams.Proposal)
	if err != nil {
		return PrepareResult{}, err
	}

	var pFragment client.ProposalFragment
	if err := json.Unmarshal(proposal.Snapshot, &pFragment); err != nil {
		return PrepareResult{}, err
	}

	typedData := a.typedSignDataBuilder.Build(prepareParams, &pFragment)
	typedDataJSON, err := json.Marshal(typedData)
	if err != nil {
		return PrepareResult{}, err
	}

	return PrepareResult{
		TypedData: string(typedDataJSON),
	}, nil
}

func (a *ActionService) validateVotingPower(validateParams ValidateParams, pFragment *client.ProposalFragment) (validateVPResult, error) {
	params := snapshot.GetVotingPowerParams{
		Address:    validateParams.Voter,
		Network:    pFragment.Network,
		Strategies: convertStrategies(pFragment.Strategies),
		Snapshot:   getSnapshot(pFragment.Snapshot),
		Space:      pFragment.Space.ID,
		Delegation: false, // TODO: add delegation
	}

	votingPowerResp, err := a.snapshotSDK.GetVotingPower(context.Background(), params)
	if err != nil {
		return validateVPResult{}, err
	}

	if votingPowerResp.Result.VP == 0 {
		return validateVPResult{
			validation: validation{
				ok: false,
			},
		}, nil
	}

	return validateVPResult{
		validation: validation{
			ok: true,
		},
		votingPower: votingPowerResp.Result.VP,
	}, nil
}

func (a *ActionService) validateVote(validateParams ValidateParams, pFragment *client.ProposalFragment) (validateVoteResult, error) {
	if pFragment.Validation.Name == "any" {
		return validateVoteResult{
			validation: validation{ok: true},
		}, nil
	}

	fParams := pFragment.Validation.Params
	if pFragment.Validation.Name == "basic" {
		_, ok := fParams["strategies"]
		if !ok {
			fParams["strategies"] = pFragment.Strategies
		}
	}

	params := snapshot.ValidationParams{
		Validation: pFragment.Validation.Name,
		Author:     validateParams.Voter,
		Space:      pFragment.Space.ID,
		Network:    pFragment.Network,
		Snapshot:   getSnapshot(pFragment.Snapshot),
		Params:     fParams,
	}

	validate, err := a.snapshotSDK.Validate(context.Background(), params)
	if err != nil {
		return validateVoteResult{}, err
	}

	// FIXME: build messages and code
	return validateVoteResult{
		validation: validation{ok: validate.Result},
	}, nil
}

func convertStrategies(strategies []*client.StrategyFragment) []snapshot.StrategyFragment {
	var result []snapshot.StrategyFragment
	for _, strategy := range strategies {
		params := make(map[string]interface{})
		if strategy.Params != nil {
			params = strategy.Params
		}
		result = append(result, snapshot.StrategyFragment{
			Name:    strategy.Name,
			Network: strategy.Network,
			Params:  params,
		})
	}

	return result
}

func getSnapshot(snapshot *string) any {
	if snapshot == nil {
		return latestSnapshot
	}

	numSnapshot, err := strconv.Atoi(*snapshot)
	if err != nil {
		log.Error().Err(err).Msg("failed to convert snapshot to int")

		return latestSnapshot
	}

	return numSnapshot
}

func buildErrorValidation(v validation) ValidateResult {
	return ValidateResult{
		OK: false,
		ValidationError: &ValidationError{
			Message: v.errorMsg,
			Code:    v.errorCode,
		},
	}
}
