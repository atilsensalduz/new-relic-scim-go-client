package newrelicscim

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Client is a struct for interacting with the New Relic SCIM API.
//
// It has the following fields:
//  - BaseUrl: the base URL for the SCIM API, including the version number
//  - ApiToken: the API token for authenticating with the SCIM API
//  - HttpClient: an HTTP client with a timeout of 20 seconds, used for making requests to the SCIM API
type Client struct {
	BaseUrl    string
	ApiToken   string
	HttpClient *http.Client
}

// NewClient generates a new NewRelicSCIMClient for interacting with the New Relic SCIM API.
//
// It takes in an API token for authentication and returns a pointer to a new Client struct. The Client struct
// contains the following fields:
//  - BaseUrl: the base URL for the SCIM API, including the version number
//  - ApiToken: the API token for authenticating with the SCIM API
//  - HttpClient: an HTTP client with a timeout of 20 seconds, used for making requests to the SCIM API
//
// The client can be used to make requests to the SCIM API, such as retrieving or updating user information.
func NewClient(apiToken string) *Client {
	h := &http.Client{
		Timeout: 20 * time.Second,
	}

	return &Client{
		BaseUrl:    "https://scim-provisioning.service.newrelic.com/scim/v2/",
		ApiToken:   apiToken,
		HttpClient: h,
	}
}

// doRequest is a helper function that sends an HTTP request and returns the response body as a slice of bytes.
//
// It takes in a pointer to an HTTP request and adds the necessary headers for authenticating with the New Relic SCIM API
// using the client's API token. The function then makes the request and reads the response body into a slice of bytes.
// If the request or response encounters an error or the response status code is not in the 2xx range, an error is returned.
// Otherwise, the response body is returned as a slice of bytes.
func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", "Bearer "+c.ApiToken)
	req.Header.Set("content-type", "application/json")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if !((resp.StatusCode >= 200) && (resp.StatusCode <= 299)) {
		return nil, fmt.Errorf("error body: %s\nstatus Code: %d", body, resp.StatusCode)
	}

	return body, nil
}
