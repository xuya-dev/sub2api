package handler

import (
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

type LeaderboardHandler struct {
	leaderboardService *service.LeaderboardService
	checkinService     *service.CheckinService
}

func NewLeaderboardHandler(leaderboardService *service.LeaderboardService, checkinService *service.CheckinService) *LeaderboardHandler {
	return &LeaderboardHandler{
		leaderboardService: leaderboardService,
		checkinService:     checkinService,
	}
}

func (h *LeaderboardHandler) GetBalanceLeaderboard(c *gin.Context) {
	page, pageSize := response.ParsePagination(c)

	result, err := h.leaderboardService.GetBalanceLeaderboard(c.Request.Context(), page, pageSize)
	if err != nil {
		response.InternalError(c, "Failed to get balance leaderboard")
		return
	}

	response.Paginated(c, result.Entries, result.Total, page, pageSize)
}

func (h *LeaderboardHandler) GetConsumptionLeaderboard(c *gin.Context) {
	period := c.DefaultQuery("period", "daily")
	if period != "daily" && period != "weekly" && period != "monthly" {
		response.BadRequest(c, "Invalid period, must be daily, weekly or monthly")
		return
	}

	page, pageSize := response.ParsePagination(c)

	result, err := h.leaderboardService.GetConsumptionLeaderboard(c.Request.Context(), period, page, pageSize)
	if err != nil {
		response.InternalError(c, "Failed to get consumption leaderboard")
		return
	}

	response.Paginated(c, result.Entries, result.Total, page, pageSize)
}

func (h *LeaderboardHandler) GetCheckinLeaderboard(c *gin.Context) {
	page, pageSize := response.ParsePagination(c)

	result, err := h.leaderboardService.GetCheckinLeaderboard(c.Request.Context(), page, pageSize)
	if err != nil {
		response.InternalError(c, "Failed to get checkin leaderboard")
		return
	}

	response.Paginated(c, result.Entries, result.Total, page, pageSize)
}
