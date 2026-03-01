package examples

import (
	"fmt"

	"github.com/authorizerdev/authorizer-go"
)

// UpdateProfileExample demonstrates how to use UpdateProfile function of authorizer sdk
func UpdateProfileExample() {
	c, err := authorizer.NewAuthorizerClient(ClientID, AuthorizerURL, "", nil)
	if err != nil {
		panic(err)
	}

	loginRes, err := c.Login(&authorizer.LoginRequest{
		Email:    &TestEmail,
		Password: "Abc@123",
	})
	if err != nil {
		panic(err)
	}

	res, err := c.UpdateProfile(&authorizer.UpdateProfileRequest{
		FamilyName: authorizer.NewStringRef("test"),
	}, map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", authorizer.StringValue(loginRes.AccessToken)),
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Message)
}
