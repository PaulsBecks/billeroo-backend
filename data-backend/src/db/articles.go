package db

import (
	"context"
	"fmt"
	"strconv"

	"billeroo.de/data-backend/src/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CUSTOMERS

func FindArticlesByUserId(d *mongo.Database, userId string) ([]bson.M, error) {
	articleCollection := d.Collection("articles")

	var result []bson.M

	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return result, err
	}

	filter := bson.M{"userId": id}
	cursor, err := articleCollection.Find(context.Background(), filter)
	if err != nil {
		return result, err
	}
	err = cursor.All(context.Background(), &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func UpdateArticleById(d *mongo.Database, uId string, cId string, article bson.M) (bson.M, error) {
	var result bson.M

	userId, err := primitive.ObjectIDFromHex(uId)
	if err != nil {
		return result, err
	}

	articleId, err := primitive.ObjectIDFromHex(cId)
	if err != nil {
		return result, err
	}

	delete(article, "_id")
	delete(article, "userId")

	fmt.Println(article)

	filter := bson.M{"userId": userId, "_id": articleId}
	articleCollection := d.Collection("articles")
	_, err = articleCollection.UpdateOne(context.Background(), filter, bson.M{"$set": article})

	if err != nil {
		return result, err
	}

	err = articleCollection.FindOne(context.Background(), filter).Decode(&result)

	if err != nil {
		return result, err
	}

	return result, nil
}

func CreateArticle(d *mongo.Database, uId string, article bson.M) (bson.M, error) {
	var result bson.M

	userId, err := primitive.ObjectIDFromHex(uId)
	if err != nil {
		return result, err
	}

	article["userId"] = userId

	articleCollection := d.Collection("articles")
	res, err := articleCollection.InsertOne(context.Background(), article)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	id := res.InsertedID

	filter := bson.M{"_id": id}
	err = articleCollection.FindOne(context.Background(), filter).Decode(&result)

	return result, err
}

func FindOrCreateArticleByWHArticleId(database *mongo.Database, uId string, article models.WebhookLineItem) (bson.M, error) {
	articlesCollection := database.Collection("articles")

	result := bson.M{}

	userId, err := primitive.ObjectIDFromHex(uId)
	if err != nil {
		return result, err
	}

	err = articlesCollection.FindOne(context.Background(), bson.M{"userId": userId, "whArticleId": strconv.Itoa(article.Product_id)}).Decode(&result)
	if err != nil {
		result = bson.M{}

		result["name"] = article.Name
		result["isbn"] = ""
		result["price"] = article.Price
		result["amount"] = 100
		result["whArticleId"] = strconv.Itoa(article.Product_id)

		result, err = CreateArticle(database, uId, result)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
	}

	result["toBeSend"] = article.Quantity
	result["toBePayed"] = article.Quantity

	return result, nil
}

func FindOrCreateArticlesByWPlineItems(database *mongo.Database, userId string, articles []models.WebhookLineItem) ([]bson.M, error) {
	results := []bson.M{}

	for _, a := range articles {
		res, err := FindOrCreateArticleByWHArticleId(database, userId, a)
		if err != nil {
			return nil, nil
		}

		results = append(results, res)
	}

	return results, nil
}
