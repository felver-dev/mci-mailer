package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mcicare/mci-mailer/internal/domain"
	"github.com/mcicare/mci-mailer/internal/dto"
)

type emailLogRepo struct {
	db *pgxpool.Pool
}

func NewEmailLogRepository(db *pgxpool.Pool) EmailLogRepository {
	return &emailLogRepo{db: db}
}

func (r *emailLogRepo) Create(ctx context.Context, log *domain.EmailLog) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO email_logs
		 (id, api_key_id, from_address, to_addresses, cc_addresses, bcc_addresses,
		  subject, template_name, status, attempts, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		log.ID, log.ApiKeyID, log.FromAddress, log.ToAddresses,
		log.CcAddresses, log.BccAddresses, log.Subject, log.TemplateName,
		log.Status, log.Attempts, log.CreatedAt,
	)
	return err
}

func (r *emailLogRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.EmailStatus, errMsg *string, attempts int) error {
	_, err := r.db.Exec(ctx,
		`UPDATE email_logs SET status = $1, error_msg = $2, attempts = $3 WHERE id = $4`,
		status, errMsg, attempts, id,
	)
	return err
}

func (r *emailLogRepo) MarkSent(ctx context.Context, id uuid.UUID, attempts int) error {
	_, err := r.db.Exec(ctx,
		`UPDATE email_logs SET status = 'sent', error_msg = NULL, attempts = $1, sent_at = NOW() WHERE id = $2`,
		attempts, id,
	)
	return err
}

func (r *emailLogRepo) FindAll(ctx context.Context, filter dto.LogFilter) ([]domain.EmailLog, int, error) {
	conditions := []string{"1=1"}
	args := []any{}
	i := 1

	if filter.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", i))
		args = append(args, filter.Status)
		i++
	}
	if filter.ApiKeyID != "" {
		conditions = append(conditions, fmt.Sprintf("api_key_id = $%d", i))
		args = append(args, filter.ApiKeyID)
		i++
	}
	if filter.From != "" {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", i))
		args = append(args, filter.From)
		i++
	}
	if filter.To != "" {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", i))
		args = append(args, filter.To)
		i++
	}

	where := strings.Join(conditions, " AND ")

	var total int
	err := r.db.QueryRow(ctx, fmt.Sprintf("SELECT COUNT(*) FROM email_logs WHERE %s", where), args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	if filter.PageSize <= 0 {
		filter.PageSize = 50
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.PageSize

	args = append(args, filter.PageSize, offset)
	query := fmt.Sprintf(
		`SELECT id, api_key_id, from_address, to_addresses, cc_addresses, bcc_addresses,
		        subject, template_name, status, error_msg, attempts, sent_at, created_at
		 FROM email_logs WHERE %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`,
		where, i, i+1,
	)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []domain.EmailLog
	for rows.Next() {
		var l domain.EmailLog
		if err := rows.Scan(
			&l.ID, &l.ApiKeyID, &l.FromAddress, &l.ToAddresses,
			&l.CcAddresses, &l.BccAddresses, &l.Subject, &l.TemplateName,
			&l.Status, &l.ErrorMsg, &l.Attempts, &l.SentAt, &l.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		logs = append(logs, l)
	}
	return logs, total, nil
}
