package admin

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

type TransferAdminHandler struct {
	transferService *service.BalanceTransferService
}

func NewTransferAdminHandler(transferService *service.BalanceTransferService) *TransferAdminHandler {
	return &TransferAdminHandler{transferService: transferService}
}

func (h *TransferAdminHandler) ListTransfers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	filter := &service.TransferFilter{
		Status:       c.DefaultQuery("status", ""),
		TransferType: c.DefaultQuery("transfer_type", ""),
	}
	if uidStr := c.Query("user_id"); uidStr != "" {
		uid, _ := strconv.ParseInt(uidStr, 10, 64)
		if uid > 0 {
			filter.UserID = &uid
		}
	}
	if startStr := c.Query("start_time"); startStr != "" {
		if t, err := time.Parse(time.RFC3339, startStr); err == nil {
			filter.StartTime = t
		}
	}
	if endStr := c.Query("end_time"); endStr != "" {
		if t, err := time.Parse(time.RFC3339, endStr); err == nil {
			filter.EndTime = t
		}
	}
	records, total, err := h.transferService.GetAllTransfers(c.Request.Context(), filter, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": records, "total": total, "page": page, "page_size": pageSize})
}

func (h *TransferAdminHandler) FreezeTransfer(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		Reason string `json:"reason"`
	}
	_ = c.ShouldBindJSON(&req)
	adminID := getUserIDFromContext(c)
	if err := h.transferService.FreezeTransfer(c.Request.Context(), adminID, id); err != nil {
		writeAppError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "transfer frozen"})
}

func (h *TransferAdminHandler) RevokeTransfer(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		Reason string `json:"reason"`
	}
	_ = c.ShouldBindJSON(&req)
	adminID := getUserIDFromContext(c)
	if err := h.transferService.RevokeTransfer(c.Request.Context(), adminID, id, req.Reason); err != nil {
		writeAppError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "transfer revoked"})
}

func (h *TransferAdminHandler) BatchDistribute(c *gin.Context) {
	var req struct {
		Targets []service.BatchDistributeTarget `json:"targets" binding:"required"`
		Memo    *string                          `json:"memo"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	adminID := getUserIDFromContext(c)
	records, err := h.transferService.BatchDistribute(c.Request.Context(), adminID, req.Targets, req.Memo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": records, "count": len(records)})
}

func (h *TransferAdminHandler) GetFeeStats(c *gin.Context) {
	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -30)
	if startStr := c.Query("start_time"); startStr != "" {
		if t, err := time.Parse(time.RFC3339, startStr); err == nil {
			startTime = t
		}
	}
	if endStr := c.Query("end_time"); endStr != "" {
		if t, err := time.Parse(time.RFC3339, endStr); err == nil {
			endTime = t
		}
	}
	stats, err := h.transferService.GetFeeStats(c.Request.Context(), startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *TransferAdminHandler) ListRedPackets(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	records, total, err := h.transferService.GetAllRedPackets(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": records, "total": total, "page": page, "page_size": pageSize})
}

func getUserIDFromContext(c *gin.Context) int64 {
	subject, ok := middleware.GetAuthSubjectFromContext(c)
	if !ok {
		return 0
	}
	return subject.UserID
}

func writeAppError(c *gin.Context, err error) {
	var appErr *infraerrors.ApplicationError
	if errors.As(err, &appErr) {
		c.JSON(int(appErr.Code), gin.H{"error": appErr.Message, "code": appErr.Reason})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
}
