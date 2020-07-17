package db

import (
	"context"
	"fmt"
	"strconv"

	"billeroo.de/data-backend/src/models"

	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CUSTOMERS

func FindSubscriptionsByUserId(d *mongo.Database, userId string) ([]bson.M, error) {
	subscriptionCollection := d.Collection("subscriptions")

	var result []bson.M

	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return result, err
	}

	filter := bson.M{"userId": id}
	cursor, err := subscriptionCollection.Find(context.Background(), filter)
	if err != nil {
		return result, err
	}
	err = cursor.All(context.Background(), &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func FindMostRecentSubscriptionByUserId(d *mongo.Database, userId string) (bson.M, error) {
	subscriptionCollection := d.Collection("subscriptions")

	var result bson.M

	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return result, err
	}
	options := options.FindOne()
	options.SetSort(bson.M{"_id": -1})
	filter := bson.M{"userId": id, "deleted": bson.M{"$exists": false}}
	err = subscriptionCollection.FindOne(context.Background(), filter, options).Decode(&result)

	return result, nil
}

func UpdateSubscriptionById(d *mongo.Database, uId string, cId string, subscription bson.M) (bson.M, error) {
	var result bson.M

	userId, err := primitive.ObjectIDFromHex(uId)
	if err != nil {
		return result, err
	}

	subscriptionId, err := primitive.ObjectIDFromHex(cId)
	if err != nil {
		return result, err
	}

	delete(subscription, "_id")
	delete(subscription, "userId")

	fmt.Println(subscription)

	filter := bson.M{"userId": userId, "_id": subscriptionId}
	subscriptionCollection := d.Collection("subscriptions")
	_, err = subscriptionCollection.UpdateOne(context.Background(), filter, bson.M{"$set": subscription})

	if err != nil {
		return result, err
	}

	err = subscriptionCollection.FindOne(context.Background(), filter).Decode(&result)

	if err != nil {
		return result, err
	}

	return result, nil
}

func CreateSubscription(d *mongo.Database, uId string, subscription bson.M) (bson.M, error) {
	var result bson.M

	userId, err := primitive.ObjectIDFromHex(uId)
	if err != nil {
		return result, err
	}

	subscription["userId"] = userId

	subscriptionCollection := d.Collection("subscriptions")
	res, err := subscriptionCollection.InsertOne(context.Background(), subscription)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	id := res.InsertedID

	filter := bson.M{"_id": id}
	err = subscriptionCollection.FindOne(context.Background(), filter).Decode(&result)

	return result, err
}

func FindOrCreateSubscriptionByWHSubscriptionId(database *mongo.Database, uId string, subscription models.WebhookLineItem) (bson.M, error) {
	subscriptionsCollection := database.Collection("subscriptions")

	result := bson.M{}

	userId, err := primitive.ObjectIDFromHex(uId)
	if err != nil {
		return result, err
	}

	err = subscriptionsCollection.FindOne(context.Background(), bson.M{"userId": userId, "whSubscriptionId": strconv.Itoa(subscription.Product_id)}).Decode(&result)
	if err != nil {
		result = bson.M{}

		result["name"] = subscription.Name
		result["isbn"] = ""
		result["price"] = subscription.Price
		result["amount"] = 100
		result["whSubscriptionId"] = strconv.Itoa(subscription.Product_id)

		result, err = CreateSubscription(database, uId, result)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
	}

	result["toBeSend"] = subscription.Quantity
	result["toBePayed"] = subscription.Quantity

	return result, nil
}

func FindOrCreateSubscriptionsByWPlineItems(database *mongo.Database, userId string, subscriptions []models.WebhookLineItem) ([]bson.M, error) {
	results := []bson.M{}

	for _, a := range subscriptions {
		res, err := FindOrCreateSubscriptionByWHSubscriptionId(database, userId, a)
		if err != nil {
			return nil, nil
		}

		results = append(results, res)
	}

	return results, nil
}
