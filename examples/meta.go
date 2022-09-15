package examples

import (
	"authorizer-go"
	"fmt"
)

// GetMetaDataExample demonstrates how to use GetMetaData function of authorizer sdk
func GetMetaDataExample() {
	c, err := authorizer.NewAuthorizerClient(ClientID, AuthorizerURL, "", nil)
	if err != nil {
		panic(err)
	}

	res, err := c.GetMetaData()
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Version, res.ClientID)
}
