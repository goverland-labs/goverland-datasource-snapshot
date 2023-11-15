package voting

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"

	"github.com/goverland-labs/sdk-snapshot-go/client"
	"github.com/goverland-labs/sdk-snapshot-go/snapshot"

	"github.com/goverland-labs/datasource-snapshot/internal/db"
)

const (
	latestSnapshot = "latest"
)

type ActionService struct {
	snapshotSDK snapshotSDK

	proposalGetter       proposalGetter
	typedSignDataBuilder *TypedSignDataBuilder
	preparedVoteStorage  preparedVoteStorage
}

func NewActionService(snapshotSDK snapshotSDK, proposalGetter proposalGetter, typedSignDataBuilder *TypedSignDataBuilder, preparedVoteStorage preparedVoteStorage) *ActionService {
	return &ActionService{snapshotSDK: snapshotSDK, proposalGetter: proposalGetter, typedSignDataBuilder: typedSignDataBuilder, preparedVoteStorage: preparedVoteStorage}
}

func (a *ActionService) Validate(validateParams ValidateParams) (ValidateResult, error) {
	log.Info().Fields(map[string]any{
		"req": validateParams,
	}).Msg("got validation request")

	proposal, err := a.proposalGetter.GetByID(validateParams.Proposal)
	if err != nil {
		return ValidateResult{}, fmt.Errorf("failed to get proposal: %w", err)
	}

	var pFragment client.ProposalFragment
	if err := json.Unmarshal(proposal.Snapshot, &pFragment); err != nil {
		return ValidateResult{}, fmt.Errorf("failed to unmarshal proposal snapshot: %w", err)
	}

	vResult, err := a.validateVote(validateParams, &pFragment)
	if err != nil {
		return ValidateResult{}, fmt.Errorf("failed to validate vote: %w", err)
	}
	if !vResult.ok {
		return vResult.toValidateResult(), nil
	}

	vpResult, err := a.validateVotingPower(validateParams, &pFragment)
	if err != nil {
		return ValidateResult{}, fmt.Errorf("failed to validate voting power: %w", err)
	}

	return vpResult.toValidateResult(), nil
}

func (a *ActionService) Prepare(prepareParams PrepareParams) (db.PreparedVote, error) {
	proposal, err := a.proposalGetter.GetByID(prepareParams.Proposal)
	if err != nil {
		return db.PreparedVote{}, fmt.Errorf("failed to get proposal: %w", err)
	}

	var pFragment client.ProposalFragment
	if err := json.Unmarshal(proposal.Snapshot, &pFragment); err != nil {
		return db.PreparedVote{}, fmt.Errorf("failed to unmarshal proposal snapshot: %w", err)
	}

	checksumVoter := common.HexToAddress(prepareParams.Voter).Hex()
	typedData := a.typedSignDataBuilder.Build(checksumVoter, prepareParams.Reason, prepareParams.Choice, &pFragment)
	typedDataJSON, err := json.Marshal(typedData)
	if err != nil {
		return db.PreparedVote{}, fmt.Errorf("failed to marshal typed data: %w", err)
	}

	preparedVote := db.PreparedVote{
		Voter:     checksumVoter,
		Proposal:  pFragment.ID,
		TypedData: string(typedDataJSON),
	}
	err = a.preparedVoteStorage.Create(&preparedVote)
	if err != nil {
		return db.PreparedVote{}, fmt.Errorf("failed to create prepared vote: %w", err)
	}

	return preparedVote, nil
}

func (a *ActionService) Vote(voteParams VoteParams) (SuccessVote, error) {
	preparedVote, err := a.preparedVoteStorage.Get(voteParams.ID)
	if err != nil {
		return SuccessVote{}, fmt.Errorf("failed to get prepared vote: %w", err)
	}

	var typedData TypedData
	if err := json.Unmarshal([]byte(preparedVote.TypedData), &typedData); err != nil {
		return SuccessVote{}, fmt.Errorf("failed to unmarshal prepared vote typed data: %w", err)
	}

	voteResult, err := a.snapshotSDK.Vote(context.Background(), snapshot.VoteParams{
		Address: preparedVote.Voter,
		Sig:     voteParams.Sig,
		Data:    typedData,
	})
	if err != nil {
		return SuccessVote{}, fmt.Errorf("failed to vote: %w", err)
	}

	return SuccessVote{
		ID:   voteResult.ID,
		IPFS: voteResult.IPFS,
		Relayer: Relayer{
			Address: voteResult.Relayer.Address,
			Receipt: voteResult.Relayer.Receipt,
		},
	}, nil
}

func (a *ActionService) validateVotingPower(validateParams ValidateParams, pFragment *client.ProposalFragment) (validateVPResult, error) {
	snapshotNum := getSnapshot(pFragment.Snapshot)
	params := snapshot.GetVotingPowerParams{
		Address:    validateParams.Voter,
		Network:    pFragment.Network,
		Strategies: convertStrategies(pFragment.Strategies),
		Snapshot:   snapshotNum,
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
				ok:              false,
				ValidationError: NoVotingPowerErr.WithVars(fmt.Sprintf("%v", snapshotNum)),
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
	if pFragment.Validation == nil {
		return validateVoteResult{
			validation: validation{ok: true},
		}, nil
	}

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

	return a.resolveValidationResult(validate, pFragment)
}

func (a *ActionService) resolveValidationResult(validate snapshot.ValidationResponse, pFragment *client.ProposalFragment) (validateVoteResult, error) {
	if validate.Result {
		return validateVoteResult{
			validation: validation{ok: true},
		}, nil
	}

	if pFragment.Validation.Name == "basic" {
		return validateVoteResult{
			validation: validation{
				ok:              false,
				ValidationError: BasicValidationErr,
			},
		}, nil
	}

	if pFragment.Validation.Name == "passport-gated" {
		return validateVoteResult{
			validation: validation{
				ok:              false,
				ValidationError: PassportGatedValidationErr.WithVars(getPassportGatedParams(pFragment)),
			},
		}, nil
	}

	return validateVoteResult{
		validation: validation{
			ok:              false,
			ValidationError: NotValidErr,
		},
	}, nil
}

func getPassportGatedParams(fragment *client.ProposalFragment) (string, string, string) {
	const (
		defaultOperator       = "one"
		defaultStamps         = ""
		defaultScoreThreshold = "0"
	)

	if fragment.Validation == nil {
		return defaultOperator, defaultStamps, defaultScoreThreshold
	}

	if fragment.Validation.Params == nil {
		return defaultOperator, defaultStamps, defaultScoreThreshold
	}

	operator, ok := fragment.Validation.Params["operator"].(string)
	if !ok {
		operator = defaultOperator
	}

	stamps, ok := fragment.Validation.Params["stamps"].([]string)
	if !ok {
		stamps = []string{}
	}

	scoreThresholdResult := defaultScoreThreshold
	scoreThreshold, ok := fragment.Validation.Params["scoreThreshold"].(float64)
	if ok {
		scoreThresholdResult = fmt.Sprintf("%.2f", scoreThreshold)
	}

	return operator, strings.Join(stamps, ", "), scoreThresholdResult
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
