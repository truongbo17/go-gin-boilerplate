package mail

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/smtp"
	"strings"
	"time"

	"github.com/truongbo17/go-gin-boilerplate/config"
)

type Sender struct {
	Config config.Mail
}

func (sender Sender) Send(ctx context.Context, to, subject, body string) error {
	address := net.JoinHostPort(sender.Config.Host, sender.Config.Port)
	connection, err := (&net.Dialer{Timeout: 5 * time.Second}).DialContext(ctx, "tcp", address)
	if err != nil {
		return fmt.Errorf("connect SMTP: %w", err)
	}
	defer connection.Close()
	deadline := time.Now().Add(10 * time.Second)
	if ctxDeadline, ok := ctx.Deadline(); ok && ctxDeadline.Before(deadline) {
		deadline = ctxDeadline
	}
	if err := connection.SetDeadline(deadline); err != nil {
		return err
	}
	if sender.Config.TLSMode == "implicit" {
		secure := tls.Client(connection, &tls.Config{ServerName: sender.Config.Host, MinVersion: tls.VersionTLS12})
		if err := secure.HandshakeContext(ctx); err != nil {
			return fmt.Errorf("SMTP TLS handshake: %w", err)
		}
		connection = secure
	}
	client, err := smtp.NewClient(connection, sender.Config.Host)
	if err != nil {
		return fmt.Errorf("create SMTP client: %w", err)
	}
	defer client.Close()
	if sender.Config.TLSMode == "starttls" {
		if err := client.StartTLS(&tls.Config{ServerName: sender.Config.Host, MinVersion: tls.VersionTLS12}); err != nil {
			return fmt.Errorf("SMTP STARTTLS: %w", err)
		}
	}
	if sender.Config.Username != "" {
		if err := client.Auth(smtp.PlainAuth("", sender.Config.Username, sender.Config.Password, sender.Config.Host)); err != nil {
			return fmt.Errorf("authenticate SMTP: %w", err)
		}
	}
	if err := client.Mail(sender.Config.From); err != nil {
		return err
	}
	if err := client.Rcpt(to); err != nil {
		return err
	}
	writer, err := client.Data()
	if err != nil {
		return err
	}
	message := "From: " + sender.Config.From + "\r\nTo: " + to + "\r\nSubject: " + subject + "\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n" + strings.ReplaceAll(body, "\n", "\r\n")
	if _, err := io.WriteString(writer, message); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}
	// DATA has been accepted. A failed QUIT cannot unsend the message, so
	// retrying here could deliver it twice.
	_ = client.Quit()
	return nil
}
