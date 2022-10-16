package main

import (
	"fga-final/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("can't open .env file ", err)
	}

	dbConfig := database.PGConfig{
		Host:     os.Getenv("POSTGRES_ADDR"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Database: os.Getenv("POSTGRES_DB"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
	}

	_, err = database.StartDB(dbConfig)
	if err != nil {
		log.Fatalln("can't open database ", err)
	}
}
