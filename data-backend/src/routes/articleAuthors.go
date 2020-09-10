package routes

import (
	"fmt"
	"log"

	"billeroo.de/data-backend/src/db"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func PostArticleAuthor(database *mongo.Database) func(ctx *gin.Context) {
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

		articleAuthorId, okArticleAuthor := (*bsonData)["_id"].(string)

		if okArticleAuthor {
			var result bson.M
			result, err := db.UpdateArticleAuthorById(database, requestUser.Id, articleAuthorId, *bsonData)
			if err != nil {
				fmt.Println(err)
				ctx.Status(500)
				return
			}
			ctx.JSON(200, gin.H{"body": result})
		} else {
			articleId, okArticle := (*bsonData)["articleId"].(string)
			authorId, okAuthor := (*bsonData)["authorId"].(string)

			if !okAuthor || !okArticle {
				fmt.Println("Error 404: No valid article and author provided.")
				ctx.Status(404)
				return
			}
			var result bson.M

			result, err := db.CreateArticleAuthor(database, requestUser.Id, authorId, articleId, *bsonData)
			if err != nil {
				fmt.Println(err)
				ctx.Status(500)
				return
			}
			ctx.JSON(200, gin.H{"body": result})
		}

	}
}

func GetArticleAuthors(database *mongo.Database) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		requestUser, ok := getRequestUserFromContext(ctx)
		if !ok {
			ctx.Status(401)
			return
		}

		articleAuthors, err := db.FindArticleAuthorsByUserId(database, requestUser.Id)

		if err != nil {
			log.Fatal(err)
			ctx.Status(500)
			return
		}

		ctx.JSON(200, gin.H{"body": articleAuthors})
		return
	}
}
