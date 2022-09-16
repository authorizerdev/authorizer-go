package examples

import (
	"fmt"

	"github.com/authorizerdev/authorizer-go"
)

// MagicLinkLoginExample demonstrates how to use MagicLinkLogin function of authorizer sdk
func MagicLinkLoginExample() {
	c, err := authorizer.NewAuthorizerClient(ClientID, AuthorizerURL, "", nil)
	if err != nil {
		panic(err)
	}

	res, err := c.MagicLinkLogin(&authorizer.MagicLinkLoginInput{
		Email: "test@yopmail.com",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Message)
}
