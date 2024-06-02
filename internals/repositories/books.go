package repositories

import (
	"context"
	"golbry/internals/database"

	"github.com/jackc/pgx/v5"
)

type BookRepository struct {
	service *database.Service
}

type Book struct {
	Id     uint   `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   uint16 `json:"year"`
}

type DeleteBook struct {
	Id uint `json:"id"`
}

func NewBookRepository(service *database.Service) BookRepository {
	return BookRepository{service: service}
}

func (br BookRepository) GetById(id uint) (book Book, err error) {
	row := br.service.DB.QueryRow(context.Background(), "SELECT * FROM books WHERE id=$1", id)
	err = row.Scan(&book.Id, &book.Title, &book.Author, &book.Year)

	return
}

func (br BookRepository) GetAll() (books []Book, err error) {
	rows, err := br.service.DB.Query(context.Background(), "SELECT * FROM books")
	if err != nil {
		return
	}
	defer rows.Close()

	books, err = pgx.CollectRows(rows, pgx.RowToStructByName[Book])
	return
}

func (br BookRepository) InsertOne(book Book) (id uint, err error) {
	err = br.service.DB.QueryRow(context.Background(), `
		INSERT INTO books (
			title, author, year
		) VALUES (
			$1, $2, $3
		) RETURNING id
	`, book.Title, book.Author, book.Year).Scan(&id)

	return
}

func (br BookRepository) DeleteById(id uint) (err error) {
	_, err = br.service.DB.Exec(context.Background(), "DELETE FROM books WHERE id=$1", id)
	return
}
