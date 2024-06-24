package handlers

import (
	"ballot-box/internal/modules/votes/models"
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CastVote(db *sql.DB) gin.HandlerFunc {
    return func (c *gin.Context) {
        var vote models.Vote

        if err := c.ShouldBindJSON(&vote); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        voterID := c.GetString("userID")

        query := `
            INSERT INTO votes (election_id, contestant_id, voter_id)
            VALUES ($1, $2, $3)
            RETURNING id
        `

        vote.VoterID = "a623f292-5532-411d-9809-ebfab8e80523"
        log.Println(voterID)

        err := db.QueryRow(
            query,
            vote.ElectionID,
            vote.ContestantID,
            vote.VoterID,
        ).Scan(&vote.ID)

        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not cast vote", "message": err.Error()})
            return
        }

        c.JSON(http.StatusOK, vote)
    }
}
