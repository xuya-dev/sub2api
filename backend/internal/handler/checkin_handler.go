package handler

import (
	"strconv"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

type luckCheckinRequest struct {
	BetAmount float64 `json:"bet_amount" binding:"required,gt=0"`
}

type CheckinHandler struct {
	checkinService  *service.CheckinService
	blindboxService *service.BlindBoxService
}

func NewCheckinHandler(checkinService *service.CheckinService, blindboxService *service.BlindBoxService) *CheckinHandler {
	return &CheckinHandler{
		checkinService:  checkinService,
		blindboxService: blindboxService,
	}
}

func (h *CheckinHandler) Checkin(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	result, err := h.checkinService.Checkin(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, result)
}

func (h *CheckinHandler) LuckCheckin(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req luckCheckinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: bet_amount is required and must be greater than 0")
		return
	}

	result, err := h.checkinService.LuckCheckin(c.Request.Context(), subject.UserID, req.BetAmount)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, result)
}

func (h *CheckinHandler) GetStatus(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	status, err := h.checkinService.GetStatus(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, status)
}

func (h *CheckinHandler) GetCalendar(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	result, err := h.checkinService.GetCalendar(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, result)
}

func (h *CheckinHandler) GetBlindboxRecords(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	page := 1
	pageSize := 20
	if p, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil && p > 0 {
		page = p
	}
	if ps, err := strconv.Atoi(c.DefaultQuery("page_size", "20")); err == nil && ps > 0 && ps <= 100 {
		pageSize = ps
	}

	result, err := h.blindboxService.GetUserRecords(c.Request.Context(), subject.UserID, page, pageSize)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, result)
}
