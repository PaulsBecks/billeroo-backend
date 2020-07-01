package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)

func getRequestUserFromContext(ctx *gin.Context) (user, bool) {
	var requestUser user
	if u, ok := ctx.Get("user"); ok {
		if _u, _ok := u.(user); _ok {
			return _u, true
		}
		return requestUser, false
	}
	return requestUser, false
}

func routes(app *gin.Engine) *gin.Engine {
	client := connectToClient()
	database := client.Database("billeroo")

	app.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"greetings": "Hi there!"})
	})

	app.GET("/data", func(ctx *gin.Context) {
		requestUser, ok := getRequestUserFromContext(ctx)
		if !ok {
			ctx.Status(401)
			return
		}
		result, err := findDataByUserId(database, requestUser.id)
		if err != nil {
			log.Fatal(err)
			ctx.Status(500)
			return
		}
		ctx.JSON(200, gin.H{"body": result})
		return
	})

	app.POST("/data", func(ctx *gin.Context) {
		requestUser, ok := getRequestUserFromContext(ctx)
		if !ok {
			ctx.Status(401)
			return
		}
		data, ok := ctx.Get("data")
		bsonData, bsonOk := data.(*bson.M)
		if !ok || !bsonOk {
			ctx.Status(404)
			return
		}

		result, err := updateDataByUserId(database, requestUser.id, *bsonData)
		if err != nil {
			fmt.Println(err)
			ctx.Status(500)
			return
		}
		ctx.JSON(200, gin.H{"body": result})
		return
	})
	return app
}
