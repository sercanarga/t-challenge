package internal

import (
	"fmt"
	"log"
	"net/http"
	"t-challenge/internal/middleware"
	"t-challenge/internal/routes"
)

type Server struct{}

func (s *Server) SetupRoutes(mux *http.ServeMux) {
	routes.Health(mux)

	routes.Login(mux)
	routes.Register(mux)

	routes.MyAccounts(mux)
	routes.Sent(mux)

}

func (s *Server) SetupMiddleware(handler http.Handler) http.Handler {
	// @note: in a live project more middleware will be needed. like cors, logger, rate limiter, etc.
	panicHandler := middleware.RecoverPanic(handler)
	return panicHandler
}

func (s *Server) StartServer(handler http.Handler, port string) {
	fmt.Printf("server is starting on port: %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
