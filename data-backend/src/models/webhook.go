package models

type WebhookData struct {
	Number         string              `json:"number"`
	Total          string              `json:"total"`
	Total_shipping string              `json:"total_shipping"`
	Billing        WebhookLocationData `json:"billing"`
	Line_items     []WebhookLineItem   `json:"line_items"`
	Shipping       WebhookLocationData `json:"shipping"`
	Customer_id    int                 `json:"customer_id"`
}

type WebhookLocationData struct {
	Address_1  string `json:"address_1"`
	Address_2  string `json:"address_2"`
	City       string `json:"city"`
	Company    string `json:"company"`
	Country    string `json:"country"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Postcode   string `json:"postcode"`
	State      string `json:"state"`
}

type WebhookLineItem struct {
	Id         int     `json:"id"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	Product_id int     `json:"product_id"`
	Quantity   int     `json:"quantity"`
}
