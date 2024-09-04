package web

import (
	"errors"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"goworkwebook/webook003/internal/domain"
	"goworkwebook/webook003/internal/service"
	"net/http"
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
	u.s.POST("/users/login", u.Login)
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
		MaxAge: 10,
		//MaxAge: 3600 * 24 * 7,
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
