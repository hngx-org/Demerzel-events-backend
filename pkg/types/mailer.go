package types

type MailSendParam struct {
	Sender    string
	Subject   string
	Body      string
	Recipient []string
}

type Mailer interface {
	Send(params *MailSendParam) error
}
