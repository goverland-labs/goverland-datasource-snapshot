package metrics

import (
	"net/http"
	"strconv"
)

type HeaderWatcher struct {
	name string
}

func NewHeaderWatcher(name string) *HeaderWatcher {
	return &HeaderWatcher{
		name: name,
	}
}

func (m *HeaderWatcher) RoundTrip(r *http.Request) (*http.Response, error) {
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		return nil, err
	}

	if data, ok := resp.Header["X-Api-Key-Remaining"]; ok && len(data) > 0 {
		if val, err := strconv.ParseFloat(data[0], 64); err == nil {
			CollectSnapshotKeyState(m.name, "remaining_value", val)
		}
	}

	if data, ok := resp.Header["X-Api-Key-Limit"]; ok && len(data) > 0 {
		if val, err := strconv.ParseFloat(data[0], 64); err == nil {
			CollectSnapshotKeyState(m.name, "limit_value", val)
		}
	}

	return resp, nil
}
