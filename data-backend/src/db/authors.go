package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CUSTOMERS

func FindAuthorsByUserId(d *mongo.Database, userId string) ([]bson.M, error) {
	authorCollection := d.Collection("authors")

	var result []bson.M

	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return result, err
	}

	filter := bson.M{"userId": id}
	cursor, err := authorCollection.Find(context.Background(), filter)
	if err != nil {
		return result, err
	}
	err = cursor.All(context.Background(), &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func UpdateAuthorById(d *mongo.Database, uId string, cId string, author bson.M) (bson.M, error) {
	var result bson.M

	userId, err := primitive.ObjectIDFromHex(uId)
	if err != nil {
		return result, err
	}

	authorId, err := primitive.ObjectIDFromHex(cId)
	if err != nil {
		return result, err
	}

	delete(author, "_id")
	delete(author, "userId")

	filter := bson.M{"userId": userId, "_id": authorId}
	authorCollection := d.Collection("authors")
	_, err = authorCollection.UpdateOne(context.Background(), filter, bson.M{"$set": author})

	if err != nil {
		return result, err
	}

	err = authorCollection.FindOne(context.Background(), filter).Decode(&result)

	if err != nil {
		return result, err
	}

	return result, nil
}

func CreateAuthor(d *mongo.Database, uId string, author bson.M) (bson.M, error) {
	var result bson.M

	userId, err := primitive.ObjectIDFromHex(uId)
	if err != nil {
		return result, err
	}

	author["userId"] = userId

	authorCollection := d.Collection("authors")
	res, err := authorCollection.InsertOne(context.Background(), author)
	id := res.InsertedID

	filter := bson.M{"_id": id}
	err = authorCollection.FindOne(context.Background(), filter).Decode(&result)

	return result, err
}
