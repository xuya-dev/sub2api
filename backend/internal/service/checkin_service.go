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
	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	"github.com/Wei-Shaw/sub2api/internal/pkg/timezone"
)

const (
	CheckinTypeNormal = "normal"
	CheckinTypeLuck   = "luck"
)

var (
	ErrCheckinDisabled     = infraerrors.Forbidden("CHECKIN_DISABLED", "check-in feature is not enabled")
	ErrCheckinLuckDisabled = infraerrors.Forbidden("CHECKIN_LUCK_DISABLED", "luck check-in feature is not enabled")
	ErrAlreadyCheckedIn    = infraerrors.Conflict("ALREADY_CHECKED_IN", "you have already checked in today")
	ErrCheckinNotAllowed   = infraerrors.Forbidden("CHECKIN_NOT_ALLOWED", "check-in is not allowed for your account")
	ErrInvalidBetAmount    = infraerrors.BadRequest("INVALID_BET_AMOUNT", "bet amount must be greater than 0 and not exceed your balance")
)

type CheckinResult struct {
	RewardAmount float64        `json:"reward_amount"`
	StreakDays   int            `json:"streak_days"`
	CheckedAt    string         `json:"checked_at"`
	CheckinType  string         `json:"checkin_type"`
	BetAmount    float64        `json:"bet_amount,omitempty"`
	Multiplier   float64        `json:"multiplier,omitempty"`
	Blindbox     *BlindboxResult `json:"blindbox,omitempty"`
}

type CheckinStatus struct {
	Enabled              bool     `json:"enabled"`
	LuckEnabled          bool     `json:"luck_enabled"`
	BlindboxEnabled      bool     `json:"blindbox_enabled"`
	BlindboxTriggerType  string   `json:"blindbox_trigger_type,omitempty"`
	BlindboxInterval     int      `json:"blindbox_interval,omitempty"`
	CanCheckin           bool     `json:"can_checkin"`
	StreakDays           int      `json:"streak_days"`
	TodayReward          *float64 `json:"today_reward,omitempty"`
	TodayCheckinType     string   `json:"today_checkin_type,omitempty"`
	TodayMultiplier      *float64 `json:"today_multiplier,omitempty"`
	MinReward            float64  `json:"min_reward"`
	MaxReward            float64  `json:"max_reward"`
	MinMultiplier        float64  `json:"min_multiplier"`
	MaxMultiplier        float64  `json:"max_multiplier"`
	Balance              float64  `json:"balance"`
}

type CheckinService struct {
	entClient            *dbent.Client
	userRepo             UserRepository
	redeemCodeRepo       RedeemCodeRepository
	settingService       *SettingService
	billingCacheService  *BillingCacheService
	authCacheInvalidator APIKeyAuthCacheInvalidator
	blindboxService      *BlindBoxService
}

func NewCheckinService(
	entClient *dbent.Client,
	userRepo UserRepository,
	redeemCodeRepo RedeemCodeRepository,
	settingService *SettingService,
	billingCacheService *BillingCacheService,
	authCacheInvalidator APIKeyAuthCacheInvalidator,
	blindboxService *BlindBoxService,
) *CheckinService {
	return &CheckinService{
		entClient:            entClient,
		userRepo:             userRepo,
		redeemCodeRepo:       redeemCodeRepo,
		settingService:       settingService,
		billingCacheService:  billingCacheService,
		authCacheInvalidator: authCacheInvalidator,
		blindboxService:      blindboxService,
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

	today := timezone.Today()
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
			CheckinType:  existing.CheckinType,
			BetAmount:    existing.BetAmount,
			Multiplier:   existing.Multiplier,
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
		SetCheckinType(CheckinTypeNormal).
		Save(txCtx)
	if err != nil {
		return nil, fmt.Errorf("create checkin record: %w", err)
	}

	if err := s.userRepo.UpdateBalance(txCtx, userID, rewardAmount); err != nil {
		return nil, fmt.Errorf("update user balance: %w", err)
	}

	s.createAuditRecord(txCtx, userID, rewardAmount, AdjustmentTypeCheckin, 0, 0)

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	s.invalidateCaches(ctx, userID)

	result := &CheckinResult{
		RewardAmount: rewardAmount,
		StreakDays:   streakDays,
		CheckedAt:    todayDate,
		CheckinType:  CheckinTypeNormal,
	}

	if s.blindboxService != nil && s.blindboxService.ShouldTriggerBlindbox(ctx, userID, streakDays) {
		blindboxResult, err := s.blindboxService.Draw(ctx, userID, streakDays)
		if err == nil && blindboxResult != nil {
			result.Blindbox = blindboxResult
		} else if err != nil {
			logger.LegacyPrintf("service.checkin", "blindbox draw failed for user %d: %v", userID, err)
		}
	}

	return result, nil
}

func (s *CheckinService) LuckCheckin(ctx context.Context, userID int64, betAmount float64) (*CheckinResult, error) {
	if !s.settingService.IsCheckinLuckEnabled(ctx) {
		return nil, ErrCheckinLuckDisabled
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if user.Status != StatusActive {
		return nil, ErrCheckinNotAllowed
	}

	if betAmount <= 0 || betAmount > user.Balance {
		return nil, ErrInvalidBetAmount
	}

	if user.Balance <= 0 {
		return nil, ErrInvalidBetAmount
	}

	today := timezone.Today()
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
			CheckinType:  existing.CheckinType,
			BetAmount:    existing.BetAmount,
			Multiplier:   existing.Multiplier,
		}, nil
	}
	if err != nil && !dbent.IsNotFound(err) {
		return nil, fmt.Errorf("query checkin: %w", err)
	}

	minMultiplier, maxMultiplier := s.settingService.GetCheckinLuckMultiplierRange(ctx)
	multiplier := minMultiplier + rand.Float64()*(maxMultiplier-minMultiplier)
	multiplier = math.Round(multiplier*100) / 100

	rewardAmount := math.Round(betAmount*(multiplier-1)*100) / 100

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
		SetCheckinType(CheckinTypeLuck).
		SetBetAmount(betAmount).
		SetMultiplier(multiplier).
		Save(txCtx)
	if err != nil {
		return nil, fmt.Errorf("create checkin record: %w", err)
	}

	if rewardAmount != 0 {
		if err := s.userRepo.UpdateBalance(txCtx, userID, rewardAmount); err != nil {
			return nil, fmt.Errorf("update user balance: %w", err)
		}
	}

	s.createAuditRecord(txCtx, userID, rewardAmount, AdjustmentTypeCheckinLuck, multiplier, betAmount)

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	s.invalidateCaches(ctx, userID)

	result := &CheckinResult{
		RewardAmount: rewardAmount,
		StreakDays:   streakDays,
		CheckedAt:    todayDate,
		CheckinType:  CheckinTypeLuck,
		BetAmount:    betAmount,
		Multiplier:   multiplier,
	}

	if s.blindboxService != nil && s.blindboxService.ShouldTriggerBlindbox(ctx, userID, streakDays) {
		blindboxResult, err := s.blindboxService.Draw(ctx, userID, streakDays)
		if err == nil && blindboxResult != nil {
			result.Blindbox = blindboxResult
		} else if err != nil {
			logger.LegacyPrintf("service.checkin", "blindbox draw failed for user %d (luck): %v", userID, err)
		}
	}

	return result, nil
}

func (s *CheckinService) GetStatus(ctx context.Context, userID int64) (*CheckinStatus, error) {
	normalEnabled := s.settingService.IsCheckinEnabled(ctx)
	luckEnabled := s.settingService.IsCheckinLuckEnabled(ctx)
	minReward, maxReward := s.settingService.GetCheckinBalanceRange(ctx)
	minMultiplier, maxMultiplier := s.settingService.GetCheckinLuckMultiplierRange(ctx)

	anyEnabled := normalEnabled || luckEnabled

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	if !anyEnabled {
		return &CheckinStatus{
			Enabled:             normalEnabled,
			LuckEnabled:         luckEnabled,
			BlindboxEnabled:     s.settingService.IsCheckinBlindboxEnabled(ctx),
			CanCheckin:          false,
			MinReward:           minReward,
			MaxReward:           maxReward,
			MinMultiplier:       minMultiplier,
			MaxMultiplier:       maxMultiplier,
			Balance:             user.Balance,
		}, nil
	}

	today := timezone.Today()

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
		Enabled:             normalEnabled,
		LuckEnabled:         luckEnabled,
		BlindboxEnabled:     s.settingService.IsCheckinBlindboxEnabled(ctx),
		BlindboxTriggerType: s.settingService.GetCheckinBlindboxTriggerType(ctx),
		BlindboxInterval:    s.settingService.GetCheckinBlindboxInterval(ctx),
		CanCheckin:          true,
		MinReward:           minReward,
		MaxReward:           maxReward,
		MinMultiplier:       minMultiplier,
		MaxMultiplier:       maxMultiplier,
		Balance:             user.Balance,
	}

	if todayCheckin != nil {
		status.CanCheckin = false
		status.StreakDays = todayCheckin.StreakDays
		reward := todayCheckin.RewardAmount
		status.TodayReward = &reward
		status.TodayCheckinType = todayCheckin.CheckinType
		if todayCheckin.CheckinType == CheckinTypeLuck {
			multiplier := todayCheckin.Multiplier
			status.TodayMultiplier = &multiplier
		}
	} else {
		status.StreakDays = s.calculateStreak(ctx, userID, today)
	}

	return status, nil
}

func sameDate(a, b time.Time) bool {
	aY, aM, aD := a.Date()
	bY, bM, bD := b.Date()
	return aY == bY && aM == bM && aD == bD
}

func (s *CheckinService) calculateStreak(ctx context.Context, userID int64, today time.Time) int {
	yesterday := today.AddDate(0, 0, -1)

	lastCheckin, err := s.entClient.Checkin.
		Query().
		Where(
			checkin.UserID(userID),
			checkin.CheckinDateLT(today),
		).
		Order(dbent.Desc(checkin.FieldCheckinDate)).
		First(ctx)
	if err != nil {
		return 1
	}

	if sameDate(lastCheckin.CheckinDate, yesterday) {
		return lastCheckin.StreakDays + 1
	}

	return 1
}

func (s *CheckinService) createAuditRecord(txCtx context.Context, userID int64, rewardAmount float64, adjType string, multiplier float64, betAmount float64) {
	code, err := GenerateRedeemCode()
	if err != nil {
		return
	}
	now := time.Now()
	adjustmentRecord := &RedeemCode{
		Code:       code,
		Type:       adjType,
		Value:      rewardAmount,
		Status:     StatusUsed,
		UsedBy:     &userID,
		UsedAt:     &now,
		Multiplier: multiplier,
		BetAmount:  betAmount,
	}
	if createErr := s.redeemCodeRepo.Create(txCtx, adjustmentRecord); createErr != nil {
		logger.LegacyPrintf("service.checkin", "failed to create checkin redeem code record: %v", createErr)
	}
}

type CheckinCalendarDay struct {
	Date        string  `json:"date"`
	CheckedIn   bool    `json:"checked_in"`
	RewardType  string  `json:"reward_type,omitempty"`
	RewardValue float64 `json:"reward_value,omitempty"`
	StreakDays  int     `json:"streak_days,omitempty"`
}

type CheckinCalendar struct {
	Days []CheckinCalendarDay `json:"days"`
}

func (s *CheckinService) GetCalendar(ctx context.Context, userID int64) (*CheckinCalendar, error) {
	today := timezone.Today()
	startDate := today.AddDate(0, 0, -29)

	records, err := s.entClient.Checkin.
		Query().
		Where(
			checkin.UserID(userID),
			checkin.CheckinDateGTE(startDate),
			checkin.CheckinDateLTE(today),
		).
		Order(dbent.Asc(checkin.FieldCheckinDate)).
		All(ctx)
	if err != nil && !dbent.IsNotFound(err) {
		return nil, fmt.Errorf("query calendar: %w", err)
	}

	recordMap := make(map[string]*dbent.Checkin, len(records))
	for i := range records {
		key := records[i].CheckinDate.Format("2006-01-02")
		recordMap[key] = records[i]
	}

	days := make([]CheckinCalendarDay, 30)
	for i := 0; i < 30; i++ {
		d := startDate.AddDate(0, 0, i)
		key := d.Format("2006-01-02")
		days[i] = CheckinCalendarDay{Date: key}

		if rec, ok := recordMap[key]; ok {
			days[i].CheckedIn = true
			days[i].RewardType = rec.CheckinType
			days[i].RewardValue = rec.RewardAmount
			days[i].StreakDays = rec.StreakDays
		}
	}

	return &CheckinCalendar{Days: days}, nil
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
