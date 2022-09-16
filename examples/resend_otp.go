package examples

import (
	"fmt"

	"github.com/authorizerdev/authorizer-go"
)

// ResendOTPExample demonstrates how to use ResendOTP function of authorizer sdk
func ResendOTPExample() {
	c, err := authorizer.NewAuthorizerClient(ClientID, AuthorizerURL, "", nil)
	if err != nil {
		panic(err)
	}

	res, err := c.ResendOTP(&authorizer.ResendOTPInput{
		Email: "test@yopmail.com",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Message)
}
