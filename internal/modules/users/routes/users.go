package routes

import (
	"ballot-box/internal/modules/users/handlers"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func Setup(router *gin.Engine, db *sql.DB) {
	router.POST("/admin/register", handlers.RegisterAdmin(db))
	router.POST("/voter/register", handlers.Register(db))	

	router.POST("/admin/login", handlers.LoginAdmin(db))
	router.POST("/voter/login", handlers.Login(db))
}