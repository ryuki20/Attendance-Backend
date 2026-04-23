package handler

import (
	"net/http"
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
