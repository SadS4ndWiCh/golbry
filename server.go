package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"golbry/database"
	"golbry/repositories"
	"golbry/utils"
	"log"
	"net/http"
	"strconv"
)

func main() {
	db := database.ConnectDB()
	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /books", func(w http.ResponseWriter, r *http.Request) {
		bookRepository := repositories.NewBookRepository(db)
		books, err := bookRepository.GetAll()
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		json, err := json.Marshal(books)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Write(json)
	})

	mux.HandleFunc("GET /books/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "Invalid id", http.StatusBadRequest)
			return
		}

		bookRepository := repositories.NewBookRepository(db)
		book, err := bookRepository.GetById(uint(id))
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		json, err := json.Marshal(book)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Write(json)
	})

	mux.HandleFunc("POST /books", func(w http.ResponseWriter, r *http.Request) {
		var newBook repositories.Book

		if err := utils.DecodeJSONBody(w, r, &newBook); err != nil {
			var malformedRequest *utils.MalformedRequest

			if errors.As(err, &malformedRequest) {
				http.Error(w, malformedRequest.Error(), malformedRequest.Status)
			} else {
				log.Print(err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}

			return
		}

		bookRepository := repositories.NewBookRepository(db)
		if _, err := bookRepository.InsertOne(newBook); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(201)
	})

	mux.HandleFunc("DELETE /books", func(w http.ResponseWriter, r *http.Request) {
		var deleteBook repositories.DeleteBook

		if err := utils.DecodeJSONBody(w, r, &deleteBook); err != nil {
			var malformedRequest utils.MalformedRequest

			if errors.As(err, &malformedRequest) {
				http.Error(w, malformedRequest.Error(), malformedRequest.Status)
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}

			return
		}

		bookRepository := repositories.NewBookRepository(db)
		if err := bookRepository.DeleteById(deleteBook.Id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	fmt.Println("🪅 Server is running at :3000")

	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Fatal(err)
	}
}
