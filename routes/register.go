package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"example.com/events-api/models"
	"github.com/gin-gonic/gin"
)

func registerForEvent(context *gin.Context) {
	userId := context.GetInt64("userId")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid event id",
		})
		return
	}

	event, err := models.GetEventById(eventId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Failed to fetch event with id %v", eventId),
		})
		return
	}

	isRegistered, err := event.IsUserRegistered(userId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to check if user is already registered",
		})
		return
	}

	if isRegistered {
		context.JSON(http.StatusConflict, gin.H{
			"message": "User is already registered for this event",
		})
		return
	}

	err = event.Register(userId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Failed to register user with id %v to event with id %v", userId, eventId),
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Successfully registered",
	})
}

func cancelRegistration(context *gin.Context) {
	userId := context.GetInt64("userId")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid event id",
		})
		return
	}

	event, err := models.GetEventById(eventId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Failed to fetch event with id %v", eventId),
		})
		return
	}

	isRegistered, err := event.IsUserRegistered(userId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to check if user is already registered",
		})
		return
	}

	if !isRegistered {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "User is not registered for this event",
		})
		return
	}

	err = event.CancelRegistration(userId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to cancel registration",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Registration canceled",
	})
}
