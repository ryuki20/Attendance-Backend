package handler

import (
	"net/http"
	"regexp"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/myuto/attendance-backend/internal/domain/entity"
	"github.com/myuto/attendance-backend/internal/usecase"
)

type AttendanceHandler struct {
	attendanceUseCase usecase.AttendanceUseCase
}

func NewAttendanceHandler(attendanceUseCase usecase.AttendanceUseCase) *AttendanceHandler {
	return &AttendanceHandler{attendanceUseCase: attendanceUseCase}
}

type attendanceResponse struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	Date      string  `json:"date"`
	ClockIn   *string `json:"clock_in"`
	ClockOut  *string `json:"clock_out"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

func toAttendanceResponse(a *entity.Attendance) attendanceResponse {
	res := attendanceResponse{
		ID:        a.ID,
		UserID:    a.UserID,
		Date:      a.Date.Format("2006-01-02"),
		CreatedAt: a.CreatedAt.Format(time.RFC3339),
		UpdatedAt: a.UpdatedAt.Format(time.RFC3339),
	}
	if a.ClockIn != nil {
		s := a.ClockIn.Format(time.RFC3339)
		res.ClockIn = &s
	}
	if a.ClockOut != nil {
		s := a.ClockOut.Format(time.RFC3339)
		res.ClockOut = &s
	}
	return res
}

func (h *AttendanceHandler) GetAttendances(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	yearMonth := c.QueryParam("year_month")
	if yearMonth == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid year_month format (use YYYY-MM)",
		})
	}

	matched, _ := regexp.MatchString(`^\d{4}-(0[1-9]|1[0-2])$`, yearMonth)
	if !matched {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid year_month format (use YYYY-MM)",
		})
	}

	attendances, err := h.attendanceUseCase.GetAttendancesByMonth(c.Request().Context(), userID, yearMonth)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal server error",
		})
	}

	res := make([]attendanceResponse, len(attendances))
	for i, a := range attendances {
		res[i] = toAttendanceResponse(a)
	}
	return c.JSON(http.StatusOK, res)
}

func (h *AttendanceHandler) ClockIn(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	attendance, err := h.attendanceUseCase.ClockIn(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, toAttendanceResponse(attendance))
}

func (h *AttendanceHandler) ClockOut(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	attendance, err := h.attendanceUseCase.ClockOut(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, toAttendanceResponse(attendance))
}
