package routes

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"billeroo.de/data-backend/src/db"
	"billeroo.de/data-backend/src/mail"
	"billeroo.de/data-backend/src/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func PostWebhook(database *mongo.Database) func(ctx *gin.Context) {
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

		webhookId, ok := (*bsonData)["_id"].(string)

		if ok {
			var result bson.M
			result, err := db.UpdateWebhookById(database, requestUser.Id, webhookId, *bsonData)
			if err != nil {
				fmt.Println(err)
				ctx.Status(500)
				return
			}
			ctx.JSON(200, gin.H{"body": result})
		} else {
			var result bson.M

			const charset = "abcdefghijklmnopqrstuvwxyz" +
				"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

			webhook := bson.M{
				"secret": stringWithCharset(20, charset),
			}

			result, err := db.CreateWebhook(database, requestUser.Id, webhook)
			if err != nil {
				fmt.Println(err)
				ctx.Status(500)
				return
			}
			ctx.JSON(200, gin.H{"body": result})
		}

	}
}

func GetWebhooks(database *mongo.Database) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		requestUser, ok := getRequestUserFromContext(ctx)
		if !ok {
			ctx.Status(401)
			return
		}

		webhooks, err := db.FindWebhooksByUserId(database, requestUser.Id)
		fmt.Println(requestUser)
		if err != nil {
			log.Fatal(err)
			ctx.Status(500)
			return
		}

		ctx.JSON(200, gin.H{"body": webhooks})
		return
	}
}

func ReceiveWebhook(database *mongo.Database) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		webhookId := ctx.Param("webhookId")

		webhook, err := db.FindWebhookById(database, webhookId)
		if err != nil {
			fmt.Println("Unable to find webhook id %s", webhookId)
			fmt.Println(err.Error())
			ctx.Status(http.StatusNotFound)
			return
		}

		userId := webhook["userId"].(primitive.ObjectID).Hex()

		user, err := db.FindUserById(database, userId)

		if err != nil {
			fmt.Println(err.Error())
			ctx.Status(http.StatusInternalServerError)
			return
		}

		company, err := db.FindCompanyByUserId(database, userId)

		var wd models.WebhookData
		err = ctx.ShouldBindJSON(&wd)

		if err != nil {
			fmt.Println(err.Error())
			ctx.Status(http.StatusOK)
			return
		}

		dateTime := time.Now()
		invoice := bson.M{}

		invoice["userId"] = userId
		invoice["company"] = company
		invoice["orderDate"] = dateTime
		invoice["paymentDate"] = nil
		invoice["invoiceDate"] = dateTime
		invoice["shippingDate"] = dateTime
		invoiceNumber := "" + strconv.Itoa(dateTime.Year())
		if dateTime.Month() < 10 {
			invoiceNumber = invoiceNumber + "0"
		}
		invoiceNumber = invoiceNumber + strconv.Itoa((int(dateTime.Month())))

		count, err := strconv.Atoi(wd.Number)
		if err != nil {
			invoice["invoiceNumber"] = wd.Number
		} else {
			if count < 10 {
				invoiceNumber = invoiceNumber + "0"
			}
			if count < 100 {
				invoiceNumber = invoiceNumber + "0"
			}
			invoice["invoiceNumber"] = invoiceNumber + wd.Number
		}
		invoice["payed"] = false
		portoGross, err := strconv.ParseFloat(wd.Total_shipping, 64)
		if err != nil {
			fmt.Println(err.Error())
			ctx.Status(http.StatusBadRequest)
			return
		}

		invoice["porto"] = fmt.Sprintf("%.2f", (portoGross/107) * 100)
		invoice["totalPrice"] = wd.Total
		fmt.Println("Porto: $s", invoice["porto"])
		
		customer, err := db.FindOrCreateCustomerByWPcustomerId(database, userId, strconv.Itoa(wd.Customer_id), wd.Billing, wd.Shipping)
		if err != nil {
			fmt.Println(err.Error())
			ctx.Status(http.StatusBadRequest)
			return
		}

		articles, err := db.FindOrCreateArticlesByWPlineItems(database, userId, wd.Line_items)

		if err != nil {
			fmt.Println(err.Error())
			ctx.Status(http.StatusBadRequest)
			return
		}

		invoice["customer"] = customer
		invoice["articles"] = articles
		invoice["services"] = bson.A{}

		invoice, err = db.CreateInvoice(database, userId, invoice)

		if err != nil {
			fmt.Println(err.Error())
			ctx.Status(http.StatusInternalServerError)
			return
		}

		invoiceId := invoice["_id"].(primitive.ObjectID).Hex()

		mail.NewInvoiceEmail(user.Email, invoiceId)

		ctx.JSON(http.StatusCreated, gin.H{"body": invoice})
		return
	}
}
