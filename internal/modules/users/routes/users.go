package routes

import (
	"database/sql"
	"voting-platform/internal/modules/users/handlers"

	"github.com/gin-gonic/gin"
)

func Setup(router *gin.Engine, db *sql.DB) {
	router.POST("/voter/register", handlers.Register(db))
	router.POST("/voter/login", handlers.Login(db))

}