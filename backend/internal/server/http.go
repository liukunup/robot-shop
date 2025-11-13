package server

import (
	"backend/docs"
	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/pkg/jwt"
	"backend/pkg/log"
	"backend/pkg/server/http"
	"backend/web"
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"mime"
	"path/filepath"
	"strings"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var serverStart = time.Now().UTC()

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

	// 前端 SPA 静态资源与回退处理
	if distFS, err := fs.Sub(web.Assets(), "dist"); err != nil {
		logger.Error("Failed to load frontend assets: " + err.Error())
	} else {
		// 预读 index.html 避免 gin Static 的 301 异常重定向问题
		indexBytes, readErr := fs.ReadFile(distFS, "index.html")
		if readErr != nil {
			logger.Error("index.html missing: " + readErr.Error())
		} else {
			// 根路径直接返回 index.html 内容
			s.GET("/", func(c *gin.Context) {
				etag := calcETag(indexBytes)
				setCacheHeaders(c, "index.html", etag)
				// 304 短路
				if matchIfNoneMatch(c, etag) {
					c.Status(304)
					return
				}
				c.Data(200, "text/html; charset=utf-8", indexBytes)
			})
		}

		// 统一处理前端静态文件与 SPA 回退
		s.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			// API 前缀 -> JSON 404
			if len(path) >= 3 && path[:3] == "/v1" {
				c.JSON(404, gin.H{"error": "API endpoint not found"})
				return
			}
			// swagger 前缀放行（交由已有路由处理）
			if len(path) >= 8 && path[:8] == "/swagger" {
				return
			}

			cleanPath := path
			if len(cleanPath) > 1 && cleanPath[0] == '/' {
				cleanPath = cleanPath[1:]
			}

			// 存在真实静态文件则返回
			if f, openErr := distFS.Open(cleanPath); openErr == nil {
				if info, statErr := f.Stat(); statErr == nil && !info.IsDir() {
					data, _ := fs.ReadFile(distFS, cleanPath)
					ext := filepath.Ext(cleanPath)
					ct := mime.TypeByExtension(ext)
					if ct == "" {
						ct = "application/octet-stream"
					}
					etag := calcETag(data)
					setCacheHeaders(c, cleanPath, etag)
					if matchIfNoneMatch(c, etag) {
						c.Status(304)
						return
					}
					c.Data(200, ct, data)
					return
				}
			}
			// 回退到 index.html（确保前端路由正常）
			if len(indexBytes) > 0 {
				etag := calcETag(indexBytes)
				setCacheHeaders(c, "index.html", etag)
				if matchIfNoneMatch(c, etag) {
					c.Status(304)
					return
				}
				c.Data(200, "text/html; charset=utf-8", indexBytes)
				return
			}
			c.String(500, "index.html not loaded")
		})
	}

	// 健康检查路由
	s.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"time":    time.Now().UTC().Format(time.RFC3339),
			"version": "v1",
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

// --- 工具函数 ---

func calcETag(data []byte) string {
	sha := sha256.Sum256(data)
	return "\"sha256-" + hex.EncodeToString(sha[:]) + "\"" // 使用完整哈希作为强ETag
}

func matchIfNoneMatch(c *gin.Context, etag string) bool {
	if inm := c.GetHeader("If-None-Match"); inm != "" {
		return inm == etag
	}
	return false
}

func setCacheHeaders(c *gin.Context, name, etag string) {
	// 根据文件名是否包含hash决定是否长期缓存
	cc := "no-cache"
	if looksHashed(name) {
		cc = "public, max-age=31536000, immutable"
	}
	c.Header("Cache-Control", cc)
	c.Header("ETag", etag)
	c.Header("Last-Modified", serverStart.Format(time.RFC1123))
}

func looksHashed(name string) bool {
	// 简单判断：文件名中是否包含 8+ 位的hex片段
	for _, part := range strings.Split(name, ".") {
		if len(part) >= 8 {
			allHex := true
			for _, r := range part {
				if (r < '0' || r > '9') && (r < 'a' || r > 'f') && (r < 'A' || r > 'F') {
					allHex = false
					break
				}
			}
			if allHex {
				return true
			}
		}
	}
	return false
}
