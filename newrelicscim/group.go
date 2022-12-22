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

// Group represents a group in the New Relic SCIM API.
//
// It has the following fields:
//  - Schemas: a slice of strings containing the SCIM schema URIs that define the attributes of the group
//  - DisplayName: the name of the group, which is used to identify it in the New Relic user interface
type Group struct {
	Schemas     []string `json:"schemas"`
	DisplayName string   `json:"displayName"`
}

// GroupResponse represents a response from the New Relic SCIM API for a group creation or update request.
//
// It has the following fields:
//  - Schemas: a slice of strings containing the SCIM schema URIs that define the attributes of the group
//  - ID: the unique identifier for the group, assigned by the New Relic SCIM API
//  - DisplayName: the name of the group, which is used to identify it in the New Relic user interface
//  - Meta: metadata about the group, including the resource type, creation date, and last modification date
//  - Members: a slice of interfaces representing the members of the group (typically user resources)
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

// GroupErrorResponse represents an error response from the New Relic SCIM API for a group creation or update request.
//
// It has the following fields:
//  - Schemas: a slice of strings containing the SCIM schema URIs that define the attributes of the group error response
//  - ScimType: a string indicating the type of error that occurred
//  - Detail: a string describing the error in more detail
//  - Status: a string indicating the HTTP status code of the error
type GroupErrorResponse struct {
	Schemas  []string `json:"schemas"`
	ScimType string   `json:"scimType"`
	Detail   string   `json:"detail"`
	Status   string   `json:"status"`
}

// GroupsResponse represents a response from the New Relic SCIM API for a group list request.
//
// It has the following fields:
//  - TotalResults: an integer indicating the total number of groups that match the list request
//  - Schemas: a slice of strings containing the SCIM schema URIs that define the attributes of the group list response
//  - Resources: a slice of structs representing the groups that match the list request, each with the fields described in the GroupResponse struct
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
		Members []struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"members"`
	} `json:"Resources"`
}

// UpdateGroup represents a request to update a group in the New Relic SCIM API using the patch operation.
//
// It has the following fields:
//  - Schemas: a slice of strings containing the SCIM schema URIs that define the attributes of the update request
//  - Operations: a slice of structs representing the patch operations to be performed on the group, such as adding or
//    removing members or changing the group name
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

// fill_defaults is a helper function that sets default values for a Group struct if they are not already present.
//
// It sets the Schemas field to the SCIM schema URI for groups if it is empty. This is required by the SCIM API when
// creating or updating a group.
func (g *Group) fill_defaults() {

	// setting default values
	// if no values present
	if len(g.Schemas) == 0 {
		g.Schemas = []string{"urn:ietf:params:scim:schemas:core:2.0:Group"}
	}

}

// fill_defaults is a helper function that sets default values for an UpdateGroup struct if they are not already present.
//
// It sets the Schemas field to the SCIM schema URI for patch operations if it is empty. This is required by the SCIM API
// when patching a group.
func (ug *UpdateGroup) fill_defaults() {

	// setting default values
	// if no values present
	if len(ug.Schemas) == 0 {
		ug.Schemas = []string{"urn:ietf:params:scim:api:messages:2.0:PatchOp"}
	}

}

// CreateGroup is a function that creates a new group in the New Relic SCIM API using the provided group name.
//
// It takes the following arguments:
//  - ctx: a context for cancelling or timing out the request
//  - groupName: the name of the group to be created
//
// It returns the following values:
//  - groupResponse: a GroupResponse struct containing the details of the created group if the operation was successful
//  - groupErrorResponse: a GroupErrorResponse struct containing details of the error if the operation was not successful
//  - err: an error value if there was an issue with the request or response
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

// UpdateGroup is a function that updates an existing group in the New Relic SCIM API using the provided group name.
//
// It takes the following arguments:
//  - ctx: a context for cancelling or timing out the request
//  - groupName: the new name of the group to be updated
//
// It returns the following values:
//  - groupResponse: a GroupResponse struct containing the details of the updated group if the operation was successful
//  - groupErrorResponse: a GroupErrorResponse struct containing details of the error if the operation was not successful
//  - err: an error value if there was an issue with the request or response
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

// GroupList is a function that retrieves a list of groups from the New Relic SCIM API.
//
// It takes the following arguments:
//  - ctx: a context for cancelling or timing out the request
//
// It returns the following values:
//  - groupsResponse: a GroupsResponse struct containing the details of the retrieved groups if the operation was successful
//  - groupErrorResponse: a GroupErrorResponse struct containing details of the error if the operation was not successful
//  - err: an error value if there was an issue with the request or response
func (c *Client) GroupList(ctx context.Context) (groupsResponse GroupsResponse, groupErrorResponse GroupErrorResponse, err error) {
	// Construct the full URL for the request
	fullUrl := fmt.Sprintf("%s%s", c.BaseUrl, groupPath)

	// Create a new HTTP GET request
	req, err := http.NewRequest("Get", fullUrl, nil)
	if err != nil {
		return groupsResponse, groupErrorResponse, err
	}

	// Send the request and get the response
	resp, err := c.doRequest(req)
	if err != nil {
		return groupsResponse, groupErrorResponse, err
	}

	// Unmarshal the response into a GroupsResponse struct
	if err := json.Unmarshal(resp, &groupsResponse); err != nil {
		return groupsResponse, groupErrorResponse, err
	}

	// If the response is an error, unmarshal it into a GroupErrorResponse struct
	if groupsResponse.Schemas[0] == "urn:ietf:params:scim:api:messages:2.0:Error" {
		if err := json.Unmarshal(resp, &groupErrorResponse); err != nil {
			return groupsResponse, groupErrorResponse, err
		}
	}

	return groupsResponse, groupErrorResponse, nil
}

// GetGroupByID fetches a group by its ID using the SCIM API.
//
// It takes the following arguments:
//  - ctx: the context for the request
//  - groupID: the ID of the group to fetch
//
// It returns the following values:
//  - groupsResponse: a GroupsResponse struct containing the group information if the request is successful
//  - groupErrorResponse: a GroupErrorResponse struct containing the error information if there is an error with the request
//  - err: an error if there is any issue with the request or response
func (c *Client) GetGroupByID(ctx context.Context, groupID string) (groupsResponse GroupsResponse, groupErrorResponse GroupErrorResponse, err error) {

	// Construct the full URL for the request
	fullUrl := fmt.Sprintf("%s%s/%s", c.BaseUrl, groupPath, groupID)

	// Create a new HTTP GET request
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return groupsResponse, groupErrorResponse, err
	}

	// Send the request and get the response
	resp, err := c.doRequest(req)
	if err != nil {
		return groupsResponse, groupErrorResponse, err
	}

	// Unmarshal the response into a GroupsResponse struct
	if err := json.Unmarshal(resp, &groupsResponse); err != nil {
		return groupsResponse, groupErrorResponse, err
	}

	// If the response is an error, unmarshal it into a GroupErrorResponse struct
	if groupsResponse.Schemas[0] == "urn:ietf:params:scim:api:messages:2.0:Error" {
		if err := json.Unmarshal(resp, &groupErrorResponse); err != nil {
			return groupsResponse, groupErrorResponse, err
		}
	}

	return groupsResponse, groupErrorResponse, nil
}

// GetGroupByName is a function that retrieves a group by its name using the New Relic SCIM API.
//
// It takes the following arguments:
//  - ctx: a context for cancelling or timing out the request
//  - groupName: the name of the group to retrieve
//
// It returns the following values:
//  - groupsResponse: a GroupsResponse struct containing the details of the retrieved group if the operation was successful
//  - groupErrorResponse: a GroupErrorResponse struct containing details of the error if the operation was not successful
//  - err: an error value if there was an issue with the request or response
func (c *Client) GetGroupByName(ctx context.Context, groupName string) (groupsResponse GroupsResponse, groupErrorResponse GroupErrorResponse, err error) {
	// Construct the full URL for the request
	fullUrl := fmt.Sprintf("%s%s", c.BaseUrl, groupPath)

	// Create a new HTTP GET request
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return groupsResponse, groupErrorResponse, err
	}

	// Add the filter parameter to the request URL to filter the results by group name
	q := req.URL.Query()
	filter := fmt.Sprintf(`displayName eq "%s"`, groupName)
	q.Add("filter", filter)
	req.URL.RawQuery = q.Encode()

	// Send the request and get the response
	resp, err := c.doRequest(req)
	if err != nil {
		return groupsResponse, groupErrorResponse, err
	}

	// Unmarshal the response into a GroupsResponse struct
	if err := json.Unmarshal(resp, &groupsResponse); err != nil {
		return groupsResponse, groupErrorResponse, err
	}

	// If the response is an error, unmarshal it into a GroupErrorResponse struct
	if groupsResponse.Schemas[0] == "urn:ietf:params:scim:api:messages:2.0:Error" {
		if err := json.Unmarshal(resp, &groupErrorResponse); err != nil {
			return groupsResponse, groupErrorResponse, err
		}
	}

	return groupsResponse, groupErrorResponse, nil
}

// GroupMemberOps is a function that performs an operation on a group member in the New Relic SCIM API.
//
// It takes the following arguments:
//  - ctx: a context for cancelling or timing out the request
//  - groupID: the ID of the group to perform the operation on
//  - userID: the ID of the user to perform the operation on
//  - operation: the operation to perform on the group member (e.g. "add", "remove")
//
// It returns the following values:
//  - groupResponse: a GroupResponse struct containing the details of the modified group if the operation was successful
//  - groupErrorResponse: a GroupErrorResponse struct containing details of the error if the operation was not successful
//  - err: an error value if there was an issue with the request or response
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
