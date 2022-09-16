package examples

import (
	"fmt"

	"github.com/authorizerdev/authorizer-go"
)

// GetSessionExample demonstrates how to use GetSession function of authorizer sdk
func GetSessionExample() {
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

	res, err := c.GetSession(&authorizer.SessionQueryInput{
		Roles: []*string{authorizer.NewStringRef("test")},
	}, map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", authorizer.StringValue(loginRes.AccessToken)),
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(authorizer.StringValue(res.Message))
}
