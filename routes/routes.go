package routes

import (
	"example.com/events-api/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	// Events routes
	authGroup := server.
		Group("/").
		Use(middleware.Auth)
	server.GET("/events", getEvents)
	server.GET("/events/:id", getEventById)
	authGroup.POST("/events", createEvent)
	authGroup.PUT("/events/:id", updateEvent)
	authGroup.PATCH("/events/:id", patchEvent)
	authGroup.DELETE("/events/:id", deleteEvent)
	authGroup.POST("/events/:id/register", registerForEvent)
	authGroup.DELETE("/events/:id/register", cancelRegistration)

	// Users routes
	server.POST("/signup", signup)
	server.POST("/login", login)
}
