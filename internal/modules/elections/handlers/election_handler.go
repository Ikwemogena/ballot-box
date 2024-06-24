package handlers

import (
	contestant_model "ballot-box/internal/modules/contestants/models"
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

func GetElection(db *sql.DB) gin.HandlerFunc {
	return func (c *gin.Context) {
		electionID := c.Param("election_id")

		query := `
			SELECT id, title, description, start_time, end_time, created_by
			FROM elections
			WHERE id = $1
		`

		var election models.Election

		err := db.QueryRow(query, electionID).Scan(
			&election.ID,
			&election.Title,
			&election.Description,
			&election.StartTime,
			&election.EndTime,
			&election.CreatedBy,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get election"})
			return
		}

		c.JSON(http.StatusOK, election)
	}
}

func GetElectionContestants(db *sql.DB) gin.HandlerFunc {
	return func (c *gin.Context) {
		electionID := c.Param("election_id")

		query := `
			SELECT id, election_id, name
			FROM contestants
			WHERE election_id = $1
		`

		rows, err := db.Query(query, electionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get contestants"})
			return
		}
		defer rows.Close()

		contestants := []contestant_model.Contestant{}


		for rows.Next() {
			var contestant contestant_model.Contestant
			err := rows.Scan(&contestant.ID, &contestant.Name, &contestant.ElectionID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get contestants"})
				return
			}
			contestants = append(contestants, contestant)
		}

		c.JSON(http.StatusOK, contestants)
	}
}