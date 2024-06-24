package routes

import (
	"ballot-box/internal/middleware"
	"ballot-box/internal/modules/contestants/handlers"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func Setup(router *gin.Engine, db *sql.DB) {

	auth := router.Group("/contestant")

	auth.Use(middleware.AuthMiddleware())
	auth.Use(middleware.AdminOnlyMiddleware())
	{
		auth.POST("/add", handlers.AddContestant(db))
	}
}