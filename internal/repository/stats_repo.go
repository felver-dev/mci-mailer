package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mcicare/mci-mailer/internal/dto"
)

type pgxStatsRepository struct {
	db *pgxpool.Pool
}

func NewStatsRepository(db *pgxpool.Pool) StatsRepository {
	return &pgxStatsRepository{db: db}
}

func (r *pgxStatsRepository) GetOverview(ctx context.Context) (*dto.StatsOverview, error) {
	o := &dto.StatsOverview{}

	err := r.db.QueryRow(ctx, `
		SELECT
			COUNT(*)                                                        AS total_apps,
			COUNT(*) FILTER (WHERE is_active)                              AS active_apps
		FROM api_keys
	`).Scan(&o.TotalApps, &o.ActiveApps)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&o.TotalUsers)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, `
		SELECT
			COUNT(*)                                                         AS total,
			COUNT(*) FILTER (WHERE status = 'sent')                         AS sent,
			COUNT(*) FILTER (WHERE status = 'failed')                       AS failed,
			COUNT(*) FILTER (WHERE status = 'sent' AND created_at >= CURRENT_DATE)                       AS today,
			COUNT(*) FILTER (WHERE status = 'sent' AND created_at >= DATE_TRUNC('month', NOW()))         AS this_month
		FROM email_logs
	`).Scan(&o.TotalEmails, &o.SentEmails, &o.FailedEmails, &o.SentToday, &o.SentThisMonth)
	if err != nil {
		return nil, err
	}

	if o.TotalEmails > 0 {
		o.SuccessRate = float64(o.SentEmails) / float64(o.TotalEmails) * 100
	}

	return o, nil
}

func (r *pgxStatsRepository) GetPerApp(ctx context.Context) ([]dto.AppStats, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			ak.id::TEXT,
			ak.name,
			COUNT(el.id)                                    AS total,
			COUNT(el.id) FILTER (WHERE el.status = 'sent') AS sent,
			COUNT(el.id) FILTER (WHERE el.status = 'failed') AS failed,
			u.name AS created_by
		FROM api_keys ak
		LEFT JOIN email_logs el ON el.api_key_id = ak.id
		LEFT JOIN users u ON u.id = ak.created_by_user_id
		GROUP BY ak.id, ak.name, u.name
		ORDER BY total DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []dto.AppStats
	for rows.Next() {
		var s dto.AppStats
		if err := rows.Scan(&s.AppID, &s.AppName, &s.Total, &s.Sent, &s.Failed, &s.CreatedByName); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}
