package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"sync"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

var (
	ErrTransferDisabled       = infraerrors.Forbidden("TRANSFER_DISABLED", "transfer feature is disabled")
	ErrTransferSelf           = infraerrors.BadRequest("TRANSFER_SELF", "cannot transfer to yourself")
	ErrTransferAmountInvalid  = infraerrors.BadRequest("TRANSFER_AMOUNT_INVALID", "invalid transfer amount")
	ErrTransferInsufficient   = infraerrors.BadRequest("TRANSFER_INSUFFICIENT", "insufficient balance")
	ErrTransferDailyLimit     = infraerrors.Forbidden("TRANSFER_DAILY_LIMIT", "daily transfer limit exceeded")
	ErrTransferDailyCount     = infraerrors.Forbidden("TRANSFER_DAILY_COUNT", "daily transfer count limit exceeded")
	ErrTransferReceiverNotFound = infraerrors.NotFound("RECEIVER_NOT_FOUND", "receiver not found")
	ErrTransferNotFound       = infraerrors.NotFound("TRANSFER_NOT_FOUND", "transfer not found")
	ErrTransferAlreadyFrozen  = infraerrors.BadRequest("TRANSFER_ALREADY_FROZEN", "transfer already frozen")
	ErrTransferAlreadyRevoked = infraerrors.BadRequest("TRANSFER_ALREADY_REVOKED", "transfer already revoked")
	ErrRedPacketDisabled      = infraerrors.Forbidden("REDPACKET_DISABLED", "red packet feature is disabled")
	ErrRedPacketNotFound      = infraerrors.NotFound("REDPACKET_NOT_FOUND", "red packet not found")
	ErrRedPacketExpired       = infraerrors.BadRequest("REDPACKET_EXPIRED", "red packet has expired")
	ErrRedPacketExhausted     = infraerrors.BadRequest("REDPACKET_EXHAUSTED", "red packet has been fully claimed")
	ErrRedPacketAlreadyClaimed = infraerrors.BadRequest("REDPACKET_ALREADY_CLAIMED", "you have already claimed this red packet")
	ErrRedPacketSelfClaim     = infraerrors.BadRequest("REDPACKET_SELF_CLAIM", "cannot claim your own red packet")
	ErrRedPacketCountInvalid  = infraerrors.BadRequest("REDPACKET_COUNT_INVALID", "invalid red packet count")
)

type BalanceTransferService struct {
	transferRepo   BalanceTransferRepository
	redPacketRepo  BalanceRedPacketRepository
	userRepo       UserRepository
	settingService *SettingService
	claimLocks     sync.Map
}

func NewBalanceTransferService(
	transferRepo BalanceTransferRepository,
	redPacketRepo BalanceRedPacketRepository,
	userRepo UserRepository,
	settingService *SettingService,
) *BalanceTransferService {
	return &BalanceTransferService{
		transferRepo:   transferRepo,
		redPacketRepo:  redPacketRepo,
		userRepo:       userRepo,
		settingService: settingService,
	}
}

func (s *BalanceTransferService) getTransferSettings(ctx context.Context) *TransferSettings {
	settings, err := s.settingService.GetAllSettings(ctx)
	if err != nil {
		return &TransferSettings{}
	}
	return &TransferSettings{
		Enabled:              settings.TransferEnabled,
		FeeRate:              settings.TransferFeeRate,
		MinAmount:            settings.TransferMinAmount,
		MaxAmount:            settings.TransferMaxAmount,
		DailyLimit:           settings.TransferDailyLimit,
		DailyCountLimit:      settings.TransferDailyCountLimit,
		VIPFeeExempt:         settings.TransferVIPFeeExempt,
		RedPacketEnabled:     settings.RedPacketEnabled,
		RedPacketMaxCount:    settings.RedPacketMaxCount,
		RedPacketExpireHours: settings.RedPacketExpireHours,
	}
}

func (s *BalanceTransferService) Transfer(ctx context.Context, senderID, receiverID int64, amount float64, memo *string) (*BalanceTransferRecord, error) {
	cfg := s.getTransferSettings(ctx)
	if !cfg.Enabled {
		return nil, ErrTransferDisabled
	}
	if senderID == receiverID {
		return nil, ErrTransferSelf
	}
	amount = math.Round(amount*1e8) / 1e8
	if amount < cfg.MinAmount || (cfg.MaxAmount > 0 && amount > cfg.MaxAmount) || amount <= 0 {
		return nil, ErrTransferAmountInvalid
	}
	dailyTotal, dailyCount, err := s.transferRepo.GetDailyTransferTotal(ctx, senderID)
	if err != nil {
		return nil, fmt.Errorf("check daily limit: %w", err)
	}
	if cfg.DailyLimit > 0 && dailyTotal+amount > cfg.DailyLimit {
		return nil, ErrTransferDailyLimit
	}
	if cfg.DailyCountLimit > 0 && dailyCount >= cfg.DailyCountLimit {
		return nil, ErrTransferDailyCount
	}
	receiver, err := s.userRepo.GetByID(ctx, receiverID)
	if err != nil {
		return nil, ErrTransferReceiverNotFound
	}
	if receiver == nil {
		return nil, ErrTransferReceiverNotFound
	}
	feeRate := cfg.FeeRate
	fee := math.Round(amount*feeRate*1e8) / 1e8
	if fee < 0 {
		fee = 0
	}
	grossAmount := amount + fee
	sender, err := s.userRepo.GetByID(ctx, senderID)
	if err != nil {
		return nil, fmt.Errorf("get sender: %w", err)
	}
	if sender.Balance < grossAmount {
		return nil, ErrTransferInsufficient
	}
	var record *BalanceTransferRecord
	if err := s.transferRepo.RunInTx(ctx, func(txCtx context.Context) error {
		if err := s.userRepo.DeductBalance(txCtx, senderID, grossAmount); err != nil {
			return fmt.Errorf("deduct sender balance: %w", err)
		}
		if err := s.userRepo.UpdateBalance(txCtx, receiverID, amount); err != nil {
			return fmt.Errorf("credit receiver balance: %w", err)
		}
		record = &BalanceTransferRecord{
			SenderID:     senderID,
			ReceiverID:   receiverID,
			Amount:       amount,
			Fee:          fee,
			FeeRate:      feeRate,
			GrossAmount:  grossAmount,
			TransferType: "direct",
			Status:       "completed",
			Memo:         memo,
			CreatedAt:    time.Now(),
		}
		return s.transferRepo.Create(txCtx, record)
	}); err != nil {
		return nil, err
	}
	return record, nil
}

func (s *BalanceTransferService) ValidateTransfer(ctx context.Context, senderID, receiverID int64, amount float64) (fee float64, feeRate float64, err error) {
	cfg := s.getTransferSettings(ctx)
	if !cfg.Enabled {
		return 0, 0, ErrTransferDisabled
	}
	if senderID == receiverID {
		return 0, 0, ErrTransferSelf
	}
	amount = math.Round(amount*1e8) / 1e8
	if amount < cfg.MinAmount || (cfg.MaxAmount > 0 && amount > cfg.MaxAmount) || amount <= 0 {
		return 0, 0, ErrTransferAmountInvalid
	}
	feeRate = cfg.FeeRate
	fee = math.Round(amount*feeRate*1e8) / 1e8
	return fee, feeRate, nil
}

func (s *BalanceTransferService) GetHistory(ctx context.Context, userID int64, role string, page, pageSize int) ([]*BalanceTransferRecord, int, error) {
	return s.transferRepo.ListByUserExcludeType(ctx, userID, role, "redpacket", page, pageSize)
}

func (s *BalanceTransferService) GetAllTransfers(ctx context.Context, filter *TransferFilter, page, pageSize int) ([]*BalanceTransferRecord, int, error) {
	records, total, err := s.transferRepo.ListAll(ctx, filter, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	userIDs := make(map[int64]struct{})
	for _, r := range records {
		userIDs[r.SenderID] = struct{}{}
		userIDs[r.ReceiverID] = struct{}{}
	}
	emails := make(map[int64]string)
	for uid := range userIDs {
		if u, err := s.userRepo.GetByID(ctx, uid); err == nil {
			emails[uid] = u.Email
		}
	}
	for _, r := range records {
		r.SenderEmail = emails[r.SenderID]
		r.ReceiverEmail = emails[r.ReceiverID]
	}
	return records, total, nil
}

func (s *BalanceTransferService) FreezeTransfer(ctx context.Context, adminID, transferID int64) error {
	record, err := s.transferRepo.GetByID(ctx, transferID)
	if err != nil {
		return ErrTransferNotFound
	}
	if record.Status == "frozen" {
		return ErrTransferAlreadyFrozen
	}
	if record.Status == "revoked" {
		return ErrTransferAlreadyRevoked
	}
	now := time.Now()
	return s.transferRepo.UpdateStatus(ctx, transferID, "frozen", &now, &adminID, nil)
}

func (s *BalanceTransferService) RevokeTransfer(ctx context.Context, adminID, transferID int64, reason string) error {
	record, err := s.transferRepo.GetByID(ctx, transferID)
	if err != nil {
		return ErrTransferNotFound
	}
	if record.Status == "revoked" {
		return ErrTransferAlreadyRevoked
	}
	return s.transferRepo.RunInTx(ctx, func(txCtx context.Context) error {
		if err := s.userRepo.DeductBalance(txCtx, record.ReceiverID, record.Amount); err != nil {
			return fmt.Errorf("deduct receiver balance: %w", err)
		}
		if err := s.userRepo.UpdateBalance(txCtx, record.SenderID, record.GrossAmount); err != nil {
			return fmt.Errorf("return sender balance: %w", err)
		}
		return s.transferRepo.UpdateStatus(txCtx, transferID, "revoked", record.FrozenAt, &adminID, &reason)
	})
}

func (s *BalanceTransferService) BatchDistribute(ctx context.Context, adminID int64, targets []BatchDistributeTarget, memo *string) ([]*BalanceTransferRecord, error) {
	var records []*BalanceTransferRecord
	err := s.transferRepo.RunInTx(ctx, func(txCtx context.Context) error {
		for _, t := range targets {
			if t.Amount <= 0 || t.UserID <= 0 {
				continue
			}
			if _, err := s.userRepo.GetByID(txCtx, t.UserID); err != nil {
				continue
			}
			if err := s.userRepo.UpdateBalance(txCtx, t.UserID, t.Amount); err != nil {
				return fmt.Errorf("update balance for user %d: %w", t.UserID, err)
			}
			record := &BalanceTransferRecord{
				SenderID:     adminID,
				ReceiverID:   t.UserID,
				Amount:       t.Amount,
				Fee:          0,
				FeeRate:      0,
				GrossAmount:  t.Amount,
				TransferType: "batch",
				Status:       "completed",
				Memo:         memo,
				CreatedAt:    time.Now(),
			}
			if err := s.transferRepo.Create(txCtx, record); err != nil {
				return fmt.Errorf("create transfer record for user %d: %w", t.UserID, err)
			}
			records = append(records, record)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (s *BalanceTransferService) GetFeeStats(ctx context.Context, startTime, endTime time.Time) ([]*DailyFeeStat, error) {
	return s.transferRepo.GetFeeStats(ctx, startTime, endTime)
}

func (s *BalanceTransferService) GetLeaderboard(ctx context.Context, period string, limit int) ([]*TransferRankEntry, error) {
	now := time.Now()
	var start time.Time
	switch period {
	case "week":
		start = now.AddDate(0, 0, -7)
	case "month":
		start = now.AddDate(0, -1, 0)
	default:
		start = now.AddDate(0, 0, -1)
	}
	return s.transferRepo.GetLeaderboard(ctx, start, now, limit, "amount")
}

type BatchDistributeTarget struct {
	UserID int64   `json:"user_id"`
	Amount float64 `json:"amount"`
}

func (s *BalanceTransferService) CreateRedPacket(ctx context.Context, senderID int64, totalAmount float64, count int, redPacketType string, memo *string) (*RedPacketRecord, error) {
	cfg := s.getTransferSettings(ctx)
	if !cfg.Enabled || !cfg.RedPacketEnabled {
		return nil, ErrRedPacketDisabled
	}
	if count <= 0 || count > cfg.RedPacketMaxCount {
		return nil, ErrRedPacketCountInvalid
	}
	totalAmount = math.Round(totalAmount*1e8) / 1e8
	minRequired := float64(count) * 0.01
	if totalAmount < minRequired {
		return nil, infraerrors.BadRequest("REDPACKET_AMOUNT_TOO_SMALL", fmt.Sprintf("minimum amount for %d packets is %.2f", count, minRequired))
	}
	feeRate := cfg.FeeRate
	fee := math.Round(totalAmount*feeRate*1e8) / 1e8
	grossAmount := totalAmount + fee
	sender, err := s.userRepo.GetByID(ctx, senderID)
	if err != nil {
		return nil, fmt.Errorf("get sender: %w", err)
	}
	if sender.Balance < grossAmount {
		return nil, ErrTransferInsufficient
	}
	code, err := generateRedPacketCode()
	if err != nil {
		return nil, fmt.Errorf("generate code: %w", err)
	}
	expireHours := cfg.RedPacketExpireHours
	if expireHours <= 0 {
		expireHours = 24
	}
	var rp *RedPacketRecord
	if err := s.transferRepo.RunInTx(ctx, func(txCtx context.Context) error {
		if err := s.userRepo.DeductBalance(txCtx, senderID, grossAmount); err != nil {
			return fmt.Errorf("deduct sender balance: %w", err)
		}
		rp = &RedPacketRecord{
			SenderID:        senderID,
			TotalAmount:     totalAmount,
			TotalCount:      count,
			RemainingAmount: totalAmount,
			RemainingCount:  count,
			RedPacketType:   redPacketType,
			Fee:             fee,
			FeeRate:         feeRate,
			Code:            code,
			Status:          "active",
			Memo:            memo,
			ExpireAt:        time.Now().Add(time.Duration(expireHours) * time.Hour),
			CreatedAt:       time.Now(),
		}
		return s.redPacketRepo.Create(txCtx, rp)
	}); err != nil {
		return nil, err
	}
	return rp, nil
}

func (s *BalanceTransferService) ClaimRedPacket(ctx context.Context, userID int64, code string) (*RedPacketClaimRecord, error) {
	cfg := s.getTransferSettings(ctx)
	if !cfg.Enabled || !cfg.RedPacketEnabled {
		return nil, ErrRedPacketDisabled
	}
	rp, err := s.redPacketRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, ErrRedPacketNotFound
	}
	if rp.SenderID == userID {
		return nil, ErrRedPacketSelfClaim
	}
	if rp.Status != "active" {
		if rp.Status == "expired" {
			return nil, ErrRedPacketExpired
		}
		return nil, ErrRedPacketExhausted
	}
	if time.Now().After(rp.ExpireAt) {
		return nil, ErrRedPacketExpired
	}
	claimed, err := s.redPacketRepo.HasClaimed(ctx, rp.ID, userID)
	if err != nil {
		return nil, fmt.Errorf("check claimed: %w", err)
	}
	if claimed {
		return nil, ErrRedPacketAlreadyClaimed
	}

	lockKey := fmt.Sprintf("rp:%d", rp.ID)
	actual, _ := s.claimLocks.LoadOrStore(lockKey, &sync.Mutex{})
	mu := actual.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()

	freshRp, err := s.redPacketRepo.GetByID(ctx, rp.ID)
	if err != nil {
		return nil, ErrRedPacketNotFound
	}
	if freshRp.Status != "active" || freshRp.RemainingCount <= 0 || freshRp.RemainingAmount <= 0 {
		return nil, ErrRedPacketExhausted
	}

	amount := s.calculateClaimAmount(freshRp)
	if amount <= 0 {
		return nil, ErrRedPacketExhausted
	}
	remainingCount := freshRp.RemainingCount
	var claimRecord *RedPacketClaimRecord
	if err := s.transferRepo.RunInTx(ctx, func(txCtx context.Context) error {
		if err := s.redPacketRepo.DecrementClaim(txCtx, freshRp.ID, amount); err != nil {
			return ErrRedPacketExhausted
		}
		if err := s.userRepo.UpdateBalance(txCtx, userID, amount); err != nil {
			return fmt.Errorf("credit balance: %w", err)
		}
		claimRecord = &RedPacketClaimRecord{
			RedPacketID: freshRp.ID,
			UserID:      userID,
			Amount:      amount,
			CreatedAt:   time.Now(),
		}
		transferRecord := &BalanceTransferRecord{
			SenderID:     freshRp.SenderID,
			ReceiverID:   userID,
			Amount:       amount,
			Fee:          0,
			FeeRate:      0,
			GrossAmount:  amount,
			TransferType: "redpacket",
			Status:       "completed",
			RedpacketID:  &freshRp.ID,
			CreatedAt:    time.Now(),
		}
		if err := s.transferRepo.Create(txCtx, transferRecord); err != nil {
			return fmt.Errorf("create transfer record: %w", err)
		}
		claimRecord.TransferID = &transferRecord.ID
		if err := s.redPacketRepo.CreateClaim(txCtx, claimRecord); err != nil {
			return fmt.Errorf("create claim record: %w", err)
		}
		if remainingCount <= 1 {
			return s.redPacketRepo.MarkExhausted(txCtx, freshRp.ID)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return claimRecord, nil
}

func (s *BalanceTransferService) GetRedPacketDetail(ctx context.Context, redPacketID int64) (*RedPacketRecord, []*RedPacketClaimRecord, error) {
	rp, err := s.redPacketRepo.GetByID(ctx, redPacketID)
	if err != nil {
		return nil, nil, ErrRedPacketNotFound
	}
	claims, err := s.redPacketRepo.GetClaims(ctx, redPacketID)
	if err != nil {
		claims = []*RedPacketClaimRecord{}
	}
	for _, c := range claims {
		if u, err := s.userRepo.GetByID(ctx, c.UserID); err == nil {
			c.UserEmail = u.Email
		}
	}
	return rp, claims, nil
}

func (s *BalanceTransferService) GetMyRedPackets(ctx context.Context, senderID int64, page, pageSize int) ([]*RedPacketRecord, int, error) {
	return s.redPacketRepo.ListBySender(ctx, senderID, page, pageSize)
}

func (s *BalanceTransferService) ExpireRedPackets(ctx context.Context) error {
	rps, err := s.redPacketRepo.ListActiveExpired(ctx)
	if err != nil {
		return err
	}
	for _, rp := range rps {
		_ = s.transferRepo.RunInTx(ctx, func(txCtx context.Context) error {
			remaining, err := s.redPacketRepo.ReturnRemaining(txCtx, rp.ID, rp.SenderID)
			if err != nil {
				return err
			}
			if remaining > 0 {
				return s.userRepo.UpdateBalance(txCtx, rp.SenderID, remaining)
			}
			return nil
		})
	}
	return nil
}

func (s *BalanceTransferService) GetAllRedPackets(ctx context.Context, page, pageSize int) ([]*RedPacketRecord, int, error) {
	return s.redPacketRepo.ListAll(ctx, page, pageSize)
}

func (s *BalanceTransferService) calculateClaimAmount(rp *RedPacketRecord) float64 {
	if rp.RemainingCount <= 0 || rp.RemainingAmount <= 0 {
		return 0
	}
	if rp.RedPacketType == "equal" {
		return math.Round(rp.RemainingAmount/float64(rp.RemainingCount)*1e8) / 1e8
	}
	if rp.RemainingCount == 1 {
		return math.Round(rp.RemainingAmount*1e8) / 1e8
	}
	maxAllowed := rp.RemainingAmount - float64(rp.RemainingCount-1)*0.01
	upperBound := maxAllowed / float64(rp.RemainingCount) * 2
	if upperBound <= 0.01 {
		return 0.01
	}
	n, err := rand.Int(rand.Reader, big.NewInt(int64(upperBound*100)))
	if err != nil {
		return math.Round(maxAllowed/float64(rp.RemainingCount)*1e8) / 1e8
	}
	amount := float64(n.Int64())/100 + 0.01
	if amount < 0.01 {
		amount = 0.01
	}
	if amount > maxAllowed {
		amount = maxAllowed
	}
	return math.Round(amount*1e8) / 1e8
}

func generateRedPacketCode() (string, error) {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (s *BalanceTransferService) GetTransferStats(ctx context.Context, userID int64) (sent float64, received float64, feePaid float64, err error) {
	return s.transferRepo.GetUserTransferStats(ctx, userID)
}

type UserSearchResult struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func (s *BalanceTransferService) SearchUsers(ctx context.Context, query string) ([]*UserSearchResult, error) {
	if query == "" {
		return nil, nil
	}
	users, _, err := s.userRepo.ListWithFilters(ctx, pagination.PaginationParams{Page: 1, PageSize: 10}, UserListFilters{Search: query})
	if err != nil {
		return nil, err
	}
	var results []*UserSearchResult
	for _, u := range users {
		results = append(results, &UserSearchResult{ID: u.ID, Email: u.Email, Username: u.Username})
	}
	return results, nil
}
