package main

import (
	"fmt"
	"log"

	"github.com/myuto/attendance-backend/internal/infrastructure/database"
	"github.com/myuto/attendance-backend/internal/infrastructure/repository"
	"github.com/myuto/attendance-backend/internal/infrastructure/router"
	"github.com/myuto/attendance-backend/internal/interface/handler"
	"github.com/myuto/attendance-backend/internal/interface/middleware"
	"github.com/myuto/attendance-backend/internal/usecase"
	"github.com/myuto/attendance-backend/pkg/config"
)

func main() {
	// 設定の読み込み
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// データベース接続
	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Successfully connected to database")

	// リポジトリの初期化
	userRepo := repository.NewUserRepository(db)
	attendanceRepo := repository.NewAttendanceRepository(db)

	// ユースケースの初期化
	authUseCase := usecase.NewAuthUseCase(userRepo, cfg.JWT.Secret, cfg.JWT.Expiration)
	attendanceUseCase := usecase.NewAttendanceUseCase(attendanceRepo)
	adminUseCase := usecase.NewAdminUseCase(userRepo)

	// ハンドラーの初期化
	authHandler := handler.NewAuthHandler(authUseCase)
	attendanceHandler := handler.NewAttendanceHandler(attendanceUseCase)
	adminHandler := handler.NewAdminHandler(adminUseCase)

	// ミドルウェアの初期化
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWT.Secret)

	// ルーターの初期化
	r := router.NewRouter(authHandler, attendanceHandler, adminHandler, authMiddleware, cfg.CORS.AllowOrigins)
	e := r.Setup()

	// サーバーの起動
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Starting server on %s (environment: %s)", addr, cfg.Server.Env)
	if err := e.Start(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
