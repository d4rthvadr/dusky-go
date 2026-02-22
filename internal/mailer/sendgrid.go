package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

const (
	fromName = "Dusky Team"
)

type SendGridMailer struct {
	fromEmail  string
	client     *sendgrid.Client
	maxRetries int
}

func NewSendGridMailer(apiKey string, fromEmail string, maxRetries int) (*SendGridMailer, error) {

	if apiKey == "" {
		return nil, fmt.Errorf("sendgrid api key should not be null")
	}
	client := sendgrid.NewSendClient(apiKey)
	return &SendGridMailer{
		fromEmail:  fromEmail,
		client:     client,
		maxRetries: maxRetries,
	}, nil
}

// parseTemplate parses the email template and returns the body and subject as strings.
// It uses the Go html/template package to execute the template with the provided data.
func (m *SendGridMailer) parseTemplate(templateFile string, data any) (bodyStr string, subjectStr string, errStr error) {

	tmpl, err := template.ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse template: %w", err)
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return "", "", fmt.Errorf("failed to execute subject template: %w", err)
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return "", "", fmt.Errorf("failed to execute body template: %w", err)
	}

	return body.String(), subject.String(), nil
}

// Send constructs and sends an email using the SendGrid API. It takes the template file, recipient's username and email, data for template execution, and a flag for sandbox mode.
func (m *SendGridMailer) Send(templateFile, username, email string, data any, isSandbox bool) error {

	from := mail.NewEmail(fromName, m.fromEmail)
	to := mail.NewEmail(username, email)

	// template parsing and building
	body, subject, err := m.parseTemplate(templateFile, data)
	if err != nil {
		return err
	}

	message := m.buildMessage(from, to, subject, body, isSandbox)

	return m.sendEmailWithRetry(message, isSandbox)

}

// buildMessage constructs the email message using the SendGrid mail helper.
func (m *SendGridMailer) buildMessage(from, to *mail.Email, subject, body string, isSandbox bool) *mail.SGMailV3 {
	message := mail.NewSingleEmail(from, subject, to, "", body)

	message.SetMailSettings(
		&mail.MailSettings{
			SandboxMode: &mail.Setting{
				Enable: &isSandbox,
			},
		},
	)

	return message
}

// sendEmailWithRetry attempts to send the email and retries if it fails, with a backoff strategy.
func (m *SendGridMailer) sendEmailWithRetry(message *mail.SGMailV3, isSandbox bool) error {

	// If we're in sandbox mode, we won't actually send the email, so we can skip the retry logic.
	if isSandbox {
		return nil
	}

	isSent := false
	var errSend error

	for i := 0; i < m.maxRetries; i++ {
		response, err := m.client.Send(message)
		if err != nil {
			errSend = err

			// Implement a backoff strategy before retrying
			backoffDuration := time.Duration(i+1) * time.Second
			time.Sleep(backoffDuration)
			continue
		}

		if response.StatusCode >= 500 {
			errSend = fmt.Errorf("server error: %s", response.Body)
			continue
		}

		if response.StatusCode >= 200 && response.StatusCode < 300 {
			isSent = true
			break
		}

	}

	if isSent {
		return nil
	}
	return fmt.Errorf("failed to send email: %w", errSend)
}
