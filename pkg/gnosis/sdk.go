package gnosis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

const (
	baseURL = "https://delegate-api.gnosisguild.org/api/v1"
)

type SDK struct {
	commonHTTP *commonHTTP
}

func NewSDK() *SDK {
	return &SDK{
		commonHTTP: &commonHTTP{},
	}
}

func (s *SDK) GetTopDelegates(ctx context.Context, req TopDelegatesRequest) (TopDelegatesResponse, error) {
	url := fmt.Sprintf("%s/%s/pin/top-delegates", baseURL, req.Dao)
	queryParams := map[string]string{
		"by":     req.By,
		"limit":  strconv.Itoa(req.Limit),
		"offset": strconv.Itoa(req.Offset),
	}

	jsonResp, err := s.commonHTTP.doPostHttpRequest(ctx, url, queryParams, strategyBodyRequest{
		Strategy: req.Strategy,
	})
	if err != nil {
		return TopDelegatesResponse{}, err
	}

	var resp TopDelegatesResponse
	if err := json.Unmarshal(jsonResp, &resp); err != nil {
		return TopDelegatesResponse{}, err
	}

	return resp, nil
}

func (s *SDK) GetDelegateProfile(ctx context.Context, req DelegateProfileRequest) (DelegateProfileResponse, error) {
	url := fmt.Sprintf("%s/%s/pin/%s", baseURL, req.Dao, req.Address)

	jsonResp, err := s.commonHTTP.doPostHttpRequest(ctx, url, map[string]string{}, strategyBodyRequest{
		Strategy: req.Strategy,
	})
	if err != nil {
		return DelegateProfileResponse{}, err
	}

	var resp DelegateProfileResponse
	if err := json.Unmarshal(jsonResp, &resp); err != nil {
		return DelegateProfileResponse{}, err
	}

	return resp, nil
}
