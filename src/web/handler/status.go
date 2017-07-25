package handler

import (
    "github.com/gin-gonic/gin"
    "fmt"
)

func StatusHandler(context *gin.Context) {
    fmt.Println(context.Request.Header.Get("token"))
    context.String(200, "ok")
}