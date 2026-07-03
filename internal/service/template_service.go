package service

import (
	"bytes"
	"context"
	"errors"
	"html/template"
	"time"

	"github.com/google/uuid"
	"github.com/mcicare/mci-mailer/internal/domain"
	"github.com/mcicare/mci-mailer/internal/dto"
	"github.com/mcicare/mci-mailer/internal/repository"
)

type TemplateService struct {
	repo repository.TemplateRepository
}

func NewTemplateService(repo repository.TemplateRepository) *TemplateService {
	return &TemplateService{repo: repo}
}

func (s *TemplateService) Create(ctx context.Context, req dto.CreateTemplateRequest) (*dto.TemplateResponse, error) {
	existing, _ := s.repo.FindByName(ctx, req.Name)
	if existing != nil {
		return nil, errors.New("template with this name already exists")
	}

	t := &domain.Template{
		ID:        uuid.New(),
		Name:      req.Name,
		Subject:   req.Subject,
		HtmlBody:  req.HtmlBody,
		TextBody:  req.TextBody,
		Variables: req.Variables,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, t); err != nil {
		return nil, err
	}
	return toTemplateResponse(t), nil
}

func (s *TemplateService) List(ctx context.Context) ([]dto.TemplateResponse, error) {
	templates, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]dto.TemplateResponse, len(templates))
	for i, t := range templates {
		t := t
		result[i] = *toTemplateResponse(&t)
	}
	return result, nil
}

func (s *TemplateService) Update(ctx context.Context, name string, req dto.UpdateTemplateRequest) (*dto.TemplateResponse, error) {
	t, err := s.repo.FindByName(ctx, name)
	if err != nil {
		return nil, errors.New("template not found")
	}
	if req.Subject != "" {
		t.Subject = req.Subject
	}
	if req.HtmlBody != "" {
		t.HtmlBody = req.HtmlBody
	}
	if req.TextBody != "" {
		t.TextBody = req.TextBody
	}
	if err := s.repo.Update(ctx, t); err != nil {
		return nil, err
	}
	return toTemplateResponse(t), nil
}

func (s *TemplateService) Delete(ctx context.Context, name string) error {
	if _, err := s.repo.FindByName(ctx, name); err != nil {
		return errors.New("template not found")
	}
	return s.repo.Delete(ctx, name)
}

func (s *TemplateService) Render(ctx context.Context, name string, variables map[string]any) (subject, html, text string, err error) {
	t, err := s.repo.FindByName(ctx, name)
	if err != nil {
		return "", "", "", errors.New("template not found: " + name)
	}

	subject, err = renderString(t.Subject, variables)
	if err != nil {
		return
	}
	html, err = renderString(t.HtmlBody, variables)
	if err != nil {
		return
	}
	if t.TextBody != "" {
		text, err = renderString(t.TextBody, variables)
	}
	return
}

func renderString(tmpl string, data map[string]any) (string, error) {
	t, err := template.New("").Parse(tmpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func toTemplateResponse(t *domain.Template) *dto.TemplateResponse {
	return &dto.TemplateResponse{
		ID:        t.ID.String(),
		Name:      t.Name,
		Subject:   t.Subject,
		HtmlBody:  t.HtmlBody,
		TextBody:  t.TextBody,
		Variables: t.Variables,
		CreatedAt: t.CreatedAt.Format(time.RFC3339),
		UpdatedAt: t.UpdatedAt.Format(time.RFC3339),
	}
}
