package db

import (
	"context"

	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CUSTOMERS

func FindArticleAuthorsByUserId(d *mongo.Database, userId string) ([]bson.M, error) {
	articleAuthorCollection := d.Collection("articleAuthors")

	var result []bson.M

	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return result, err
	}

	filter := bson.M{"userId": id}
	cursor, err := articleAuthorCollection.Find(context.Background(), filter)
	if err != nil {
		return result, err
	}
	err = cursor.All(context.Background(), &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func UpdateArticleAuthorById(d *mongo.Database, uId string, cId string, articleAuthor bson.M) (bson.M, error) {
	var result bson.M

	userId, err := primitive.ObjectIDFromHex(uId)
	if err != nil {
		return result, err
	}

	articleAuthorId, err := primitive.ObjectIDFromHex(cId)
	if err != nil {
		return result, err
	}

	delete(articleAuthor, "authorId")
	delete(articleAuthor, "articleId")
	delete(articleAuthor, "_id")
	delete(articleAuthor, "userId")

	filter := bson.M{"userId": userId, "_id": articleAuthorId}
	articleAuthorCollection := d.Collection("articleAuthors")
	_, err = articleAuthorCollection.UpdateOne(context.Background(), filter, bson.M{"$set": articleAuthor})

	if err != nil {
		return result, err
	}

	err = articleAuthorCollection.FindOne(context.Background(), filter).Decode(&result)

	if err != nil {
		return result, err
	}

	return result, nil
}

func CreateArticleAuthor(d *mongo.Database, uId string, auId string, arId string, articleAuthor bson.M) (bson.M, error) {
	var result bson.M

	userId, err := primitive.ObjectIDFromHex(uId)
	if err != nil {
		return result, err
	}

	articleId, err := primitive.ObjectIDFromHex(arId)
	if err != nil {
		return result, err
	}

	authorId, err := primitive.ObjectIDFromHex(auId)
	if err != nil {
		return result, err
	}

	articleAuthor["userId"] = userId
	articleAuthor["authorId"] = authorId
	articleAuthor["articleId"] = articleId

	articleAuthorCollection := d.Collection("articleAuthors")
	res, err := articleAuthorCollection.InsertOne(context.Background(), articleAuthor)
	id := res.InsertedID

	filter := bson.M{"_id": id}
	err = articleAuthorCollection.FindOne(context.Background(), filter).Decode(&result)

	return result, err
}
