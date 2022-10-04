package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/khuchuz/go-clean-architecture-sql/auth/app/database"
	"github.com/khuchuz/go-clean-architecture-sql/auth/controllers"
	"github.com/khuchuz/go-clean-architecture-sql/auth/services"
	authrepo "github.com/khuchuz/go-clean-architecture-sql/auth/services/repository"
	authusecase "github.com/khuchuz/go-clean-architecture-sql/auth/services/usecase"
)

type App struct {
	httpServer *http.Server
	authUC     services.UseCase
}

func NewApp() *App {
	db := database.SetupDatabase()

	userRepo := authrepo.InitUserRepositorySQL(db)

	return &App{
		authUC: authusecase.NewAuthUseCase(
			userRepo,
			"hash_salt",
			[]byte("signing_key"),
			86400,
		),
	}
}

func (a *App) Run(port string) error {
	// To Disable debug
	//gin.SetMode(gin.ReleaseMode)
	// Init gin handler
	router := gin.Default()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	// Set up http handlers
	controllers.RegisterHTTPEndpoints(router, a.authUC)

	// API endpoints
	authMiddleware := controllers.NewAuthMiddleware(a.authUC)
	_ = router.Group("/api", authMiddleware)

	// HTTP Server
	a.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}
