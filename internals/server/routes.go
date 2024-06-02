package server

import (
	"net/http"

	books "golbry/internals/http"
)

func (s *Server) Bootstrap() *http.ServeMux {
	api := http.NewServeMux()

	bookHandler := books.NewBookHandler(s.service)
	api.HandleFunc("GET /books", bookHandler.GetAll)
	api.HandleFunc("GET /books/{id}", bookHandler.GetById)
	api.HandleFunc("POST /books", bookHandler.InsertOne)
	api.HandleFunc("DELETE /books", bookHandler.DeleteById)

	return api
}
