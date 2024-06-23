package handlers

import (
	"ballot-box/internal/modules/elections/models"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateElection(db *sql.DB) gin.HandlerFunc {
	return func (c *gin.Context) {
		
		var election models.Election

		if err := c.ShouldBindJSON(&election); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID := c.GetString("userID")
    	userRole := c.GetString("role")

		election.CreatedBy = userID

		if userRole != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Admin access required"})
			return
    	}	


		query := `
			INSERT INTO elections (title, description, start_time, end_time, created_by)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		`


		err := db.QueryRow(
			query,
			election.Title,
			election.Description,
			election.StartTime,
			election.EndTime,
			election.CreatedBy,
		).Scan(&election.ID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create election"})
			return
		}

		c.JSON(http.StatusOK, election)
	}
}