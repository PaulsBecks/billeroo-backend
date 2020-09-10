package routes

import (
	"fmt"
	"log"

	"billeroo.de/data-backend/src/db"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func PostAuthor(database *mongo.Database) func(ctx *gin.Context) {
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

		authorId, ok := (*bsonData)["_id"].(string)

		if ok {
			var result bson.M
			result, err := db.UpdateAuthorById(database, requestUser.Id, authorId, *bsonData)
			if err != nil {
				fmt.Println(err)
				ctx.Status(500)
				return
			}
			ctx.JSON(200, gin.H{"body": result})
		} else {
			var result bson.M

			result, err := db.CreateAuthor(database, requestUser.Id, *bsonData)
			if err != nil {
				fmt.Println(err)
				ctx.Status(500)
				return
			}
			ctx.JSON(200, gin.H{"body": result})
		}

	}
}

func GetAuthors(database *mongo.Database) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		requestUser, ok := getRequestUserFromContext(ctx)
		if !ok {
			ctx.Status(401)
			return
		}

		authors, err := db.FindAuthorsByUserId(database, requestUser.Id)
		fmt.Println(requestUser)
		if err != nil {
			log.Fatal(err)
			ctx.Status(500)
			return
		}

		ctx.JSON(200, gin.H{"body": authors})
		return
	}
}
