package examples

import (
	"fmt"

	"github.com/authorizerdev/authorizer-go"
)

// LoginExample demonstrates how to use Login function of authorizer sdk
func LoginExample() {
	c, err := authorizer.NewAuthorizerClient(ClientID, AuthorizerURL, "", nil)
	if err != nil {
		panic(err)
	}

	res, err := c.Login(&authorizer.LoginInput{
		Email:    &TestEmail,
		Password: "Abc@123",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(authorizer.StringValue(res.Message))
}
