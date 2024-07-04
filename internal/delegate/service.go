package delegate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetDelegates(req GetDelegatesRequest) ([]Delegate, error) {
	url := fmt.Sprintf("https://delegate-api.gnosisguild.org/api/v1/%s/delegates", req.Dao)
	queryParams := map[string]string{
		"by":     req.By,
		"limit":  strconv.Itoa(req.Limit),
		"offset": strconv.Itoa(req.Offset),
	}

	jsonResp, err := s.doPostHttpRequest(url, queryParams, GnosisTopDelegatesBodyRequest{
		Strategy: req.Strategy,
	})
	if err != nil {
		return nil, err
	}

	var resp GnosisTopDelegatesResponse
	if err := json.Unmarshal(jsonResp, &resp); err != nil {
		return nil, err
	}

	delegates := make([]Delegate, 0, len(resp.Delegates))
	for _, d := range resp.Delegates {
		delegates = append(delegates, Delegate{
			Address:              d.Address,
			DelegatorCount:       d.DelegatorCount,
			PercentOfDelegators:  d.PercentOfDelegators,
			VotingPower:          d.VotingPower,
			PercentOfVotingPower: d.PercentOfVotingPower,
		})
	}

	return delegates, nil
}

func (s *Service) doPostHttpRequest(url string, queryParams map[string]string, bodyParams any) ([]byte, error) {
	reqParams, err := json.Marshal(bodyParams)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest("POST", url, bytes.NewBuffer(reqParams))
	if err != nil {
		return nil, err
	}

	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")

	q := r.URL.Query()
	for k, v := range queryParams {
		q.Add(k, v)
	}
	r.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed: %s", string(content))
	}

	return content, nil
}
