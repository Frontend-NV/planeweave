package handler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/planeweave/planeweave/internal/config"
	"github.com/planeweave/planeweave/internal/middleware"
	"github.com/planeweave/planeweave/internal/service"
)

type Handlers struct {
	Auth       *AuthHandler
	Project    *ProjectHandler
	Invite     *InviteHandler
	Task       *TaskHandler
	Dependency *DependencyHandler
	WS         *WSHandler
}

func NewRouter(cfg config.Config, auth *service.AuthService, handlers Handlers) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSOrigins,
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	r.GET("/health", Health)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", handlers.Auth.Register)
		authGroup.POST("/login", handlers.Auth.Login)
		authGroup.GET("/me", middleware.Auth(auth.Tokens()), handlers.Auth.Me)
	}

	api := r.Group("/")
	api.Use(middleware.Auth(auth.Tokens()))
	{
		api.GET("/projects", handlers.Project.List)
		api.POST("/projects", handlers.Project.Create)

		project := api.Group("/projects/:projectId")
		{
			project.GET("", handlers.Project.Get)
			project.PATCH("", handlers.Project.Patch)
			project.DELETE("", handlers.Project.Delete)
			project.GET("/members", handlers.Project.Members)
			project.POST("/invites", handlers.Invite.Create)
			project.GET("/tasks", handlers.Task.List)
			project.POST("/tasks", handlers.Task.Create)
			project.GET("/dependencies", handlers.Dependency.List)
			project.POST("/dependencies", handlers.Dependency.Create)
		}

		api.POST("/invites/:token/accept", handlers.Invite.Accept)
		api.PATCH("/tasks/:taskId", handlers.Task.Patch)
		api.DELETE("/tasks/:taskId", handlers.Task.Delete)
		api.DELETE("/dependencies/:dependencyId", handlers.Dependency.Delete)
	}

	r.GET("/ws", handlers.WS.Serve)
	return r
}
