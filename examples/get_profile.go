package examples

import (
	"fmt"

	"github.com/authorizerdev/authorizer-go"
)

// GetProfileExample demonstrates how to use GetProfile function of authorizer sdk
func GetProfileExample() {
	c, err := authorizer.NewAuthorizerClient(ClientID, AuthorizerURL, "", nil)
	if err != nil {
		panic(err)
	}

	loginRes, err := c.Login(&authorizer.LoginInput{
		Email:    &TestEmail,
		Password: "Abc@123",
	})
	if err != nil {
		panic(err)
	}

	res, err := c.GetProfile(map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", authorizer.StringValue(loginRes.AccessToken)),
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Email)
}
