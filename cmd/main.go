package main

import (
	"log"
	"voting-platform/database"
	voterRoutes "voting-platform/internal/modules/users/routes"
	voteRoutes "voting-platform/internal/modules/votes/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer db.Close()

	r := gin.Default()

	voterRoutes.Setup(r, db)
	voteRoutes.Setup(r, db)

	r.Run(":8080")
}