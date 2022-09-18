# authorizer-go

Golang SDK for [authorizer.dev](https://authorizer.dev) server. This SDK will be handy to add API middleware where you can authorize your users. It will also empower you to perform various auth operations on authorizer server.

For detailed explanation of each functions check official [docs](https://pkg.go.dev/github.com/authorizerdev/authorizer-go)

## Getting Started

**Pre-requisite**: You will need an authorizer instance up and running. Checkout how you can host your instance in the [docs](https://docs.authorizer.dev/deployment)

Follow the steps here to install authorizer-go in your golang project and use the methods of SDK to protect/authorize your APIs

Once you have deployed authorizer instance. Get `Client ID` from your authorizer instance dashboard

![client_id](https://res.cloudinary.com/dcfpom7fo/image/upload/v1663437088/Authorizer/client_id_ptjsvc.png)

### Step 1: Install authorizer-go SDK

Run the following command to download authorizer-go SDK

```bash
go get github.com/authorizerdev/authorizer-go
```

### Step 2: Initialize authorizer client

**Required Parameters**

| Key             | Type | Required | Description                                                                                                     |
| --------------- | --------------------------------------------------------------------------------------------------------------- |
| `clientID`      | `string` | `true` | Your unique client identifier obtained from authorizer dashboard                                                |
| `authorizerURL` | `string` |`true` |Authorizer server URL                                                                                           |
| `redirectURL`   | `string`| `false` |Default URL to which you would like to redirect the user in case of successful signup / login / forgot password |
| `extraHeaders` | `map[string]string` | `false` | set of headers that you would like to pass with each request |

__Example__

```go
defaultHeaders := map[string]string{}

authorizerClient, err := authorizer.NewAuthorizerClient("YOUR_CLIENT_ID", "YOUR_AUHTORIZER_URL", "OPTIONAL_REDIRECT_URL", defaultHeaders)
if err != nil {
    panic(err)
}
```

### Step 3: Access all the SDK methods using authorizer client instance, initialized on step 2

__Example__

```go
response, err := authorizerClient.Login(&authorizer.LoginInput{
    Email:    "test@yopmail.com",
    Password: "Abc@123",
})
if err != nil {
    panic(err)
}
```