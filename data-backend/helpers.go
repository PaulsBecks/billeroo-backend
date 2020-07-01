package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Load env variables from .env file
func getEnvVariable(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error when loading env")
	}

	return os.Getenv(key)
}

// Mongo DB

func connectToClient() *mongo.Client {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(getEnvVariable("MONGO_DB_URI")))
	if err != nil {
		panic(err)
	}
	//defer client.Disconnect(ctx)
	return client
}
