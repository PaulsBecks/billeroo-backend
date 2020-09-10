package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CUSTOMERS

func FindCompanyByUserId(d *mongo.Database, userId string) (bson.M, error) {
	companyCollection := d.Collection("companies")

	var result bson.M

	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return result, err
	}

	filter := bson.M{"userId": id}
	err = companyCollection.FindOne(context.Background(), filter).Decode(&result)

	if err != nil {
		return result, err
	}

	return result, nil
}

func UpdateCompanyById(d *mongo.Database, uId string, cId string, company bson.M) (bson.M, error) {
	var result bson.M

	userId, err := primitive.ObjectIDFromHex(uId)
	if err != nil {
		return result, err
	}

	companyId, err := primitive.ObjectIDFromHex(cId)
	if err != nil {
		return result, err
	}

	delete(company, "_id")
	delete(company, "userId")

	filter := bson.M{"userId": userId, "_id": companyId}
	companyCollection := d.Collection("companies")
	_, err = companyCollection.UpdateOne(context.Background(), filter, bson.M{"$set": company})

	if err != nil {
		return result, err
	}

	err = companyCollection.FindOne(context.Background(), filter).Decode(&result)

	if err != nil {
		return result, err
	}

	return result, nil
}
