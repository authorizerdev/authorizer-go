package examples

import (
	"authorizer-go"
	"fmt"
)

// ForgotPasswordInputExample demonstrates how to use ForgotPassword function of authorizer skd
func ForgotPasswordInputExample() {
	c, err := authorizer.NewAuthorizerClient(ClientID, AuthorizerURL, "", nil)
	if err != nil {
		panic(err)
	}

	res, err := c.ForgotPassword(&authorizer.ForgotPasswordInput{
		Email: "test@yopmail.com",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Message)
}
