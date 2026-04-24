package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/myuto/attendance-backend/internal/interface/handler"
	custommw "github.com/myuto/attendance-backend/internal/interface/middleware"
)

type Router struct {
	echo              *echo.Echo
	authHandler       *handler.AuthHandler
	attendanceHandler *handler.AttendanceHandler
	adminHandler      *handler.AdminHandler
	authMiddleware    *custommw.AuthMiddleware
	corsOrigins       []string
}

func NewRouter(
	authHandler *handler.AuthHandler,
	attendanceHandler *handler.AttendanceHandler,
	adminHandler *handler.AdminHandler,
	authMiddleware *custommw.AuthMiddleware,
	corsOrigins []string,
) *Router {
	return &Router{
		echo:              echo.New(),
		authHandler:       authHandler,
		attendanceHandler: attendanceHandler,
		adminHandler:      adminHandler,
		authMiddleware:    authMiddleware,
		corsOrigins:       corsOrigins,
	}
}

func (r *Router) Setup() *echo.Echo {
	// Validator
	r.echo.Validator = custommw.NewCustomValidator()

	// Middleware
	r.echo.Use(middleware.Logger())
	r.echo.Use(middleware.Recover())
	r.echo.Use(custommw.CORS(r.corsOrigins))

	// Health check
	r.echo.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	// API v1
	v1 := r.echo.Group("/api/v1")

	// Auth routes (public)
	auth := v1.Group("/auth")
	auth.POST("/register", r.authHandler.Register)
	auth.POST("/login", r.authHandler.Login)

	// Protected routes
	protected := v1.Group("")
	protected.Use(r.authMiddleware.Authenticate)

	// Auth - Me endpoint
	protected.GET("/auth/me", r.authHandler.Me)

	// Attendance routes (authenticated users)
	attendance := protected.Group("/attendances")
	attendance.GET("", r.attendanceHandler.GetAttendances)
	attendance.POST("/clock-in", r.attendanceHandler.ClockIn)
	attendance.POST("/clock-out", r.attendanceHandler.ClockOut)

	// Admin routes (admin only)
	admin := protected.Group("/admin")
	admin.Use(r.authMiddleware.AdminOnly)
	admin.GET("/employees", r.adminHandler.GetEmployees)
	admin.GET("/employees/:id", r.adminHandler.GetEmployee)

	return r.echo
}
