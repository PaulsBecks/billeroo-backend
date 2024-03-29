package routes

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"billeroo.de/data-backend/src/db"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func PostInvoice(database *mongo.Database) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		requestUser, ok := getRequestUserFromContext(ctx)
		fmt.Println(requestUser)
		if !ok {
			ctx.Status(401)
			return
		}

		bsonData := &bson.M{}
		err := ctx.Bind(bsonData)
		if err != nil {
			ctx.Status(400)
			return
		}

		invoiceId, ok := (*bsonData)["_id"].(string)

		if ok {
			var result bson.M
			result, err := db.UpdateInvoiceById(database, requestUser.Id, invoiceId, *bsonData)
			if err != nil {
				fmt.Println(err)
				ctx.Status(500)
				return
			}
			ctx.JSON(200, gin.H{"body": result})
		} else {
			var result bson.M

			result, err := db.CreateInvoice(database, requestUser.Id, *bsonData)
			if err != nil {
				fmt.Println(err)
				ctx.Status(500)
				return
			}
			ctx.JSON(200, gin.H{"body": result})
		}

	}
}

func GetInvoices(database *mongo.Database) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		requestUser, ok := getRequestUserFromContext(ctx)
		if !ok {
			ctx.Status(401)
			return
		}

		limitStr := ctx.DefaultQuery("limit", "-1")
		limit, err := strconv.Atoi(limitStr)

		if err != nil {
			log.Println("Limit set but not an integer")
			ctx.Status(http.StatusBadRequest)
			return
		}

		invoices, err := db.FindInvoicesByUserId(database, requestUser.Id, int64(limit))
		fmt.Println(requestUser)
		if err != nil {
			log.Fatal(err)
			ctx.Status(500)
			return
		}

		ctx.JSON(200, gin.H{"body": invoices})
		return
	}
}
