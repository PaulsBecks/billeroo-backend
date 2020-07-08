package mail

import (
	"fmt"
)

func NewInvoiceEmail(to string, newInvoiceId string) {
	title := "Eine neue Rechnung wurde angelegt"
	body := "Hallo,\n eine neue Rechnung wurde automatisch angelegt. \n\n Du findest die Rechnung unter https://billeroo.de/invoices/" + newInvoiceId + "\n\n" + "Viele Grüße,\n Billeroo Service"

	err := sendMail(to, title, body)
	if err != nil {
		fmt.Println(err.Error())
	}
}
