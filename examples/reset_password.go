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

	token := "token obtained via forgot password email"
	res, err := c.ResetPassword(&authorizer.ResetPasswordRequest{
		Token:           &token,
		Password:        "new password",
		ConfirmPassword: "new password",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Message)
}
