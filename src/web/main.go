package main

import (
    "github.com/gin-gonic/gin"
    "web/handler"
)

func main() {
    r := gin.Default()
    r.GET("/ping", handler.SimpleHandler)
    r.Run()
}
