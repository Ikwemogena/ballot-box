package handlers

import (
	"ballot-box/internal/modules/users/models"
	"ballot-box/internal/utils/auth"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterAdmin(db *sql.DB) gin.HandlerFunc  {
	return func(c *gin.Context) {
        var user models.User

        if err := c.ShouldBindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

		if UserExists(db, user.Username, user.Email) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
			return
		}

		hashedPassword := hashPassword(user.Password)
		user.Password = hashedPassword

		user.Role = "admin"
		query := "INSERT INTO users (username, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id"
		user.Role = "admin"
        err := db.QueryRow(query,
            user.Username, user.Email, user.Password, user.Role).Scan(&user.ID)

        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating the user."})
            return
        }

        c.JSON(http.StatusOK, user)
    }
}

func LoginAdmin(db *sql.DB) gin.HandlerFunc {
	return func (c *gin.Context) {

		var loginRequest struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hashedPassword := hashPassword(loginRequest.Password)

		var user struct {
            ID       string    `json:"id"`
            Username string `json:"username"`
            Email    string `json:"email"`
            Role     string `json:"role"`
        }

		query := "SELECT id, username, email, role FROM users WHERE username=$1 AND password=$2 AND role='admin'"
		err := db.QueryRow(query, loginRequest.Username, hashedPassword).Scan(&user.ID, &user.Username, &user.Email, &user.Role)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed"})
			}
			return
		}

		token, err := auth.GenerateJWT(user.ID, user.Username, user.Role)
	
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"accessToken": token, "user": user})

	}
}

func Register(db *sql.DB) gin.HandlerFunc  {

	return func(c *gin.Context) {
        var user models.User

        if err := c.ShouldBindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

		if UserExists(db, user.Username, user.Email) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
			return
		}

		hashedPassword := hashPassword(user.Password)
		user.Password = hashedPassword

		query := "INSERT INTO users (username, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id"
		user.Role = "voter"
        err := db.QueryRow(query,
            user.Username, user.Email, user.Password, user.Role).Scan(&user.ID)

        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating the user."})
			log.Println("Error creating the user", err)
            return
        }

        c.JSON(http.StatusOK, user)
    }
}

func Login(db *sql.DB) gin.HandlerFunc {
	return func (c *gin.Context) {

		var loginRequest struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hashedPassword := hashPassword(loginRequest.Password)

		var user struct {
            ID       string    `json:"id"`
            Username string `json:"username"`
            Email    string `json:"email"`
            Role     string `json:"role"`
        }

		query := "SELECT id, username, email, role FROM users WHERE username=$1 AND password=$2 AND role='voter'"
		err := db.QueryRow(query, loginRequest.Username, hashedPassword).Scan(&user.ID, &user.Username, &user.Email, &user.Role)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed"})
			}
			return
		}

		token, err := auth.GenerateJWT(user.ID, user.Username, user.Role)
	
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"accessToken": token, "user": user})

	}
}

func hashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	
	hashed := hash.Sum(nil)
	return hex.EncodeToString(hashed)
}

func GetAllUsers(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []models.User

		rows, err := db.Query("SELECT * FROM users")

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		defer rows.Close()

		for rows.Next() {
			var user models.User
			if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			users = append(users, user)
		}
	}
}

func UserExists(db *sql.DB, username, email string) bool {
	var user models.User
	query := "SELECT id, username, email, password, role FROM users WHERE username=$1 OR email=$2"
	err := db.QueryRow(query, username, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		return false
	}
	return true
}