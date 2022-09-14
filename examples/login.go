package examples

import (
	"authorizer-go"
	"fmt"
)

// LoginExample demonstrates how to use Login function of authorizer sdk
func LoginExample() {
	c, err := authorizer.NewAuthorizerClient(ClientID, AuthorizerURL, "", nil)
	if err != nil {
		panic(err)
	}

	res, err := c.Login(&authorizer.LoginRequest{
		Email:    "test@yopmail.com",
		Password: "Abc@123",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(authorizer.StringValue(res.Message))
}
