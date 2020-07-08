package db

import (
	"context"
	"fmt"

	"billeroo.de/data-backend/src/models"
	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CUSTOMERS

func FindCustomersByUserId(d *mongo.Database, userId string) ([]bson.M, error) {
	customerCollection := d.Collection("customers")

	var result []bson.M

	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return result, err
	}

	filter := bson.M{"userId": id}
	cursor, err := customerCollection.Find(context.Background(), filter)
	if err != nil {
		return result, err
	}
	err = cursor.All(context.Background(), &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func UpdateCustomerById(d *mongo.Database, uId string, cId string, customer bson.M) (bson.M, error) {
	var result bson.M

	userId, err := primitive.ObjectIDFromHex(uId)
	if err != nil {
		return result, err
	}

	customerId, err := primitive.ObjectIDFromHex(cId)
	if err != nil {
		return result, err
	}

	// remove these keys if present, otherwise they will object will be gone
	delete(customer, "_id")
	delete(customer, "userId")

	filter := bson.M{"userId": userId, "_id": customerId}
	customerCollection := d.Collection("customers")
	_, err = customerCollection.UpdateOne(context.Background(), filter, bson.M{"$set": customer})

	if err != nil {
		return result, err
	}

	err = customerCollection.FindOne(context.Background(), filter).Decode(&result)

	if err != nil {
		return result, err
	}

	return result, nil
}

func CreateCustomer(d *mongo.Database, uId string, customer bson.M) (bson.M, error) {
	var result bson.M

	userId, err := primitive.ObjectIDFromHex(uId)
	if err != nil {
		return result, err
	}

	customer["userId"] = userId

	customerCollection := d.Collection("customers")
	res, err := customerCollection.InsertOne(context.Background(), customer)
	id := res.InsertedID

	filter := bson.M{"_id": id}
	err = customerCollection.FindOne(context.Background(), filter).Decode(&result)

	return result, err
}

func FindOrCreateCustomerByWPcustomerId(database *mongo.Database, uId string, whCustomerId string, billing models.WebhookLocationData, shipping models.WebhookLocationData) (bson.M, error) {
	customersCollection := database.Collection("customers")

	customer := bson.M{}

	userId, err := primitive.ObjectIDFromHex(uId)
	if err != nil {
		return customer, err
	}

	err = customersCollection.FindOne(context.Background(), bson.M{"userId": userId, "whCustomerId": whCustomerId}).Decode(&customer)
	if err == nil {
		return customer, err
	}

	customer["discount"] = 0
	billingName := billing.First_name + " " + billing.Last_name + "<br/>"
	billingCompany := billing.Company + "<br/>"
	billingAddress1 := billing.Address_1 + "<br/>"
	billingAddress2 := billing.Address_2 + "<br/>"
	billingPostCodeCity := billing.Postcode + " " + billing.City
	customer["invoiceAddress"] = fmt.Sprintf("<p>%s  %s  %s %s %s</p>", billingName, billingCompany, billingAddress1, billingAddress2, billingPostCodeCity)

	customer["name"] = billing.First_name + " " + billing.Last_name

	shippingName := shipping.First_name + " " + shipping.Last_name
	shippingCompany := shipping.Company + "<br/>"
	shippingAddress1 := shipping.Address_1
	shippingAddress2 := shipping.Address_2
	shippingPostCodeCity := shipping.Postcode + " " + shipping.City

	customer["shippingAddress"] = fmt.Sprintf("<p>%s • %s  %s %s • %s</p>", shippingName, shippingCompany, shippingAddress1, shippingAddress2, shippingPostCodeCity)
	customer["ust"] = 5
	customer["whCustomerId"] = whCustomerId

	customer, err = CreateCustomer(database, uId, customer)

	return customer, err
}
