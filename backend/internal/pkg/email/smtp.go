package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"
)

// Sender sends a plain-text email to one or more recipients.
type Sender interface {
	Send(ctx context.Context, to []string, subject, body string) error
}

// SMTPSender delivers email over SMTP with optional STARTTLS and AUTH.
type SMTPSender struct {
	host     string
	port     string
	user     string
	password string
	from     string
	timeout  time.Duration
}

// NewSMTPSender creates a new SMTPSender.
// port defaults to 587; from defaults to user when empty.
func NewSMTPSender(host, port, user, password, from string) *SMTPSender {
	if port == "" {
		port = "587"
	}
	if from == "" {
		from = user
	}
	return &SMTPSender{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		from:     from,
		timeout:  15 * time.Second,
	}
}

// Send delivers an email. STARTTLS is negotiated when the server supports it.
func (s *SMTPSender) Send(ctx context.Context, to []string, subject, body string) error {
	if len(to) == 0 {
		return nil
	}

	addr := net.JoinHostPort(s.host, s.port)
	dialer := &net.Dialer{Timeout: s.timeout}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("email: dial %s: %w", addr, err)
	}

	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		return fmt.Errorf("email: smtp client: %w", err)
	}
	defer client.Close()

	if ok, _ := client.Extension("STARTTLS"); ok {
		if err := client.StartTLS(&tls.Config{ServerName: s.host}); err != nil {
			return fmt.Errorf("email: starttls: %w", err)
		}
	}

	if s.user != "" {
		if ok, _ := client.Extension("AUTH"); ok {
			if err := client.Auth(smtp.PlainAuth("", s.user, s.password, s.host)); err != nil {
				return fmt.Errorf("email: auth: %w", err)
			}
		}
	}

	if err := client.Mail(s.from); err != nil {
		return fmt.Errorf("email: mail from: %w", err)
	}
	for _, rcpt := range to {
		if err := client.Rcpt(rcpt); err != nil {
			return fmt.Errorf("email: rcpt %s: %w", rcpt, err)
		}
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("email: data: %w", err)
	}
	msg := strings.Join([]string{
		"From: " + s.from,
		"To: " + strings.Join(to, ", "),
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
	}, "\r\n")
	if _, err := w.Write([]byte(msg)); err != nil {
		return fmt.Errorf("email: write: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("email: data close: %w", err)
	}

	return client.Quit()
}
