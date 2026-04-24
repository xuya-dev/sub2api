package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

type BalanceTransferHandler struct {
	transferService *service.BalanceTransferService
}

func NewBalanceTransferHandler(transferService *service.BalanceTransferService) *BalanceTransferHandler {
	return &BalanceTransferHandler{transferService: transferService}
}

func getUserID(c *gin.Context) int64 {
	subject, ok := middleware.GetAuthSubjectFromContext(c)
	if !ok {
		return 0
	}
	return subject.UserID
}

func (h *BalanceTransferHandler) Transfer(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req struct {
		ReceiverID int64   `json:"receiver_id" binding:"required"`
		Amount     float64 `json:"amount" binding:"required,gt=0"`
		Memo       *string `json:"memo"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	record, err := h.transferService.Transfer(c.Request.Context(), userID, req.ReceiverID, req.Amount, req.Memo)
	if err != nil {
		WriteAppError(c, err)
		return
	}
	c.JSON(http.StatusOK, record)
}

func (h *BalanceTransferHandler) ValidateTransfer(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req struct {
		ReceiverID int64   `json:"receiver_id" binding:"required"`
		Amount     float64 `json:"amount" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fee, feeRate, err := h.transferService.ValidateTransfer(c.Request.Context(), userID, req.ReceiverID, req.Amount)
	if err != nil {
		WriteAppError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"fee": fee, "fee_rate": feeRate})
}

func (h *BalanceTransferHandler) GetHistory(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	role := c.DefaultQuery("role", "all")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	records, total, err := h.transferService.GetHistory(c.Request.Context(), userID, role, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": records, "total": total, "page": page, "page_size": pageSize})
}

func (h *BalanceTransferHandler) GetStats(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	sent, received, feePaid, err := h.transferService.GetTransferStats(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total_sent": sent, "total_received": received, "total_fee_paid": feePaid})
}

func (h *BalanceTransferHandler) CreateRedPacket(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req struct {
		TotalAmount   float64 `json:"total_amount" binding:"required,gt=0"`
		Count         int     `json:"count" binding:"required,gt=0"`
		RedPacketType string  `json:"redpacket_type"`
		Memo          *string `json:"memo"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.RedPacketType == "" {
		req.RedPacketType = "equal"
	}
	rp, err := h.transferService.CreateRedPacket(c.Request.Context(), userID, req.TotalAmount, req.Count, req.RedPacketType, req.Memo)
	if err != nil {
		WriteAppError(c, err)
		return
	}
	c.JSON(http.StatusOK, rp)
}

func (h *BalanceTransferHandler) ClaimRedPacket(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	claim, err := h.transferService.ClaimRedPacket(c.Request.Context(), userID, req.Code)
	if err != nil {
		WriteAppError(c, err)
		return
	}
	c.JSON(http.StatusOK, claim)
}

func (h *BalanceTransferHandler) GetRedPacketDetail(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	rp, claims, err := h.transferService.GetRedPacketDetail(c.Request.Context(), id)
	if err != nil {
		WriteAppError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"redpacket": rp, "claims": claims})
}

func (h *BalanceTransferHandler) GetMyRedPackets(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	records, total, err := h.transferService.GetMyRedPackets(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": records, "total": total, "page": page, "page_size": pageSize})
}

func (h *BalanceTransferHandler) GetLeaderboard(c *gin.Context) {
	period := c.DefaultQuery("period", "day")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	entries, err := h.transferService.GetLeaderboard(c.Request.Context(), period, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entries)
}

func (h *BalanceTransferHandler) SearchUsers(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	q := c.Query("q")
	results, err := h.transferService.SearchUsers(c.Request.Context(), q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if results == nil {
		results = []*service.UserSearchResult{}
	}
	c.JSON(http.StatusOK, results)
}

func GetUserIDAware(c *gin.Context) int64 {
	subject, ok := middleware.GetAuthSubjectFromContext(c)
	if !ok {
		return 0
	}
	return subject.UserID
}

func WriteAppError(c *gin.Context, err error) {
	var appErr *infraerrors.ApplicationError
	if errors.As(err, &appErr) {
		c.JSON(int(appErr.Code), gin.H{"error": appErr.Message, "code": appErr.Reason})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
}
