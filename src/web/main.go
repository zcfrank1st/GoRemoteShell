package main

import (
    "github.com/gin-gonic/gin"
    "web/handler"
)

func main() {
    r := gin.Default()

    r.GET("/status", handler.StatusHandler)

    v1 := r.Group("/openapi/v1")
    {
        v1.POST("/member", handler.MemberHandler)
        v1.POST("/display", handler.DisplayHandler)
    }

    r.Run()
}
