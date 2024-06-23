package main

import (
	"ballot-box/database"
	electionRoutes "ballot-box/internal/modules/elections/routes"
	voterRoutes "ballot-box/internal/modules/users/routes"
	voteRoutes "ballot-box/internal/modules/votes/routes"
	"log"

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
	electionRoutes.Setup(r, db)

	r.Run(":8080")
}