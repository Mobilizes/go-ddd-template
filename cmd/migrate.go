package cmd

import (
	"fmt"

	"github.com/joho/godotenv"
)

func Migrate() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("error loading .env file: %v", err)
		return
	}
}
