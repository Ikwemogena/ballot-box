package routes

import (
	"ballot-box/internal/middleware"
	"ballot-box/internal/modules/votes/handlers"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func Setup(router *gin.Engine, db *sql.DB) {
	
	auth := router.Group("/vote")

	auth.Use(middleware.AuthMiddleware())
	{
		auth.POST("/",  handlers.CastVote(db))
	}
}