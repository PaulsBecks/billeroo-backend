package main

import (
	"fmt"
	"strings"

	"billeroo.de/data-backend/src/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func middleware(app *gin.Engine) {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	app.Use(cors.New(config))
	app.Use(
		getUserFromToken,
	)

}

func getUserFromToken(ctx *gin.Context) {
	reqToken := ctx.Request.Header.Get("Authorization")
	fmt.Println(reqToken)
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) != 2 {
		ctx.Next()
		return
	}

	reqToken = splitToken[1]

	token, err := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return getEnvVariable("JWT_SECRET"), nil
	})

	var requestUser models.User

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println(claims)
		if id, _ok := claims["_id"].(string); _ok {
			requestUser.Id = id
		}
		if email, _ok := claims["email"].(string); _ok {
			requestUser.Email = email
		}
	} else {
		fmt.Println(err)
		ctx.Next()
		return
	}
	fmt.Println(requestUser)
	ctx.Set("user", requestUser)
	ctx.Next()
}
