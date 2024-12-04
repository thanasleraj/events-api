package middleware

import (
	"net/http"

	"example.com/events-api/utils"
	"github.com/gin-gonic/gin"
)

func Auth(context *gin.Context) {
	authHeader := context.Request.Header.Get("Authorization")
	userId, err := utils.VerifyJwtToken(authHeader)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	context.Set("userId", userId)

	context.Next()
}
