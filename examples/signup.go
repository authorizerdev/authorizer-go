package examples

import (
	"fmt"

	"github.com/authorizerdev/authorizer-go"
)

// SignUpExample demonstrates how to use SignUp function of authorizer sdk
func SignUpExample() {
	c, err := authorizer.NewAuthorizerClient(ClientID, AuthorizerURL, "", nil)
	if err != nil {
		panic(err)
	}

	res, err := c.SignUp(&authorizer.LoginInput{
		Email:    "test@yopmail.com",
		Password: "Abc@123",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(authorizer.StringValue(res.Message))
}
