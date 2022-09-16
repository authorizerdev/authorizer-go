package examples

import (
	"fmt"

	"github.com/authorizerdev/authorizer-go"
)

// RevokeTokenExample demonstrates how to use RevokeToken function of authorizer sdk
func RevokeTokenExample() {
	c, err := authorizer.NewAuthorizerClient(ClientID, AuthorizerURL, "", nil)
	if err != nil {
		panic(err)
	}

	res, err := c.RevokeToken(&authorizer.RevokeTokenInput{
		RefreshToken: "example refresh token",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Message)
}
