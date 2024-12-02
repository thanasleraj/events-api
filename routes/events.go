package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"example.com/events-api/models"
	"github.com/gin-gonic/gin"
)

func getEvents(context *gin.Context) {
	events, err := models.GetAllEvents()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch events",
		})

		return
	}

	context.JSON(http.StatusOK, events)
}

func getEventById(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid event id",
		})
		return
	}

	event, err := models.GetEventById(id)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Failed to fetch event with id %v", id),
		})
		return
	}

	context.JSON(http.StatusOK, event)
}

func createEvent(context *gin.Context) {
	var event models.Event

	err := context.ShouldBindJSON(&event)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to parse request body",
		})

		return
	}

	event.UserID = 1

	err = event.Save()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create event",
		})

		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Successfully created event",
		"event":   event,
	})
}

func updateEvent(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid event id",
		})
		return
	}

	existingEvent, err := models.GetEventById(id)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Failed to fetch event with id %v", id),
		})
		return
	}

	var updatedEvent models.Event
	err = context.ShouldBindJSON(&updatedEvent)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to parse request body",
		})
		return
	}

	updatedEvent.ID = existingEvent.ID
	updatedEvent.UserID = existingEvent.UserID

	err = updatedEvent.Update()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Failed to update event with id %v", id),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Successfully updated event",
		"event":   updatedEvent,
	})
}

func patchEvent(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid event id",
		})
		return
	}

	event, err := models.GetEventById(id)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Failed to fetch event with id %v", id),
		})
		return
	}

	var patchRequest models.PatchEventRequest
	err = context.ShouldBindJSON(&patchRequest)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to parse request body",
		})
		return
	}

	err = event.Patch(patchRequest)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Failed to update event with id %v", id),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Successfully updated event",
		"event":   event,
	})
}

func deleteEvent(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid event id",
		})
		return
	}

	event, err := models.GetEventById(id)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Failed to fetch event with id %v", id),
		})
		return
	}

	err = event.Delete()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Failed to delete event with id %v", id),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Successfully deleted event",
	})
}
