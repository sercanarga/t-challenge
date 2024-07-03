package routes

import (
	"net/http"
	"t-challenge/internal/durable"
	"t-challenge/internal/middleware"
	"t-challenge/internal/response"
)

func getAccounts(uuid string) ([]response.Account, error) {
	var accounts []response.Account

	accountsResult := durable.Connection().Table("accounts").Select("accounts.account_number, balances.balance").
		Joins("left join balances on balances.account_uuid = accounts.uuid").
		Where("accounts.user_uuid = ?", uuid).Scan(&accounts)

	return accounts, accountsResult.Error
}

func MyAccounts(mux *http.ServeMux) {
	mux.HandleFunc("/my-accounts", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		uuid := r.Context().Value("uuid").(string)

		// get user accounts
		accounts, err := getAccounts(uuid)
		if err != nil {
			response.WriteResponse(w, &response.Response{
				Status:  http.StatusInternalServerError,
				Success: false,
				Message: "Internal server error",
			})
		}

		response.WriteResponse(w, &response.Response{
			Status:   http.StatusOK,
			Message:  "Accounts retrieved successfully",
			Success:  true,
			Accounts: accounts,
		})
	}))
}
