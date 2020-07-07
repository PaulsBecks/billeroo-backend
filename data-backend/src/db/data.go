package db

import (
	"context"
	"fmt"

	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindDataByUserId(d *mongo.Database, userId string) (bson.M, error) {
	dataCollection := d.Collection("data")

	var result bson.M

	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return result, err
	}
	filter := bson.M{"userId": id}
	err = dataCollection.FindOne(context.Background(), filter).Decode(&result)

	if err != nil {
		return result, err
	}
	return result, nil
}

func UpdateDataByUserId(d *mongo.Database, userId string, data bson.M) (bson.M, error) {
	dataCollection := d.Collection("data")
	fmt.Println(data, userId)
	var result bson.M

	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return result, err
	}

	filter := bson.M{"userId": id}

	_, err = dataCollection.UpdateOne(
		context.Background(),
		filter,
		bson.M{"$set": data},
	)

	if err != nil {
		return result, err
	}

	err = dataCollection.FindOne(context.Background(), filter).Decode(&result)

	if err != nil {
		return result, err
	}

	return result, nil

}
