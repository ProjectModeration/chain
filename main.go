package main

import (
	"ProjectModeration/chain/chain"
	"fmt"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	chain.StartChaining(2578983379) // flawless
}
