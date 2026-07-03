package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mcicare/mci-mailer/internal/config"
	"github.com/mcicare/mci-mailer/internal/domain"
	"github.com/mcicare/mci-mailer/internal/dto"
	"github.com/mcicare/mci-mailer/internal/repository"
	smtpClient "github.com/mcicare/mci-mailer/internal/smtp"
)

const maxRetries = 3

type MailerService struct {
	smtp        *smtpClient.Client
	logRepo     repository.EmailLogRepository
	tmplService *TemplateService
	cfg         *config.SMTPConfig
}

func NewMailerService(
	smtp *smtpClient.Client,
	logRepo repository.EmailLogRepository,
	tmplService *TemplateService,
	cfg *config.SMTPConfig,
) *MailerService {
	return &MailerService{smtp: smtp, logRepo: logRepo, tmplService: tmplService, cfg: cfg}
}

func (s *MailerService) Send(ctx context.Context, req dto.SendRequest, apiKeyID *uuid.UUID) (*dto.SendResponse, error) {
	subject, html, text, err := s.resolveContent(ctx, req)
	if err != nil {
		return nil, err
	}

	if subject == "" || (html == "" && text == "") {
		return nil, errors.New("subject and at least one of html or text body are required")
	}

	toAddresses := extractEmails(req.To)
	logEntry := &domain.EmailLog{
		ID:           uuid.New(),
		ApiKeyID:     apiKeyID,
		FromAddress:  s.resolveFrom(req),
		ToAddresses:  toAddresses,
		CcAddresses:  extractEmails(req.CC),
		BccAddresses: extractEmails(req.BCC),
		Subject:      subject,
		Status:       domain.StatusQueued,
		Attempts:     0,
		CreatedAt:    time.Now(),
	}
	if req.TemplateName != "" {
		logEntry.TemplateName = &req.TemplateName
	}

	if err := s.logRepo.Create(ctx, logEntry); err != nil {
		return nil, fmt.Errorf("failed to create log: %w", err)
	}

	msg := s.buildMessage(req, subject, html, text)

	var sendErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		sendErr = s.smtp.Send(msg)
		if sendErr == nil {
			_ = s.logRepo.MarkSent(ctx, logEntry.ID, attempt)
			return &dto.SendResponse{
				LogID:   logEntry.ID.String(),
				Status:  string(domain.StatusSent),
				Message: "email sent successfully",
			}, nil
		}
		if attempt < maxRetries {
			time.Sleep(time.Duration(attempt*2) * time.Second)
		}
	}

	errMsg := sendErr.Error()
	_ = s.logRepo.UpdateStatus(ctx, logEntry.ID, domain.StatusFailed, &errMsg, maxRetries)

	return &dto.SendResponse{
		LogID:   logEntry.ID.String(),
		Status:  string(domain.StatusFailed),
		Message: "email delivery failed after retries",
	}, sendErr
}

func (s *MailerService) resolveContent(ctx context.Context, req dto.SendRequest) (subject, html, text string, err error) {
	if req.TemplateName != "" {
		return s.tmplService.Render(ctx, req.TemplateName, req.Variables)
	}
	return req.Subject, req.Html, req.Text, nil
}

func (s *MailerService) resolveFrom(req dto.SendRequest) string {
	if req.From != nil && req.From.Email != "" {
		return req.From.Email
	}
	return s.cfg.From
}

func (s *MailerService) buildMessage(req dto.SendRequest, subject, html, text string) smtpClient.Message {
	msg := smtpClient.Message{
		To: extractEmails(req.To),
		CC:          extractEmails(req.CC),
		BCC:         extractEmails(req.BCC),
		Subject:     subject,
		Html:        html,
		Text:        text,
	}

	if req.From != nil {
		msg.FromAddress = req.From.Email
		msg.FromName = req.From.Name
	}
	if req.ReplyTo != nil {
		msg.ReplyTo = req.ReplyTo.Email
	}

	for _, att := range req.Attachments {
		decoded, err := smtpClient.DecodeBase64Attachment(att.Content)
		if err != nil {
			continue
		}
		mime := att.MimeType
		if mime == "" {
			mime = "application/octet-stream"
		}
		msg.Attachments = append(msg.Attachments, smtpClient.Attachment{
			Filename: att.Filename,
			Content:  decoded,
			MimeType: mime,
		})
	}
	return msg
}

func extractEmails(addresses []dto.EmailAddress) []string {
	result := make([]string, 0, len(addresses))
	for _, a := range addresses {
		if strings.TrimSpace(a.Name) != "" {
			result = append(result, fmt.Sprintf("%s <%s>", a.Name, a.Email))
		} else {
			result = append(result, a.Email)
		}
	}
	return result
}
