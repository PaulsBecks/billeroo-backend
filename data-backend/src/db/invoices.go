package db

import (
	"context"

	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CUSTOMERS

func FindInvoicesByUserId(d *mongo.Database, userId string) ([]bson.M, error) {
	invoiceCollection := d.Collection("invoices")

	var result []bson.M

	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return result, err
	}
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"_id": -1})

	filter := bson.M{"userId": id}

	cursor, err := invoiceCollection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return result, err
	}
	err = cursor.All(context.Background(), &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func UpdateInvoiceById(d *mongo.Database, uId string, cId string, invoice bson.M) (bson.M, error) {
	var result bson.M

	userId, err := primitive.ObjectIDFromHex(uId)
	if err != nil {
		return result, err
	}

	invoiceId, err := primitive.ObjectIDFromHex(cId)
	if err != nil {
		return result, err
	}

	delete(invoice, "_id")
	delete(invoice, "userId")

	filter := bson.M{"userId": userId, "_id": invoiceId}
	invoiceCollection := d.Collection("invoices")
	_, err = invoiceCollection.UpdateOne(context.Background(), filter, bson.M{"$set": invoice})

	if err != nil {
		return result, err
	}

	err = invoiceCollection.FindOne(context.Background(), filter).Decode(&result)

	if err != nil {
		return result, err
	}

	return result, nil
}

func CreateInvoice(d *mongo.Database, uId string, invoice bson.M) (bson.M, error) {
	var result bson.M

	userId, err := primitive.ObjectIDFromHex(uId)
	if err != nil {
		return result, err
	}

	invoice["userId"] = userId

	invoiceCollection := d.Collection("invoices")
	res, err := invoiceCollection.InsertOne(context.Background(), invoice)
	id := res.InsertedID

	filter := bson.M{"_id": id}
	err = invoiceCollection.FindOne(context.Background(), filter).Decode(&result)

	return result, err
}
