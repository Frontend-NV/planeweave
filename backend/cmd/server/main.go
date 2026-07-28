// @title           Planeweave API
// @version         1.0
// @description     REST API командной доски задач с графом зависимостей.
// @description     WebSocket (синхронизация): см. backend/docs/websocket.md
//
// @host      localhost:8080
// @BasePath  /
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT-токен: `Bearer <token>`

package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"

	_ "github.com/planeweave/planeweave/docs/swag"

	"github.com/planeweave/planeweave/internal/config"
	"github.com/planeweave/planeweave/internal/database"
	"github.com/planeweave/planeweave/internal/handler"
	"github.com/planeweave/planeweave/internal/pkg/jwt"
	"github.com/planeweave/planeweave/internal/repository"
	"github.com/planeweave/planeweave/internal/service"
	"github.com/planeweave/planeweave/internal/ws"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()
	ctx := context.Background()

	pool, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	migrationsDir := filepath.Join("db", "migrations")
	if _, err := os.Stat(migrationsDir); err != nil {
		migrationsDir = filepath.Join("backend", "db", "migrations")
	}
	if err := database.Migrate(ctx, pool, migrationsDir); err != nil {
		log.Fatal("migrate: ", err)
	}

	repo := repository.New(pool)
	tokens := jwt.NewService(cfg.JWTSecret)
	hub := ws.NewHub()

	accessSvc := service.NewAccessService(repo)
	authSvc := service.NewAuthService(repo, tokens)
	projectSvc := service.NewProjectService(repo, accessSvc)
	inviteSvc := service.NewInviteService(repo, accessSvc)
	taskSvc := service.NewTaskService(repo, accessSvc, hub)
	depSvc := service.NewDependencyService(repo, accessSvc, hub)

	router := handler.NewRouter(cfg, authSvc, handler.Handlers{
		Auth:       handler.NewAuthHandler(authSvc),
		Project:    handler.NewProjectHandler(projectSvc),
		Invite:     handler.NewInviteHandler(inviteSvc),
		Task:       handler.NewTaskHandler(taskSvc),
		Dependency: handler.NewDependencyHandler(depSvc),
		WS:         handler.NewWSHandler(tokens, projectSvc, hub),
	})

	log.Printf("listening on %s", cfg.HTTPAddr)
	log.Printf("swagger UI: http://localhost%s/swagger/index.html", cfg.HTTPAddr)
	if err := router.Run(cfg.HTTPAddr); err != nil {
		log.Fatal(err)
	}
}
