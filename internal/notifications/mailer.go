package notifications

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"

	"gopkg.in/gomail.v2"
)

type mailtrapClient struct {
	fromEmail string
	apiKey    string
}

var FS embed.FS

func NewMailTrapClient(apiKey, fromEmail string) (mailtrapClient, error) {
	if apiKey == "" {
		return mailtrapClient{}, errors.New("api key is required")
	}

	return mailtrapClient{
		fromEmail: fromEmail,
		apiKey:    apiKey,
	}, nil
}

func (m mailtrapClient) Send(templateFile NotificationMessageTemplate, recipients []string, subject string, data any) (int, error) {
	// Template parsing and building
	templatePath := filepath.Join("templates/", templateFile.String())
	fmt.Println(templatePath)

	tmpl, err := template.ParseFS(FS, "templates/"+templateFile.String())
	if err != nil {
		return -1, err
	}

	var body bytes.Buffer

	if err := tmpl.Execute(&body, data); err != nil {
		return -1, err
	}

	message := gomail.NewMessage()
	message.SetHeader("From", m.fromEmail)
	message.SetHeader("To", strings.Join(recipients, ","))
	message.SetHeader("Subject", subject)

	message.AddAlternative("text/html", body.String())

	dialer := gomail.NewDialer("live.smtp.mailtrap.io", 587, "api", m.apiKey)

	if err := dialer.DialAndSend(message); err != nil {
		return -1, err
	}

	return 200, nil
}
