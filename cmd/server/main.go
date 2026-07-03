package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mcicare/mci-mailer/internal/config"
	"github.com/mcicare/mci-mailer/internal/handler"
	"github.com/mcicare/mci-mailer/internal/middleware"
	"github.com/mcicare/mci-mailer/internal/repository"
	"github.com/mcicare/mci-mailer/internal/service"
	smtpClient "github.com/mcicare/mci-mailer/internal/smtp"
)

func main() {
	cfg := config.Load()

	// ── Database ────────────────────────────────────────────────────────────
	db, err := pgxpool.New(context.Background(), cfg.Database.URL)
	if err != nil {
		log.Fatalf("[main] failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		log.Fatalf("[main] database ping failed: %v", err)
	}
	log.Println("[main] database connected")

	// ── SMTP ────────────────────────────────────────────────────────────────
	smtp := smtpClient.NewClient(&cfg.SMTP)
	log.Printf("[main] SMTP configured → %s:%d", cfg.SMTP.Host, cfg.SMTP.Port)

	// ── Repositories ────────────────────────────────────────────────────────
	apiKeyRepo  := repository.NewApiKeyRepository(db)
	templateRepo := repository.NewTemplateRepository(db)
	emailLogRepo := repository.NewEmailLogRepository(db)

	// ── Services ─────────────────────────────────────────────────────────────
	apiKeySvc  := service.NewApiKeyService(apiKeyRepo)
	templateSvc := service.NewTemplateService(templateRepo)
	mailerSvc  := service.NewMailerService(smtp, emailLogRepo, templateSvc, &cfg.SMTP)

	// ── Handlers ─────────────────────────────────────────────────────────────
	healthHandler   := handler.NewHealthHandler(smtp)
	mailHandler     := handler.NewMailHandler(mailerSvc)
	templateHandler := handler.NewTemplateHandler(templateSvc)
	apiKeyHandler   := handler.NewApiKeyHandler(apiKeySvc)
	logHandler      := handler.NewLogHandler(emailLogRepo)

	// ── Router ───────────────────────────────────────────────────────────────
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())

	r.GET("/health", healthHandler.Check)

	v1 := r.Group("/v1")
	v1.Use(middleware.Auth(apiKeyRepo))
	v1.Use(middleware.RateLimit(10, 20)) // 10 req/s, burst 20

	// Mail
	v1.POST("/mail/send",
		middleware.RequireScope("mail:send"),
		mailHandler.Send,
	)

	// Logs
	v1.GET("/logs",
		middleware.RequireScope("logs:read"),
		logHandler.List,
	)

	// Templates
	v1.GET("/templates",
		middleware.RequireScope("templates:read"),
		templateHandler.List,
	)
	v1.POST("/templates",
		middleware.RequireScope("templates:write"),
		templateHandler.Create,
	)
	v1.PUT("/templates/:name",
		middleware.RequireScope("templates:write"),
		templateHandler.Update,
	)
	v1.DELETE("/templates/:name",
		middleware.RequireScope("templates:write"),
		templateHandler.Delete,
	)

	// API Keys (gestion interne — scope keys:manage requis)
	v1.GET("/apikeys",
		middleware.RequireScope("keys:manage"),
		apiKeyHandler.List,
	)
	v1.POST("/apikeys",
		middleware.RequireScope("keys:manage"),
		apiKeyHandler.Create,
	)
	v1.DELETE("/apikeys/:id",
		middleware.RequireScope("keys:manage"),
		apiKeyHandler.Revoke,
	)

	// ── Graceful shutdown ─────────────────────────────────────────────────────
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		log.Printf("[main] mci-mailer listening on :%s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[main] server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[main] shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("[main] forced shutdown: %v", err)
	}
	log.Println("[main] server stopped")
}
