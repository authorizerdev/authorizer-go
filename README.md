# authorizer-go

Go SDK for [authorizer.dev](https://authorizer.dev) — self-hosted authentication & authorization. Current version: **2.0.0**.

Use this SDK to add auth middleware to your Go services and to call any Authorizer API method. For full API docs see [pkg.go.dev/github.com/authorizerdev/authorizer-go](https://pkg.go.dev/github.com/authorizerdev/authorizer-go).

## Getting Started

**Pre-requisite:** A running Authorizer instance. See the [deployment guide](https://docs.authorizer.dev/deployment).

Once deployed, get your `Client ID` from the Authorizer dashboard.

### Step 1: Install the SDK

```bash
go get github.com/authorizerdev/authorizer-go
```

### Step 2: Initialize the client

**Parameters**

| Key             | Type                | Required | Description                                                                                          |
| --------------- | ------------------- | -------- | ---------------------------------------------------------------------------------------------------- |
| `clientID`      | `string`            | Yes      | Client ID from the Authorizer dashboard                                                              |
| `authorizerURL` | `string`            | Yes      | Base URL of your Authorizer instance                                                                 |
| `redirectURL`   | `string`            | No       | Default redirect URL for signup / login / forgot-password flows                                      |
| `extraHeaders`  | `map[string]string` | No       | Headers sent on every request (e.g. `Origin`)                                                       |

**Example**

```go
package main

import (
    authorizer "github.com/authorizerdev/authorizer-go"
)

func main() {
    client, err := authorizer.NewAuthorizerClient(
        "YOUR_CLIENT_ID",
        "https://your-instance.example.com",
        "https://your-app.example.com", // optional redirect URL
        map[string]string{},
    )
    if err != nil {
        panic(err)
    }
    _ = client
}
```

> **Note (Authorizer >= v2.3.0):** the server's CSRF guard requires an `Origin` header on state-changing requests. The client sends the Authorizer server's own origin by default, which always passes. If your instance restricts `ALLOWED_ORIGINS`, pass your app's origin instead via `extraHeaders`: `map[string]string{"Origin": "https://your-app.com"}`.

### Step 3: Use the SDK

```go
response, err := client.Login(&authorizer.LoginInput{
    Email:    "user@example.com",
    Password: "Abc@123",
})
if err != nil {
    panic(err)
}
fmt.Println("access_token:", response.AccessToken)
```

## Admin API

The SDK exposes admin methods for server-side user management. Admin operations require the admin secret; supply it via `extraHeaders`:

```go
adminClient, err := authorizer.NewAuthorizerClient(
    "YOUR_CLIENT_ID",
    "https://your-instance.example.com",
    "",
    map[string]string{"x-authorizer-admin-secret": "YOUR_ADMIN_SECRET"},
)
```

See [admin_methods.go](admin_methods.go) for the full list of available admin operations.

## Fine-grained authorization (FGA)

Authorizer ships an embedded [OpenFGA](https://openfga.dev) engine for relationship-based access control (ReBAC). You model your domain as object **types** with **relations** (`viewer`, `editor`, `owner`…), grant access by writing **relationship tuples**, and ask the engine whether access is allowed.

Authoring the model and tuples is an admin task — do it once in the dashboard under **Authorization**, or via the `_fga_*` admin GraphQL API. The SDK exposes only the read-side checks an application needs at request time. The subject defaults to the authenticated caller and is pinned server-side from the request headers. The optional `User` field overrides the subject but is honored only for super-admins or when it equals the caller's own token subject.

**Check permissions** — evaluates one or more relation checks in a single round trip:

```go
res, err := client.CheckPermissions(&authorizer.CheckPermissionsRequest{
    Checks: []*authorizer.PermissionCheckInput{
        {Relation: "can_view", Object: "document:1"},
        {Relation: "can_edit", Object: "document:1"},
    },
}, map[string]string{
    "Authorization": "Bearer " + token,
})
if err != nil {
    panic(err)
}
for _, r := range res.Results {
    fmt.Println(r.Relation, r.Object, r.Allowed)
}
```

**List accessible objects** — returns the ids of every object of a type the caller holds a relation on:

```go
res, err := client.ListPermissions(&authorizer.ListPermissionsRequest{
    Relation:   "can_view",
    ObjectType: "document",
}, map[string]string{"Authorization": "Bearer " + token})
if err != nil {
    panic(err)
}
fmt.Println(res.Objects) // ["document:1", "document:7", ...]
```

## API gateway example (gin)

```go
package main

import (
    "net/http"
    "strings"

    authorizer "github.com/authorizerdev/authorizer-go"
    "github.com/gin-gonic/gin"
)

func AuthorizeMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.Request.Header.Get("Authorization")
        tokenSplit := strings.Split(authHeader, " ")

        client, err := authorizer.NewAuthorizerClient(
            "YOUR_CLIENT_ID",
            "YOUR_AUTHORIZER_URL",
            "",
            map[string]string{},
        )
        if err != nil || len(tokenSplit) < 2 || tokenSplit[1] == "" {
            c.AbortWithStatusJSON(401, "unauthorized")
            return
        }

        res, err := client.ValidateJWTToken(&authorizer.ValidateJWTTokenInput{
            TokenType: authorizer.TokenTypeIDToken,
            Token:     tokenSplit[1],
        })
        if err != nil || !res.IsValid {
            c.AbortWithStatusJSON(401, "unauthorized")
            return
        }

        c.Next()
    }
}

func main() {
    router := gin.New()
    router.Use(AuthorizeMiddleware())

    router.GET("/ping", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "pong"})
    })

    router.Run(":8090")
}
```

Test the protected endpoint:

```bash
curl -H 'Authorization: Bearer JWT_TOKEN' http://localhost:8090/ping
```

---

## Release

1. Tag the commit: `git tag v<version>`
2. Push with tags: `git push origin main --tags`

The GitHub Actions release workflow handles the GitHub Release creation automatically. Go modules are versioned by tags — no additional publish step is needed.
