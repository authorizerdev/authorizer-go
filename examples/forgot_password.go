package examples

import (
	"fmt"

	"github.com/authorizerdev/authorizer-go/v2"
)

// ForgotPasswordInputExample demonstrates how to use ForgotPassword function of authorizer skd
func ForgotPasswordInputExample() {
	c, err := authorizer.NewAuthorizerClient(ClientID, AuthorizerURL, "", nil)
	if err != nil {
		panic(err)
	}

	email := "test@yopmail.com"
	res, err := c.ForgotPassword(&authorizer.ForgotPasswordRequest{
		Email: &email,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Message)
}
