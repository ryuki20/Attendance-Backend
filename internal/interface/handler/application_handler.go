package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/myuto/attendance-backend/internal/domain/entity"
	"github.com/myuto/attendance-backend/internal/usecase"
)

type ApplicationHandler struct {
	applicationUseCase usecase.ApplicationUseCase
}

func NewApplicationHandler(applicationUseCase usecase.ApplicationUseCase) *ApplicationHandler {
	return &ApplicationHandler{applicationUseCase: applicationUseCase}
}

// --- レスポンス型 ---

type employeeSummaryResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type attendanceCorrectionDetailsResponse struct {
	Date              string  `json:"date"`
	RequestedClockIn  *string `json:"requested_clock_in"`
	RequestedClockOut *string `json:"requested_clock_out"`
}

type applicationResponse struct {
	ID          string                              `json:"id"`
	Employee    employeeSummaryResponse             `json:"employee"`
	Type        string                              `json:"type"`
	Status      string                              `json:"status"`
	Reason      string                              `json:"reason"`
	Details     attendanceCorrectionDetailsResponse `json:"details"`
	CreatedAt   string                              `json:"created_at"`
	UpdatedAt   string                              `json:"updated_at"`
}

type applicationDetailResponse struct {
	applicationResponse
	ApprovedBy   *employeeSummaryResponse `json:"approved_by"`
	ApprovedAt   *string                  `json:"approved_at"`
	AdminComment *string                  `json:"admin_comment"`
}

func toApplicationResponse(app *entity.Application) applicationResponse {
	res := applicationResponse{
		ID: app.ID,
		Employee: employeeSummaryResponse{
			ID:   app.Employee.ID,
			Name: app.Employee.Name,
		},
		Type:      string(app.Type),
		Status:    string(app.Status),
		Reason:    app.Reason,
		CreatedAt: app.CreatedAt.Format(time.RFC3339),
		UpdatedAt: app.UpdatedAt.Format(time.RFC3339),
	}

	if app.Correction != nil {
		res.Details.Date = app.Correction.Date.Format("2006-01-02")
		if app.Correction.RequestedClockIn != nil {
			v := app.Correction.RequestedClockIn.Format(time.RFC3339)
			res.Details.RequestedClockIn = &v
		}
		if app.Correction.RequestedClockOut != nil {
			v := app.Correction.RequestedClockOut.Format(time.RFC3339)
			res.Details.RequestedClockOut = &v
		}
	}

	return res
}

func toApplicationDetailResponse(app *entity.Application) applicationDetailResponse {
	res := applicationDetailResponse{
		applicationResponse: toApplicationResponse(app),
	}

	if app.Approver != nil {
		approver := employeeSummaryResponse{
			ID:   app.Approver.ID,
			Name: app.Approver.Name,
		}
		res.ApprovedBy = &approver
	}

	if app.ApprovedAt != nil {
		v := app.ApprovedAt.Format(time.RFC3339)
		res.ApprovedAt = &v
	}

	res.AdminComment = app.AdminComment

	return res
}

// --- ページネーションパラメータ共通パーサー ---

func parsePaginationParams(c echo.Context) (page, perPage int, err error) {
	page = 1
	perPage = 20

	if p := c.QueryParam("page"); p != "" {
		v, e := strconv.Atoi(p)
		if e != nil || v < 1 {
			return 0, 0, errors.New("page must be a positive integer")
		}
		page = v
	}

	if pp := c.QueryParam("per_page"); pp != "" {
		v, e := strconv.Atoi(pp)
		if e != nil || v < 1 || v > 100 {
			return 0, 0, errors.New("per_page must be between 1 and 100")
		}
		perPage = v
	}

	return page, perPage, nil
}

// --- 社員向けエンドポイント ---

func (h *ApplicationHandler) GetApplications(c echo.Context) error {
	employeeID := c.Get("employee_id").(string)

	page, perPage, err := parsePaginationParams(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	var status *entity.ApplicationStatus
	if s := c.QueryParam("status"); s != "" {
		v := entity.ApplicationStatus(s)
		if v != entity.ApplicationStatusPending && v != entity.ApplicationStatusApproved && v != entity.ApplicationStatusRejected {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid status value"})
		}
		status = &v
	}

	apps, total, err := h.applicationUseCase.GetApplications(c.Request().Context(), employeeID, status, page, perPage)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	list := make([]applicationResponse, len(apps))
	for i, app := range apps {
		list[i] = toApplicationResponse(app)
	}

	return c.JSON(http.StatusOK, map[string]any{
		"applications": list,
		"total":        total,
		"page":         page,
		"per_page":     perPage,
	})
}

func (h *ApplicationHandler) CreateApplication(c echo.Context) error {
	employeeID := c.Get("employee_id").(string)

	var req struct {
		Type              string  `json:"type"`
		Date              string  `json:"date"`
		RequestedClockIn  *string `json:"requested_clock_in"`
		RequestedClockOut *string `json:"requested_clock_out"`
		Reason            string  `json:"reason"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if req.Type == "" || req.Date == "" || req.Reason == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "type, date and reason are required"})
	}
	if entity.ApplicationType(req.Type) != entity.ApplicationTypeAttendanceCorrection {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid type value"})
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid date format"})
	}

	var clockIn, clockOut *time.Time
	if req.RequestedClockIn != nil {
		t, err := time.Parse(time.RFC3339, *req.RequestedClockIn)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid requested_clock_in format"})
		}
		clockIn = &t
	}
	if req.RequestedClockOut != nil {
		t, err := time.Parse(time.RFC3339, *req.RequestedClockOut)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid requested_clock_out format"})
		}
		clockOut = &t
	}

	app, err := h.applicationUseCase.CreateAttendanceCorrectionApplication(
		c.Request().Context(), employeeID, date, clockIn, clockOut, req.Reason,
	)
	if err != nil {
		if errors.Is(err, usecase.ErrCorrectionFieldsRequired) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, toApplicationResponse(app))
}

func (h *ApplicationHandler) CancelApplication(c echo.Context) error {
	id := c.Param("id")
	employeeID := c.Get("employee_id").(string)

	app, err := h.applicationUseCase.CancelApplication(c.Request().Context(), id, employeeID)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrApplicationNotFound):
			return c.JSON(http.StatusNotFound, map[string]string{"error": "application not found"})
		case errors.Is(err, usecase.ErrApplicationNotOwned):
			return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
		case errors.Is(err, usecase.ErrApplicationNotPending):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "only pending applications can be cancelled"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}
	}

	return c.JSON(http.StatusOK, toApplicationResponse(app))
}

// --- 管理者向けエンドポイント ---

func (h *ApplicationHandler) AdminGetApplications(c echo.Context) error {
	page, perPage, err := parsePaginationParams(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	var status *entity.ApplicationStatus
	if s := c.QueryParam("status"); s != "" {
		v := entity.ApplicationStatus(s)
		if v != entity.ApplicationStatusPending && v != entity.ApplicationStatusApproved && v != entity.ApplicationStatusRejected {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid status value"})
		}
		status = &v
	}

	apps, total, err := h.applicationUseCase.AdminListApplications(c.Request().Context(), status, page, perPage)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	list := make([]applicationResponse, len(apps))
	for i, app := range apps {
		list[i] = toApplicationResponse(app)
	}

	return c.JSON(http.StatusOK, map[string]any{
		"applications": list,
		"total":        total,
		"page":         page,
		"per_page":     perPage,
	})
}

func (h *ApplicationHandler) AdminGetApplication(c echo.Context) error {
	id := c.Param("id")

	app, err := h.applicationUseCase.AdminGetApplication(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrApplicationNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "application not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, toApplicationDetailResponse(app))
}

func (h *ApplicationHandler) ApproveApplication(c echo.Context) error {
	id := c.Param("id")
	approvedBy := c.Get("employee_id").(string)

	var req struct {
		AdminComment *string `json:"admin_comment"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	app, err := h.applicationUseCase.ApproveApplication(c.Request().Context(), id, approvedBy, req.AdminComment)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrApplicationNotFound):
			return c.JSON(http.StatusNotFound, map[string]string{"error": "application not found"})
		case errors.Is(err, usecase.ErrApplicationNotPending):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "only pending applications can be approved"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}
	}

	return c.JSON(http.StatusOK, toApplicationDetailResponse(app))
}

func (h *ApplicationHandler) RejectApplication(c echo.Context) error {
	id := c.Param("id")
	approvedBy := c.Get("employee_id").(string)

	var req struct {
		AdminComment *string `json:"admin_comment"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	app, err := h.applicationUseCase.RejectApplication(c.Request().Context(), id, approvedBy, req.AdminComment)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrApplicationNotFound):
			return c.JSON(http.StatusNotFound, map[string]string{"error": "application not found"})
		case errors.Is(err, usecase.ErrApplicationNotPending):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "only pending applications can be rejected"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}
	}

	return c.JSON(http.StatusOK, toApplicationDetailResponse(app))
}
