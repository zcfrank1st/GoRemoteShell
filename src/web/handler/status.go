package handler

import (
    "github.com/gin-gonic/gin"
)

func StatusHandler(context *gin.Context) {
    context.String(200, "ok")
}