package vegaapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/daniel1302/vega-assistant/types"
)

func httpCall(requestURL string, endpoint string) ([]byte, error) {
	apiURL := fmt.Sprintf(
		"%s/%s",
		strings.TrimRight(requestURL, "/"),
		strings.TrimLeft(endpoint, "/"),
	)
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	client := http.DefaultClient
	client.Timeout = 5 * time.Second

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making http request: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"invalid response code for url '%s': expected %d, got %d",
			apiURL,
			http.StatusOK,
			res.StatusCode,
		)
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %w", err)
	}

	return resBody, nil
}

func getRestApiResponse(urls []string, endpoint string, result interface{}) error {
	var (
		errs error
		resp []byte
		err  error
	)
	for _, url := range urls {
		resp, err = httpCall(url, endpoint)
		if err != nil {
			if errs == nil {
				errs = err
			} else {
				errs = errors.Join(errs, err)
			}

			continue
		}
		break
	}

	if resp == nil {
		return fmt.Errorf("failed to get http response for the vega api: %w", errs)
	}

	if err := json.Unmarshal(resp, result); err != nil {
		return fmt.Errorf("failed to unmarshal api response: %w", err)
	}

	return nil
}

func Statistics(urls []string) (*types.VegaStatistics, error) {
	result := types.VegaStatistics{}
	if err := getRestApiResponse(urls, "/statistics", &result); err != nil {
		return nil, fmt.Errorf("failed to get vega statistics: %w", err)
	}

	return &result, nil
}

func Snapshots(urls []string) (*types.CoreSnapshots, error) {
	result := types.CoreSnapshots{}
	if err := getRestApiResponse(urls, "/api/v2/snapshots", &result); err != nil {
		return nil, fmt.Errorf("failed to get core snapshots: %w", err)
	}

	return &result, nil
}

func NetworkHistorySegments(urls []string) (*types.NetworkHistorySegments, error) {
	result := types.NetworkHistorySegments{}
	if err := getRestApiResponse(urls, "/api/v2/networkhistory/segments", &result); err != nil {
		return nil, fmt.Errorf("failed to get network history segments: %w", err)
	}

	return &result, nil
}
