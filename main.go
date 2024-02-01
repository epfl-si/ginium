package main

import (
    "context"
    "net/http"
    "strings"

    "golang.org/x/oauth2"
    "github.com/coreos/go-oidc/v3/oidc"
    "github.com/gin-gonic/gin"
)

func TokenAuthMiddleware(provider *oidc.Provider) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get the bearer token from the request header
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
            c.Abort()
            return
        }

        // Parse and verify the token
        bearerToken := strings.Split(tokenString, "Bearer ")[1]

        _, err := provider.UserInfo(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: bearerToken}))
	if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to get userinfo: "+err.Error()})
            c.Abort()
	    return
	}

        // Token is valid, continue with the request
        c.Next()
    }
}

func main() {
    // Initialize the Gin router
    router := gin.Default()

    // No X-Forwarded-For business (yet)
    router.SetTrustedProxies(nil)

    router.LoadHTMLFiles("index.html")

    router.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.html", gin.H{})
    })

    // Define a route handler for a toy public API
    router.GET("/api/v1/public", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Hello, Gin!",
        })
    })

    // Create an OpenID Connect provider
    provider, err := oidc.NewProvider(context.Background(), "http://localhost:8080/realms/golang")
    if err != nil {
        panic(err)
    }

    clientID := "hello_golang"

    router.GET("/login", func(c *gin.Context) {
	config := oauth2.Config{
		ClientID:     clientID,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  "http://localhost:8088/",
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

        nonce := "samesamebutdifferent"
        c.Redirect(http.StatusFound, config.AuthCodeURL(nonce, oidc.Nonce(nonce)))
    })

    // Define a route handler that is protected by OpenID-Connect
    router.GET("/api/v1/private", TokenAuthMiddleware(provider), func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "I is serious cat. This is serious data.",
        })
    })

    // Run the server on port 8088
    router.Run(":8088")
}
