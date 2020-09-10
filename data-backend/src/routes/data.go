package routes

import (
	"fmt"
	"log"

	"billeroo.de/data-backend/src/db"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetData(database *mongo.Database) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		requestUser, ok := getRequestUserFromContext(ctx)
		fmt.Println(requestUser)
		if !ok {
			ctx.Status(401)
			return
		}
		result, err := db.FindDataByUserId(database, requestUser.Id)
		if err != nil {
			log.Fatal(err)
			ctx.Status(500)
			return
		}
		ctx.JSON(200, gin.H{"body": result})
		return
	}
}

func PostData(database *mongo.Database) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		requestUser, ok := getRequestUserFromContext(ctx)
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

		result, err := db.UpdateDataByUserId(database, requestUser.Id, *bsonData)
		if err != nil {
			fmt.Println(err)
			ctx.Status(500)
			return
		}
		ctx.JSON(200, gin.H{"body": result})
		return
	}
}
