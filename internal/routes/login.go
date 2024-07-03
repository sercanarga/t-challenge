package routes

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"os"
	"t-challenge/internal/durable"
	"t-challenge/internal/model"
	"t-challenge/internal/response"
	"time"
)

var privateKey *rsa.PrivateKey

type loginBody struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func init() {
	privateKeyBytes, err := os.ReadFile("cert/private_key.pem")
	if err != nil {
		log.Fatal("could not read private key")
	}

	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		log.Fatal("could not parse private key")
	}
}

func Login(mux *http.ServeMux) {
	mux.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		var lBody loginBody
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &lBody)

		validate := validator.New()
		err := validate.Struct(lBody)
		if err != nil {
			response.WriteResponse(w, &response.Response{
				Status:  http.StatusBadRequest,
				Success: false,
				Message: "Bad request",
			})
			return
		}

		var user model.User
		result := durable.Connection().Model(&model.User{}).Where("email = ?", lBody.Email).First(&user)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			response.WriteResponse(w, &response.Response{
				Status:  http.StatusForbidden,
				Success: false,
				Message: "Username or password is wrong!",
			})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(lBody.Password))
		if err != nil {
			response.WriteResponse(w, &response.Response{
				Status:  http.StatusForbidden,
				Success: false,
				Message: "Username or password is wrong!",
			})
			return
		}

		token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
			"exp":  time.Now().Add((time.Hour * 24) * 7).Unix(), // 1 week
			"uuid": user.UUID,
		}).SignedString(privateKey)
		if err != nil {
			response.WriteResponse(w, &response.Response{
				Status:  http.StatusInternalServerError,
				Success: false,
				Message: "Token signing failed",
			})
			return
		}

		response.WriteResponse(w, &response.Response{
			Status:  http.StatusOK,
			Success: true,
			Message: "Login successful",
			Token:   &token,
		})
	})
}
