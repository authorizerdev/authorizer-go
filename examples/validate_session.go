package examples

import (
	"fmt"

	"github.com/authorizerdev/authorizer-go"
)

// ValidateSessionExample demonstrates how to use ValidateSession function of authorizer sdk
func ValidateSessionExample() {
	c, err := authorizer.NewAuthorizerClient(ClientID, AuthorizerURL, "", nil)
	if err != nil {
		panic(err)
	}
	_, err = c.Login(&authorizer.LoginInput{
		Email:    &TestEmail,
		Password: "Abc@123",
	})
	if err != nil {
		panic(err)
	}
	res, err := c.ValidateSession(&authorizer.ValidateSessionInput{
		Cookie: "", // TODO set cookie here
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(res.IsValid)
}
