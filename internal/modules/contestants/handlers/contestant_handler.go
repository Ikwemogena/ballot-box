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