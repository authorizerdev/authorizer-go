package examples

import (
	"fmt"

	"github.com/authorizerdev/authorizer-go"
)

// ResendVerifyEmailExample demonstrates how to use ResendVerifyEmail function of authorizer sdk
func ResendVerifyEmailExample() {
	c, err := authorizer.NewAuthorizerClient(ClientID, AuthorizerURL, "", nil)
	if err != nil {
		panic(err)
	}

	res, err := c.ResendVerifyEmail(&authorizer.ResendVerifyEmailRequest{
		Email: "test@yopmail.com",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Message)
}
