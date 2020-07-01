package main

import (
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)

func middleware(app *gin.Engine) {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	app.Use(cors.New(config))
	app.Use(
		getUserFromToken,
		parseJSONfromBody,
	)

}

func getUserFromToken(ctx *gin.Context) {
	reqToken := ctx.Request.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) != 2 {
		ctx.Next()
	}

	reqToken = splitToken[1]

	token, err := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return getEnvVariable("JWT_SECRET"), nil
	})

	var requestUser user

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if id, _ok := claims["_id"].(string); _ok {
			requestUser.id = id
		}
		if email, _ok := claims["email"].(string); _ok {
			requestUser.email = email
		}
	} else {
		fmt.Println(err)
		ctx.Next()
		return
	}

	ctx.Set("user", requestUser)
	ctx.Next()
}

func parseJSONfromBody(ctx *gin.Context) {
	data := &bson.M{}
	ctx.Bind(data)
	fmt.Println("Data", data)
	ctx.Set("data", data)
	ctx.Next()
}
