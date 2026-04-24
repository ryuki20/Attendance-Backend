package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/myuto/attendance-backend/internal/domain/entity"
	"github.com/myuto/attendance-backend/internal/usecase"
)

type AuthHandler struct {
	authUseCase usecase.AuthUseCase
}

func NewAuthHandler(authUseCase usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUseCase: authUseCase}
}

type RegisterRequest struct {
	Email    string              `json:"email" validate:"required,email"`
	Password string              `json:"password" validate:"required,min=6"`
	Name     string              `json:"name" validate:"required"`
	Role     entity.EmployeeRole `json:"role" validate:"omitempty,oneof=admin employee"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token    string           `json:"token"`
	Employee *entity.Employee `json:"employee"`
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	if req.Role == "" {
		req.Role = entity.RoleEmployee
	}

	employee, err := h.authUseCase.Register(c.Request().Context(), req.Email, req.Password, req.Name, req.Role)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, employee)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	token, employee, err := h.authUseCase.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, AuthResponse{
		Token:    token,
		Employee: employee,
	})
}

func (h *AuthHandler) Me(c echo.Context) error {
	employeeID, ok := c.Get("employee_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"employee_id": employeeID,
		"role":        c.Get("role"),
	})
}
