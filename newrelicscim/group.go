package newrelicscim

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const groupPath = "Groups"

type Group struct {
	Schemas     []string `json:"schemas"`
	DisplayName string   `json:"displayName"`
}

type GroupResponse struct {
	Schemas     []string `json:"schemas"`
	ID          string   `json:"id"`
	DisplayName string   `json:"displayName"`
	Meta        struct {
		ResourceType string    `json:"resourceType"`
		Created      time.Time `json:"created"`
		LastModified time.Time `json:"lastModified"`
	} `json:"meta"`
	Members []interface{} `json:"members"`
}

type GroupErrorResponse struct {
	Schemas  []string `json:"schemas"`
	ScimType string   `json:"scimType"`
	Detail   string   `json:"detail"`
	Status   string   `json:"status"`
}

type GroupsResponse struct {
	TotalResults int      `json:"totalResults"`
	Schemas      []string `json:"schemas"`
	Resources    []struct {
		Schemas     []string `json:"schemas"`
		ID          string   `json:"id"`
		DisplayName string   `json:"displayName"`
		Meta        struct {
			ResourceType string    `json:"resourceType"`
			Created      time.Time `json:"created"`
			LastModified time.Time `json:"lastModified"`
		} `json:"meta"`
		Members []interface{} `json:"members"`
	} `json:"Resources"`
}

type UpdateGroup struct {
	Schemas    []string `json:"schemas"`
	Operations []struct {
		Op    string `json:"op"`
		Path  string `json:"path"`
		Value []struct {
			Value string `json:"value"`
		} `json:"value"`
	} `json:"Operations"`
}

func (g *Group) fill_defaults() {

	// setting default values
	// if no values present
	if len(g.Schemas) == 0 {
		g.Schemas = []string{"urn:ietf:params:scim:schemas:core:2.0:Group"}
	}

}

func (ug *UpdateGroup) fill_defaults() {

	// setting default values
	// if no values present
	if len(ug.Schemas) == 0 {
		ug.Schemas = []string{"urn:ietf:params:scim:api:messages:2.0:PatchOp"}
	}

}

func (c *Client) CreateGroup(ctx context.Context, groupName string) (groupResponse GroupResponse, groupErrorResponse GroupErrorResponse, err error) {
	fullUrl := fmt.Sprintf("%s%s", c.BaseUrl, groupPath)
	group := Group{
		DisplayName: groupName,
	}
	group.fill_defaults()

	//Encode the data
	postBody, _ := json.Marshal(group)
	requestBody := bytes.NewBuffer(postBody)

	req, err := http.NewRequest("POST", fullUrl, requestBody)
	if err != nil {
		return groupResponse, groupErrorResponse, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return groupResponse, groupErrorResponse, err
	}
	if err := json.Unmarshal(resp, &groupResponse); err != nil {
		return groupResponse, groupErrorResponse, err
	}
	if groupResponse.Schemas[0] == "urn:ietf:params:scim:api:messages:2.0:Error" {
		if err := json.Unmarshal(resp, &groupErrorResponse); err != nil {
			return groupResponse, groupErrorResponse, err
		}

	}

	return groupResponse, groupErrorResponse, nil
}

func (c *Client) UpdateGroup(ctx context.Context, groupName string) (groupResponse GroupResponse, groupErrorResponse GroupErrorResponse, err error) {
	fullUrl := fmt.Sprintf("%s%s", c.BaseUrl, groupPath)
	group := Group{
		DisplayName: groupName,
	}
	group.fill_defaults()

	//Encode the data
	postBody, _ := json.Marshal(group)
	requestBody := bytes.NewBuffer(postBody)

	req, err := http.NewRequest("PUT", fullUrl, requestBody)
	if err != nil {
		return groupResponse, groupErrorResponse, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return groupResponse, groupErrorResponse, err
	}
	if err := json.Unmarshal(resp, &groupResponse); err != nil {
		return groupResponse, groupErrorResponse, err
	}
	if groupResponse.Schemas[0] == "urn:ietf:params:scim:api:messages:2.0:Error" {
		if err := json.Unmarshal(resp, &groupErrorResponse); err != nil {
			return groupResponse, groupErrorResponse, err
		}

	}

	return groupResponse, groupErrorResponse, nil
}

func (c *Client) GroupList(ctx context.Context) (groupsResponse GroupsResponse, groupErrorResponse GroupErrorResponse, err error) {
	fullUrl := fmt.Sprintf("%s%s", c.BaseUrl, groupPath)
	req, err := http.NewRequest("Get", fullUrl, nil)
	if err != nil {
		return groupsResponse, groupErrorResponse, err
	}
	resp, err := c.doRequest(req)
	if err != nil {
		return groupsResponse, groupErrorResponse, err
	}
	if err := json.Unmarshal(resp, &groupsResponse); err != nil {
		return groupsResponse, groupErrorResponse, err
	}
	if groupsResponse.Schemas[0] == "urn:ietf:params:scim:api:messages:2.0:Error" {
		if err := json.Unmarshal(resp, &groupErrorResponse); err != nil {
			return groupsResponse, groupErrorResponse, err
		}

	}
	return groupsResponse, groupErrorResponse, nil
}

func (c *Client) GetGroupByID(ctx context.Context, groupID string) (groupsResponse GroupsResponse, groupErrorResponse GroupErrorResponse,err error) {

	fullUrl := fmt.Sprintf("%s%s/%s", c.BaseUrl, groupPath, groupID)

	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return groupsResponse, groupErrorResponse, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return groupsResponse, groupErrorResponse, err
	}
	if err := json.Unmarshal(resp, &groupsResponse); err != nil {
		return groupsResponse, groupErrorResponse, err
	}

	if groupsResponse.Schemas[0] == "urn:ietf:params:scim:api:messages:2.0:Error" {
		if err := json.Unmarshal(resp, &groupErrorResponse); err != nil {
			return groupsResponse, groupErrorResponse, err
		}

	}
	return groupsResponse, groupErrorResponse, nil
}

func (c *Client) GetGroupByName(ctx context.Context, groupName string) (groupsResponse GroupsResponse, groupErrorResponse GroupErrorResponse, err error) {

	fullUrl := fmt.Sprintf("%s%s", c.BaseUrl, groupPath)

	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return groupsResponse, groupErrorResponse, err
	}
	q := req.URL.Query()
	filter := fmt.Sprintf(`displayName eq "%s"`, groupName)
	fmt.Println(filter)
	q.Add("filter", filter)
	req.URL.RawQuery = q.Encode()

	fmt.Println(req.URL.String())

	resp, err := c.doRequest(req)
	if err != nil {
		return groupsResponse, groupErrorResponse, err
	}
	if err := json.Unmarshal(resp, &groupsResponse); err != nil {
		return groupsResponse, groupErrorResponse, err
	}

	if groupsResponse.Schemas[0] == "urn:ietf:params:scim:api:messages:2.0:Error" {
		if err := json.Unmarshal(resp, &groupErrorResponse); err != nil {
			return groupsResponse, groupErrorResponse, err
		}

	}
	return groupsResponse, groupErrorResponse, nil}

func (c *Client) GroupMemberOps(ctx context.Context, groupID string, userID string, operation string) (groupResponse GroupResponse, groupErrorResponse GroupErrorResponse, err error) {

	fullUrl := fmt.Sprintf("%s%s/%s", c.BaseUrl, groupPath, groupID)
	//Encode the data
	updateGroup := UpdateGroup{
		Operations: []struct {
			Op    string "json:\"op\""
			Path  string "json:\"path\""
			Value []struct {
				Value string "json:\"value\""
			} "json:\"value\""
		}{
			{Op: operation, Path: "members", Value: []struct {
				Value string "json:\"value\""
			}{{Value: userID}}},
		},
	}
	updateGroup.fill_defaults()

	putBody, _ := json.Marshal(updateGroup)
	requestBody := bytes.NewBuffer(putBody)

	req, err := http.NewRequest("PATCH", fullUrl, requestBody)
	if err != nil {
		return groupResponse, groupErrorResponse, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return groupResponse, groupErrorResponse, err
	}
	if err := json.Unmarshal(resp, &groupResponse); err != nil {
		return groupResponse, groupErrorResponse, err
	}
	if groupResponse.Schemas[0] == "urn:ietf:params:scim:api:messages:2.0:Error" {
		if err := json.Unmarshal(resp, &groupErrorResponse); err != nil {
			return groupResponse, groupErrorResponse, err
		}

	}

	return groupResponse, groupErrorResponse, nil
}

func (c *Client) AddUserToGroup(ctx context.Context, groupID string, userID string) (groupResponse GroupResponse, groupErrorResponse GroupErrorResponse, err error) {
	return c.GroupMemberOps(ctx, groupID, userID, "Add")
}

func (c *Client) RemoveUserToGroup(ctx context.Context, groupID string, userID string) (groupResponse GroupResponse, groupErrorResponse GroupErrorResponse, err error) {
	return c.GroupMemberOps(ctx, groupID, userID, "Remove")
}

func (c *Client) DeleteGroup(ctx context.Context, groupID string) (err error) {

	fullUrl := fmt.Sprintf("%s%s/%s", c.BaseUrl, groupPath, groupID)

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


