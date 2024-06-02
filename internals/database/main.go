package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	DB *pgxpool.Pool
}

var (
	dbName   = os.Getenv("DB_DATABASE")
	username = os.Getenv("DB_USERNAME")
	password = os.Getenv("DB_PASSWORD")
	host     = os.Getenv("DB_HOST")
	port     = os.Getenv("DB_PORT")
	instance *Service
)

func New() *Service {
	if instance != nil {
		return instance
	}

	databaseUrl := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", username, password, host, port, dbName)
	conn, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	instance = &Service{
		DB: conn,
	}

	return instance
}

func (s *Service) Close() {
	log.Println("Database closed")
	s.DB.Close()
}
