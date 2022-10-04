package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	authhttp "github.com/khuchuz/go-clean-architecture-sql/auth/delivery"
	itface "github.com/khuchuz/go-clean-architecture-sql/auth/itface"
	authrepo "github.com/khuchuz/go-clean-architecture-sql/auth/repository"
	authusecase "github.com/khuchuz/go-clean-architecture-sql/auth/usecase"
	"github.com/khuchuz/go-clean-architecture-sql/models"
)

func main() {

	app := NewApp()

	if err := app.Run("8000"); err != nil {
		log.Fatalf("%s", err.Error())
	}
}

type App struct {
	httpServer *http.Server
	authUC     itface.UseCase
}

func NewApp() *App {
	db := SetupDatabase()

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
	// Init gin handler
	router := gin.Default()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	// Set up http handlers
	authhttp.RegisterHTTPEndpoints(router, a.authUC)

	// API endpoints
	authMiddleware := authhttp.NewAuthMiddleware(a.authUC)
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

func SetupDatabase() *gorm.DB {
	var DBHost string = "127.0.0.1"
	var DBPort string = "3306"
	var DBUser string = "root"
	var DBPass string = ""
	var DBName string = "go_clean_architecture"

	DSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", DBUser, DBPass, DBHost, DBPort, DBName)

	db, err := gorm.Open(mysql.Open(DSN), &gorm.Config{})
	db.AutoMigrate(&models.User{})
	if err != nil {
		panic(err)
	}
	return db
}
