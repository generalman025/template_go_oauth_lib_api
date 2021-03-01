package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/generalman025/template_go_util_lib_api/rest_errors"
	"github.com/go-resty/resty/v2"
)

const (
	// baseURL          = "http://localhost:8080" // for localhost
	baseURL          = "http://host.docker.internal:8080" // for docker
	headerXPublic    = "X-Public"
	headerXClientId  = "X-Client-Id"
	headerXCallerId  = "X-Caller-Id"
	paramAccessToken = "access_token"
)

var (
	// oauthRestClient = rest.RequestBuilder{
	// 	BaseURL: "http://localhost:8080",
	// 	Timeout: 200 * time.Millisecond,
	// }

	oauthRestClient = resty.New().R().EnableTrace()
)

type oauthClient struct {
}

type oauthInterface interface {
}

type accessToken struct {
	Id       string `json:"id"`
	UserId   int64  `json:"user_id"`
	ClientId int64  `json:"client_id"`
}

func init() {

}

func IsPublic(request *http.Request) bool {
	if request == nil {
		return true
	}
	return request.Header.Get(headerXPublic) == "true"
}

func GetCallerId(request *http.Request) int64 {
	if request == nil {
		return 0
	}
	callerId, err := strconv.ParseInt(request.Header.Get(headerXCallerId), 10, 64)
	if err != nil {
		return 0
	}

	return callerId
}

func GetClientId(request *http.Request) int64 {
	if request == nil {
		return 0
	}
	clientId, err := strconv.ParseInt(request.Header.Get(headerXClientId), 10, 64)
	if err != nil {
		return 0
	}

	return clientId
}

func AuthenticateRequest(request *http.Request) rest_errors.RestErr {
	if request == nil {
		return nil
	}

	cleanRequest(request)

	accessTokenId := strings.TrimSpace(request.URL.Query().Get(paramAccessToken))
	if accessTokenId == "" {
		return nil
	}

	at, err := getAccessToken(accessTokenId)
	if err != nil {
		if err.Status() == http.StatusNotFound {
			return nil
		}
		return err
	}

	request.Header.Add(headerXClientId, fmt.Sprintf("%v", at.ClientId))
	request.Header.Add(headerXCallerId, fmt.Sprintf("%v", at.UserId))

	return nil
}

func cleanRequest(request *http.Request) {
	if request == nil {
		return
	}

	request.Header.Del(headerXClientId)
	request.Header.Del(headerXCallerId)
}

func getAccessToken(accessTokenId string) (*accessToken, rest_errors.RestErr) {
	response, _ := oauthRestClient.Get(fmt.Sprintf("%s/oauth/access_token/%s", baseURL, accessTokenId))
	if response == nil {
		return nil, rest_errors.NewInternalServerError("invalid resclient response when trying to get access token", rest_errors.NewError("error when trying to get access token"))
	}
	if response.StatusCode() > 299 {
		var restErr rest_errors.RestErr
		if err := json.Unmarshal(response.Body(), restErr); err != nil {
			return nil, rest_errors.NewInternalServerError("invalid error interface when trying to get access token", err)
		}
		if response.StatusCode() == 404 {
			return nil, restErr
		}
		return nil, restErr
	}

	var at accessToken
	if err := json.Unmarshal(response.Body(), &at); err != nil {
		return nil, rest_errors.NewInternalServerError("error when trying to unmarshal access token response", err)
	}
	return &at, nil
}

// func getAccessToken(accessTokenId string) (*accessToken, rest_errors.RestErr) {
// 	response := oauthRestClient.Get(fmt.Sprintf("/oauth/access_token/%s", accessTokenId))
// 	if response == nil || response.Response == nil {
// 		return nil, rest_errors.NewInternalServerError("invalid resclient response when trying to get access token", rest_errors.NewError("error when trying to get access token"))
// 	}
// 	if response.StatusCode > 299 {
// 		var restErr rest_errors.RestErr
// 		if err := json.Unmarshal(response.Bytes(), restErr); err != nil {
// 			return nil, rest_errors.NewInternalServerError("invalid error interface when trying to get access token", err)
// 		}
// 		if response.StatusCode == 404 {
// 			return nil, restErr
// 		}
// 		return nil, restErr
// 	}

// 	var at accessToken
// 	if err := json.Unmarshal(response.Bytes(), &at); err != nil {
// 		return nil, rest_errors.NewInternalServerError("error when trying to unmarshal access token response", err)
// 	}
// 	return &at, nil
// }
