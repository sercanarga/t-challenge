package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	// Common
	Status  int    `json:"status"`
	Success bool   `json:"success"`
	Message string `json:"message"`

	// Auth
	Token *string `json:"token,omitempty"`

	// Accounts
	Accounts []Account `json:"accounts,omitempty"`
}

type Account struct {
	AccountNumber string  `json:"account_number"`
	Balance       float64 `json:"balance"`
}

func WriteResponse(w http.ResponseWriter, response *Response) {
	jsonBody, _ := json.Marshal(response)

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(response.Status)
	_, err := w.Write(jsonBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
