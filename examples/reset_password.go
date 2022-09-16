package examples

import (
	"fmt"

	"github.com/authorizerdev/authorizer-go"
)

// ResetPasswordExample demonstrates how to use ResetPassword function of authorizer sdk
func ResetPasswordExample() {
	c, err := authorizer.NewAuthorizerClient(ClientID, AuthorizerURL, "", nil)
	if err != nil {
		panic(err)
	}

	res, err := c.ResetPassword(&authorizer.ResetPasswordInput{
		Token:           "token obtained via forgot password email",
		Password:        "new password",
		ConfirmPassword: "new password",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Message)
}
