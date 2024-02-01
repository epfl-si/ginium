package main

import "github.com/gin-gonic/gin"

func main() {
    // Initialize the Gin router
    router := gin.Default()

    // No X-Forwarded-For business (yet)
    router.SetTrustedProxies(nil)

    // Define a route handler for the root path
    router.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Hello, Gin!",
        })
    })

    // Run the server on port 8088
    router.Run(":8088")
}
