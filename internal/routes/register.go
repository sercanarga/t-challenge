package routes

import (
	"encoding/json"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"io"
	"net/http"
	"t-challenge/internal/durable"
	"t-challenge/internal/model"
	"t-challenge/internal/response"
)

type registerBody struct {
	Name     string `json:"name" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /register", func(w http.ResponseWriter, r *http.Request) {
		var rBody registerBody
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &rBody)

		validate := validator.New()
		err := validate.Struct(rBody)
		if err != nil {
			response.WriteResponse(w, &response.Response{
				Status:  http.StatusBadRequest,
				Success: false,
				Message: "Bad request",
			})
			return
		}

		// Check if user already exists
		if durable.Connection().Where("email = ?", rBody.Email).First(&model.User{}).RowsAffected != 0 {
			response.WriteResponse(w, &response.Response{
				Status:  http.StatusForbidden,
				Success: false,
				Message: "User already exists",
			})
			return
		}

		hashedPassword, err := durable.HashPassword(rBody.Password)
		if err != nil {
			response.WriteResponse(w, &response.Response{
				Status:  http.StatusInternalServerError,
				Success: false,
				Message: "Internal server error",
			})
			return
		}

		// create user
		userUUID := uuid.New()
		newUser := &model.User{
			UUID:     userUUID.String(),
			Name:     rBody.Name,
			Email:    rBody.Email,
			Password: hashedPassword,
			Status:   true,
		}
		if err := durable.Connection().Create(newUser).Error; err != nil {
			response.WriteResponse(w, &response.Response{
				Status:  http.StatusInternalServerError,
				Success: false,
				Message: "Internal server error",
			})
			return
		}

		// create account
		// @todo: unique status will be check.
		// known issue: created account number can throw conflict problems in the future.
		account := &model.Account{
			UUID:          (uuid.New()).String(),
			UserUUID:      userUUID.String(),
			AccountNumber: durable.GenerateAccountNumber(),
		}
		if err := durable.Connection().Create(account).Error; err != nil {
			response.WriteResponse(w, &response.Response{
				Status:  http.StatusInternalServerError,
				Success: false,
				Message: "Internal server error",
			})
			return
		}

		// create balance
		balance := &model.Balance{
			UUID:        (uuid.New()).String(),
			AccountUUID: account.UUID,
			Balance:     durable.GenerateAmount(), // initial balance
		}
		if err := durable.Connection().Create(balance).Error; err != nil {
			response.WriteResponse(w, &response.Response{
				Status:  http.StatusInternalServerError,
				Success: false,
				Message: "Internal server error",
			})
			return
		}

		response.WriteResponse(w, &response.Response{
			Status:  http.StatusCreated,
			Success: true,
			Message: "User created",
		})
	})
}
