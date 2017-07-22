package handler

import (
    "github.com/gin-gonic/gin"
)

func SimpleHandler(context *gin.Context) {
    context.JSON(200, gin.H{
        "message": "pong",
    })
}