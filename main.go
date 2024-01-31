package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	uri := os.Getenv("MONGO_CONN_STR")
	if uri == "" {
		log.Fatal("MONGO_CONN_STR not found in environments")
	}
	client, err := initStorage(uri)

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("DBName not found in environments")
	}

	userCollection := client.Database(dbName).Collection("user")

	storage := NewMongoStorage(userCollection)
	service := NewServiceImpl(storage)

	server := NewAPIServer(":3000", service)

	// run server
	log.Printf("Listen on Port : %d", 3000)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}
