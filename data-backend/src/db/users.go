package db

import (
	"context"

	"billeroo.de/data-backend/src/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindUserById(database *mongo.Database, userId string) (models.User, error) {
	usersCollection := database.Collection("users")

	var user models.User

	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return user, err
	}

	err = usersCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)

	return user, err
}
