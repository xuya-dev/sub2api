package service

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	dbuser "github.com/Wei-Shaw/sub2api/ent/user"
	"github.com/Wei-Shaw/sub2api/internal/pkg/timezone"
)

type LeaderboardEntry struct {
	Rank          int     `json:"rank"`
	Username      string  `json:"username"`
	Value         float64 `json:"value"`
	ExtraInt      int     `json:"extra_int,omitempty"`
	ExtraInt2     int     `json:"extra_int2,omitempty"`
	ExtraFloat    float64 `json:"extra_float,omitempty"`
	ExtraDate     string  `json:"extra_date,omitempty"`
}

type LeaderboardResult struct {
	Entries []LeaderboardEntry `json:"items"`
	Total   int64              `json:"total"`
}

type LeaderboardService struct {
	entClient *dbent.Client
	db        *sql.DB
}

func NewLeaderboardService(entClient *dbent.Client, db *sql.DB) *LeaderboardService {
	return &LeaderboardService{
		entClient: entClient,
		db:        db,
	}
}

func (s *LeaderboardService) GetBalanceLeaderboard(ctx context.Context, page, pageSize int) (*LeaderboardResult, error) {
	offset := (page - 1) * pageSize

	total, err := s.entClient.User.Query().
		Where(
			dbuser.DeletedAtIsNil(),
			dbuser.StatusEQ(StatusActive),
			dbuser.RoleNEQ("admin"),
			dbuser.BalanceGT(0),
		).
		Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("count users: %w", err)
	}

	users, err := s.entClient.User.Query().
		Where(
			dbuser.DeletedAtIsNil(),
			dbuser.StatusEQ(StatusActive),
			dbuser.RoleNEQ("admin"),
			dbuser.BalanceGT(0),
		).
		Order(dbent.Desc(dbuser.FieldBalance)).
		Offset(offset).
		Limit(pageSize).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}

	entries := make([]LeaderboardEntry, 0, len(users))
	for i, u := range users {
		var checkinCount int
		s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM checkins WHERE user_id = $1`, u.ID).Scan(&checkinCount)
		entries = append(entries, LeaderboardEntry{
			Rank:     offset + i + 1,
			Username: maskUsername(u.Username, u.Email),
			Value:    math.Round(u.Balance*100) / 100,
			ExtraInt: checkinCount,
		})
	}

	return &LeaderboardResult{Entries: entries, Total: int64(total)}, nil
}

func (s *LeaderboardService) GetConsumptionLeaderboard(ctx context.Context, period string, page, pageSize int) (*LeaderboardResult, error) {
	today := timezone.Today()
	var startTime time.Time
	switch period {
	case "daily":
		startTime = today
	case "weekly":
		startTime = today.AddDate(0, 0, -7)
	case "monthly":
		startTime = today.AddDate(0, 0, -30)
	default:
		startTime = today
	}

	offset := (page - 1) * pageSize

	countQuery := `
		SELECT COUNT(*) FROM (
			SELECT ul.user_id
			FROM usage_logs ul
			INNER JOIN users u ON ul.user_id = u.id AND u.deleted_at IS NULL
			WHERE ul.created_at >= $1 AND u.status = 'active' AND u.role != 'admin'
			GROUP BY ul.user_id
			HAVING SUM(ul.actual_cost) > 0
		) sub
	`
	var total int64
	if err := s.db.QueryRowContext(ctx, countQuery, startTime).Scan(&total); err != nil {
		return nil, fmt.Errorf("count consumption: %w", err)
	}

	dataQuery := `
		SELECT u.username, u.email, COALESCE(SUM(ul.actual_cost), 0) as total_cost, COUNT(*) as request_count
		FROM usage_logs ul
		INNER JOIN users u ON ul.user_id = u.id AND u.deleted_at IS NULL
		WHERE ul.created_at >= $1 AND u.status = 'active' AND u.role != 'admin'
		GROUP BY ul.user_id, u.username, u.email
		HAVING SUM(ul.actual_cost) > 0
		ORDER BY total_cost DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := s.db.QueryContext(ctx, dataQuery, startTime, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("query consumption: %w", err)
	}
	defer rows.Close()

	entries := make([]LeaderboardEntry, 0)
	rank := offset
	for rows.Next() {
		rank++
		var username, email string
		var totalCost float64
		var requestCount int
		if err := rows.Scan(&username, &email, &totalCost, &requestCount); err != nil {
			return nil, fmt.Errorf("scan consumption row: %w", err)
		}
		entries = append(entries, LeaderboardEntry{
			Rank:     rank,
			Username: maskUsername(username, email),
			Value:    math.Round(totalCost*100) / 100,
			ExtraInt: requestCount,
		})
	}

	return &LeaderboardResult{Entries: entries, Total: total}, nil
}

func (s *LeaderboardService) GetCheckinLeaderboard(ctx context.Context, page, pageSize int) (*LeaderboardResult, error) {
	today := timezone.Today()
	yesterday := today.AddDate(0, 0, -1)

	offset := (page - 1) * pageSize

	countQuery := `
		SELECT COUNT(*) FROM (
			SELECT c.user_id
			FROM checkins c
			INNER JOIN (
				SELECT user_id, MAX(checkin_date) as max_date
				FROM checkins
				GROUP BY user_id
			) latest ON c.user_id = latest.user_id AND c.checkin_date = latest.max_date
			INNER JOIN users u ON c.user_id = u.id AND u.deleted_at IS NULL
			WHERE c.checkin_date >= $1 AND u.status = 'active' AND u.role != 'admin'
		) sub
	`
	var total int64
	if err := s.db.QueryRowContext(ctx, countQuery, yesterday).Scan(&total); err != nil {
		return nil, fmt.Errorf("count checkin: %w", err)
	}

	dataQuery := `
		SELECT u.username, u.email, c.streak_days, c.reward_amount,
			(SELECT COUNT(*) FROM checkins WHERE user_id = c.user_id) as total_checkins,
			(SELECT MAX(checkin_date) FROM checkins WHERE user_id = c.user_id) as last_date
		FROM checkins c
		INNER JOIN (
			SELECT user_id, MAX(checkin_date) as max_date
			FROM checkins
			GROUP BY user_id
		) latest ON c.user_id = latest.user_id AND c.checkin_date = latest.max_date
		INNER JOIN users u ON c.user_id = u.id AND u.deleted_at IS NULL
		WHERE c.checkin_date >= $1 AND u.status = 'active' AND u.role != 'admin'
		ORDER BY c.streak_days DESC, c.checkin_date DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := s.db.QueryContext(ctx, dataQuery, yesterday, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("query checkin: %w", err)
	}
	defer rows.Close()

	entries := make([]LeaderboardEntry, 0)
	rank := offset
	for rows.Next() {
		rank++
		var username, email string
		var streakDays int
		var rewardAmount float64
		var totalCheckins int
		var lastDate time.Time
		if err := rows.Scan(&username, &email, &streakDays, &rewardAmount, &totalCheckins, &lastDate); err != nil {
			return nil, fmt.Errorf("scan checkin row: %w", err)
		}
		entries = append(entries, LeaderboardEntry{
			Rank:       rank,
			Username:   maskUsername(username, email),
			Value:      float64(streakDays),
			ExtraInt:   totalCheckins,
			ExtraFloat: math.Round(rewardAmount*100) / 100,
			ExtraDate:  lastDate.Format("2006-01-02"),
		})
	}

	return &LeaderboardResult{Entries: entries, Total: total}, nil
}

func maskUsername(username, email string) string {
	if username != "" {
		return username
	}
	if email == "" {
		return "user"
	}
	if len(email) <= 3 {
		return email
	}
	return email[:3] + "***"
}
