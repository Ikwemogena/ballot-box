package routes

import (
	"ballot-box/internal/middleware"
	"ballot-box/internal/modules/elections/handlers"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func Setup(router *gin.Engine, db *sql.DB) {

	auth := router.Group("/election")

	auth.Use(middleware.AuthMiddleware())
	auth.Use(middleware.AdminOnlyMiddleware())
	{
		auth.POST("/create", handlers.CreateElection(db))
		auth.GET("/:election_id", handlers.GetElectionDetails(db))
		auth.POST("/:election_id/add-position", handlers.AddPositionToElection(db))
		auth.GET("/:election_id/contestants", handlers.GetElectionContestants(db))
		
	}
}