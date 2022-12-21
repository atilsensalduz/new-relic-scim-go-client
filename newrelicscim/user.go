package newrelicscim

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const userPath = "Users"

type User struct {
	Schemas  []string `json:"schemas"`
	UserName string   `json:"userName"`
	Name     struct {
		FamilyName string `json:"familyName"`
		GivenName  string `json:"givenName"`
	} `json:"name"`
	Emails []struct {
		Primary bool   `json:"primary"`
		Value   string `json:"value"`
	} `json:"emails"`
	Active   bool   `json:"active"`
	Timezone string `json:"timezone"`
}

func (u *User) fill_defaults() {

	// setting default values
	// if no values present
	if len(u.Schemas) == 0 {
		u.Schemas = []string{"urn:ietf:params:scim:schemas:core:2.0:User"}
	}
	if u.Timezone == "" {
		u.Timezone = "Europe/Istanbul"
	}
}

type UserResponse struct {
	Schemas    []string `json:"schemas"`
	ID         string   `json:"id"`
	ExternalID string   `json:"externalId"`
	UserName   string   `json:"userName"`
	Name       struct {
		FamilyName string `json:"familyName"`
		GivenName  string `json:"givenName"`
	} `json:"name"`
	Emails []struct {
		Value   string `json:"value"`
		Primary bool   `json:"primary"`
	} `json:"emails"`
	Timezone string `json:"timezone"`
	Active   bool   `json:"active"`
	Meta     struct {
		ResourceType string    `json:"resourceType"`
		Created      time.Time `json:"created"`
		LastModified time.Time `json:"lastModified"`
	} `json:"meta"`
	Groups []interface{} `json:"groups"`
}

type UserErrorResponse struct {
	Schemas  []string `json:"schemas"`
	ScimType string   `json:"scimType"`
	Detail   string   `json:"detail"`
	Status   string   `json:"status"`
}

type UsersResponse struct {
	TotalResults int      `json:"totalResults"`
	Schemas      []string `json:"schemas"`
	Resources    []struct {
		Schemas    []string    `json:"schemas"`
		ID         string      `json:"id"`
		ExternalID interface{} `json:"externalId"`
		UserName   string      `json:"userName"`
		Name       struct {
			FamilyName string `json:"familyName"`
			GivenName  string `json:"givenName"`
		} `json:"name"`
		Emails []struct {
			Value   string `json:"value"`
			Primary bool   `json:"primary"`
		} `json:"emails"`
		Timezone string `json:"timezone"`
		Active   bool   `json:"active"`
		Meta     struct {
			ResourceType string    `json:"resourceType"`
			Created      time.Time `json:"created"`
			LastModified time.Time `json:"lastModified"`
		} `json:"meta"`
		Groups []struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"groups"`
	} `json:"Resources"`
}

type UserTypeBody struct {
	Schemas                                         []string `json:"schemas"`
	UrnIetfParamsScimSchemasExtensionNewrelic21User struct {
		NrUserType string `json:"nrUserType"`
	} `json:"urn:ietf:params:scim:schemas:extension:newrelic:2.1:User"`
}

func (u *UserTypeBody) fill_defaults() {

	// setting default values
	// if no values present
	if len(u.Schemas) == 0 {
		u.Schemas = []string{"urn:ietf:params:scim:schemas:core:2.0:User", "urn:ietf:params:scim:schemas:extension:newrelic:2.0:User"}
	}
}

func (c *Client) UserList(ctx context.Context) (usersResponse UsersResponse, userErrorResponse UserErrorResponse, err error) {
	fullUrl := fmt.Sprintf("%s%s", c.BaseUrl, userPath)
	req, err := http.NewRequest("Get", fullUrl, nil)
	if err != nil {
		return usersResponse, userErrorResponse, err
	}
	resp, err := c.doRequest(req)
	if err != nil {
		return usersResponse, userErrorResponse, err
	}
	if err := json.Unmarshal(resp, &usersResponse); err != nil {
		return usersResponse, userErrorResponse, err
	}
	if usersResponse.Schemas[0] == "urn:ietf:params:scim:api:messages:2.0:Error" {
		if err := json.Unmarshal(resp, &userErrorResponse); err != nil {
			return usersResponse, userErrorResponse, err
		}

	}

	return usersResponse, userErrorResponse, nil
}

func (c *Client) GetUserByID(ctx context.Context, userID string) (userResponse UserResponse, userErrorResponse UserErrorResponse, err error) {
	fullUrl := fmt.Sprintf("%s%s/%s", c.BaseUrl, userPath, userID)
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return userResponse, userErrorResponse, err
	}
	resp, err := c.doRequest(req)
	if err != nil {
		return userResponse, userErrorResponse, err
	}
	if err := json.Unmarshal(resp, &userResponse); err != nil {
		return userResponse, userErrorResponse, err
	}
	if userResponse.Schemas[0] == "urn:ietf:params:scim:api:messages:2.0:Error" {
		if err := json.Unmarshal(resp, &userErrorResponse); err != nil {
			return userResponse, userErrorResponse, err
		}

	}

	return userResponse, userErrorResponse, nil
}

func (c *Client) GetUserByName(ctx context.Context, userName string) (userResponse UserResponse, userErrorResponse UserErrorResponse, err error) {

	fullUrl := fmt.Sprintf("%s%s", c.BaseUrl, userPath)

	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return userResponse, userErrorResponse, err
	}
	q := req.URL.Query()
	filter := fmt.Sprintf(`userName eq "%s"`, userName)
	fmt.Println(filter)
	q.Add("filter", filter)
	req.URL.RawQuery = q.Encode()

	fmt.Println(req.URL.String())

	resp, err := c.doRequest(req)
	if err != nil {
		return userResponse, userErrorResponse, err
	}
	if err := json.Unmarshal(resp, &userResponse); err != nil {
		return userResponse, userErrorResponse, err
	}

	if userResponse.Schemas[0] == "urn:ietf:params:scim:api:messages:2.0:Error" {
		if err := json.Unmarshal(resp, &userErrorResponse); err != nil {
			return userResponse, userErrorResponse, err
		}

	}

	return userResponse, userErrorResponse, nil
}

func (c *Client) CreateUser(ctx context.Context, user User) (userResponse UserResponse, userErrorResponse UserErrorResponse, err error) {

	fullUrl := fmt.Sprintf("%s%s", c.BaseUrl, userPath)
	user.fill_defaults()
	//Encode the data
	postBody, _ := json.Marshal(user)
	responseBody := bytes.NewBuffer(postBody)

	req, err := http.NewRequest("POST", fullUrl, responseBody)
	if err != nil {
		return userResponse, userErrorResponse, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return userResponse, userErrorResponse, err
	}
	if err := json.Unmarshal(resp, &userResponse); err != nil {
		return userResponse, userErrorResponse, err
	}
	if userResponse.Schemas[0] == "urn:ietf:params:scim:api:messages:2.0:Error" {
		if err := json.Unmarshal(resp, &userErrorResponse); err != nil {
			return userResponse, userErrorResponse, err
		}

	}

	return userResponse, userErrorResponse, nil
}

func (c *Client) UpdateUser(ctx context.Context, userID string, user User) (userResponse UserResponse, userErrorResponse UserErrorResponse, err error) {

	fullUrl := fmt.Sprintf("%s%s/%s", c.BaseUrl, userPath, userID)
	//Encode the data
	user.fill_defaults()
	postBody, _ := json.Marshal(user)
	responseBody := bytes.NewBuffer(postBody)

	req, err := http.NewRequest("PUT", fullUrl, responseBody)
	if err != nil {
		return userResponse, userErrorResponse, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return userResponse, userErrorResponse, err
	}
	if err := json.Unmarshal(resp, &userResponse); err != nil {
		return userResponse, userErrorResponse, err
	}
	if userResponse.Schemas[0] == "urn:ietf:params:scim:api:messages:2.0:Error" {
		if err := json.Unmarshal(resp, &userErrorResponse); err != nil {
			return userResponse, userErrorResponse, err
		}

	}

	return userResponse, userErrorResponse, nil
}

func (c *Client) DeleteUser(ctx context.Context, userID string) (err error) {

	fullUrl := fmt.Sprintf("%s%s/%s", c.BaseUrl, userPath, userID)

	req, err := http.NewRequest("DELETE", fullUrl, nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}
	return nil
}

type UserType int64

const (
	Full UserType = iota
	Core
	Basic
)

func (u UserType) String() string {
	switch u {
	case Full:
		return "Full User"
	case Core:
		return "Core User"
	case Basic:
		return "Basic User"
	}
	return "unknown"
}

func (c *Client) ChangeUserType(ctx context.Context, userID string, userType UserType) (userResponse UserResponse, userErrorResponse UserErrorResponse, err error) {

	fullUrl := fmt.Sprintf("%s%s/%s", c.BaseUrl, userPath, userID)
	userTypeBody := UserTypeBody{
		UrnIetfParamsScimSchemasExtensionNewrelic21User: struct {
			NrUserType string "json:\"nrUserType\""
		}{NrUserType: userType.String()},
	}
	//Encode the data
	userTypeBody.fill_defaults()
	putBody, _ := json.Marshal(userTypeBody)
	responseBody := bytes.NewBuffer(putBody)

	req, err := http.NewRequest("PUT", fullUrl, responseBody)
	if err != nil {
		return userResponse, userErrorResponse, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return userResponse, userErrorResponse, err
	}
	if err := json.Unmarshal(resp, &userResponse); err != nil {
		return userResponse, userErrorResponse, err
	}
	if userResponse.Schemas[0] == "urn:ietf:params:scim:api:messages:2.0:Error" {
		if err := json.Unmarshal(resp, &userErrorResponse); err != nil {
			return userResponse, userErrorResponse, err
		}

	}

	return userResponse, userErrorResponse, nil
}
