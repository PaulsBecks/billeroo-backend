package mail

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"os"

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

func sendMail(_to string, title string, body string) error {

	smtpServer := getEnvVariable("EMAIL_PROVIDER")
	user := getEnvVariable("EMAIL_USER")
	password := getEnvVariable("EMAIL_PASSWORD")

	auth := smtp.PlainAuth(
		"",
		user,
		password,
		smtpServer,
	)

	from := mail.Address{Name: "Billeroo Service", Address: getEnvVariable("SERVICE_EMAIL")}
	to := mail.Address{Name: "", Address: _to}

	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = to.String()
	header["Subject"] = title
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(
		smtpServer+":587",
		auth,
		from.Address,
		[]string{to.Address},
		[]byte(message),
		//[]byte("This is the email body."),
	)

	return err
}
