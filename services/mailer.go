package services

import (
	"demerzel-events/dependencies/mailersend"
	"demerzel-events/internal/models"
	"demerzel-events/pkg/helpers"
	"demerzel-events/pkg/types"
	"fmt"
	"os"
)

func SendEventSubscriptionEmail(user *models.User, event *models.EventResponse) error {
	body := `
	<h3 style="text-align: center;">Hello %s.</h3>
	<p style="text-align: center;">You have successfully subscribed to the event "%s".</p>
	<p style="text-align: center;">You can add the event to your calendar by clicking <a href="%s">here</a></p>`

	startDate, _ := helpers.FormatDateTimeStr(event.StartDate, event.StartTime)
	endDate, _ := helpers.FormatDateTimeStr(event.EndDate, event.EndTime)

	calendarUrl := "https://calendar.google.com/calendar/render?action=TEMPLATE&text=%s&details=%s&dates=%s/%s&location=%s"
	calendarUrl = fmt.Sprintf(calendarUrl, event.Title, event.Description, startDate, endDate, event.Location)

	params := types.MailSendParam{
		Recipient: []string{user.Email},
		Sender:    os.Getenv("MAIL_FROM"),
		Subject:   "Event subscription successful",
		Body:      fmt.Sprintf(body, user.Name, event.Title, calendarUrl),
	}

	return sendEmail(mailersend.Ms, &params)
}

func SendEventUnsubscriptionEmail(user *models.User, event *models.EventResponse) error {
	body := `
	<h3 style="text-align: center;">Hello %s.</h3>
    <p style="text-align: center;">You have successfully unsubscribed from the event "%s"</p>`

	params := types.MailSendParam{
		Recipient: []string{user.Email},
		Sender:    os.Getenv("MAIL_FROM"),
		Subject:   "Event unsubscription successful",
		Body:      fmt.Sprintf(body, user.Name, event.Title),
	}

	return sendEmail(mailersend.Ms, &params)
}

func sendEmail(mailer types.Mailer, params *types.MailSendParam) error {
	return mailer.Send(params)
}
