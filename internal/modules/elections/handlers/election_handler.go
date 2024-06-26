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


func AddPositionToElection(db *sql.DB) gin.HandlerFunc {
	return func (c *gin.Context) {
		electionID := c.Param("election_id")

		var position models.Position

		if err := c.ShouldBindJSON(&position); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		query := `
			INSERT INTO positions (name, election_id)
			VALUES ($1, $2)
			RETURNING id
		`

		err := db.QueryRow(query, position.Name, electionID).Scan(&position.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add position to election"})
			return
		}

		c.JSON(http.StatusOK, position)
	}

}

func GetElectionDetails(db *sql.DB) gin.HandlerFunc {
	return func (c *gin.Context) {
		electionID := c.Param("election_id")


		electionQuery := `
			SELECT id, title, description, start_time, end_time, created_by
			FROM elections
			WHERE id = $1
		`

		var election models.Election

		err := db.QueryRow(electionQuery, electionID).Scan(
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

		positonsQuery := `
			SELECT id, name, election_id, created_by, created_at 
			FROM positions 
			WHERE election_id = $1
		`

		rows, err := db.Query(positonsQuery, electionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get positions"})
			return
		}

		defer rows.Close()
		
		var positions []models.Position

		for rows.Next() {
			var position models.Position
			err := rows.Scan(&position.ID, &position.Name, &position.ElectionID, &position.CreatedBy, &position.CreatedAt)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not scan positions"})
				return
			}

			contestantsQuery := `
				SELECT id, name, position_id, created_at 
				FROM contestants 
				WHERE position_id = $1
			`
			contestantRows, err := db.Query(contestantsQuery, position.ID)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying contestants"})
				return
			}

			defer contestantRows.Close()

			var contestants []models.Contestant

			for contestantRows.Next() {
				var contestant models.Contestant
				err := contestantRows.Scan(&contestant.ID, &contestant.Name, &contestant.PositionID, &contestant.CreatedAt)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning contestant"})
					return
				}
				contestants = append(contestants, contestant)
			}
			position.Contestants = contestants
			positions = append(positions, position)
		}

		election.Positions = positions
		c.JSON(http.StatusOK, election)
	}
}

func GetElectionContestants(db *sql.DB) gin.HandlerFunc {
	return func (c *gin.Context) {
		electionID := c.Param("election_id")

		query := `
			SELECT id, name, election_id
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
