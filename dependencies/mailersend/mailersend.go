package mailersend

import (
	"context"
	"demerzel-events/pkg/types"
	"github.com/mailersend/mailersend-go"
	"os"
	"time"
)

var Ms *Mailersend

type Mailersend struct {
	ms *mailersend.Mailersend
}

func Initialize() {
	apiKey := os.Getenv("MAIL_API_KEY")

	Ms = &Mailersend{
		ms: mailersend.NewMailersend(apiKey),
	}
}

func (ms *Mailersend) Send(param *types.MailSendParam) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	from := mailersend.From{Email: param.Sender}
	var recipients []mailersend.Recipient
	for _, recipient := range param.Recipient {
		recipients = append(recipients, mailersend.Recipient{Email: recipient})
	}

	message := ms.ms.Email.NewMessage()

	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetSubject(param.Subject)
	message.SetHTML(param.Body)
	message.SetText(param.Body)

	_, err := ms.ms.Email.Send(ctx, message)
	return err
}
