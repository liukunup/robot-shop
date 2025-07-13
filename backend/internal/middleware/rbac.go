package middleware

import (
	v1 "backend/api/v1"
	"backend/internal/constant"
	"backend/pkg/jwt"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(e *casbin.SyncedEnforcer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从上下文获取用户信息（假设通过 JWT 或其他方式设置）
		v, exists := ctx.Get("claims")
		if !exists {
			v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
			ctx.Abort()
			return
		}
		uid := v.(*jwt.AccessClaims).UserID
		if convertor.ToString(uid) == constant.AdminUserID {
			// 防呆设计，超管跳过API权限检查
			ctx.Next()
			return
		}

		// 获取请求的资源和操作
		sub := convertor.ToString(uid)
		obj := constant.ApiResourcePrefix + ctx.Request.URL.Path
		act := ctx.Request.Method

		// 检查权限
		allowed, err := e.Enforce(sub, obj, act)
		if err != nil {
			v1.HandleError(ctx, http.StatusForbidden, v1.ErrForbidden, nil)
			ctx.Abort()
			return
		}
		if !allowed {
			v1.HandleError(ctx, http.StatusForbidden, v1.ErrForbidden, nil)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
