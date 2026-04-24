package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/balancetransfer"
	"github.com/Wei-Shaw/sub2api/ent/predicate"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

type balanceTransferRepo struct {
	client *dbent.Client
	db     *sql.DB
}

func NewBalanceTransferRepository(client *dbent.Client, db *sql.DB) service.BalanceTransferRepository {
	return &balanceTransferRepo{client: client, db: db}
}

func (r *balanceTransferRepo) Create(ctx context.Context, t *service.BalanceTransferRecord) error {
	client := clientFromContext(ctx, r.client)
	builder := client.BalanceTransfer.Create().
		SetSenderID(t.SenderID).
		SetReceiverID(t.ReceiverID).
		SetAmount(t.Amount).
		SetFee(t.Fee).
		SetFeeRate(t.FeeRate).
		SetGrossAmount(t.GrossAmount).
		SetTransferType(t.TransferType).
		SetStatus(t.Status).
		SetCreatedAt(time.Now())
	if t.Memo != nil {
		builder.SetMemo(*t.Memo)
	}
	if t.RedpacketID != nil {
		builder.SetRedpacketID(*t.RedpacketID)
	}
	saved, err := builder.Save(ctx)
	if err != nil {
		return fmt.Errorf("create balance transfer: %w", err)
	}
	t.ID = saved.ID
	return nil
}

func (r *balanceTransferRepo) GetByID(ctx context.Context, id int64) (*service.BalanceTransferRecord, error) {
	client := clientFromContext(ctx, r.client)
	t, err := client.BalanceTransfer.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get balance transfer %d: %w", id, err)
	}
	return toTransferRecord(t), nil
}

func (r *balanceTransferRepo) UpdateStatus(ctx context.Context, id int64, status string, frozenAt *time.Time, frozenBy *int64, revokeReason *string) error {
	client := clientFromContext(ctx, r.client)
	builder := client.BalanceTransfer.UpdateOneID(id).
		SetStatus(status)
	if frozenAt != nil {
		builder.SetFrozenAt(*frozenAt)
	}
	if frozenBy != nil {
		builder.SetFrozenBy(*frozenBy)
	}
	if revokeReason != nil {
		builder.SetRevokeReason(*revokeReason)
	}
	_, err := builder.Save(ctx)
	return err
}

func (r *balanceTransferRepo) ListByUser(ctx context.Context, userID int64, role string, page, pageSize int) ([]*service.BalanceTransferRecord, int, error) {
	client := clientFromContext(ctx, r.client)
	var preds []predicate.BalanceTransfer
	switch role {
	case "sender":
		preds = append(preds, balancetransfer.SenderID(userID))
	case "receiver":
		preds = append(preds, balancetransfer.ReceiverID(userID))
	default:
		preds = append(preds, balancetransfer.Or(
			balancetransfer.SenderID(userID),
			balancetransfer.ReceiverID(userID),
		))
	}
	query := client.BalanceTransfer.Query().Where(preds...).Order(dbent.Desc(balancetransfer.FieldCreatedAt))
	return r.queryWithPagination(ctx, query, page, pageSize)
}

func (r *balanceTransferRepo) ListByUserExcludeType(ctx context.Context, userID int64, role, excludeType string, page, pageSize int) ([]*service.BalanceTransferRecord, int, error) {
	client := clientFromContext(ctx, r.client)
	var preds []predicate.BalanceTransfer
	switch role {
	case "sender":
		preds = append(preds, balancetransfer.SenderID(userID))
	case "receiver":
		preds = append(preds, balancetransfer.ReceiverID(userID))
	default:
		preds = append(preds, balancetransfer.Or(
			balancetransfer.SenderID(userID),
			balancetransfer.ReceiverID(userID),
		))
	}
	if excludeType != "" {
		preds = append(preds, balancetransfer.TransferTypeNEQ(excludeType))
	}
	query := client.BalanceTransfer.Query().Where(preds...).Order(dbent.Desc(balancetransfer.FieldCreatedAt))
	return r.queryWithPagination(ctx, query, page, pageSize)
}

func (r *balanceTransferRepo) ListAll(ctx context.Context, filter *service.TransferFilter, page, pageSize int) ([]*service.BalanceTransferRecord, int, error) {
	client := clientFromContext(ctx, r.client)
	var preds []predicate.BalanceTransfer
	if filter != nil {
		if filter.Status != "" {
			preds = append(preds, balancetransfer.StatusEQ(filter.Status))
		}
		if filter.TransferType != "" {
			preds = append(preds, balancetransfer.TransferTypeEQ(filter.TransferType))
		}
		if filter.UserID != nil {
			preds = append(preds, balancetransfer.Or(
				balancetransfer.SenderID(*filter.UserID),
				balancetransfer.ReceiverID(*filter.UserID),
			))
		}
		if !filter.StartTime.IsZero() {
			preds = append(preds, balancetransfer.CreatedAtGTE(filter.StartTime))
		}
		if !filter.EndTime.IsZero() {
			preds = append(preds, balancetransfer.CreatedAtLTE(filter.EndTime))
		}
	}
	query := client.BalanceTransfer.Query().Where(preds...).Order(dbent.Desc(balancetransfer.FieldCreatedAt))
	return r.queryWithPagination(ctx, query, page, pageSize)
}

func (r *balanceTransferRepo) GetDailyTransferTotal(ctx context.Context, userID int64) (float64, int, error) {
	startOfDay := time.Now().Truncate(24 * time.Hour)
	var total float64
	var count int
	err := r.db.QueryRowContext(ctx,
		"SELECT COALESCE(SUM(gross_amount),0), COALESCE(COUNT(*),0) FROM balance_transfers WHERE sender_id = $1 AND status != 'revoked' AND created_at >= $2",
		userID, startOfDay,
	).Scan(&total, &count)
	return total, count, err
}

func (r *balanceTransferRepo) GetFeeStats(ctx context.Context, startTime, endTime time.Time) ([]*service.DailyFeeStat, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT DATE(created_at) as day, COALESCE(SUM(fee),0) as total_fee, COUNT(*) as count FROM balance_transfers WHERE status = 'completed' AND created_at >= $1 AND created_at < $2 GROUP BY DATE(created_at) ORDER BY day",
		startTime, endTime,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var stats []*service.DailyFeeStat
	for rows.Next() {
		var s service.DailyFeeStat
		if err := rows.Scan(&s.Date, &s.TotalFee, &s.Count); err != nil {
			return nil, err
		}
		stats = append(stats, &s)
	}
	return stats, nil
}

func (r *balanceTransferRepo) GetLeaderboard(ctx context.Context, startTime, endTime time.Time, limit int, orderBy string) ([]*service.TransferRankEntry, error) {
	col := "SUM(bt.amount)"
	if orderBy == "count" {
		col = "COUNT(*)"
	}
	query := fmt.Sprintf(
		"SELECT u.id, u.email, COALESCE(SUM(bt.amount),0) as total_amount, COUNT(*) as total_count FROM balance_transfers bt JOIN users u ON u.id = bt.sender_id WHERE bt.status = 'completed' AND bt.created_at >= $1 AND bt.created_at < $2 GROUP BY u.id, u.email ORDER BY %s DESC LIMIT $3",
		col,
	)
	rows, err := r.db.QueryContext(ctx, query, startTime, endTime, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var entries []*service.TransferRankEntry
	rank := 1
	for rows.Next() {
		var e service.TransferRankEntry
		if err := rows.Scan(&e.UserID, &e.Email, &e.TotalAmount, &e.TotalCount); err != nil {
			return nil, err
		}
		e.Rank = rank
		rank++
		entries = append(entries, &e)
	}
	return entries, nil
}

func (r *balanceTransferRepo) queryWithPagination(ctx context.Context, query *dbent.BalanceTransferQuery, page, pageSize int) ([]*service.BalanceTransferRecord, int, error) {
	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	offset := (&pagination.PaginationParams{Page: page, PageSize: pageSize}).Offset()
	items, err := query.Offset(offset).Limit(pageSize).All(ctx)
	if err != nil {
		return nil, 0, err
	}
	records := make([]*service.BalanceTransferRecord, len(items))
	for i, item := range items {
		records[i] = toTransferRecord(item)
	}
	return records, total, nil
}

func (r *balanceTransferRepo) RunInTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	txCtx := dbent.NewTxContext(ctx, tx)
	defer func() { _ = tx.Rollback() }()
	if err := fn(txCtx); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *balanceTransferRepo) GetUserTransferStats(ctx context.Context, userID int64) (sent, received, feePaid float64, err error) {
	err = r.db.QueryRowContext(ctx,
		`SELECT
			COALESCE((SELECT SUM(amount) FROM balance_transfers WHERE sender_id = $1 AND status != 'revoked' AND transfer_type != 'redpacket'), 0),
			COALESCE((SELECT SUM(amount) FROM balance_transfers WHERE receiver_id = $1 AND status != 'revoked' AND transfer_type != 'redpacket'), 0),
			COALESCE((SELECT SUM(fee) FROM balance_transfers WHERE sender_id = $1 AND status != 'revoked' AND transfer_type != 'redpacket'), 0)`,
		userID,
	).Scan(&sent, &received, &feePaid)
	return
}

func toTransferRecord(t *dbent.BalanceTransfer) *service.BalanceTransferRecord {
	return &service.BalanceTransferRecord{
		ID:           t.ID,
		SenderID:     t.SenderID,
		ReceiverID:   t.ReceiverID,
		Amount:       t.Amount,
		Fee:          t.Fee,
		FeeRate:      t.FeeRate,
		GrossAmount:  t.GrossAmount,
		TransferType: t.TransferType,
		Status:       t.Status,
		Memo:         t.Memo,
		RedpacketID:  t.RedpacketID,
		FrozenAt:     t.FrozenAt,
		FrozenBy:     t.FrozenBy,
		RevokeReason: t.RevokeReason,
		CreatedAt:    t.CreatedAt,
	}
}
