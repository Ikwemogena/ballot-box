package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("secret-key")

type Claims struct {
    UserID int `json:"user_id"`
	UserRole string `json:"user_role"`
	UserUsername string `json:"user_username"`
    jwt.Claims
}

func GenerateJWT(id, username, role string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, 
        jwt.MapClaims{ 
        "username": username, 
		"role": role,
		"id": id,
        "exp": time.Now().Add(time.Hour * 24).Unix(), 
        })

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
      return secretKey, nil
   })
  
   if err != nil {
      return nil, err
   }
  
   if !token.Valid {
      return nil, err
   }
  
   return token, nil
}
