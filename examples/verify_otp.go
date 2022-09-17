package examples

import (
	"fmt"

	"github.com/authorizerdev/authorizer-go"
)

// VerifyOTPExample demonstrates how to use VerifyOTP function of authorizer sdk
func VerifyOTPExample() {
	c, err := authorizer.NewAuthorizerClient(ClientID, AuthorizerURL, "", nil)
	if err != nil {
		panic(err)
	}

	res, err := c.VerifyOTP(&authorizer.VerifyOTPInput{
		OTP:   "test",
		Email: "test@yopmail.com",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(authorizer.StringValue(res.AccessToken))
}
