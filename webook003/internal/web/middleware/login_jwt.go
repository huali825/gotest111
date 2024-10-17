package middleware

import (
	"encoding/gob"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	ijwt "goworkwebook/webook003/internal/web/jwt"
	"net/http"
	"time"
)

type LoginJWTMiddlewareBuilder struct {
	ijwt.Handler
	paths []string
}

func NewLoginJWTMiddlewareBuilder(hdl ijwt.Handler) *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{
		Handler: hdl,
	}
}

func (m *LoginJWTMiddlewareBuilder) IgnorePaths(path string) *LoginJWTMiddlewareBuilder {
	m.paths = append(m.paths, path)
	return m
}

func (m *LoginJWTMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	// 注册当前时间
	gob.Register(time.Now())

	return func(ctx *gin.Context) {
		for _, path := range m.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		// 根据约定，token 在 Authorization 头部
		// Bearer XXXX

		tokenStr := m.ExtractToken(ctx)

		var uc ijwt.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return ijwt.JWTKey, nil
		})
		if err != nil {
			// token 不对，token 是伪造的
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid {
			// token 解析出来了，但是 token 可能是非法的，或者过期了的
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		err = m.CheckSession(ctx, uc.Ssid)
		if err != nil {
			// token 无效或者 redis 有问题
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//if uc.UserAgent != ctx.GetHeader("User-Agent") {
		//	// 使用user agent 防止 token 被劫持
		//	// 后期我们讲到了监控告警的时候，这个地方要埋点
		//	// 能够进来这个分支的，大概率是攻击者
		//	ctx.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}

		//expireTime := uc.ExpiresAt
		// 不判定都可以
		//if expireTime.Before(time.Now()) {
		//	ctx.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}
		// 剩余过期时间 < 50s 就要刷新
		//if expireTime.Sub(time.Now()) < time.Second*50 {
		//	println("refresh token")
		//	uc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute * 1))
		//	tokenStr, err = token.SignedString(web.JWTKey)
		//	ctx.Header("x-jwt-token", tokenStr)
		//	if err != nil {
		//		// 这边不要中断，因为仅仅是过期时间没有刷新，但是用户是登录了的
		//		log.Println(err)
		//	}
		//}

		ctx.Set("user", uc)
	}
}
