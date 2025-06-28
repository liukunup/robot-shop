package server

import (
	apiV1 "backend/api/v1"
	"backend/docs"
	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/pkg/jwt"
	"backend/pkg/log"
	"backend/pkg/server/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewHTTPServer(
	logger *log.Logger,
	conf *viper.Viper,
	jwt *jwt.JWT,
	e *casbin.SyncedEnforcer,
	userHandler *handler.UserHandler,
	roleHandler *handler.RoleHandler,
	menuHandler *handler.MenuHandler,
	apiHandler *handler.ApiHandler,
	robotHandler *handler.RobotHandler,
) *http.Server {
	gin.SetMode(gin.DebugMode)
	s := http.NewServer(
		gin.Default(),
		logger,
		http.WithServerHost(conf.GetString("http.host")),
		http.WithServerPort(conf.GetInt("http.port")),
	)

	// swagger doc
	docs.SwaggerInfo.BasePath = "/"
	s.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerfiles.Handler,
		//ginSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", conf.GetInt("app.http.port"))),
		ginSwagger.DefaultModelsExpandDepth(-1),
		ginSwagger.PersistAuthorization(true),
	))

	s.Use(
		middleware.CORSMiddleware(),
		middleware.ResponseLogMiddleware(logger),
		middleware.RequestLogMiddleware(logger),
		//middleware.SignMiddleware(log),
	)
	s.GET("/", func(ctx *gin.Context) {
		logger.WithContext(ctx).Info("hello")
		apiV1.HandleSuccess(ctx, map[string]interface{}{
			":)": "Thank you for using nunu!",
		})
	})

	v1 := s.Group("/v1")
	{
		// No route group has permission
		noAuthRouter := v1.Group("/")
		{
			noAuthRouter.POST("/register", userHandler.Register)
			noAuthRouter.POST("/login", userHandler.Login)
		}

		// Strict permission routing group
		strictAuthRouter := v1.Group("/").Use(middleware.StrictAuth(jwt, logger), middleware.AuthMiddleware(e))
		{
			// User
			strictAuthRouter.GET("/admin/users", userHandler.ListUsers)
			strictAuthRouter.POST("/admin/users", userHandler.UserCreate)
			strictAuthRouter.PUT("/admin/users/:id", userHandler.UserUpdate)
			strictAuthRouter.DELETE("/admin/users/:id", userHandler.UserDelete)
			strictAuthRouter.GET("/users/me", userHandler.GetCurrentUser)
			strictAuthRouter.GET("/users/me/permission", userHandler.GetUserPermission)
			strictAuthRouter.GET("/users/me/menus", userHandler.GetUserMenu)

			// Role
			strictAuthRouter.GET("/admin/roles", roleHandler.ListRoles)
			strictAuthRouter.POST("/admin/roles", roleHandler.RoleCreate)
			strictAuthRouter.PUT("/admin/roles/:id", roleHandler.RoleUpdate)
			strictAuthRouter.DELETE("/admin/roles/:id", roleHandler.RoleDelete)
			strictAuthRouter.GET("/admin/roles/permission", roleHandler.GetRolePermission)
			strictAuthRouter.PUT("/admin/roles/permission", roleHandler.UpdateRolePermission)

			// Menu
			strictAuthRouter.GET("/admin/menus", menuHandler.ListMenus)
			strictAuthRouter.POST("/admin/menus", menuHandler.MenuCreate)
			strictAuthRouter.PUT("/admin/menus/:id", menuHandler.MenuUpdate)
			strictAuthRouter.DELETE("/admin/menus/:id", menuHandler.MenuDelete)

			// API
			strictAuthRouter.GET("/admin/apis", apiHandler.ListApis)
			strictAuthRouter.POST("/admin/apis", apiHandler.ApiCreate)
			strictAuthRouter.PUT("/admin/apis/:id", apiHandler.ApiUpdate)
			strictAuthRouter.DELETE("/admin/apis/:id", apiHandler.ApiDelete)

			// Robot
			strictAuthRouter.GET("/robots", robotHandler.ListRobots)
			strictAuthRouter.POST("/robots", robotHandler.RobotCreate)
			strictAuthRouter.GET("/robots/:id", robotHandler.GetRobot)
			strictAuthRouter.PUT("/robots/:id", robotHandler.RobotUpdate)
			strictAuthRouter.DELETE("/robots/:id", robotHandler.RobotDelete)
		}
	}

	return s
}
