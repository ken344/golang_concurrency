package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"zin-golang_concurrency/calculator"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		return
	}
	fmt.Println(os.Getenv("GO_ENV"))
	fmt.Println(calculator.Offset)
	fmt.Println(calculator.Sum(1, 2))
	fmt.Println(calculator.Multiply(1, 2))
}
