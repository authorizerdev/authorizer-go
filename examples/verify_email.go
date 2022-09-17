package examples

import (
	"fmt"

	"github.com/authorizerdev/authorizer-go"
)

// VerifyEmailExample demonstrates how to use VerifyEmail function of authorizer sdk
func VerifyEmailExample() {
	c, err := authorizer.NewAuthorizerClient(ClientID, AuthorizerURL, "", nil)
	if err != nil {
		panic(err)
	}

	res, err := c.VerifyEmail(&authorizer.VerifyEmailInput{
		Token: "token sent via email",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(authorizer.StringValue(res.AccessToken))
}
