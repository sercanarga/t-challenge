package routes

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"t-challenge/internal/durable"
	"t-challenge/internal/middleware"
	"time"

	"t-challenge/internal/model"
	"t-challenge/internal/response"
)

type sentBody struct {
	SenderAccountNumber   string  `json:"senderAccountNumber"`
	ReceiverAccountNumber string  `json:"receiverAccountNumber"`
	Amount                float64 `json:"amount"`
}

func Sent(mux *http.ServeMux) {
	mux.HandleFunc("POST /sent", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		userid := r.Context().Value("uuid").(string)

		var rBody sentBody
		if err := json.NewDecoder(r.Body).Decode(&rBody); err != nil {
			response.WriteResponse(w, &response.Response{
				Status:  http.StatusBadRequest,
				Success: false,
				Message: "Invalid request",
			})
			return
		}

		// fmt.Println(rBody.SenderAccountNumber, rBody.ReceiverAccountNumber, rBody.Amount, uuid)

		db := durable.Connection()
		tx := db.Begin()

		var senderAccount, receiverAccount model.Account
		var senderBalance, receiverBalance model.Balance
		if err := tx.Where("account_number = ? AND user_uuid = ?", rBody.SenderAccountNumber, userid).First(&senderAccount).Error; err != nil {
			tx.Rollback()
			response.WriteResponse(w, &response.Response{
				Status:  http.StatusNotFound,
				Success: false,
				Message: "Sender account not found",
			})
			return
		}

		if err := tx.Where("account_number = ? AND user_uuid != ?", rBody.ReceiverAccountNumber, userid).First(&receiverAccount).Error; err != nil {
			tx.Rollback()
			response.WriteResponse(w, &response.Response{
				Status:  http.StatusNotFound,
				Success: false,
				Message: "Receiver account not found",
			})
			return
		}

		if err := tx.Where("account_uuid = ?", senderAccount.UUID).First(&senderBalance).Error; err != nil {
			tx.Rollback()
			response.WriteResponse(w, &response.Response{
				Status:  http.StatusInternalServerError,
				Success: false,
				Message: "Sender balance not found",
			})
			return
		}

		if err := tx.Where("account_uuid = ?", receiverAccount.UUID).First(&receiverBalance).Error; err != nil {
			tx.Rollback()
			response.WriteResponse(w, &response.Response{
				Status:  http.StatusInternalServerError,
				Success: false,
				Message: "Receiver balance not found",
			})
			return
		}

		if senderBalance.Balance < rBody.Amount {
			tx.Rollback()
			response.WriteResponse(w, &response.Response{
				Status:  http.StatusBadRequest,
				Success: false,
				Message: "Insufficient funds",
			})
			return
		}

		senderBalance.Balance -= rBody.Amount
		receiverBalance.Balance += rBody.Amount

		if err := tx.Save(&senderBalance).Error; err != nil || tx.Save(&receiverBalance).Error != nil {
			tx.Rollback()
			response.WriteResponse(w, &response.Response{
				Status:  http.StatusInternalServerError,
				Success: false,
				Message: "Transaction failed",
			})
			return
		}

		newTransaction := model.Transaction{
			UUID:                (uuid.New()).String(),
			AccountUUID:         senderAccount.UUID,
			ReceiverAccountUUID: receiverAccount.UUID,
			Amount:              rBody.Amount,
			CreatedAt:           time.Now(), //@note: not required because autoCreateTime is set
		}

		if err := tx.Create(&newTransaction).Error; err != nil {
			tx.Rollback()
			response.WriteResponse(w, &response.Response{
				Status:  http.StatusInternalServerError,
				Success: false,
				Message: "Failed to record transaction",
			})
			return
		}

		tx.Commit()
		response.WriteResponse(w, &response.Response{
			Status:  http.StatusOK,
			Success: true,
			Message: "Transaction successful",
		})
	}))
}
