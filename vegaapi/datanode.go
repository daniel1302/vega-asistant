package vegaapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/daniel1302/vega-assistant/types"
	"github.com/daniel1302/vega-assistant/utils"
	"github.com/hashicorp/go-multierror"
)

const (
	healthyBlocksThreshold = 500
	defaultTimeout         = 5 * time.Second
)

type NetworkAPI struct {
	httpClient *http.Client
	apiREST    []string
}

func NewNetworkAPI(apiREST []string, safeOnly bool, client *http.Client) (*NetworkAPI, error) {
	if len(apiREST) < 1 {
		return nil, fmt.Errorf("at least one api rest endpoint required to create NetworkAPI client")
	}

	if client == nil {
		client = newDefaultHTTPClient()
	}

	safeApiREST := []string{}
	if safeOnly {
		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
		defer cancel()
		latestStatistics, err := getLatestStatistics(ctx, client, apiREST)
		if err != nil {
			return nil, fmt.Errorf("failed to get network statistics for the network head: %w", err)
		}

		for _, endpoint := range apiREST {
			if isRESTEndpointHealthy(ctx, client, latestStatistics.BlockHeight, endpoint) {
				safeApiREST = append(safeApiREST, endpoint)
			}
		}

		if len(safeApiREST) < 1 {
			return nil, fmt.Errorf("not found any healthy endpoint for the network")
		}
		apiREST = safeApiREST
	}

	return &NetworkAPI{
		httpClient: client,
		apiREST:    apiREST,
	}, nil
}

func (n *NetworkAPI) HealthyEndpoints(ctx context.Context, endpoints []types.EndpointWithVegaREST) ([]string, error) {
	latestStatistics, err := getLatestStatistics(ctx, n.httpClient, n.apiREST)
	if err != nil {
		return nil, fmt.Errorf("failed to get network statistics for the network head: %w", err)
	}

	res := []string{}
	for _, endpoint := range endpoints {
		if isRESTEndpointHealthy(ctx, n.httpClient, latestStatistics.BlockHeight, endpoint.REST) {
			res = append(res, endpoint.Endpoint)
		}
	}

	return res, nil
}

func (n *NetworkAPI) Statistics(ctx context.Context) (*types.VegaStatistics, error) {
	var resErr error

	for _, endpoint := range n.apiREST {
		stats, err := getStatistics(ctx, n.httpClient, endpoint)
		if err != nil {
			resErr = multierror.Append(resErr, err)
			continue
		}

		return stats, nil
	}

	return nil, resErr
}

func (n *NetworkAPI) Snapshots(ctx context.Context) (*types.CoreSnapshots, error) {
	if len(n.apiREST) < 1 {
		return nil, fmt.Errorf("failed to get statistics for network: no endpoint available")
	}
	var resErr error
	for _, endpoint := range n.apiREST {
		res, err := n.getSnapshots(ctx, endpoint)
		if err != nil {
			resErr = multierror.Append(resErr, err)
			continue
		}

		return res, nil
	}

	return nil, resErr
}

func (n *NetworkAPI) NetworkHistorySegments(ctx context.Context, networkHight uint64) (*types.NetworkHistorySegments, error) {
	const segmentThreshold = 350

	if len(n.apiREST) < 1 {
		return nil, fmt.Errorf("failed to get statistics for network: no endpoint available")
	}

	// type NetworkHistorySegment struct {
	// 	FromHeight       string `json:"fromHeight"`
	// 	ToHeight         string `json:"toHeight"`
	// 	HistorySegmentId string `json:"historySegmentId"`
	// }

	// type NetworkHistorySegments struct {
	// 	Segments []NetworkHistorySegment `json:"segments"`
	// }

	var resErr error
	for _, endpoint := range n.apiREST {
		res, err := n.getNetworkHistorySegments(ctx, endpoint)

		if err != nil {
			resErr = multierror.Append(resErr, err)
			continue
		}

		// Make sure there is segment close to the current network head block
		foundHeadCloseSegment := false
		for _, segment := range res.Segments {
			if utils.MustUint64(segment.ToHeight) < networkHight-segmentThreshold {
				continue
			}

			foundHeadCloseSegment = true
			break
		}
		if !foundHeadCloseSegment {
			resErr = multierror.Append(resErr, fmt.Errorf("the latest data-node segment for node %s not found", endpoint))
			continue
		}

		return res, nil
	}

	return nil, resErr
}

func getLatestStatistics(ctx context.Context, httpClient *http.Client, restEndpoints []string) (*types.VegaStatistics, error) {
	if len(restEndpoints) < 1 {
		return nil, fmt.Errorf("no rest endpoint passed")
	}

	var latestStatistics *types.VegaStatistics

	for _, endpoint := range restEndpoints {
		statistics, err := utils.RetryReturn(3, 500*time.Millisecond, func() (*types.VegaStatistics, error) {
			return getStatistics(ctx, httpClient, endpoint)
		})

		if err != nil {
			// TODO: Maybe we can think about logging
			continue
		}

		if latestStatistics == nil || latestStatistics.BlockHeight < statistics.BlockHeight {
			latestStatistics = statistics
		}
	}

	if latestStatistics == nil {
		return nil, fmt.Errorf("all endpoints are unhealthy")
	}

	return latestStatistics, nil
}

func isRESTEndpointHealthy(ctx context.Context, httpClient *http.Client, networkHeadHeight uint64, restURL string) bool {
	statistics, err := utils.RetryReturn(3, 500*time.Millisecond, func() (*types.VegaStatistics, error) {
		return getStatistics(ctx, httpClient, restURL)
	})

	if err != nil {
		//	logger.Info(fmt.Sprintf("The %s endpoint unhealthy: failed to get statistics endpoint", restURL), zap.Error(err))
		return false
	}

	headBlocksDiff := networkHeadHeight - statistics.BlockHeight
	if statistics.BlockHeight < networkHeadHeight && headBlocksDiff > healthyBlocksThreshold {
		// logger.Sugar().Infof(
		// 	"The %s endpoint unhealthy: core height(%d) is %d behind the network head(%d), only %d blocks lag allowed",
		// 	restURL,
		// 	statistics.BlockHeight,
		// 	headBlocksDiff,
		// 	networkHeadHeight,
		// 	healthyBlocksThreshold,
		// )
		return false
	}

	if statistics.DataNodeHeight > 0 {
		blocksDiff := statistics.BlockHeight - statistics.DataNodeHeight
		if statistics.DataNodeHeight < statistics.BlockHeight && blocksDiff > healthyBlocksThreshold {
			// logger.Sugar().Infof(
			// 	"The %s endpoint unhealthy: data node is %d blocks behind core, only %d blocks lag allowed",
			// 	restURL,
			// 	blocksDiff,
			// 	healthyBlocksThreshold,
			// )
			return false
		}
	}

	// We do not check time diff here, because we want run test even if the network is not producing blocks.
	// 		It can give us extra information
	// timeDiff := statistics.CurrentTime.Sub(statistics.VegaTime)
	// if timeDiff > HealthyTimeThreshold {
	// 	logger.Sugar().Infof(
	// 		"The %s endpoint unhealthy: time lag is %s, only %s allowed",
	// 		restURL,
	// 		timeDiff.String(),
	// 		HealthyTimeThreshold.String(),
	// 	)
	// 	return false
	// }

	// logger.Sugar().Infof("The %s endpoint is healthy", restURL)

	return true
}

func getStatistics(ctx context.Context, httpClient *http.Client, restURL string) (*types.VegaStatistics, error) {
	statisticsURL := fmt.Sprintf("%s/statistics", strings.TrimRight(restURL, "/"))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, statisticsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for %s: %w", statisticsURL, err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send get query to the statistics endpoint: %w", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read statistics response body: %w", err)
	}

	rawResult := &types.VegaRawStatistics{}
	if err := json.Unmarshal(body, rawResult); err != nil {
		return nil, fmt.Errorf("failed to unmarshal statistics response: %w", err)
	}

	blockHeight, err := strconv.ParseUint(rawResult.Statistics.BlockHeight, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed parse block height: %w", err)
	}

	currentTime, err := time.Parse(time.RFC3339Nano, rawResult.Statistics.CurrentTime)
	if err != nil {
		return nil, fmt.Errorf("failed to parse current time: %w", err)
	}

	vegaTime, err := time.Parse(time.RFC3339Nano, rawResult.Statistics.VegaTime)
	if err != nil {
		return nil, fmt.Errorf("failed to parse vega time: %w", err)
	}

	dataNodeHeight := uint64(0)
	if dataNodeHeightStr := resp.Header.Get("x-block-height"); dataNodeHeightStr != "" {
		dataNodeHeight, err = strconv.ParseUint(dataNodeHeightStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed parse data node block height: %w", err)
		}
	}

	result := &types.VegaStatistics{
		BlockHeight:    blockHeight,
		DataNodeHeight: dataNodeHeight,
		CurrentTime:    currentTime,
		VegaTime:       vegaTime,

		ChainID:    rawResult.Statistics.ChainID,
		AppVersion: rawResult.Statistics.AppVersion,
	}

	return result, nil
}

func newDefaultHTTPClient() *http.Client {
	return http.DefaultClient
}

func (n *NetworkAPI) httpCall(req *http.Request, result any) error {
	res, err := n.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making http request: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"invalid response code: expected %d, got %d",
			http.StatusOK,
			res.StatusCode,
		)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("could not read response body: %w", err)
	}

	if resBody == nil {
		return fmt.Errorf("failed to get http response for the vega api: %w", err)
	}

	if err := json.Unmarshal(resBody, result); err != nil {
		return fmt.Errorf("failed to unmarshal api response: %w", err)
	}

	return nil
}

func (n *NetworkAPI) getSnapshots(ctx context.Context, endpoint string) (*types.CoreSnapshots, error) {
	result := types.CoreSnapshots{}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/v2/snapshots", endpoint), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request snapshots api request for %s: %w", endpoint, err)
	}

	if err := n.httpCall(req, &result); err != nil {
		return nil, fmt.Errorf("failed to get core snapshots: %w", err)
	}

	return &result, nil
}

func (n *NetworkAPI) getNetworkHistorySegments(ctx context.Context, endpoint string) (*types.NetworkHistorySegments, error) {
	result := types.NetworkHistorySegments{}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/v2/networkhistory/segments", endpoint), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request history segments api request for %s: %w", endpoint, err)
	}

	if err := n.httpCall(req, &result); err != nil {
		return nil, fmt.Errorf("failed to get network history segments: %w", err)
	}

	return &result, nil
}
