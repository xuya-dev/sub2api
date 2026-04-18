package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type MonitoringService struct {
	db *sql.DB
}

func NewMonitoringService(db *sql.DB) *MonitoringService {
	return &MonitoringService{db: db}
}

type GroupHealth struct {
	GroupID        int64  `json:"group_id"`
	GroupName      string `json:"group_name"`
	TotalAccounts  int    `json:"total_accounts"`
	ActiveAccounts int    `json:"active_accounts"`
	ErrorAccounts  int    `json:"error_accounts"`
	RateLimited    int    `json:"rate_limited"`
	Overload       int    `json:"overload"`
	Disabled       int    `json:"disabled"`
}

type ModelLatency struct {
	Model           string  `json:"model"`
	RequestCount    int     `json:"request_count"`
	SuccessCount    int     `json:"success_count"`
	ErrorCount      int     `json:"error_count"`
	AvgLatencyMs    float64 `json:"avg_latency_ms"`
	P50LatencyMs    float64 `json:"p50_latency_ms"`
	P95LatencyMs    float64 `json:"p95_latency_ms"`
	P99LatencyMs    float64 `json:"p99_latency_ms"`
	AvgFirstTokenMs float64 `json:"avg_first_token_ms"`
}

type GroupModelStats struct {
	GroupID      int64  `json:"group_id"`
	GroupName    string `json:"group_name"`
	Model        string `json:"model"`
	RequestCount int    `json:"request_count"`
	SuccessCount int    `json:"success_count"`
	ErrorCount   int    `json:"error_count"`
	AvgLatencyMs float64 `json:"avg_latency_ms"`
	P50LatencyMs float64 `json:"p50_latency_ms"`
	P95LatencyMs float64 `json:"p95_latency_ms"`
	AvgTTFT      float64 `json:"avg_ttft"`
}

type ErrorAccount struct {
	AccountID     int64  `json:"account_id"`
	AccountName   string `json:"account_name"`
	GroupName     string `json:"group_name"`
	Status        string `json:"status"`
	ErrorMessage  string `json:"error_message"`
	RateLimitedAt string `json:"rate_limited_at,omitempty"`
	OverloadUntil string `json:"overload_until,omitempty"`
}

type MonitoringOverview struct {
	Groups         []GroupHealth     `json:"groups"`
	GroupModels    []GroupModelStats `json:"group_models"`
	ModelLatencies []ModelLatency    `json:"model_latencies"`
	ErrorAccounts  []ErrorAccount    `json:"error_accounts"`
	TotalRequests  int64             `json:"total_requests_today"`
	AvgLatencyMs   float64           `json:"avg_latency_ms_today"`
}

func (s *MonitoringService) GetOverview(ctx context.Context) (*MonitoringOverview, error) {
	overview := &MonitoringOverview{}

	if err := s.queryGroupHealth(ctx, overview); err != nil {
		log.Printf("[Monitoring] queryGroupHealth error: %v", err)
		return nil, fmt.Errorf("query group health: %w", err)
	}
	if err := s.queryGroupModelStats(ctx, overview); err != nil {
		log.Printf("[Monitoring] queryGroupModelStats error: %v", err)
		return nil, fmt.Errorf("query group model stats: %w", err)
	}
	if err := s.queryModelLatency(ctx, overview); err != nil {
		log.Printf("[Monitoring] queryModelLatency error: %v", err)
		return nil, fmt.Errorf("query model latency: %w", err)
	}
	if err := s.queryErrorAccounts(ctx, overview); err != nil {
		log.Printf("[Monitoring] queryErrorAccounts error: %v", err)
		return nil, fmt.Errorf("query error accounts: %w", err)
	}
	if err := s.queryTodaySummary(ctx, overview); err != nil {
		log.Printf("[Monitoring] queryTodaySummary error: %v", err)
		return nil, fmt.Errorf("query today summary: %w", err)
	}

	return overview, nil
}

func (s *MonitoringService) queryGroupHealth(ctx context.Context, overview *MonitoringOverview) error {
	now := time.Now().UTC()
	query := `
		SELECT
			g.id,
			COALESCE(g.name, ''),
			COUNT(a.id),
			COUNT(a.id) FILTER (WHERE a.status = 'active' AND a.schedulable = true
				AND (a.rate_limited_at IS NULL OR a.rate_limit_reset_at <= $1)
				AND (a.overload_until IS NULL OR a.overload_until <= $1)
				AND (a.temp_unschedulable_until IS NULL OR a.temp_unschedulable_until <= $1)),
			COUNT(a.id) FILTER (WHERE a.status = 'error'),
			COUNT(a.id) FILTER (WHERE a.rate_limited_at IS NOT NULL AND a.rate_limit_reset_at > $1),
			COUNT(a.id) FILTER (WHERE a.overload_until IS NOT NULL AND a.overload_until > $1),
			COUNT(a.id) FILTER (WHERE a.status = 'disabled' OR a.schedulable = false)
		FROM groups g
		LEFT JOIN accounts a ON a.group_accounts = g.id AND a.deleted_at IS NULL
		WHERE g.deleted_at IS NULL
		GROUP BY g.id, g.name
		ORDER BY g.name`

	rows, err := s.db.QueryContext(ctx, query, now)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var g GroupHealth
		if err := rows.Scan(&g.GroupID, &g.GroupName, &g.TotalAccounts, &g.ActiveAccounts, &g.ErrorAccounts, &g.RateLimited, &g.Overload, &g.Disabled); err != nil {
			return err
		}
		overview.Groups = append(overview.Groups, g)
	}
	return rows.Err()
}

func (s *MonitoringService) queryGroupModelStats(ctx context.Context, overview *MonitoringOverview) error {
	since := time.Now().UTC().Add(-24 * time.Hour)
	query := `
		SELECT
			g.id,
			COALESCE(g.name, ''),
			u.model,
			COUNT(*) as cnt,
			COUNT(*) FILTER (WHERE u.output_tokens > 0) as success_cnt,
			COUNT(*) FILTER (WHERE u.output_tokens = 0) as error_cnt,
			COALESCE(AVG(u.duration_ms) FILTER (WHERE u.duration_ms IS NOT NULL AND u.duration_ms > 0), 0)::float8,
			COALESCE(percentile_cont(0.5) WITHIN GROUP (ORDER BY CASE WHEN u.duration_ms IS NOT NULL AND u.duration_ms > 0 THEN u.duration_ms END), 0)::float8,
			COALESCE(percentile_cont(0.95) WITHIN GROUP (ORDER BY CASE WHEN u.duration_ms IS NOT NULL AND u.duration_ms > 0 THEN u.duration_ms END), 0)::float8,
			COALESCE(AVG(u.first_token_ms) FILTER (WHERE u.first_token_ms IS NOT NULL AND u.first_token_ms > 0), 0)::float8
		FROM usage_logs u
		JOIN groups g ON u.group_id = g.id
		WHERE u.created_at >= $1
		  AND g.deleted_at IS NULL
		GROUP BY g.id, g.name, u.model
		HAVING COUNT(*) > 0
		ORDER BY g.name, cnt DESC`

	rows, err := s.db.QueryContext(ctx, query, since)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var m GroupModelStats
		if err := rows.Scan(&m.GroupID, &m.GroupName, &m.Model, &m.RequestCount, &m.SuccessCount, &m.ErrorCount, &m.AvgLatencyMs, &m.P50LatencyMs, &m.P95LatencyMs, &m.AvgTTFT); err != nil {
			return err
		}
		overview.GroupModels = append(overview.GroupModels, m)
	}
	return rows.Err()
}

func (s *MonitoringService) queryModelLatency(ctx context.Context, overview *MonitoringOverview) error {
	since := time.Now().UTC().Add(-24 * time.Hour)
	query := `
		SELECT
			model,
			COUNT(*) as cnt,
			COUNT(*) FILTER (WHERE output_tokens > 0) as success_cnt,
			COUNT(*) FILTER (WHERE output_tokens = 0) as error_cnt,
			COALESCE(AVG(duration_ms) FILTER (WHERE duration_ms IS NOT NULL AND duration_ms > 0), 0)::float8,
			COALESCE(percentile_cont(0.5) WITHIN GROUP (ORDER BY CASE WHEN duration_ms IS NOT NULL AND duration_ms > 0 THEN duration_ms END), 0)::float8,
			COALESCE(percentile_cont(0.95) WITHIN GROUP (ORDER BY CASE WHEN duration_ms IS NOT NULL AND duration_ms > 0 THEN duration_ms END), 0)::float8,
			COALESCE(percentile_cont(0.99) WITHIN GROUP (ORDER BY CASE WHEN duration_ms IS NOT NULL AND duration_ms > 0 THEN duration_ms END), 0)::float8,
			COALESCE(AVG(first_token_ms) FILTER (WHERE first_token_ms IS NOT NULL AND first_token_ms > 0), 0)::float8
		FROM usage_logs
		WHERE created_at >= $1
		GROUP BY model
		ORDER BY cnt DESC
		LIMIT 30`

	rows, err := s.db.QueryContext(ctx, query, since)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var m ModelLatency
		if err := rows.Scan(&m.Model, &m.RequestCount, &m.SuccessCount, &m.ErrorCount, &m.AvgLatencyMs, &m.P50LatencyMs, &m.P95LatencyMs, &m.P99LatencyMs, &m.AvgFirstTokenMs); err != nil {
			return err
		}
		overview.ModelLatencies = append(overview.ModelLatencies, m)
	}
	return rows.Err()
}

func (s *MonitoringService) queryErrorAccounts(ctx context.Context, overview *MonitoringOverview) error {
	query := `
		SELECT
			a.id,
			COALESCE(a.name, ''),
			COALESCE(g.name, ''),
			a.status,
			COALESCE(a.error_message, ''),
			a.rate_limited_at,
			a.overload_until
		FROM accounts a
		LEFT JOIN groups g ON a.group_accounts = g.id
		WHERE a.deleted_at IS NULL
		  AND (a.status = 'error'
		       OR (a.rate_limited_at IS NOT NULL AND a.rate_limit_reset_at > NOW())
		       OR (a.overload_until IS NOT NULL AND a.overload_until > NOW()))
		ORDER BY
			CASE a.status WHEN 'error' THEN 0 ELSE 1 END,
			a.rate_limited_at DESC NULLS LAST
		LIMIT 50`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var e ErrorAccount
		var rateLimitedAt, overloadUntil sql.NullTime
		if err := rows.Scan(&e.AccountID, &e.AccountName, &e.GroupName, &e.Status, &e.ErrorMessage, &rateLimitedAt, &overloadUntil); err != nil {
			return err
		}
		if rateLimitedAt.Valid {
			e.RateLimitedAt = rateLimitedAt.Time.Format("2006-01-02 15:04:05")
		}
		if overloadUntil.Valid {
			e.OverloadUntil = overloadUntil.Time.Format("2006-01-02 15:04:05")
		}
		overview.ErrorAccounts = append(overview.ErrorAccounts, e)
	}
	return rows.Err()
}

func (s *MonitoringService) queryTodaySummary(ctx context.Context, overview *MonitoringOverview) error {
	todayStart := time.Now().UTC().Truncate(24 * time.Hour)
	query := `
		SELECT
			COUNT(*),
			COALESCE(AVG(duration_ms), 0)::float8
		FROM usage_logs
		WHERE created_at >= $1
		  AND duration_ms IS NOT NULL`

	return s.db.QueryRowContext(ctx, query, todayStart).Scan(&overview.TotalRequests, &overview.AvgLatencyMs)
}
