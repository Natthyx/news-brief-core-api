package external_services

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"

	"github.com/RealEskalate/G6-NewsBrief/internal/domain/contract"
)

// smtp attribute
type EmailService struct {
	Host        string
	Port        string
	Username    string
	AppPassword string
	From        string
}

// EmailService factory
func NewEmailService(host, port, username, appPassword, from string) *EmailService {
	return &EmailService{
		Host:        host,
		Port:        port,
		Username:    username,
		AppPassword: appPassword,
		From:        from,
	}
}

// make sure EmailService implements contract.IEmailService.go
var _ contract.IEmailService = (*EmailService)(nil)

// send email method
func (es *EmailService) SendEmail(ctx context.Context, to, subject, body string) error {
	fromAddr := es.From
	if strings.TrimSpace(fromAddr) == "" {
		fromAddr = es.Username
	}
	if strings.TrimSpace(fromAddr) == "" {
		return fmt.Errorf("email sender not configured: From/Username is empty")
	}

	// Build message with minimal MIME headers for inbox placement
	msg := []byte(
		fmt.Sprintf(
			"From: %s\r\n"+
				"To: %s\r\n"+
				"Subject: %s\r\n"+
				"MIME-Version: 1.0\r\n"+
				"Content-Type: text/plain; charset=\"UTF-8\"\r\n"+
				"\r\n"+
				"%s\r\n",
			fromAddr, to, subject, body,
		),
	)

	addr := net.JoinHostPort(es.Host, es.Port)

	// Establish connection
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("smtp dial failed to %s: %w", addr, err)
	}
	defer conn.Close()

	c, err := smtp.NewClient(conn, es.Host)
	if err != nil {
		return fmt.Errorf("smtp new client failed: %w", err)
	}
	defer c.Quit()

	// Upgrade to TLS if supported
	if ok, _ := c.Extension("STARTTLS"); ok {
		config := &tls.Config{ServerName: es.Host}
		if err = c.StartTLS(config); err != nil {
			return fmt.Errorf("smtp STARTTLS failed: %w", err)
		}
	}

	// Authenticate if server supports it and creds provided
	if es.Username != "" && es.AppPassword != "" {
		if ok, _ := c.Extension("AUTH"); ok {
			auth := smtp.PlainAuth("", es.Username, es.AppPassword, es.Host)
			if err = c.Auth(auth); err != nil {
				return fmt.Errorf("smtp auth failed: %w", err)
			}
		}
	}

	// Set the sender and recipient
	if err = c.Mail(fromAddr); err != nil {
		return fmt.Errorf("smtp MAIL FROM failed: %w", err)
	}
	if err = c.Rcpt(to); err != nil {
		return fmt.Errorf("smtp RCPT TO failed for %s: %w", to, err)
	}

	// Send the email body
	wc, err := c.Data()
	if err != nil {
		return fmt.Errorf("smtp DATA start failed: %w", err)
	}
	_, werr := wc.Write(msg)
	cerr := wc.Close()
	if werr != nil {
		return fmt.Errorf("smtp write failed: %w", werr)
	}
	if cerr != nil {
		return fmt.Errorf("smtp close failed: %w", cerr)
	}

	return nil
}
