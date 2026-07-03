package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mcicare/mci-mailer/internal/domain"
)

type templateRepo struct {
	db *pgxpool.Pool
}

func NewTemplateRepository(db *pgxpool.Pool) TemplateRepository {
	return &templateRepo{db: db}
}

func (r *templateRepo) Create(ctx context.Context, t *domain.Template) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO templates (id, name, subject, html_body, text_body, variables, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		t.ID, t.Name, t.Subject, t.HtmlBody, t.TextBody, t.Variables, t.CreatedAt, t.UpdatedAt,
	)
	return err
}

func (r *templateRepo) FindByName(ctx context.Context, name string) (*domain.Template, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, name, subject, html_body, text_body, variables, created_at, updated_at
		 FROM templates WHERE name = $1`,
		name,
	)
	t := &domain.Template{}
	err := row.Scan(&t.ID, &t.Name, &t.Subject, &t.HtmlBody, &t.TextBody, &t.Variables, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (r *templateRepo) FindAll(ctx context.Context) ([]domain.Template, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, name, subject, html_body, text_body, variables, created_at, updated_at
		 FROM templates ORDER BY name ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []domain.Template
	for rows.Next() {
		var t domain.Template
		if err := rows.Scan(&t.ID, &t.Name, &t.Subject, &t.HtmlBody, &t.TextBody, &t.Variables, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		templates = append(templates, t)
	}
	return templates, nil
}

func (r *templateRepo) Update(ctx context.Context, t *domain.Template) error {
	_, err := r.db.Exec(ctx,
		`UPDATE templates SET subject = $1, html_body = $2, text_body = $3, updated_at = NOW()
		 WHERE name = $4`,
		t.Subject, t.HtmlBody, t.TextBody, t.Name,
	)
	return err
}

func (r *templateRepo) Delete(ctx context.Context, name string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM templates WHERE name = $1`, name)
	return err
}
