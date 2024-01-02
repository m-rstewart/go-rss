package main

import "github.com/joho/godotenv"

func main() {
	godotenv.Load()
	port := godotenv.Load("PORT")
}
