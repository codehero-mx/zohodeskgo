package zohodeskgo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/go-resty/resty/v2"
)

type zoho struct {
	clientID        string
	clientSecret    string
	refreshToken    string
	pathAccessToken string
	organizationID  uint64
}

// struct to general response error
type respError struct {
	ErrorCode string `json:"errorCode"`
	Message   string `json:"message"`
}

// struct to capture the OAuth token from zoho rest API
type oAuthToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// struct to oauth error
type oAuthError struct {
	Error string `json:"error"`
}

const (
	hostURL    = "https://desk.zoho.com/api/v1"
	oauthURL   = "https://accounts.zoho.com/oauth/v2/token"
	maxRetries = 2
)

var retries = 1

func (z *zoho) doRefreshToken() (err error) {
	resp, err := resty.New().
		// SetDebug(true).
		R().
		SetQueryParams(map[string]string{
			"refresh_token": z.refreshToken,
			"client_id":     z.clientID,
			"client_secret": z.clientSecret,
			"scope":         "Desk.tickets.WRITE,Desk.basic.CREATE,Desk.tickets.READ",
			"grant_type":    "refresh_token",
		}).
		Post(oauthURL)
	if err != nil {
		return
	}

	if resp.StatusCode() != 200 {
		oAuthError := new(oAuthError)
		json.Unmarshal(resp.Body(), &oAuthError)

		fmt.Printf("%s\n", oAuthError.Error)

		err = errors.New("an error occurred while refreshing the token")
		return
	}

	oAuthToken := new(oAuthToken)
	json.Unmarshal(resp.Body(), &oAuthToken)

	file, err := os.OpenFile(z.pathAccessToken, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("%s\n", err.Error())

		err = errors.New("an error occurred while refreshing the token")
		return
	}

	defer file.Close()

	_, err = file.WriteString(oAuthToken.AccessToken)
	if err != nil {
		fmt.Printf("%s\n", err.Error())

		err = errors.New("an error occurred while refreshing the token")
		return
	}

	return
}

func (z *zoho) getAccessToken() (string, error) {
	content, err := os.ReadFile(z.pathAccessToken)
	if err != nil {
		fmt.Printf("%s\n", err.Error())

		err = errors.New("an error occurred while obtaining the access token")
		return "", err
	}

	return string(content), nil
}

func (z *zoho) wrapRequestor(reqFn func() (*resty.Response, error)) (*resty.Response, error) {
	var restyResp *resty.Response

	err := resty.Backoff(func() (*resty.Response, error) {
		resp, err := reqFn()
		if err != nil {
			return nil, err
		}

		restyResp = resp

		return resp, nil
	}, resty.RetryConditions([]resty.RetryConditionFunc{func(resp *resty.Response, e error) bool {
		if resp.StatusCode() == 200 {
			retries = 1
			return false
		}

		if retries++; maxRetries < retries {
			return false
		}

		respErr := new(respError)
		json.Unmarshal(resp.Body(), &respErr)

		if respErr.ErrorCode == "INVALID_OAUTH" {
			z.doRefreshToken()
			return true
		}

		return false
	}}))
	if err != nil {
		return nil, err
	}

	return restyResp, nil
}

func (z *zoho) requestor() (*resty.Request, error) {
	accessToken, err := z.getAccessToken()
	if err != nil {
		return nil, err
	}

	client := resty.New().
		// SetDebug(true).
		SetBaseURL(hostURL).
		R().
		SetAuthToken(accessToken).
		SetHeader("orgId", strconv.FormatUint(z.organizationID, 10))

	return client, nil
}
