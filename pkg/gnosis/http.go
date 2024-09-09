package gnosis

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/goverland-labs/goverland-datasource-snapshot/internal/metrics"
)

type commonHTTP struct {
}

func (h *commonHTTP) doPostHttpRequest(ctx context.Context, url string, queryParams map[string]string, bodyParams any) (resp []byte, err error) {
	defer func(start time.Time) {
		metrics.CollectRequestsMetric("gnosis", url, err, start)
	}(time.Now())

	reqParams, err := json.Marshal(bodyParams)
	if err != nil {
		return nil, err
	}

	log.Info().
		Str("url", url).
		Any("queryParams", queryParams).
		Str("bodyParams", string(reqParams)).
		Msg("making http request to gnosis api")

	r, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqParams))
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
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("request failed: %s", string(content))
	}

	return content, nil
}
