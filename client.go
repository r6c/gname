package gname

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	gojoin "github.com/yuchenfw/go-join"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const ApiBase = "https://api.gname.com"

func MakeApiRequest[T any](method string, endpoint string, params string, appKey string, responseType T) (T, error) {
	client := http.DefaultClient

	shangHaiLoc, err := time.LoadLocation("Asia/Shanghai")
	shangHaiNowTime := time.Now().In(shangHaiLoc)

	newEndpoint := fmt.Sprintf("%s?%s&gntime=%d", endpoint, params, shangHaiNowTime.Unix())

	sortedParams, err := gojoin.Join(ApiBase+newEndpoint, gojoin.Options{
		Sep:       "&",
		KVSep:     "=",
		Order:     gojoin.ASCII,
		URLCoding: gojoin.Encoding,
	})
	if err != nil {
		return responseType, err
	}

	signParams := sortedParams + appKey

	token := md5.Sum([]byte(signParams))
	gnToken := strings.ToUpper(fmt.Sprintf("%x", token))

	bodyStr := fmt.Sprintf("%s&gntoken=%s", sortedParams, gnToken)

	fullUrl := fmt.Sprintf("%s%s", ApiBase, endpoint)

	u, err := url.Parse(fullUrl)
	if err != nil {
		return responseType, err
	}

	req, err := http.NewRequest(method, u.String(), strings.NewReader(bodyStr))
	if err != nil {
		return responseType, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return responseType, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal("Couldn't close body")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		err = errors.New("Invalid http response status, " + string(bodyBytes))
		return responseType, err
	}

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return responseType, err
	}

	commonResponse := CommonResponse{}
	err = json.Unmarshal(result, &commonResponse)
	if err != nil {
		return responseType, err
	}

	if commonResponse.Code != 1 {
		return responseType, errors.New(commonResponse.Msg)
	}

	response := responseType
	err = json.Unmarshal(result, &response)
	if err != nil {
		return responseType, err
	}

	return response, nil
}
