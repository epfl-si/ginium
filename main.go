package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

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

    // Run the server on port 8088
    router.Run(":8088")
}
