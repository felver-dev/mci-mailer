package smtp

import (
	"encoding/base64"
	"fmt"
	"io"

	"github.com/mcicare/mci-mailer/internal/config"
	"gopkg.in/gomail.v2"
)

type Attachment struct {
	Filename string
	Content  []byte
	MimeType string
}

type Message struct {
	FromAddress string
	FromName    string
	To          []string
	CC          []string
	BCC         []string
	ReplyTo     string
	Subject     string
	Html        string
	Text        string
	Attachments []Attachment
}

type Client struct {
	cfg *config.SMTPConfig
}

func NewClient(cfg *config.SMTPConfig) *Client {
	return &Client{cfg: cfg}
}

func (c *Client) Send(msg Message) error {
	m := gomail.NewMessage()

	from := msg.FromAddress
	if from == "" {
		from = c.cfg.From
	}
	fromName := msg.FromName
	if fromName == "" {
		fromName = c.cfg.FromName
	}
	m.SetAddressHeader("From", from, fromName)

	toAddresses := make([]string, len(msg.To))
	for i, addr := range msg.To {
		toAddresses[i] = addr
	}
	m.SetHeader("To", toAddresses...)

	if len(msg.CC) > 0 {
		m.SetHeader("Cc", msg.CC...)
	}
	if len(msg.BCC) > 0 {
		m.SetHeader("Bcc", msg.BCC...)
	}
	if msg.ReplyTo != "" {
		m.SetHeader("Reply-To", msg.ReplyTo)
	}

	m.SetHeader("Subject", msg.Subject)

	if msg.Html != "" && msg.Text != "" {
		m.SetBody("text/plain", msg.Text)
		m.AddAlternative("text/html", msg.Html)
	} else if msg.Html != "" {
		m.SetBody("text/html", msg.Html)
	} else {
		m.SetBody("text/plain", msg.Text)
	}

	for _, att := range msg.Attachments {
		att := att
		m.Attach(att.Filename, gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := w.Write(att.Content)
			return err
		}), gomail.SetHeader(map[string][]string{
			"Content-Type": {att.MimeType},
		}))
	}

	d := gomail.NewDialer(c.cfg.Host, c.cfg.Port, c.cfg.User, c.cfg.Password)
	d.SSL = false // TLS via STARTTLS on port 587

	return d.DialAndSend(m)
}

func (c *Client) Ping() error {
	d := gomail.NewDialer(c.cfg.Host, c.cfg.Port, c.cfg.User, c.cfg.Password)
	d.SSL = false
	conn, err := d.Dial()
	if err != nil {
		return fmt.Errorf("smtp ping failed: %w", err)
	}
	return conn.Close()
}

func DecodeBase64Attachment(b64 string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(b64)
}
