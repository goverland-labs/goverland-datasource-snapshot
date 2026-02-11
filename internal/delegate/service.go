package delegate

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/goverland-labs/goverland-datasource-snapshot/internal/helpers"
	"github.com/goverland-labs/goverland-datasource-snapshot/pkg/gnosis"
)

type Service struct {
	gnosisSDK   *gnosis.SDK
	snapshotSDK snapshotSDK
}

func NewService(gnosisSDK *gnosis.SDK, snapshotSDK snapshotSDK) *Service {
	return &Service{
		gnosisSDK:   gnosisSDK,
		snapshotSDK: snapshotSDK,
	}
}

func (s *Service) GetDelegates(ctx context.Context, req GetDelegatesParams) (DelegatesWrapper, error) {
	if len(req.Addresses) == 1 {
		return s.searchDelegateProfile(ctx, req)
	}

	topDelegatesReq := gnosis.TopDelegatesRequest{
		Dao:      req.Dao,
		Strategy: req.Strategy,
		By:       helpers.ValurOrDefault(req.By, "power"),
		Limit:    req.Limit,
		Offset:   req.Offset,
	}

	topDelegatesResp, err := s.gnosisSDK.GetTopDelegates(ctx, topDelegatesReq)
	if err != nil {
		return DelegatesWrapper{}, err
	}

	delegates := make([]Delegate, 0, len(topDelegatesResp.Delegates))
	for _, d := range topDelegatesResp.Delegates {
		delegates = append(delegates, Delegate{
			Address:              d.Address,
			DelegatorCount:       d.DelegatorCount,
			PercentOfDelegators:  basisPointToPercentage(d.PercentOfDelegators),
			VotingPower:          d.VotingPower,
			PercentOfVotingPower: basisPointToPercentage(d.PercentOfVotingPower),
		})
	}

	return DelegatesWrapper{
		Delegates: delegates,
		Total:     topDelegatesResp.Pagination.Total,
	}, nil
}

func (s *Service) GetDelegateProfile(ctx context.Context, req GetDelegateProfileParams) (DelegateProfile, error) {
	delegateProfileReq := gnosis.DelegateProfileRequest{
		Dao:      req.Dao,
		Strategy: req.Strategy,
		Address:  req.Address,
	}

	delegateProfileResp, err := s.gnosisSDK.GetDelegateProfile(ctx, delegateProfileReq)
	if err != nil {
		return DelegateProfile{}, fmt.Errorf("failed to get delegation profile: %w", err)
	}

	var expiration *time.Time
	delegates := make([]ProfileDelegateItem, 0, len(delegateProfileResp.DelegateTree))
	for _, d := range delegateProfileResp.DelegateTree {
		delegates = append(delegates, ProfileDelegateItem{
			Address:        d.Delegate,
			Weight:         d.Weight,
			DelegatedPower: d.DelegatedPower,
		})

		if d.ExpirationUnixTime != 0 {
			expiration = helpers.Ptr(time.Unix(int64(d.ExpirationUnixTime), 0))
		}
	}

	profile := DelegateProfile{
		Address:              delegateProfileResp.Address,
		VotingPower:          delegateProfileResp.VotingPower,
		IncomingPower:        delegateProfileResp.IncomingPower,
		OutgoingPower:        delegateProfileResp.OutgoingPower,
		PercentOfVotingPower: basisPointToPercentage(delegateProfileResp.PercentOfVotingPower),
		PercentOfDelegators:  basisPointToPercentage(delegateProfileResp.PercentOfDelegators),
		Delegates:            delegates,
		Expiration:           expiration,
	}

	return profile, nil
}

func (s *Service) searchDelegateProfile(ctx context.Context, req GetDelegatesParams) (DelegatesWrapper, error) {
	delegateProfileReq := gnosis.DelegateProfileRequest{
		Dao:      req.Dao,
		Strategy: req.Strategy,
		Address:  req.Addresses[0],
	}

	delegateProfileResp, err := s.gnosisSDK.GetDelegateProfile(ctx, delegateProfileReq)
	if err != nil {
		return DelegatesWrapper{}, fmt.Errorf("failed to get delegation profile: %w", err)
	}

	delegates := []Delegate{
		{
			Address:              delegateProfileResp.Address,
			DelegatorCount:       int32(len(delegateProfileResp.Delegators)),
			PercentOfDelegators:  basisPointToPercentage(delegateProfileResp.PercentOfDelegators),
			VotingPower:          delegateProfileResp.VotingPower,
			PercentOfVotingPower: basisPointToPercentage(delegateProfileResp.PercentOfVotingPower),
		},
	}

	return DelegatesWrapper{
		Delegates: delegates,
		Total:     1,
	}, nil
}

func (s *Service) GetDelegatesStatement(ctx context.Context, spaceID string, addresses []string) ([]Statement, error) {
	converted := make([]string, 0, len(addresses))
	for _, addr := range addresses {
		converted = append(converted, strings.ToLower(addr))
	}

	list, err := s.snapshotSDK.ListStatements(ctx, spaceID, converted)
	if err != nil {
		return nil, fmt.Errorf("s.snapshotSDK.ListStatements: %w", err)
	}

	statements := make([]Statement, 0, len(addresses))
	for _, info := range list {
		var address, statement string

		if val := info.GetDelegate(); val != nil {
			address = strings.ToLower(*val)
		}
		if val := info.GetStatement(); val != nil {
			statement = *val
		}

		if address == "" {
			continue
		}

		statements = append(statements, Statement{
			Address:   address,
			Statement: statement,
		})
	}

	return statements, nil
}
