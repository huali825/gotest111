package web

import (
	"errors"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"goworkwebook/webook003/internal/domain"
	"goworkwebook/webook003/internal/service"
	"net/http"
	"time"
)

type UserHandler struct {
	s                *gin.Engine
	svc              *service.UserService
	emailRegexRxp    *regexp.Regexp
	passwordRegexRxp *regexp.Regexp
}

func NewUserHandler(server *gin.Engine, svc *service.UserService) *UserHandler {
	const (
		emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	)
	return &UserHandler{
		s:                server,
		svc:              svc,
		emailRegexRxp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegexRxp: regexp.MustCompile(passwordRegexPattern, regexp.None),
	}
}

func (u *UserHandler) RegisterRoutes() {
	u.s.POST("/users/signup", u.SignUp)
	u.s.POST("/users/login", u.LoginJWT)
	//u.s.POST("/users/login", u.Login)
	u.s.POST("/users/edit", u.Edit)
	u.s.GET("/users/profile", u.Profile)

}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req SignUpReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	err := u.emailPasswordFormat(ctx, req.Email, req.Password, req.ConfirmPassword)
	if err != nil {
		return
	}

	// 调用一下 svc 的方法
	err = u.svc.SignUp(ctx, domain.DMUser{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrUserDuplicateEmail) {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}

	ctx.String(http.StatusOK, "注册成功")

	//context.String(http.StatusOK, "Hello, this is signup")
}

func (u *UserHandler) LoginJWT(ctx *gin.Context) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	DmUser, err := u.svc.Login(ctx, req.Email, req.Password)

	switch err {
	case nil:
		// 创建用户声明
		uc := UserClaims{
			Uid:       DmUser.Id,                   // 用户ID
			UserAgent: ctx.GetHeader("User-Agent"), //   使用user agent 防止 token 被劫持
			RegisteredClaims: jwt.RegisteredClaims{
				// 1 分钟过期
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 1)), // 过期时间
			},
		}
		// 创建JWT
		token := jwt.NewWithClaims(jwt.SigningMethodHS512, uc)
		// 签名JWT 加密
		tokenStr, err := token.SignedString(JWTKey)
		if err != nil {
			// 签名失败，返回系统错误
			ctx.String(http.StatusOK, "系统错误")
		}
		// 设置JWT头部
		ctx.Header("x-jwt-token", tokenStr)
		// 返回登录成功
		ctx.String(http.StatusOK, "登录成功")
	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "用户名或者密码不对")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
}

var JWTKey = []byte("k6CswdUm77WKcbM68UQUuxVsHSpTCwgK")

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	UserAgent string
}

func (u *UserHandler) Login(context *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := context.ShouldBindJSON(&req); err != nil {
		context.String(http.StatusOK, "系统错误")
		return
	}

	user, err := u.svc.Login(context, req.Email, req.Password)
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		context.String(http.StatusOK, "用户名或密码错误")
		return
	}
	if err != nil {
		context.String(http.StatusOK, "系统错误")
		return
	}

	//登录成功后设置userid的值
	sess := sessions.Default(context) //获取session
	sess.Set("userId", user.Id)       //设置session
	sess.Options(sessions.Options{
		// 设置cookie的路径
		Path: "",
		// 设置cookie的域名
		Domain: "",
		// 设置cookie的最大有效期，单位为秒
		MaxAge: 20, //单位秒
		//MaxAge: 3600 * 24 * 7, //一周过期
		// 设置cookie是否只在https协议下有效
		Secure: false,
		// 设置cookie是否只能通过http协议访问
		HttpOnly: false,
		// 设置cookie的SameSite属性，0表示不限制
		SameSite: 0,
	})
	err = sess.Save() //保存session
	if err != nil {
		context.String(http.StatusOK, "系统错误")
		return
	}

	//登录成功返回 状态200
	context.String(http.StatusOK, "登录成功")
	return
}

func (u *UserHandler) Edit(context *gin.Context) {
	context.String(http.StatusOK, "未实现!!!!")
	return
}

func (u *UserHandler) Profile(context *gin.Context) {
	context.String(http.StatusOK, "这是 profile")
	return
}

// 检验邮箱密码格式
func (u *UserHandler) emailPasswordFormat(
	ctx *gin.Context, email string, password string, password2 string) error {
	isEmail, err := u.emailRegexRxp.MatchString(email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return errors.New("系统错误")
	}
	if !isEmail {
		ctx.String(http.StatusOK, "邮箱格式错误")
		return errors.New("邮箱格式错误")
	}
	if password != password2 {
		ctx.String(http.StatusOK, "两次密码不一致")
		return errors.New("两次密码不一致")
	}
	isPassword, err := u.passwordRegexRxp.MatchString(password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return errors.New("系统错误")
	}
	if !isPassword {
		ctx.String(http.StatusOK, "密码格式错误")
		return errors.New("密码格式错误")
	}
	return nil
}
