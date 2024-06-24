package handlers

import (
	"ballot-box/internal/modules/contestants/models"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddContestant(db *sql.DB) gin.HandlerFunc {
	return func (c *gin.Context) {
		var contestant models.Contestant

		if err := c.ShouldBindJSON(&contestant); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		query := `
			INSERT INTO contestants (election_id, name)
			VALUES ($1, $2)
			RETURNING id
		`

		err := db.QueryRow(
			query,
			contestant.ElectionID,
			contestant.Name,
		).Scan(&contestant.ID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add contestant"})
			return
		}

		c.JSON(http.StatusOK, contestant)
	}
}

func GetAllContestants(db *sql.DB) gin.HandlerFunc {
	return func (c *gin.Context) {
		var contestants []models.Contestant

		query := `
			SELECT id, election_id, name
			FROM contestants
		`

		rows, err := db.Query(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get contestants"})
			return
		}

		for rows.Next() {
			var contestant models.Contestant

			err := rows.Scan(
				&contestant.ID,
				&contestant.ElectionID,
				&contestant.Name,
			)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get contestants"})
				return
			}

			contestants = append(contestants, contestant)
		}

		c.JSON(http.StatusOK, contestants)
	}
}

// check if the contestant exists in the election
func GetContestant(db *sql.DB) gin.HandlerFunc {
	return func (c *gin.Context) {
		var contestant models.Contestant

		query := `
			SELECT id, election_id, name
			FROM contestants
			WHERE id = $1
		`

		err := db.QueryRow(
			query,
			c.Param("id"),
		).Scan(
			&contestant.ID,
			&contestant.ElectionID,
			&contestant.Name,
		)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Contestant not found"})
			return
		}

		c.JSON(http.StatusOK, contestant)
	}
}