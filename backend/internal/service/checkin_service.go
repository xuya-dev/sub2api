package service

import (
	"context"
	"fmt"
	"math"
	"math/rand/v2"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/checkin"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

var (
	ErrCheckinDisabled   = infraerrors.Forbidden("CHECKIN_DISABLED", "check-in feature is not enabled")
	ErrAlreadyCheckedIn  = infraerrors.Conflict("ALREADY_CHECKED_IN", "you have already checked in today")
	ErrCheckinNotAllowed = infraerrors.Forbidden("CHECKIN_NOT_ALLOWED", "check-in is not allowed for your account")
)

type CheckinResult struct {
	RewardAmount float64 `json:"reward_amount"`
	StreakDays   int     `json:"streak_days"`
	CheckedAt    string  `json:"checked_at"`
}

type CheckinStatus struct {
	CanCheckin  bool     `json:"can_checkin"`
	StreakDays  int      `json:"streak_days"`
	TodayReward *float64 `json:"today_reward,omitempty"`
	MinReward   float64  `json:"min_reward"`
	MaxReward   float64  `json:"max_reward"`
}

type CheckinService struct {
	entClient            *dbent.Client
	userRepo             UserRepository
	settingService       *SettingService
	billingCacheService  *BillingCacheService
	authCacheInvalidator APIKeyAuthCacheInvalidator
}

func NewCheckinService(
	entClient *dbent.Client,
	userRepo UserRepository,
	settingService *SettingService,
	billingCacheService *BillingCacheService,
	authCacheInvalidator APIKeyAuthCacheInvalidator,
) *CheckinService {
	return &CheckinService{
		entClient:            entClient,
		userRepo:             userRepo,
		settingService:       settingService,
		billingCacheService:  billingCacheService,
		authCacheInvalidator: authCacheInvalidator,
	}
}

func (s *CheckinService) Checkin(ctx context.Context, userID int64) (*CheckinResult, error) {
	if !s.settingService.IsCheckinEnabled(ctx) {
		return nil, ErrCheckinDisabled
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if user.Status != StatusActive {
		return nil, ErrCheckinNotAllowed
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)
	todayDate := today.Format("2006-01-02")

	existing, err := s.entClient.Checkin.
		Query().
		Where(
			checkin.UserID(userID),
			checkin.CheckinDateEQ(today),
		).
		Only(ctx)
	if err == nil && existing != nil {
		reward := existing.RewardAmount
		return &CheckinResult{
			RewardAmount: reward,
			StreakDays:   existing.StreakDays,
			CheckedAt:    todayDate,
		}, nil
	}
	if err != nil && !dbent.IsNotFound(err) {
		return nil, fmt.Errorf("query checkin: %w", err)
	}

	minReward, maxReward := s.settingService.GetCheckinBalanceRange(ctx)
	rewardAmount := minReward + rand.Float64()*(maxReward-minReward)
	rewardAmount = math.Round(rewardAmount*100) / 100

	streakDays := s.calculateStreak(ctx, userID, today)

	tx, err := s.entClient.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	txCtx := dbent.NewTxContext(ctx, tx)

	_, err = tx.Client().Checkin.
		Create().
		SetUserID(userID).
		SetCheckinDate(today).
		SetRewardAmount(rewardAmount).
		SetStreakDays(streakDays).
		Save(txCtx)
	if err != nil {
		return nil, fmt.Errorf("create checkin record: %w", err)
	}

	if err := s.userRepo.UpdateBalance(txCtx, userID, rewardAmount); err != nil {
		return nil, fmt.Errorf("update user balance: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	s.invalidateCaches(ctx, userID)

	return &CheckinResult{
		RewardAmount: rewardAmount,
		StreakDays:   streakDays,
		CheckedAt:    todayDate,
	}, nil
}

func (s *CheckinService) GetStatus(ctx context.Context, userID int64) (*CheckinStatus, error) {
	enabled := s.settingService.IsCheckinEnabled(ctx)
	minReward, maxReward := s.settingService.GetCheckinBalanceRange(ctx)

	if !enabled {
		return &CheckinStatus{
			CanCheckin: false,
			StreakDays: 0,
			MinReward:  minReward,
			MaxReward:  maxReward,
		}, nil
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)

	todayCheckin, err := s.entClient.Checkin.
		Query().
		Where(
			checkin.UserID(userID),
			checkin.CheckinDateEQ(today),
		).
		Only(ctx)
	if err != nil && !dbent.IsNotFound(err) {
		return nil, fmt.Errorf("query today checkin: %w", err)
	}

	status := &CheckinStatus{
		CanCheckin: true,
		MinReward:  minReward,
		MaxReward:  maxReward,
	}

	if todayCheckin != nil {
		status.CanCheckin = false
		status.StreakDays = todayCheckin.StreakDays
		reward := todayCheckin.RewardAmount
		status.TodayReward = &reward
	} else {
		status.StreakDays = s.calculateStreak(ctx, userID, today)
	}

	return status, nil
}

func (s *CheckinService) calculateStreak(ctx context.Context, userID int64, today time.Time) int {
	yesterday := today.AddDate(0, 0, -1)

	lastCheckin, err := s.entClient.Checkin.
		Query().
		Where(
			checkin.UserID(userID),
			checkin.CheckinDateLTE(today),
		).
		Order(dbent.Desc(checkin.FieldCheckinDate)).
		First(ctx)
	if err != nil {
		return 1
	}

	if lastCheckin.CheckinDate.Equal(yesterday) {
		return lastCheckin.StreakDays + 1
	}

	return 1
}

func (s *CheckinService) invalidateCaches(ctx context.Context, userID int64) {
	if s.authCacheInvalidator != nil {
		s.authCacheInvalidator.InvalidateAuthCacheByUserID(ctx, userID)
	}
	if s.billingCacheService != nil {
		go func() {
			cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = s.billingCacheService.InvalidateUserBalance(cacheCtx, userID)
		}()
	}
}
