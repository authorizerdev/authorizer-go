package examples

import (
	"fmt"

	"github.com/authorizerdev/authorizer-go"
)

// ValidateJWTTokenExample demonstrates how to use ValidateJWTToken function of authorizer sdk
func ValidateJWTTokenExample() {
	c, err := authorizer.NewAuthorizerClient(ClientID, AuthorizerURL, "", nil)
	if err != nil {
		panic(err)
	}

	loginRes, err := c.Login(&authorizer.LoginInput{
		Email:    "test@yopmail.com",
		Password: "Abc@123",
	})
	if err != nil {
		panic(err)
	}

	res, err := c.ValidateJWTToken(&authorizer.ValidateJWTTokenInput{
		TokenType: authorizer.TokenTypeAccessToken,
		Token:     authorizer.StringValue(loginRes.AccessToken),
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(res.IsValid)
}
