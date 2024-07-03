package middleware

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recovered from panic: %v", err)
				fmt.Printf("Recovered from panic: %v \n", err)
				fmt.Printf("Stack trace:\n%s\n\n", string(debug.Stack()))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
