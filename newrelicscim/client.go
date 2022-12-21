package newrelicscim

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Client struct {
	BaseUrl    string
	ApiToken   string
	HttpClient *http.Client
}

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
