package routes

import (
	"encoding/json"
	"errors"
	"golbry/repositories"
	"golbry/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BookHandler struct {
	bookRepository repositories.BookRepository
}

func NewBookHandler(db *pgxpool.Pool) BookHandler {
	return BookHandler{bookRepository: repositories.NewBookRepository(db)}
}

func (bh *BookHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	books, err := bh.bookRepository.GetAll()
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
}

func (bh *BookHandler) GetById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	book, err := bh.bookRepository.GetById(uint(id))
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
}

func (bh *BookHandler) InsertOne(w http.ResponseWriter, r *http.Request) {
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

	if _, err := bh.bookRepository.InsertOne(newBook); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(201)
}

func (bh *BookHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
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

	if err := bh.bookRepository.DeleteById(deleteBook.Id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
