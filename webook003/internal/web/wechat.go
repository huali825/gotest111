package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"goworkwebook/webook003/internal/service"
	"goworkwebook/webook003/internal/service/oauth2/wechat"
	ijwt "goworkwebook/webook003/internal/web/jwt"

	//导入一个名为shortuuid的Go库，该库用于生成简洁的UUID。UUID（通用唯一标识符）
	uuid "github.com/lithammer/shortuuid/v4"
	"net/http"
)

type OAuth2WechatHandler struct {
	svc     wechat.Service
	userSvc service.UserService
	ijwt.Handler
	key             []byte
	stateCookieName string
}

func NewOAuth2WechatHandler(
	svc wechat.Service,
	hdl ijwt.Handler,
	userSvc service.UserService) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:             svc,
		userSvc:         userSvc,
		key:             []byte("k6CswdUm77WKcbM68UQUuxVsHSpTCwgB"),
		stateCookieName: "jwt-state",
		Handler:         hdl,
	}
}

func (o *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	g.GET("/authurl", o.Auth2URL)
	g.Any("/callback", o.Callback)
}

// Auth2URL 构建扫码登录url
func (o *OAuth2WechatHandler) Auth2URL(ctx *gin.Context) {
	//导入一个名为short uuid的Go库，该库用于生成简洁的UUID。UUID（通用唯一标识符）
	state := uuid.New()
	//调用svc的AuthURL方法，传入state，获取跳转URL
	val, err := o.svc.AuthURL(ctx, state)
	if err != nil {
		//如果获取跳转URL失败，返回错误信息
		ctx.JSON(http.StatusOK, Result{
			Msg:  "构造跳转URL失败",
			Code: 5,
		})
		return
	}
	//调用setStateCookie方法，将state存入cookie中
	err = o.setStateCookie(ctx, state)
	if err != nil {
		//如果存入cookie失败，返回错误信息
		ctx.JSON(http.StatusOK, Result{
			Msg:  "服务器异常",
			Code: 5,
		})
	}
	//返回跳转URL
	ctx.JSON(http.StatusOK, Result{
		Data: val,
	})
}

// Callback 函数用于处理微信OAuth2回调请求
func (o *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	// 验证state参数，防止CSRF攻击
	err := o.verifyState(ctx)
	if err != nil {
		// 如果验证失败，返回错误信息
		ctx.JSON(http.StatusOK, Result{
			Msg:  "非法请求",
			Code: 4,
		})
		return
	}
	// 获取code参数，用于获取微信用户信息
	code := ctx.Query("code")
	// 获取state参数，用于防止CSRF攻击
	// state := ctx.Query("state")
	// 使用code参数获取微信用户信息
	wechatInfo, err := o.svc.VerifyCode(ctx, code)
	if err != nil {
		// 如果获取失败，返回错误信息
		ctx.JSON(http.StatusOK, Result{
			Msg:  "授权码有误",
			Code: 4,
		})
		return
	}
	// 根据微信用户信息查找或创建用户
	u, err := o.userSvc.FindOrCreateByWechat(ctx, wechatInfo)
	if err != nil {
		// 如果查找或创建失败，返回错误信息
		ctx.JSON(http.StatusOK, Result{
			Msg:  "系统错误",
			Code: 5,
		})
		return
	}
	// 设置登录token
	err = o.SetLoginToken(ctx, u.Id)
	if err != nil {
		// 如果设置失败，返回错误信息
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	// 返回成功信息
	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
	return
}

func (o *OAuth2WechatHandler) verifyState(ctx *gin.Context) error {
	state := ctx.Query("state")
	ck, err := ctx.Cookie(o.stateCookieName)
	if err != nil {
		return fmt.Errorf("无法获得 cookie %w", err)
	}
	var sc StateClaims
	_, err = jwt.ParseWithClaims(ck, &sc, func(token *jwt.Token) (interface{}, error) {
		return o.key, nil
	})
	if err != nil {
		return fmt.Errorf("解析 token 失败 %w", err)
	}
	if state != sc.State {
		// state 不匹配，有人搞你
		return fmt.Errorf("state 不匹配")
	}
	return nil
}

func (o *OAuth2WechatHandler) setStateCookie(ctx *gin.Context,
	state string) error {
	claims := StateClaims{
		State: state,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(o.key)
	if err != nil {

		return err
	}
	ctx.SetCookie(o.stateCookieName, tokenStr,
		600, "/oauth2/wechat/callback",
		"", false, true)
	return nil
}

type StateClaims struct {
	jwt.RegisteredClaims
	State string
}
