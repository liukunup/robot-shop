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
		http.WithCertFiles(conf.GetString("http.cert_file"), conf.GetString("http.key_file")),
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
			noAuthRouter.POST("/reset-password", userHandler.ResetPassword)
			noAuthRouter.POST("/refresh-token", userHandler.RefreshToken)
		}

		// Non-strict permission routing group
		noStrictAuthRouter := v1.Group("/").Use(middleware.NoStrictAuth(jwt, logger))
		{
			// User
			noStrictAuthRouter.GET("/users/:id", userHandler.GetUserByID)
		}

		// Strict permission routing group
		strictAuthRouter := v1.Group("/").Use(middleware.StrictAuth(jwt, logger), middleware.AuthMiddleware(e))
		{
			// User
			strictAuthRouter.GET("/users/profile", userHandler.GetProfile)
			strictAuthRouter.PUT("/users/profile", userHandler.UpdateProfile)
			strictAuthRouter.PUT("/users/profile/avatar", userHandler.UploadAvatar)
			strictAuthRouter.GET("/users/menu", userHandler.GetMenu)
			strictAuthRouter.PUT("/users/password", userHandler.UpdatePassword)

			// Admin User
			strictAuthRouter.GET("/admin/users", userHandler.ListUsers)
			strictAuthRouter.POST("/admin/users", userHandler.CreateUser)
			strictAuthRouter.PUT("/admin/users/:id", userHandler.UpdateUser)
			strictAuthRouter.DELETE("/admin/users/:id", userHandler.DeleteUser)

			// Admin Role
			strictAuthRouter.GET("/admin/roles", roleHandler.ListRoles)
			strictAuthRouter.POST("/admin/roles", roleHandler.CreateRole)
			strictAuthRouter.PUT("/admin/roles/:id", roleHandler.UpdateRole)
			strictAuthRouter.DELETE("/admin/roles/:id", roleHandler.DeleteRole)
			strictAuthRouter.GET("/roles", roleHandler.ListAllRoles)
			// Admin Role Permission
			strictAuthRouter.GET("/admin/roles/permissions", roleHandler.GetRolePermissions)
			strictAuthRouter.PUT("/admin/roles/permissions", roleHandler.UpdateRolePermissions)

			// Admin Menu
			strictAuthRouter.GET("/admin/menus", menuHandler.ListMenus)
			strictAuthRouter.POST("/admin/menus", menuHandler.CreateMenu)
			strictAuthRouter.PUT("/admin/menus/:id", menuHandler.UpdateMenu)
			strictAuthRouter.DELETE("/admin/menus/:id", menuHandler.DeleteMenu)

			// Admin API
			strictAuthRouter.GET("/admin/apis", apiHandler.ListApis)
			strictAuthRouter.POST("/admin/apis", apiHandler.CreateApi)
			strictAuthRouter.PUT("/admin/apis/:id", apiHandler.UpdateApi)
			strictAuthRouter.DELETE("/admin/apis/:id", apiHandler.DeleteApi)

			// Robot
			strictAuthRouter.GET("/robots", robotHandler.ListRobots)
			strictAuthRouter.POST("/robots", robotHandler.CreateRobot)
			strictAuthRouter.PUT("/robots/:id", robotHandler.UpdateRobot)
			strictAuthRouter.DELETE("/robots/:id", robotHandler.DeleteRobot)
			strictAuthRouter.GET("/robots/:id", robotHandler.GetRobot)
		}
	}

	return s
}
