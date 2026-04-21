package handler

import (
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"
	"github.com/myuto/attendance-backend/internal/usecase"
)

type AttendanceHandler struct {
	attendanceUseCase usecase.AttendanceUseCase
}

func NewAttendanceHandler(attendanceUseCase usecase.AttendanceUseCase) *AttendanceHandler {
	return &AttendanceHandler{attendanceUseCase: attendanceUseCase}
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

	return c.JSON(http.StatusOK, attendances)
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

	return c.JSON(http.StatusOK, attendance)
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

	return c.JSON(http.StatusOK, attendance)
}

