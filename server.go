package main

import (
	"fmt"
	"golbry/database"
	"golbry/routes"
	"log"
	"net/http"
)

func main() {
	db := database.ConnectDB()
	defer db.Close()

	mux := http.NewServeMux()

	bookHandler := routes.NewBookHandler(db)
	mux.HandleFunc("GET /books", bookHandler.GetAll)
	mux.HandleFunc("GET /books/{id}", bookHandler.GetById)
	mux.HandleFunc("POST /books", bookHandler.InsertOne)
	mux.HandleFunc("DELETE /books", bookHandler.DeleteById)

	fmt.Println("🪅 Server is running at :3000")

	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Fatal(err)
	}
}
