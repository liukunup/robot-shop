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
			strictAuthRouter.GET("/user", userHandler.GetCurrentUser)
			strictAuthRouter.GET("/admin/users", userHandler.ListUsers)
			strictAuthRouter.POST("/admin/user", userHandler.UserCreate)
			strictAuthRouter.PUT("/admin/user", userHandler.UserUpdate)
			strictAuthRouter.DELETE("/admin/user", userHandler.UserDelete)
			strictAuthRouter.GET("/admin/user/permissions", userHandler.GetUserPermissions)

			// Role
			strictAuthRouter.GET("/admin/roles", userHandler.ListRoles)
			strictAuthRouter.POST("/admin/role", userHandler.RoleCreate)
			strictAuthRouter.PUT("/admin/role", userHandler.RoleUpdate)
			strictAuthRouter.DELETE("/admin/role", userHandler.RoleDelete)
			strictAuthRouter.GET("/admin/role/permissions", userHandler.GetRolePermissions)
			strictAuthRouter.PUT("/admin/role/permission", userHandler.UpdateRolePermission)

			// Menu
			strictAuthRouter.GET("/menu", userHandler.GetCurrentMenu)
			strictAuthRouter.GET("/admin/menus", userHandler.ListMenus)
			strictAuthRouter.POST("/admin/menu", userHandler.MenuCreate)
			strictAuthRouter.PUT("/admin/menu", userHandler.MenuUpdate)
			strictAuthRouter.DELETE("/admin/menu", userHandler.MenuDelete)

			// API
			strictAuthRouter.GET("/admin/apis", userHandler.ListApis)
			strictAuthRouter.POST("/admin/api", userHandler.ApiCreate)
			strictAuthRouter.PUT("/admin/api", userHandler.ApiUpdate)
			strictAuthRouter.DELETE("/admin/api", userHandler.ApiDelete)

			// Robot
			strictAuthRouter.GET("/robots", robotHandler.GetRobotList)
			strictAuthRouter.POST("/robots", robotHandler.CreateRobot)
			strictAuthRouter.GET("/robots/:id", robotHandler.GetRobot)
			strictAuthRouter.PUT("/robots/:id", robotHandler.UpdateRobot)
			strictAuthRouter.DELETE("/robots/:id", robotHandler.DeleteRobot)
		}
	}

	return s
}
