package db

import (
	"context"
	"fmt"

	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindWebhooksByUserId(database *mongo.Database, uId string) ([]bson.M, error) {
	webhooksCollection := database.Collection("webhooks")

	var result []bson.M

	id, err := primitive.ObjectIDFromHex(uId)

	if err != nil {
		fmt.Println(err)
		return result, err
	}

	courser, err := webhooksCollection.Find(context.Background(), bson.M{"userId": id})

	if err != nil {
		fmt.Println(err)
		return result, err
	}

	err = courser.All(context.Background(), &result)

	return result, err
}

func UpdateWebhookById(d *mongo.Database, uId string, cId string, webhook bson.M) (bson.M, error) {
	var result bson.M

	userId, err := primitive.ObjectIDFromHex(uId)
	if err != nil {
		return result, err
	}

	webhookId, err := primitive.ObjectIDFromHex(cId)
	if err != nil {
		return result, err
	}

	delete(webhook, "_id")
	delete(webhook, "userId")

	filter := bson.M{"userId": userId, "_id": webhookId}
	webhookCollection := d.Collection("webhooks")
	_, err = webhookCollection.UpdateOne(context.Background(), filter, bson.M{"$set": webhook})

	if err != nil {
		return result, err
	}

	err = webhookCollection.FindOne(context.Background(), filter).Decode(&result)

	if err != nil {
		return result, err
	}

	return result, nil
}

func CreateWebhook(d *mongo.Database, uId string, webhook bson.M) (bson.M, error) {
	var result bson.M

	userId, err := primitive.ObjectIDFromHex(uId)
	if err != nil {
		return result, err
	}

	webhook["userId"] = userId

	webhookCollection := d.Collection("webhooks")
	res, err := webhookCollection.InsertOne(context.Background(), webhook)
	id := res.InsertedID

	filter := bson.M{"_id": id}
	err = webhookCollection.FindOne(context.Background(), filter).Decode(&result)

	return result, err
}

func FindWebhookById(database *mongo.Database, id string) (bson.M, error) {
	var result bson.M

	webhooksCollection := database.Collection("webhooks")

	webhookId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": webhookId}
	err = webhooksCollection.FindOne(context.Background(), filter).Decode(&result)

	return result, err

}
