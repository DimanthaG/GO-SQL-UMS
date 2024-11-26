
package auth

import (
    "github.com/golang-jwt/jwt/v4"
    "time"
)

var jwtKey = []byte("your_secret_key")

func GenerateJWT(email string) (string, error) {
    claims := jwt.MapClaims{
        "email": email,
        "exp":   time.Now().Add(24 * time.Hour).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}
