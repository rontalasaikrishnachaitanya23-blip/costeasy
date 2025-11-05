package main

import (
	"log"

	"github.com/chaitu35/costeasy/backend/app/config"
	"github.com/chaitu35/costeasy/backend/app/middleware"
	"github.com/chaitu35/costeasy/backend/database"

	// GL-Core imports
	glHandler "github.com/chaitu35/costeasy/backend/internal/gl-core/handler"
	glRepo "github.com/chaitu35/costeasy/backend/internal/gl-core/repository"
	glRoutes "github.com/chaitu35/costeasy/backend/internal/gl-core/routes"
	glService "github.com/chaitu35/costeasy/backend/internal/gl-core/service"

	// Settings imports
	settingsHandler "github.com/chaitu35/costeasy/backend/internal/settings/handler"
	settingsRepo "github.com/chaitu35/costeasy/backend/internal/settings/repository"
	settingsRoutes "github.com/chaitu35/costeasy/backend/internal/settings/routes"
	settingsService "github.com/chaitu35/costeasy/backend/internal/settings/service"

	// Auth imports
	authHandler "github.com/chaitu35/costeasy/backend/internal/auth/handler"
	authRepo "github.com/chaitu35/costeasy/backend/internal/auth/repository"
	authRoutes "github.com/chaitu35/costeasy/backend/internal/auth/routes"
	authService "github.com/chaitu35/costeasy/backend/internal/auth/service"

	// Pkg imports
	"github.com/chaitu35/costeasy/backend/pkg/crypto"
	"github.com/chaitu35/costeasy/backend/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	pool, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Initialize JWT util
	jwtUtil := jwt.NewJWTUtil(
		cfg.JWTSecret,

		cfg.JWTAccessExpiry,
		cfg.JWTRefreshExpiry,
	)

	// Set Gin mode based on environment
	if cfg.LogLevel == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Auth module (must be first)
		setupAuthRoutes(v1, pool, jwtUtil)

		// GL module
		setupGLRoutes(v1, pool)

		// Settings module
		setupSettingsRoutes(v1, pool, cfg)
	}

	// Start server
	addr := cfg.Host + ":" + cfg.Port

	log.Printf("Server starting on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// setupAuthRoutes initializes all Auth routes
func setupAuthRoutes(v1 *gin.RouterGroup, pool database.Pool, jwtUtil *jwt.JWTUtil) {
	// Initialize repositories
	userRepo := authRepo.NewUserRepository(pool)
	roleRepo := authRepo.NewRoleRepository(pool)
	permRepo := authRepo.NewPermissionRepository(pool)
	refreshTokenRepo := authRepo.NewRefreshTokenRepository(pool)

	// Initialize service
	authSvc := authService.NewAuthService(userRepo, roleRepo, permRepo, refreshTokenRepo, jwtUtil)

	// Initialize handler
	handler := authHandler.NewAuthHandler(authSvc)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtUtil, authSvc)

	// Register routes
	authRoutes.RegisterAuthRoutes(v1, handler, authMiddleware)
}

// setupGLRoutes initializes all GL-Core routes
func setupGLRoutes(v1 *gin.RouterGroup, pool database.Pool) {
	gl := v1.Group("/gl")

	// Initialize Account module
	accountRepo := glRepo.NewGLAccountRepository(pool)
	accountService := glService.NewAccountService(accountRepo)
	accountHandler := glHandler.NewAccountHandler(accountService)

	// Initialize Journal Entry module
	journalEntryRepo := glRepo.NewJournalEntryRepository(pool)
	journalLineRepo := glRepo.NewJournalLineRepository(pool)
	journalEntryService := glService.NewJournalEntryService(journalEntryRepo, accountRepo)
	journalEntryHandler := glHandler.NewJournalEntryHandler(journalEntryService)

	// Initialize Import module
	importService := glService.NewImportService(pool, accountRepo, journalEntryRepo, journalLineRepo)
	importHandler := glHandler.NewImportHandler(importService)

	// Register routes
	glRoutes.RegisterAccountRoutes(gl, accountHandler)
	glRoutes.RegisterJournalEntryRoutes(gl, journalEntryHandler)
	glRoutes.RegisterImportRoutes(gl, importHandler)
}

func setupSettingsRoutes(v1 *gin.RouterGroup, pool database.Pool, cfg *config.Config) {
	settings := v1.Group("/settings")

	// Organization module
	orgRepo := settingsRepo.NewOrganizationRepository(pool)
	orgService := settingsService.NewOrganizationService(orgRepo)
	orgHandler := settingsHandler.NewOrganizationHandler(orgService)

	// Crypto service
	cryptoService, err := crypto.NewCryptoService()
	if err != nil {
		log.Fatalf("Failed to initialize crypto service: %v", err)
	}
	// Shafafiya config
	shafafiyaConfig := &settingsService.Config{
		// Add your config values here
        ShafafiyaEnvironment: cfg.Environment,
	}

	// Shafafiya module
	shafafiyaRepo := settingsRepo.NewShafafiyaRepository(pool)
	shafafiyaService := settingsService.NewShafafiyaService(shafafiyaRepo, orgRepo, cryptoService, shafafiyaConfig)
	shafafiyaHandler := settingsHandler.NewShafafiyaHandler(shafafiyaService)

	// Register routes
	settingsRoutes.RegisterOrganizationRoutes(settings, orgHandler)
	settingsRoutes.RegisterShafafiyaRoutes(settings, shafafiyaHandler)
}
