package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"golaninvite/internal/admin"
	"golaninvite/internal/auth"
	"golaninvite/internal/config"
	"golaninvite/internal/landing"
	"golaninvite/internal/middleware"
	"golaninvite/internal/packages"
	"golaninvite/internal/users"
)

type App struct {
	Config   *config.Config
	DB       *pgxpool.Pool
	Logger   *slog.Logger
	server   *http.Server

	authHandler     *auth.Handler
	usersHandler    *users.Handler
	landingHandler  *landing.Handler
	packagesHandler *packages.Handler
	adminHandler    *admin.Handler

	authSvc auth.Service
}

func New(cfg *config.Config, db *pgxpool.Pool, logger *slog.Logger) *App {
	// Repositories
	userRepo := users.NewPostgresRepository(db)
	authRepo := auth.NewPostgresRepository(db)
	landingRepo := landing.NewPostgresRepository(db)
	packagesRepo := packages.NewPostgresRepository(db)
	adminRepo := admin.NewPostgresRepository(db)

	// Services
	userSvc := users.NewService(userRepo)
	authSvc := auth.NewService(authRepo, userRepo)
	landingSvc := landing.NewService(landingRepo)
	packagesSvc := packages.NewService(packagesRepo)
	adminSvc := admin.NewService(adminRepo, packagesRepo)

	// Handlers
	authHdl := auth.NewHandler(authSvc)
	userHdl := users.NewHandler(userSvc)
	landingHdl := landing.NewHandler(landingSvc)
	packagesHdl := packages.NewHandler(packagesSvc)
	adminHdl := admin.NewHandler(adminSvc)

	return &App{
		Config:          cfg,
		DB:              db,
		Logger:          logger,
		authHandler:     authHdl,
		usersHandler:    userHdl,
		landingHandler:  landingHdl,
		packagesHandler: packagesHdl,
		adminHandler:    adminHdl,
		authSvc:         authSvc,
	}
}

func (a *App) Run(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", a.Config.AppHost, a.Config.AppPort)

	mux := http.NewServeMux()

	// 1. Health Endpoints (SSOT §77)
	mux.HandleFunc("GET /health/live", a.handleHealthLive)
	mux.HandleFunc("GET /health/ready", a.handleHealthReady)

	// 2. Public APIs (SSOT §59, §66)
	mux.HandleFunc("GET /api/v1/public/landing", a.landingHandler.HandleGetPublicLanding)
	mux.HandleFunc("GET /api/v1/public/packages", a.packagesHandler.HandleGetPublicPackages)
	mux.HandleFunc("GET /api/v1/public/templates", a.adminHandler.HandleListTemplates)

	// 3. Auth APIs (SSOT §60)
	mux.HandleFunc("POST /api/v1/auth/login", a.authHandler.HandleLogin)
	mux.HandleFunc("POST /api/v1/auth/logout", a.authHandler.HandleLogout)

	authMiddleware := middleware.RequireAuth(a.authSvc)
	adminRoleMiddleware := middleware.RequireRole("admin")

	mux.Handle("GET /api/v1/auth/session", authMiddleware(http.HandlerFunc(a.authHandler.HandleGetSession)))

	// Helper for admin protected routes
	adminProtected := func(h http.HandlerFunc) http.Handler {
		return authMiddleware(adminRoleMiddleware(h))
	}

	// 4. Admin User Management APIs
	mux.Handle("GET /api/v1/admin/overview", adminProtected(a.usersHandler.HandleAdminOverview))
	mux.Handle("GET /api/v1/admin/users", adminProtected(a.usersHandler.HandleAdminListUsers))
	mux.Handle("POST /api/v1/admin/users", adminProtected(a.usersHandler.HandleAdminCreateUser))
	mux.Handle("PATCH /api/v1/admin/users/", adminProtected(a.adminHandler.HandleUpdateUser))
	mux.Handle("POST /api/v1/admin/users/{id}/reset-password", adminProtected(a.adminHandler.HandleResetUserPassword))
	mux.Handle("DELETE /api/v1/admin/users/{id}", adminProtected(a.adminHandler.HandleDeleteUser))

	// 5. Admin Orders APIs
	mux.Handle("GET /api/v1/admin/orders", adminProtected(a.adminHandler.HandleListOrders))
	mux.Handle("PATCH /api/v1/admin/orders/{id}/status", adminProtected(a.adminHandler.HandleUpdateOrderStatus))

	// 6. Admin Invitations APIs
	mux.Handle("GET /api/v1/admin/invitations", adminProtected(a.adminHandler.HandleListInvitations))
	mux.Handle("PATCH /api/v1/admin/invitations/{id}/status", adminProtected(a.adminHandler.HandleUpdateInvitationStatus))

	// 7. Admin Templates APIs
	mux.Handle("GET /api/v1/admin/templates", adminProtected(a.adminHandler.HandleListTemplates))
	mux.Handle("POST /api/v1/admin/templates", adminProtected(a.adminHandler.HandleCreateTemplate))
	mux.Handle("DELETE /api/v1/admin/templates/{id}", adminProtected(a.adminHandler.HandleDeleteTemplate))

	// 8. Admin Packages APIs
	mux.Handle("GET /api/v1/admin/packages", adminProtected(a.adminHandler.HandleListPackages))
	mux.Handle("POST /api/v1/admin/packages", adminProtected(a.adminHandler.HandleCreatePackage))
	mux.Handle("DELETE /api/v1/admin/packages/{id}", adminProtected(a.adminHandler.HandleDeletePackage))

	// 9. Admin Domains APIs
	mux.Handle("GET /api/v1/admin/domains", adminProtected(a.adminHandler.HandleListDomains))
	mux.Handle("POST /api/v1/admin/domains/{id}/activate", adminProtected(a.adminHandler.HandleActivateDomain))

	// 10. Admin RSVP, Greetings, Transactions
	mux.Handle("GET /api/v1/admin/rsvp", adminProtected(a.adminHandler.HandleListAllRSVP))
	mux.Handle("GET /api/v1/admin/greetings", adminProtected(a.adminHandler.HandleListAllGreetings))
	mux.Handle("PATCH /api/v1/admin/greetings/{id}/status", adminProtected(a.adminHandler.HandleModerateGreeting))
	mux.Handle("GET /api/v1/admin/transactions", adminProtected(a.adminHandler.HandleListTransactions))

	// 11. Admin Packages - tambah PUT untuk edit
	mux.Handle("PUT /api/v1/admin/packages/", adminProtected(a.adminHandler.HandleUpdatePackage))

	// 12. Admin Reviews
	mux.Handle("GET /api/v1/admin/reviews", adminProtected(a.adminHandler.HandleListReviews))
	mux.Handle("PATCH /api/v1/admin/reviews/{id}/status", adminProtected(a.adminHandler.HandleModerateReview))
	mux.Handle("POST /api/v1/admin/reviews/{id}/reply", adminProtected(a.adminHandler.HandleReplyReview))

	// Middleware Chain (Request ID -> CORS)
	handler := middleware.SetRequestID(a.corsMiddleware(mux))

	a.server = &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	a.Logger.Info("server GolanInvite dimulai",
		slog.String("addr", addr),
		slog.String("env", a.Config.AppEnv),
	)

	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("app: server error: %w", err)
	}

	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	if a.server == nil {
		return nil
	}
	return a.server.Shutdown(ctx)
}

func (a *App) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			if strings.HasPrefix(origin, "http://localhost") || strings.HasPrefix(origin, "https://") {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, X-Request-ID")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *App) handleHealthLive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"live"}`))
}

func (a *App) handleHealthReady(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	w.Header().Set("Content-Type", "application/json")

	if err := a.DB.Ping(ctx); err != nil {
		a.Logger.Error("health ready: database tidak terjangkau", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"status":"not_ready","reason":"database unavailable"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ready"}`))
}
