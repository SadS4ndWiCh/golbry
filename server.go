package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"golbry/utils"
	"log"
	"net/http"
	"strconv"
)

type Book struct {
	Id     uint   `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   uint16 `json:"year"`
}

type NewBook struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   uint16 `json:"year"`
}

type DeleteBook struct {
	Id uint `json:"id"`
}

var store []Book = make([]Book, 0)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /books", func(w http.ResponseWriter, r *http.Request) {
		json, err := json.Marshal(store)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
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

		for _, book := range store {
			if book.Id == uint(id) {
				json, err := json.Marshal(book)
				if err != nil {
					http.Error(w, "Something went wrong", http.StatusInternalServerError)
					return
				}

				w.Write(json)
				return
			}
		}
	})

	mux.HandleFunc("POST /books", func(w http.ResponseWriter, r *http.Request) {
		var newBook NewBook

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

		store = append(store, Book{
			Id:     uint(len(store) + 1),
			Title:  newBook.Title,
			Author: newBook.Author,
			Year:   newBook.Year,
		})

		w.WriteHeader(201)
	})

	mux.HandleFunc("DELETE /books", func(w http.ResponseWriter, r *http.Request) {
		var deleteBook DeleteBook

		if err := utils.DecodeJSONBody(w, r, &deleteBook); err != nil {
			var malformedRequest utils.MalformedRequest

			if errors.As(err, &malformedRequest) {
				http.Error(w, malformedRequest.Error(), malformedRequest.Status)
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}

			return
		}

		updatedBooks := make([]Book, 0)
		for _, book := range store {
			if deleteBook.Id != book.Id {
				updatedBooks = append(updatedBooks, book)
			}
		}

		store = updatedBooks
	})

	fmt.Println("🪅 Server is running at :3000")

	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Fatal(err)
	}
}
