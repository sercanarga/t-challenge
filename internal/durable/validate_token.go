package durable

import (
	"crypto/rsa"
	"errors"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var publicKey *rsa.PublicKey

func init() {
	publicKeyBytes, err := os.ReadFile("cert/public_key.pem")
	if err != nil {
		log.Fatal("could not read private key")
	}

	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		log.Fatal("could not parse private key")
	}
}

func ValidateToken(r *http.Request) (string, error) {
	var uuid string

	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if authHeader == "" {
		return "", errors.New("no authorization header")
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return publicKey, nil
	})

	if token != nil && !token.Valid {
		return "", errors.New("invalid token")
	}

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if expValue, ok := claims["exp"].(float64); ok {
			if time.Unix(int64(expValue), 0).Before(time.Now()) {
				return "", errors.New("token is expired")
			}
		}

		if uuidValue, ok := claims["uuid"].(string); ok {
			uuid = uuidValue
		}
	} else {
		return "", errors.New("invalid token claims")
	}

	return uuid, nil
}
