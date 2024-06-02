package server

import (
	"fmt"
	"golbry/internals/database"
	"net/http"
	"os"
)

type Server struct {
	addr    string
	service *database.Service
}

func NewServer() *http.Server {
	port := os.Getenv("PORT")
	srv := &Server{
		addr:    fmt.Sprintf(":%s", port),
		service: database.New(),
	}

	return &http.Server{
		Addr:    srv.addr,
		Handler: srv.Bootstrap(),
	}
}
