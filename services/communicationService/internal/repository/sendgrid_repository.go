package repository

import (
	"fmt"
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendgridRepository struct {
	apiKey string
	client *sendgrid.Client
}

func NewSendgridRepository() *SendgridRepository {
	apiKey := os.Getenv("SENDGRID_ADMIN_KEY")
	if apiKey == "" {
		log.Fatal("SENDGRID_ADMIN_KEY environment variable is not set")
	}

	return &SendgridRepository{
		apiKey: apiKey,
		client: sendgrid.NewSendClient(apiKey),
	}
}

type EmailAttachment struct {
	Content     string
	Filename    string
	Type        string
	Disposition string
}

type EmailOptions struct {
	To           []string
	From         string
	Subject      string
	Text         string
	HTML         string
	Attachments  []EmailAttachment
	CC           []string
	BCC          []string
	TemplateID   string
	TemplateData map[string]interface{}
}

func (repo *SendgridRepository) SendEmail(opts EmailOptions) error {
	// Create a new SendGrid message
	message := mail.NewV3Mail()

	// Set From
	from := mail.NewEmail("Younified", os.Getenv("email_account"))
	message.SetFrom(from)

	// Set Subject
	message.Subject = opts.Subject

	// Add recipients
	personalization := mail.NewPersonalization()
	for _, recipient := range opts.To {
		to := mail.NewEmail("", recipient)
		personalization.AddTos(to)
	}

	// Add CC recipients if any
	for _, ccRecipient := range opts.CC {
		cc := mail.NewEmail("", ccRecipient)
		personalization.AddCCs(cc)
	}

	// Add BCC recipients if any
	for _, bccRecipient := range opts.BCC {
		bcc := mail.NewEmail("", bccRecipient)
		personalization.AddBCCs(bcc)
	}

	// Add content
	if opts.Text != "" {
		content := mail.NewContent("text/plain", opts.Text)
		message.AddContent(content)
	}
	if opts.HTML != "" {
		content := mail.NewContent("text/html", opts.HTML)
		message.AddContent(content)
	}

	// Add attachments
	for _, attachment := range opts.Attachments {
		a := mail.NewAttachment()
		a.SetContent(attachment.Content)
		a.SetType(attachment.Type)
		a.SetFilename(attachment.Filename)
		a.SetDisposition(attachment.Disposition)
		message.AddAttachment(a)
	}

	// Add template data if template ID is provided
	if opts.TemplateID != "" {
		message.SetTemplateID(opts.TemplateID)
		personalization.DynamicTemplateData = opts.TemplateData
	}

	// Add personalization to message
	message.AddPersonalizations(personalization)

	// Send the email
	response, err := repo.client.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	// Log response for debugging (optional)
	log.Printf("Email sent. Status Code: %d", response.StatusCode)

	return nil
}
