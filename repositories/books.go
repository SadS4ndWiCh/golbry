package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookRepository struct {
	db *pgxpool.Pool
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

func NewBookRepository(db *pgxpool.Pool) BookRepository {
	return BookRepository{db: db}
}

func (br BookRepository) GetById(id uint) (book Book, err error) {
	row := br.db.QueryRow(context.Background(), "SELECT * FROM books WHERE id=$1", id)
	err = row.Scan(&book.Id, &book.Title, &book.Author, &book.Year)

	return
}

func (br BookRepository) GetAll() (books []Book, err error) {
	rows, err := br.db.Query(context.Background(), "SELECT * FROM books")
	if err != nil {
		return
	}
	defer rows.Close()

	books, err = pgx.CollectRows(rows, pgx.RowToStructByName[Book])
	return
}

func (br BookRepository) InsertOne(book Book) (id uint, err error) {
	err = br.db.QueryRow(context.Background(), `
		INSERT INTO books (
			title, author, year
		) VALUES (
			$1, $2, $3
		) RETURNING id
	`, book.Title, book.Author, book.Year).Scan(&id)

	return
}

func (br BookRepository) DeleteById(id uint) (err error) {
	_, err = br.db.Exec(context.Background(), "DELETE FROM books WHERE id=$1", id)
	return
}
