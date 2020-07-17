package main

import (
	"billeroo.de/data-backend/src/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()
	middleware(app)

	client := connectToClient()
	database := client.Database("billeroo")

	app.GET("/data", routes.GetData(database))
	app.POST("/data", routes.PostData(database))

	app.GET("/data/customers", routes.GetCustomers(database))
	app.POST("/data/customers", routes.PostCustomer(database))

	app.GET("/data/articles", routes.GetArticles(database))
	app.POST("/data/articles", routes.PostArticle(database))

	app.GET("/data/authors", routes.GetAuthors(database))
	app.POST("/data/authors", routes.PostAuthor(database))

	app.GET("/data/invoices", routes.GetInvoices(database))
	app.POST("/data/invoices", routes.PostInvoice(database))

	app.GET("/data/services", routes.GetServices(database))
	app.POST("/data/services", routes.PostService(database))

	app.GET("/data/subscriptions", routes.GetSubscriptions(database))
	app.GET("/data/subscriptions/last", routes.GetRecentSubscription(database))
	app.POST("/data/subscriptions", routes.PostSubscription(database))

	app.GET("/data/companies", routes.GetCompany(database))
	app.POST("/data/companies", routes.PostCompany(database))

	app.GET("/data/articleAuthors", routes.GetArticleAuthors(database))
	app.POST("/data/articleAuthor", routes.PostArticleAuthor(database))

	app.GET("/data/webhooks", routes.GetWebhooks(database))
	app.POST("/data/webhooks", routes.PostWebhook(database))
	app.POST("/data/webhooks/:webhookId", routes.ReceiveWebhook(database))

	app.Run()
}
