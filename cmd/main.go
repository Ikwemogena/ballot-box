package main

import (
	"ballot-box/database"
	contestantRoutes "ballot-box/internal/modules/contestants/routes"
	electionRoutes "ballot-box/internal/modules/elections/routes"
	voterRoutes "ballot-box/internal/modules/users/routes"
	voteRoutes "ballot-box/internal/modules/votes/routes"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer db.Close()

	r := gin.Default()

	config := cors.DefaultConfig()
    config.AllowAllOrigins = true
    config.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
    config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
    config.ExposeHeaders = []string{"Content-Length"}
    config.AllowCredentials = true
    config.MaxAge = 12 * time.Hour

    r.Use(cors.New(config))

	contestantRoutes.Setup(r, db)
	voterRoutes.Setup(r, db)
	voteRoutes.Setup(r, db)
	electionRoutes.Setup(r, db)

	r.Run(":8080")
}