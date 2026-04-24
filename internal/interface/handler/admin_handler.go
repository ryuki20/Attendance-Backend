package handler

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/myuto/attendance-backend/internal/domain/entity"
	"github.com/myuto/attendance-backend/internal/usecase"
)

type AdminHandler struct {
	adminUseCase usecase.AdminUseCase
}

func NewAdminHandler(adminUseCase usecase.AdminUseCase) *AdminHandler {
	return &AdminHandler{adminUseCase: adminUseCase}
}

type userResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func toUserResponse(u *entity.User) userResponse {
	return userResponse{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		Role:      string(u.Role),
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	}
}

type listEmployeesResponse struct {
	Users   []userResponse `json:"users"`
	Total   int            `json:"total"`
	Page    int            `json:"page"`
	PerPage int            `json:"per_page"`
}

type attendancesInEmployeeResponse struct {
	Records []attendanceResponse `json:"records"`
	Total   int                  `json:"total"`
	Page    int                  `json:"page"`
	PerPage int                  `json:"per_page"`
}

type employeeDetailResponse struct {
	userResponse
	Attendances attendancesInEmployeeResponse `json:"attendances"`
}

func (h *AdminHandler) GetEmployee(c echo.Context) error {
	id := c.Param("id")

	yearMonth := c.QueryParam("year_month")
	if yearMonth == "" {
		now := time.Now()
		yearMonth = now.Format("2006-01")
	} else {
		matched, _ := regexp.MatchString(`^\d{4}-(0[1-9]|1[0-2])$`, yearMonth)
		if !matched {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "invalid year_month format (use YYYY-MM)",
			})
		}
	}

	page := 1
	perPage := 20

	if p := c.QueryParam("page"); p != "" {
		v, err := strconv.Atoi(p)
		if err != nil || v < 1 {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "page must be a positive integer",
			})
		}
		page = v
	}

	if pp := c.QueryParam("per_page"); pp != "" {
		v, err := strconv.Atoi(pp)
		if err != nil || v < 1 || v > 100 {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "per_page must be between 1 and 100",
			})
		}
		perPage = v
	}

	detail, err := h.adminUseCase.GetEmployee(c.Request().Context(), id, yearMonth, page, perPage)
	if err != nil {
		if errors.Is(err, usecase.ErrEmployeeNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "employee not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal server error",
		})
	}

	records := make([]attendanceResponse, len(detail.Attendances))
	for i, a := range detail.Attendances {
		records[i] = toAttendanceResponse(a)
	}

	return c.JSON(http.StatusOK, employeeDetailResponse{
		userResponse: toUserResponse(detail.User),
		Attendances: attendancesInEmployeeResponse{
			Records: records,
			Total:   detail.Total,
			Page:    detail.Page,
			PerPage: detail.PerPage,
		},
	})
}

func (h *AdminHandler) GetEmployees(c echo.Context) error {
	page := 1
	perPage := 20

	if p := c.QueryParam("page"); p != "" {
		v, err := strconv.Atoi(p)
		if err != nil || v < 1 {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "page must be a positive integer",
			})
		}
		page = v
	}

	if pp := c.QueryParam("per_page"); pp != "" {
		v, err := strconv.Atoi(pp)
		if err != nil || v < 1 || v > 100 {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "per_page must be between 1 and 100",
			})
		}
		perPage = v
	}

	var role *entity.UserRole
	if r := c.QueryParam("role"); r != "" {
		ur := entity.UserRole(r)
		if !ur.IsValid() {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "role must be one of admin, employee",
			})
		}
		role = &ur
	}

	users, total, err := h.adminUseCase.ListEmployees(c.Request().Context(), page, perPage, role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal server error",
		})
	}

	res := listEmployeesResponse{
		Users:   make([]userResponse, len(users)),
		Total:   total,
		Page:    page,
		PerPage: perPage,
	}
	for i, u := range users {
		res.Users[i] = toUserResponse(u)
	}

	return c.JSON(http.StatusOK, res)
}
