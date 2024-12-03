package routes

import (
	"net/http"

	"example.com/events-api/models"
	"example.com/events-api/utils"
	"github.com/gin-gonic/gin"
)

func signup(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to parse request body",
		})
		return
	}

	err = user.Save()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create user",
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Successfully created user",
	})
}

func login(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to parse request body",
		})
		return
	}

	err = user.ValidateCredentials()

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}

	token, err := utils.GenerateJwtToken(user.ID, user.Email)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to authenticate user",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged in",
		"token":   token,
	})
}
