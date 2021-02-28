package oauth

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	fmt.Println("about to start oauth tests")
	m.Run()
}

func TestOauthConstant(t *testing.T) {
	assert.EqualValues(t, "X-Public", headerXPublic)
	assert.EqualValues(t, "X-Client-Id", headerXClientId)
	assert.EqualValues(t, "X-Caller-Id", headerXCallerId)
	assert.EqualValues(t, "access_token", paramAccessToken)
}

func TestIsPublicNilRequest(t *testing.T) {
	assert.True(t, IsPublic(nil))
}

func TestIsPublicNoError(t *testing.T) {
	request := http.Request{
		Header: make(http.Header),
	}
	assert.False(t, IsPublic(&request))

	request.Header.Add("X-Public", "true")
	assert.True(t, IsPublic(&request))
}

func TestGetCallerIdNilRequest(t *testing.T) {
	assert.EqualValues(t, GetCallerId(nil), 0)
}

func TestGetCallerIdInvalidCallerFormat(t *testing.T) {

}

func TestGetCallerIdNoError(t *testing.T) {
	request := http.Request{
		Header: make(http.Header),
	}
	assert.EqualValues(t, GetCallerId(&request), 0)
}

func TestGetAccessTokenInvalidRestClientResponse(t *testing.T) {
	// Create a Resty Client
	client := resty.New()

	// Get the underlying HTTP Client and set it to Mock
	httpmock.ActivateNonDefault(client.GetClient())
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "http://localhost:8000/oauth/access_token/ABC123",
		httpmock.NewStringResponder(-1, `{}`))

	accessToken, err := getAccessToken("ABC123")
	assert.Nil(t, accessToken)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	assert.EqualValues(t, "error when trying to unmarshal access token response", err.Message())
}
