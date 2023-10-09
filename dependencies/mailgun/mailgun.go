package mailgun

import (
	"context"
	"demerzel-events/pkg/types"
	"github.com/mailgun/mailgun-go/v4"
	"os"
	"time"
)

var Mg *MailGun

type MailGun struct {
	mg mailgun.Mailgun
}

func Initialize() {
	domain := os.Getenv("MAIL_DOMAIN")
	apiKey := os.Getenv("MAIL_API_KEY")

	Mg = &MailGun{
		mg: mailgun.NewMailgun(domain, apiKey),
	}
}

func (mg *MailGun) Send(params *types.MailSendParam) error {
	message := mg.mg.NewMessage(params.Sender, params.Subject, params.Body, params.Recipient[0])
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err := mg.mg.Send(ctx, message)
	return err
}
