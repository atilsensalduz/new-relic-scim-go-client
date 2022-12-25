# new-relic-scim-go-client

[![GoDoc](https://godoc.org/github.com/atileren/new-relic-scim-go-client?status.svg)](https://godoc.org/github.com/atileren/new-relic-scim-go-client)

This repository contains a Go client library for interacting with the [New Relic SCIM API](https://docs.newrelic.com/docs/accounts/accounts/automated-user-management/automated-user-provisioning-single-sign/). The SCIM API allows you to manage users and groups in your New Relic account.

## Installation

To install the library, run:
```go
go get github.com/atileren/new-relic-scim-go-client/newrelicscim
```


## Usage

Here is a simple example of how to use the library to create a new user:

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/atileren/new-relic-scim-go-client/newrelicscim"
)

func main() {
	// Create a new client
	client := newrelicscim.NewClient("<your_api_key>")
	ctx := context.Background()
	// Create a new user
	user, userErrorResponse, err := client.CreateUser(ctx, newrelicscim.User{
		UserName: "john.doe@example.com",
		Name: newrelicscim.Name{
			FamilyName: "Doe",
			GivenName:  "John",
		},
		Emails: []newrelicscim.Email{
			{
				Primary: true,
				Value:   "john.doe@example.com",
			},
		},
	})
	if err != nil {
		log.Fatal("request error: ", err)
	}
	if userErrorResponse.Status != "" {
		log.Fatal("error detail: " + userErrorResponse.Detail)
	}

	fmt.Println("Successfully created user with ID:", user.ID)
}
```

For more detailed examples and documentation, see the [GoDoc](https://godoc.org/github.com/atileren/new-relic-scim-go-client) documentation.

## Contributing

We welcome contributions to the new-relic-scim-go-client repository. If you have an idea for a new feature or bug fix, please open an issue to discuss it. If you would like to contribute code, please fork the repository and submit a pull request.

License
This library is licensed under the MIT License. See [LICENSE](https://github.com/atileren/new-relic-scim-go-client/blob/main/LICENSE) for more details.



