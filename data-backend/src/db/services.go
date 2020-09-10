package db

import (
	"context"
	"fmt"

	"billeroo.de/data-backend/src/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CUSTOMERS

func FindServicesByUserId(d *mongo.Database, userId string) ([]bson.M, error) {
	serviceCollection := d.Collection("services")

	var result []bson.M

	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return result, err
	}

	filter := bson.M{"userId": id}
	cursor, err := serviceCollection.Find(context.Background(), filter)
	if err != nil {
		return result, err
	}
	err = cursor.All(context.Background(), &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func UpdateServiceById(d *mongo.Database, uId string, cId string, service bson.M) (bson.M, error) {
	var result bson.M

	userId, err := primitive.ObjectIDFromHex(uId)
	if err != nil {
		return result, err
	}

	serviceId, err := primitive.ObjectIDFromHex(cId)
	if err != nil {
		return result, err
	}

	// remove these keys if present, otherwise they will object will be gone
	delete(service, "_id")
	delete(service, "userId")

	filter := bson.M{"userId": userId, "_id": serviceId}
	serviceCollection := d.Collection("services")
	_, err = serviceCollection.UpdateOne(context.Background(), filter, bson.M{"$set": service})

	if err != nil {
		return result, err
	}

	err = serviceCollection.FindOne(context.Background(), filter).Decode(&result)

	if err != nil {
		return result, err
	}

	return result, nil
}

func CreateService(d *mongo.Database, uId string, service bson.M) (bson.M, error) {
	var result bson.M

	userId, err := primitive.ObjectIDFromHex(uId)
	if err != nil {
		return result, err
	}

	service["userId"] = userId

	serviceCollection := d.Collection("services")
	res, err := serviceCollection.InsertOne(context.Background(), service)
	id := res.InsertedID

	filter := bson.M{"_id": id}
	err = serviceCollection.FindOne(context.Background(), filter).Decode(&result)

	return result, err
}

func FindOrCreateServiceByWPserviceId(database *mongo.Database, uId string, whServiceId string, billing models.WebhookLocationData, shipping models.WebhookLocationData) (bson.M, error) {
	servicesCollection := database.Collection("services")

	service := bson.M{}

	userId, err := primitive.ObjectIDFromHex(uId)
	if err != nil {
		return service, err
	}

	err = servicesCollection.FindOne(context.Background(), bson.M{"userId": userId, "whServiceId": whServiceId}).Decode(&service)
	if err == nil {
		return service, err
	}

	service["discount"] = 0
	billingName := billing.First_name + " " + billing.Last_name + "<br/>"
	billingCompany := billing.Company + "<br/>"
	billingAddress1 := billing.Address_1 + "<br/>"
	billingAddress2 := billing.Address_2 + "<br/>"
	billingPostCodeCity := billing.Postcode + " " + billing.City
	service["invoiceAddress"] = fmt.Sprintf("<p>%s  %s  %s %s %s</p>", billingName, billingCompany, billingAddress1, billingAddress2, billingPostCodeCity)

	service["name"] = billing.First_name + " " + billing.Last_name

	shippingName := shipping.First_name + " " + shipping.Last_name
	shippingCompany := shipping.Company + "<br/>"
	shippingAddress1 := shipping.Address_1
	shippingAddress2 := shipping.Address_2
	shippingPostCodeCity := shipping.Postcode + " " + shipping.City

	service["shippingAddress"] = fmt.Sprintf("<p>%s • %s  %s %s • %s</p>", shippingName, shippingCompany, shippingAddress1, shippingAddress2, shippingPostCodeCity)
	service["ust"] = 5
	service["whServiceId"] = whServiceId

	service, err = CreateService(database, uId, service)

	return service, err
}
