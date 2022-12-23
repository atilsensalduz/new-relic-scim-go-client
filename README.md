# new-relic-scim-go-client

[![GoDoc](https://godoc.org/github.com/atileren/new-relic-scim-go-client?status.svg)](https://godoc.org/github.com/atileren/new-relic-scim-go-client)

This repository contains a Go client library for interacting with the [New Relic SCIM API](https://docs.newrelic.com/docs/apis/scim-api). The SCIM API allows you to manage users and groups in your New Relic account.

## Installation

To install the library, run:
```go
go get github.com/atileren/new-relic-scim-go-client
```


## Usage

Here is a simple example of how to use the library to create a new user:

```go
package main

import (
	"fmt"
	"log"

	"github.com/atileren/new-relic-scim-go-client"
)

func main() {
	// Create a new client
	client, err := newrelicscim.NewClient("<your_api_key>", "<your_account_id>")
	if err != nil {
		log.Fatal(err)
	}

	// Create a new user
	userResponse, userErrorResponse, err := client.CreateUser(newrelicscim.User{
				UserName: "john.doe@example.com",
				Name: struct {
					FamilyName string "json:\"familyName\""
					GivenName  string "json:\"givenName\""
				}{
					FamilyName: "Doe",
					GivenName:  "John",
				},
				Emails: []struct {
					Primary bool   "json:\"primary\""
					Value   string "json:\"value\""
				}{
					{
						Primary: true,
						Value:   "john.doe@example.com",
					},
				},
			})
	if userResponse.Status != "" {
		log.fatal("error detail: " + errorResp.Detail)
	}

	fmt.Println("Successfully created user with ID:", user.ID)
}
```

For more detailed examples and documentation, see the [GoDoc](https://godoc.org/github.com/atileren/new-relic-scim-go-client) documentation.

Contributing
We welcome contributions to the new-relic-scim-go-client repository. If you have an idea for a new feature or bug fix, please open an issue to discuss it. If you would like to contribute code, please fork the repository and submit a pull request.

License
This library is licensed under the MIT License. See [LICENSE](https://github.com/atileren/new-relic-scim-go-client/blob/main/LICENSE) for more details.



