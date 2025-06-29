package gname

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	gojoin "github.com/yuchenfw/go-join"
)

// libdnsZoneToDnslaDomain Strips the trailing dot from a Zone
func libdnsZoneToDnslaDomain(zone string) string {
	return strings.TrimSuffix(zone, ".")
}

const ApiBase = "https://api.gname.com"

// MakeApiRequest makes an API request using the default HTTP client
func MakeApiRequest[T any](method string, endpoint string, params string, appKey string, responseType T) (T, error) {
	return MakeApiRequestWithClient(http.DefaultClient, method, endpoint, params, appKey, responseType)
}

// MakeApiRequestWithClient makes an API request using a custom HTTP client
func MakeApiRequestWithClient[T any](client *http.Client, method string, endpoint string, params string, appKey string, responseType T) (T, error) {
	shangHaiLoc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return responseType, fmt.Errorf("failed to load timezone: %w", err)
	}
	shangHaiNowTime := time.Now().In(shangHaiLoc)

	newEndpoint := fmt.Sprintf("%s?%s&gntime=%d", endpoint, params, shangHaiNowTime.Unix())

	sortedParams, err := gojoin.Join(ApiBase+newEndpoint, gojoin.Options{
		Sep:       "&",
		KVSep:     "=",
		Order:     gojoin.ASCII,
		URLCoding: gojoin.Encoding,
	})
	if err != nil {
		return responseType, fmt.Errorf("failed to join parameters: %w", err)
	}

	signParams := sortedParams + appKey

	token := md5.Sum([]byte(signParams))
	gnToken := strings.ToUpper(fmt.Sprintf("%x", token))

	bodyStr := fmt.Sprintf("%s&gntoken=%s", sortedParams, gnToken)

	fullUrl := fmt.Sprintf("%s%s", ApiBase, endpoint)

	u, err := url.Parse(fullUrl)
	if err != nil {
		return responseType, fmt.Errorf("failed to parse URL: %w", err)
	}

	req, err := http.NewRequest(method, u.String(), strings.NewReader(bodyStr))
	if err != nil {
		return responseType, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return responseType, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Warning: failed to close response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return responseType, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return responseType, fmt.Errorf("failed to read response body: %w", err)
	}

	commonResponse := CommonResponse{}
	err = json.Unmarshal(result, &commonResponse)
	if err != nil {
		return responseType, fmt.Errorf("failed to parse common response: %w", err)
	}

	if commonResponse.Code != 1 {
		return responseType, fmt.Errorf("API error: %s (code: %d)", commonResponse.Msg, commonResponse.Code)
	}

	response := responseType
	err = json.Unmarshal(result, &response)
	if err != nil {
		return responseType, fmt.Errorf("failed to parse response: %w", err)
	}

	return response, nil
}
