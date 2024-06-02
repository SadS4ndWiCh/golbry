package main

import (
	"fmt"
	"golbry/internals/server"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	srv := server.NewServer()

	fmt.Println("🪅 Server is running")
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
