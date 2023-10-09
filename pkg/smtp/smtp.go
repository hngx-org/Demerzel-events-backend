package smtp

import (
	"demerzel-events/pkg/types"
	"fmt"
	"net/smtp"
	"os"
)

var Ml *Mailer

type Mailer struct {
	from     string
	host     string
	port     string
	password string
}

func Initialize() {
	from := os.Getenv("MAIL_FROM")
	host := os.Getenv("MAIL_HOST")
	port := os.Getenv("MAIL_PORT")
	password := os.Getenv("MAIL_PASSWORD")

	Ml = &Mailer{
		from:     from,
		host:     host,
		port:     port,
		password: password,
	}

	fmt.Println(Ml)
}

func (ml *Mailer) Send(param *types.MailSendParam) error {
	auth := smtp.PlainAuth("", ml.from, ml.password, ml.host)

	body := "Content-Type: text/html; charset=\"UTF-8\"\n\n" + param.Body
	msg := []byte("Subject: " + param.Subject + "\n" + body)

	err := smtp.SendMail(ml.host+":"+ml.port, auth, ml.from, param.Recipient, msg)
	return err
}
