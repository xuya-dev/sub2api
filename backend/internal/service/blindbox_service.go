package service

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"math/rand/v2"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/checkinblindboxrecord"
	"github.com/Wei-Shaw/sub2api/ent/checkinprizeitem"
)

const (
	BlindboxRewardBalance         = "balance"
	BlindboxRewardConcurrency     = "concurrency"
	BlindboxRewardSubscription    = "subscription"
	BlindboxRewardInvitationCode  = "invitation_code"

	RarityCommon    = "common"
	RarityRare      = "rare"
	RarityEpic      = "epic"
	RarityLegendary = "legendary"
)

type PrizeItem struct {
	ID              int64   `json:"id"`
	Name            string  `json:"name"`
	Rarity          string  `json:"rarity"`
	RewardType      string  `json:"reward_type"`
	RewardValue     float64 `json:"reward_value"`
	RewardValueMax  float64 `json:"reward_value_max"`
	SubscriptionID  *int64  `json:"subscription_id,omitempty"`
	SubscriptionDays int    `json:"subscription_days"`
	Weight          int     `json:"weight"`
	IsEnabled       bool    `json:"is_enabled"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

type BlindboxResult struct {
	PrizeName       string  `json:"prize_name"`
	Rarity          string  `json:"rarity"`
	RewardType      string  `json:"reward_type"`
	RewardValue     float64 `json:"reward_value"`
	SubscriptionDays int    `json:"subscription_days,omitempty"`
}

type BlindboxRecord struct {
	ID          int64   `json:"id"`
	PrizeName   string  `json:"prize_name"`
	Rarity      string  `json:"rarity"`
	RewardType  string  `json:"reward_type"`
	RewardValue float64 `json:"reward_value"`
	StreakDays  int     `json:"streak_days"`
	CreatedAt   string  `json:"created_at"`
}

type BlindboxRecordList struct {
	Items []BlindboxRecord `json:"items"`
	Total int64            `json:"total"`
}

type BlindBoxService struct {
	entClient         *dbent.Client
	db                *sql.DB
	settingSvc        *SettingService
	userRepo          UserRepository
	billingCache      *BillingCacheService
	subscriptionSvc   *SubscriptionService
}

func NewBlindBoxService(
	entClient *dbent.Client,
	db *sql.DB,
	settingSvc *SettingService,
	userRepo UserRepository,
	billingCache *BillingCacheService,
	subscriptionSvc *SubscriptionService,
) *BlindBoxService {
	return &BlindBoxService{
		entClient:       entClient,
		db:              db,
		settingSvc:      settingSvc,
		userRepo:        userRepo,
		billingCache:    billingCache,
		subscriptionSvc: subscriptionSvc,
	}
}

func (s *BlindBoxService) ListPrizeItems(ctx context.Context) ([]PrizeItem, error) {
	items, err := s.entClient.CheckinPrizeItem.Query().
		Where(checkinprizeitem.DeletedAtIsNil()).
		Order(dbent.Desc(checkinprizeitem.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query prize items: %w", err)
	}

	result := make([]PrizeItem, 0, len(items))
	for _, item := range items {
		result = append(result, prizeItemFromEnt(item))
	}
	return result, nil
}

type CreatePrizeItemRequest struct {
	Name             string  `json:"name" binding:"required"`
	Rarity           string  `json:"rarity" binding:"required"`
	RewardType       string  `json:"reward_type" binding:"required"`
	RewardValue      float64 `json:"reward_value"`
	RewardValueMax   float64 `json:"reward_value_max"`
	SubscriptionID   *int64  `json:"subscription_id"`
	SubscriptionDays int     `json:"subscription_days"`
	Weight           int     `json:"weight"`
	IsEnabled        *bool   `json:"is_enabled"`
}

func (s *BlindBoxService) CreatePrizeItem(ctx context.Context, req CreatePrizeItemRequest) (*PrizeItem, error) {
	builder := s.entClient.CheckinPrizeItem.Create().
		SetName(req.Name).
		SetRarity(req.Rarity).
		SetRewardType(req.RewardType).
		SetRewardValue(req.RewardValue).
		SetRewardValueMax(req.RewardValueMax).
		SetSubscriptionDays(req.SubscriptionDays).
		SetWeight(req.Weight)

	if req.SubscriptionID != nil {
		builder.SetSubscriptionID(*req.SubscriptionID)
	}
	if req.IsEnabled != nil {
		builder.SetIsEnabled(*req.IsEnabled)
	}
	if req.Weight <= 0 {
		builder.SetWeight(100)
	}

	item, err := builder.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("create prize item: %w", err)
	}

	result := prizeItemFromEnt(item)
	return &result, nil
}

type UpdatePrizeItemRequest struct {
	Name             *string  `json:"name"`
	Rarity           *string  `json:"rarity"`
	RewardType       *string  `json:"reward_type"`
	RewardValue      *float64 `json:"reward_value"`
	RewardValueMax   *float64 `json:"reward_value_max"`
	SubscriptionID   **int64  `json:"subscription_id"`
	SubscriptionDays *int     `json:"subscription_days"`
	Weight           *int     `json:"weight"`
	IsEnabled        *bool    `json:"is_enabled"`
}

func (s *BlindBoxService) UpdatePrizeItem(ctx context.Context, id int64, req UpdatePrizeItemRequest) (*PrizeItem, error) {
	builder := s.entClient.CheckinPrizeItem.UpdateOneID(id)
	if req.Name != nil {
		builder.SetName(*req.Name)
	}
	if req.Rarity != nil {
		builder.SetRarity(*req.Rarity)
	}
	if req.RewardType != nil {
		builder.SetRewardType(*req.RewardType)
	}
	if req.RewardValue != nil {
		builder.SetRewardValue(*req.RewardValue)
	}
	if req.RewardValueMax != nil {
		builder.SetRewardValueMax(*req.RewardValueMax)
	}
	if req.SubscriptionID != nil {
		builder.SetNillableSubscriptionID(*req.SubscriptionID)
	}
	if req.SubscriptionDays != nil {
		builder.SetSubscriptionDays(*req.SubscriptionDays)
	}
	if req.Weight != nil {
		builder.SetWeight(*req.Weight)
	}
	if req.IsEnabled != nil {
		builder.SetIsEnabled(*req.IsEnabled)
	}

	item, err := builder.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("update prize item: %w", err)
	}

	result := prizeItemFromEnt(item)
	return &result, nil
}

func (s *BlindBoxService) DeletePrizeItem(ctx context.Context, id int64) error {
	return s.entClient.CheckinPrizeItem.UpdateOneID(id).
		SetDeletedAt(time.Now()).
		Exec(ctx)
}

func (s *BlindBoxService) ShouldTriggerBlindbox(ctx context.Context, userID int64, streakDays int) bool {
	if !s.settingSvc.IsCheckinBlindboxEnabled(ctx) {
		return false
	}

	triggerType := s.settingSvc.GetCheckinBlindboxTriggerType(ctx)
	interval := s.settingSvc.GetCheckinBlindboxInterval(ctx)
	if interval <= 0 {
		return false
	}

	if triggerType == "total" {
		var totalCheckins int
		s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM checkins WHERE user_id = $1`, userID).Scan(&totalCheckins)
		return totalCheckins > 0 && totalCheckins%interval == 0
	}

	return streakDays > 0 && streakDays%interval == 0
}

func (s *BlindBoxService) Draw(ctx context.Context, userID int64, streakDays int) (*BlindboxResult, error) {
	items, err := s.entClient.CheckinPrizeItem.Query().
		Where(
			checkinprizeitem.IsEnabled(true),
			checkinprizeitem.DeletedAtIsNil(),
		).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query prize items: %w", err)
	}
	if len(items) == 0 {
		return nil, nil
	}

	totalWeight := 0
	for _, item := range items {
		totalWeight += item.Weight
	}

	roll := rand.IntN(totalWeight)
	cumWeight := 0
	var selected *dbent.CheckinPrizeItem
	for _, item := range items {
		cumWeight += item.Weight
		if roll < cumWeight {
			selected = item
			break
		}
	}
	if selected == nil {
		selected = items[0]
	}

	rewardValue := selected.RewardValue
	if selected.RewardType == BlindboxRewardBalance && selected.RewardValueMax > selected.RewardValue {
		rewardValue = selected.RewardValue + rand.Float64()*(selected.RewardValueMax-selected.RewardValue)
		rewardValue = math.Round(rewardValue*100) / 100
	}

	if err := s.applyReward(ctx, userID, selected, rewardValue); err != nil {
		return nil, fmt.Errorf("apply reward: %w", err)
	}

	_, err = s.entClient.CheckinBlindboxRecord.Create().
		SetUserID(userID).
		SetPrizeItemID(selected.ID).
		SetPrizeName(selected.Name).
		SetRarity(selected.Rarity).
		SetRewardType(selected.RewardType).
		SetRewardValue(rewardValue).
		SetStreakDays(streakDays).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("save blindbox record: %w", err)
	}

	return &BlindboxResult{
		PrizeName:        selected.Name,
		Rarity:           selected.Rarity,
		RewardType:       selected.RewardType,
		RewardValue:      rewardValue,
		SubscriptionDays: selected.SubscriptionDays,
	}, nil
}

func (s *BlindBoxService) applyReward(ctx context.Context, userID int64, item *dbent.CheckinPrizeItem, value float64) error {
	switch item.RewardType {
	case BlindboxRewardBalance:
		if value > 0 {
			if err := s.userRepo.UpdateBalance(ctx, userID, value); err != nil {
				return fmt.Errorf("update balance: %w", err)
			}
		}
	case BlindboxRewardConcurrency:
		_, err := s.entClient.User.UpdateOneID(userID).
			AddConcurrency(int(value)).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("update concurrency: %w", err)
		}
	case BlindboxRewardSubscription:
		if item.SubscriptionID != nil && item.SubscriptionDays > 0 {
			_, _, err := s.subscriptionSvc.AssignOrExtendSubscription(ctx, &AssignSubscriptionInput{
				UserID:       userID,
				GroupID:      *item.SubscriptionID,
				ValidityDays: item.SubscriptionDays,
				Notes:        "check-in blind box reward",
			})
			if err != nil {
				return fmt.Errorf("assign subscription: %w", err)
			}
		}
	case BlindboxRewardInvitationCode:
		// invitation code generation is informational only - record is stored
	}
	if s.billingCache != nil {
		go func() {
			cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = s.billingCache.InvalidateUserBalance(cacheCtx, userID)
		}()
	}
	return nil
}

func (s *BlindBoxService) GetUserRecords(ctx context.Context, userID int64, page, pageSize int) (*BlindboxRecordList, error) {
	offset := (page - 1) * pageSize

	total, err := s.entClient.CheckinBlindboxRecord.Query().
		Where(checkinblindboxrecord.UserID(userID)).
		Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("count records: %w", err)
	}

	records, err := s.entClient.CheckinBlindboxRecord.Query().
		Where(checkinblindboxrecord.UserID(userID)).
		Order(dbent.Desc(checkinblindboxrecord.FieldCreatedAt)).
		Offset(offset).
		Limit(pageSize).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query records: %w", err)
	}

	items := make([]BlindboxRecord, 0, len(records))
	for _, r := range records {
		items = append(items, BlindboxRecord{
			ID:          r.ID,
			PrizeName:   r.PrizeName,
			Rarity:      r.Rarity,
			RewardType:  r.RewardType,
			RewardValue: r.RewardValue,
			StreakDays:  r.StreakDays,
			CreatedAt:   r.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &BlindboxRecordList{Items: items, Total: int64(total)}, nil
}

func (s *BlindBoxService) GetStats(ctx context.Context) (map[string]interface{}, error) {
	totalItems, err := s.entClient.CheckinPrizeItem.Query().
		Where(checkinprizeitem.DeletedAtIsNil()).
		Count(ctx)
	if err != nil {
		return nil, err
	}

	enabledItems, err := s.entClient.CheckinPrizeItem.Query().
		Where(checkinprizeitem.DeletedAtIsNil(), checkinprizeitem.IsEnabled(true)).
		Count(ctx)
	if err != nil {
		return nil, err
	}

	totalDraws, err := s.entClient.CheckinBlindboxRecord.Query().Count(ctx)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_items":  totalItems,
		"enabled_items": enabledItems,
		"total_draws":  totalDraws,
	}, nil
}

func prizeItemFromEnt(item *dbent.CheckinPrizeItem) PrizeItem {
	return PrizeItem{
		ID:               item.ID,
		Name:             item.Name,
		Rarity:           item.Rarity,
		RewardType:       item.RewardType,
		RewardValue:      item.RewardValue,
		RewardValueMax:   item.RewardValueMax,
		SubscriptionID:   item.SubscriptionID,
		SubscriptionDays: item.SubscriptionDays,
		Weight:           item.Weight,
		IsEnabled:        item.IsEnabled,
		CreatedAt:        item.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        item.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
