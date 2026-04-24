package service

import (
	"context"
	"time"
)

type BalanceTransferRecord struct {
	ID             int64      `json:"id"`
	SenderID       int64      `json:"sender_id"`
	SenderEmail    string     `json:"sender_email"`
	ReceiverID     int64      `json:"receiver_id"`
	ReceiverEmail  string     `json:"receiver_email"`
	Amount         float64    `json:"amount"`
	Fee            float64    `json:"fee"`
	FeeRate        float64    `json:"fee_rate"`
	GrossAmount    float64    `json:"gross_amount"`
	TransferType   string     `json:"transfer_type"`
	Status         string     `json:"status"`
	Memo           *string    `json:"memo"`
	RedpacketID    *int64     `json:"redpacket_id"`
	FrozenAt       *time.Time `json:"frozen_at"`
	FrozenBy       *int64     `json:"frozen_by"`
	RevokeReason   *string    `json:"revoke_reason"`
	CreatedAt      time.Time  `json:"created_at"`
}

type RedPacketRecord struct {
	ID              int64      `json:"id"`
	SenderID        int64      `json:"sender_id"`
	TotalAmount     float64    `json:"total_amount"`
	TotalCount      int        `json:"total_count"`
	RemainingAmount float64    `json:"remaining_amount"`
	RemainingCount  int        `json:"remaining_count"`
	RedPacketType   string     `json:"redpacket_type"`
	Fee             float64    `json:"fee"`
	FeeRate         float64    `json:"fee_rate"`
	Code            string     `json:"code"`
	Status          string     `json:"status"`
	Memo            *string    `json:"memo"`
	ExpireAt        time.Time  `json:"expire_at"`
	CreatedAt       time.Time  `json:"created_at"`
}

type RedPacketClaimRecord struct {
	ID          int64     `json:"id"`
	RedPacketID int64     `json:"redpacket_id"`
	UserID      int64     `json:"user_id"`
	UserEmail   string    `json:"user_email"`
	Amount      float64   `json:"amount"`
	TransferID  *int64    `json:"transfer_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type TransferFilter struct {
	Status       string    `json:"status"`
	TransferType string    `json:"transfer_type"`
	UserID       *int64    `json:"user_id"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
}

type DailyFeeStat struct {
	Date     time.Time `json:"date"`
	TotalFee float64   `json:"total_fee"`
	Count    int       `json:"count"`
}

type TransferRankEntry struct {
	Rank        int     `json:"rank"`
	UserID      int64   `json:"user_id"`
	Email       string  `json:"email"`
	TotalAmount float64 `json:"total_amount"`
	TotalCount  int     `json:"total_count"`
}

type TransferSettings struct {
	Enabled                bool
	FeeRate                float64
	MinAmount              float64
	MaxAmount              float64
	DailyLimit             float64
	DailyCountLimit        int
	VIPFeeExempt           bool
	RedPacketEnabled       bool
	RedPacketMaxCount      int
	RedPacketExpireHours   int
}

type BalanceTransferRepository interface {
	Create(ctx context.Context, t *BalanceTransferRecord) error
	GetByID(ctx context.Context, id int64) (*BalanceTransferRecord, error)
	UpdateStatus(ctx context.Context, id int64, status string, frozenAt *time.Time, frozenBy *int64, revokeReason *string) error
	ListByUser(ctx context.Context, userID int64, role string, page, pageSize int) ([]*BalanceTransferRecord, int, error)
	ListByUserExcludeType(ctx context.Context, userID int64, role, excludeType string, page, pageSize int) ([]*BalanceTransferRecord, int, error)
	ListAll(ctx context.Context, filter *TransferFilter, page, pageSize int) ([]*BalanceTransferRecord, int, error)
	GetDailyTransferTotal(ctx context.Context, userID int64) (float64, int, error)
	GetFeeStats(ctx context.Context, startTime, endTime time.Time) ([]*DailyFeeStat, error)
	GetLeaderboard(ctx context.Context, startTime, endTime time.Time, limit int, orderBy string) ([]*TransferRankEntry, error)
	RunInTx(ctx context.Context, fn func(ctx context.Context) error) error
	GetUserTransferStats(ctx context.Context, userID int64) (sent, received, feePaid float64, err error)
}

type BalanceRedPacketRepository interface {
	Create(ctx context.Context, rp *RedPacketRecord) error
	GetByCode(ctx context.Context, code string) (*RedPacketRecord, error)
	GetByID(ctx context.Context, id int64) (*RedPacketRecord, error)
	DecrementClaim(ctx context.Context, id int64, amount float64) error
	MarkExhausted(ctx context.Context, id int64) error
	MarkExpired(ctx context.Context, id int64) error
	CreateClaim(ctx context.Context, claim *RedPacketClaimRecord) error
	HasClaimed(ctx context.Context, redpacketID, userID int64) (bool, error)
	GetClaims(ctx context.Context, redpacketID int64) ([]*RedPacketClaimRecord, error)
	ListBySender(ctx context.Context, senderID int64, page, pageSize int) ([]*RedPacketRecord, int, error)
	ListActiveExpired(ctx context.Context) ([]*RedPacketRecord, error)
	ListAll(ctx context.Context, page, pageSize int) ([]*RedPacketRecord, int, error)
	ReturnRemaining(ctx context.Context, id int64, senderID int64) (float64, error)
}
