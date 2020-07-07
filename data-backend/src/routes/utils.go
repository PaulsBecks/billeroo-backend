package routes

import (
	"log"
	"math/rand"
	"os"
	"time"

	"billeroo.de/data-backend/src/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Load env variables from .env file
func getEnvVariable(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error when loading env")
	}

	return os.Getenv(key)
}

func getRequestUserFromContext(ctx *gin.Context) (models.User, bool) {
	var requestUser models.User
	if u, ok := ctx.Get("user"); ok {
		if _u, _ok := u.(models.User); _ok {
			return _u, true
		}
		return requestUser, false
	}
	return requestUser, false
}

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
