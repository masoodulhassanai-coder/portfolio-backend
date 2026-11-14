package service

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"strings"

	"portfolio/internal/config"
)

type EmailService struct {
	cfg *config.Config
}

func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{cfg: cfg}
}

// sanitizeHeader removes newlines to prevent header injection
func sanitizeHeader(input string) string {
	return strings.NewReplacer("\r", "", "\n", "").Replace(input)
}

const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { background-color: #050505; color: #ededed; font-family: -apple-system, sans-serif; margin: 0; padding: 40px 20px; }
        .container { max-width: 600px; margin: 0 auto; background-color: #0a0a0a; border: 1px solid #222; border-radius: 16px; overflow: hidden; box-shadow: 0 20px 40px rgba(0,0,0,0.6); }
        .header { background: linear-gradient(180deg, #111111 0%, #0a0a0a 100%); padding: 32px 40px; border-bottom: 1px solid #222; text-align: center; }
        .header h1 { margin: 0; font-size: 24px; font-weight: 800; color: #ffffff; letter-spacing: -0.5px; }
        .header span { color: #3b82f6; }
        .content { padding: 40px; }
        .data-table { width: 100%; border-collapse: collapse; margin-bottom: 32px; }
        .data-table td { padding: 12px 0; vertical-align: top; border-bottom: 1px solid #1a1a1a; }
        .data-table td:first-child { width: 100px; font-size: 12px; font-weight: 700; color: #64748b; text-transform: uppercase; letter-spacing: 1px; padding-top: 14px; }
        .data-table td:last-child { font-size: 16px; color: #f8fafc; font-weight: 500; }
        .email-link { color: #60a5fa; text-decoration: none; font-weight: 600; }
        .badge { display: inline-block; background-color: rgba(59, 130, 246, 0.1); color: #60a5fa; border: 1px solid rgba(59, 130, 246, 0.2); padding: 6px 12px; border-radius: 6px; font-size: 12px; font-weight: 700; letter-spacing: 0.5px; }
        .message-label { font-size: 12px; font-weight: 700; color: #64748b; text-transform: uppercase; letter-spacing: 1px; margin-bottom: 12px; }
        .message-box { background-color: #111; border: 1px solid #222; border-radius: 12px; padding: 24px; font-size: 16px; color: #cbd5e1; line-height: 1.6; white-space: pre-wrap; box-shadow: inset 0 2px 4px rgba(0,0,0,0.2); }
        .footer { background-color: #050505; padding: 24px 40px; text-align: center; font-size: 12px; color: #475569; border-top: 1px solid #222; }
    </style>
</head>
<body>
    <div class="email-wrapper">
        <div class="header">
            <h1>CodeWorks<span>.</span></h1>
        </div>
        <div class="content">
            <table class="data-table">
                <tr>
                    <td>Inquiry Type</td>
                    <td><span class="badge">{{.QueryType}}</span></td>
                </tr>
                <tr>
                    <td>From</td>
                    <td>{{.Name}}</td>
                </tr>
                <tr>
                    <td>Email</td>
                    <td><a href="mailto:{{.Email}}" class="email-link">{{.Email}}</a></td>
                </tr>
            </table>
            <div>
                <div class="message-label">Message</div>
                <div class="message-box">{{.Message}}</div>
            </div>
        </div>
        <div class="footer">
            This email was sent via the CodeWorks Contact Form.
        </div>
    </div>
</body>
</html>
`

func (s *EmailService) SendContactEmail(name, email, queryType, message string) error {
	auth := smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPass, s.cfg.SMTPHost)

	tmpl, err := template.New("email").Parse(htmlTemplate)
	if err != nil {
		return err
	}

	data := struct {
		Name      string
		Email     string
		QueryType string
		Message   string
	}{
		Name:      sanitizeHeader(name),
		Email:     sanitizeHeader(email),
		QueryType: sanitizeHeader(queryType),
		Message:   message, // Body can keep newlines
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	to := []string{s.cfg.SMTPUser}
	// Sanitize subject to prevent injection
	subject := sanitizeHeader(fmt.Sprintf("New Inquiry: %s", queryType))
	
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
		"\r\n"+
		"%s\r\n", 
		to[0], subject, body.String()))

	addr := s.cfg.SMTPHost + ":" + s.cfg.SMTPPort
	err = smtp.SendMail(addr, auth, s.cfg.SMTPUser, to, msg)
	return err
}