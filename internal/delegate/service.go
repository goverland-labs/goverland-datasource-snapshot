package delegate

import (
	"context"
	"fmt"

	"github.com/goverland-labs/goverland-datasource-snapshot/pkg/gnosis"
)

type Service struct {
	gnosisSDK *gnosis.SDK
}

func NewService(gnosisSDK *gnosis.SDK) *Service {
	return &Service{gnosisSDK: gnosisSDK}
}

func (s *Service) GetDelegates(ctx context.Context, req GetDelegatesParams) ([]Delegate, error) {
	if len(req.Addresses) > 1 {
		return nil, fmt.Errorf("for now only one query address is supported")
	}

	if len(req.Addresses) == 1 {
		return s.searchDelegateProfile(req)
	}

	topDelegatesReq := gnosis.TopDelegatesRequest{
		Dao:      req.Dao,
		Strategy: req.Strategy,
		By:       req.By,
		Limit:    req.Limit,
		Offset:   req.Offset,
	}

	topDelegatesResp, err := s.gnosisSDK.GetTopDelegates(ctx, topDelegatesReq)
	if err != nil {
		return nil, err
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

	return delegates, nil
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

	profile := DelegateProfile{
		Address:              delegateProfileResp.Address,
		VotingPower:          delegateProfileResp.VotingPower,
		IncomingPower:        delegateProfileResp.IncomingPower,
		OutgoingPower:        delegateProfileResp.OutgoingPower,
		PercentOfVotingPower: basisPointToPercentage(delegateProfileResp.PercentOfVotingPower),
		PercentOfDelegators:  basisPointToPercentage(delegateProfileResp.PercentOfDelegators),
	}

	for _, d := range delegateProfileResp.DelegateTree {
		profile.Delegates = append(profile.Delegates, ProfileDelegateItem{
			Address:         d.Delegate,
			PercentOfWeight: basisPointToPercentage(d.Weight),
			DelegatedPower:  d.DelegatedPower,
		})
	}

	return profile, nil
}

func (s *Service) searchDelegateProfile(req GetDelegatesParams) ([]Delegate, error) {
	delegateProfileReq := gnosis.DelegateProfileRequest{
		Dao:      req.Dao,
		Strategy: req.Strategy,
		Address:  req.Addresses[0],
	}

	delegateProfileResp, err := s.gnosisSDK.GetDelegateProfile(nil, delegateProfileReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get delegation profile: %w", err)
	}

	return []Delegate{
		{
			Address:              delegateProfileResp.Address,
			DelegatorCount:       int32(len(delegateProfileResp.Delegators)),
			PercentOfDelegators:  basisPointToPercentage(delegateProfileResp.PercentOfDelegators),
			VotingPower:          delegateProfileResp.VotingPower,
			PercentOfVotingPower: basisPointToPercentage(delegateProfileResp.PercentOfVotingPower),
		},
	}, nil
}
