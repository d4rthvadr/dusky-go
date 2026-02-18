package mailer

import "embed"

const (
	TemplateUserInvitation = "user_invitation.tmpl"
)

//go:embed templates/*
var templateFS embed.FS

type Client interface {
	Send(templateFile, username, email string, data any, isSandbox bool) error
}
