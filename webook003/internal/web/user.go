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

const (
	emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	// 和上面比起来，用 ` 看起来就比较清爽
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	bizLogin             = "login"
)

type UserHandler struct {
	svc              service.UserService
	codeSvc          service.CodeService
	emailRegexRxp    *regexp.Regexp
	passwordRegexRxp *regexp.Regexp
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeService) *UserHandler {
	return &UserHandler{
		svc:              svc,
		codeSvc:          codeSvc,
		emailRegexRxp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegexRxp: regexp.MustCompile(passwordRegexPattern, regexp.None),
	}
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")

	ug.POST("/signup", h.SignUp)
	ug.POST("/login", h.LoginJWT)
	//ug.POST("/users/login", h.Login)
	ug.POST("/edit", h.Edit)
	ug.GET("/profile", h.Profile)

	// 手机验证码登录相关功能
	ug.POST("/login_sms/code/send", h.SendSMSLoginCode)
	ug.POST("/login_sms", h.LoginSMS)
}

func (h *UserHandler) SignUp(ctx *gin.Context) {
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
	err := h.emailPasswordFormat(ctx, req.Email, req.Password, req.ConfirmPassword)
	if err != nil {
		return
	}

	// 调用一下 svc 的方法
	err = h.svc.Signup(ctx, domain.DMUser{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrDuplicateEmail) {
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

func (h *UserHandler) LoginJWT(ctx *gin.Context) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	DmUser, err := h.svc.Login(ctx, req.Email, req.Password)

	switch err {
	case nil:
		h.setJWTToken(ctx, DmUser.Id)
		// 返回登录成功
		ctx.String(http.StatusOK, "登录成功")
	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "用户名或者密码不对")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
}

func (h *UserHandler) setJWTToken(ctx *gin.Context, uid int64) {
	// 创建用户声明
	uc := UserClaims{
		Uid:       uid,                         // 用户ID
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
}

var JWTKey = []byte("k6CswdUm77WKcbM68UQUuxVsHSpTCwgK")

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	UserAgent string
}

func (h *UserHandler) Login(context *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := context.ShouldBindJSON(&req); err != nil {
		context.String(http.StatusOK, "系统错误")
		return
	}

	user, err := h.svc.Login(context, req.Email, req.Password)
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

func (h *UserHandler) Edit(ctx *gin.Context) {

	// 嵌入一段刷新过期时间的代码
	type Req struct {
		// 改邮箱，密码，或者能不能改手机号

		Nickname string `json:"nickname"`
		// YYYY-MM-DD
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	//sess := sessions.Default(ctx)
	//sess.Get("uid")
	uc, ok := ctx.MustGet("user").(UserClaims)
	if !ok {
		//ctx.String(http.StatusOK, "系统错误")
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// 用户输入不对
	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		//ctx.String(http.StatusOK, "系统错误")
		ctx.String(http.StatusOK, "生日格式不对")
		return
	}
	err = h.svc.UpdateNonSensitiveInfo(ctx, domain.DMUser{
		Id:       uc.Uid,
		Nickname: req.Nickname,
		Birthday: birthday,
		AboutMe:  req.AboutMe,
	})
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	ctx.String(http.StatusOK, "更新成功")
}

func (h *UserHandler) Profile(ctx *gin.Context) {
	//us := ctx.MustGet("user").(UserClaims)
	//ctx.String(http.StatusOK, "这是 profile")
	// 嵌入一段刷新过期时间的代码

	uc, ok := ctx.MustGet("user").(UserClaims)
	if !ok {
		//ctx.String(http.StatusOK, "系统错误")
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	u, err := h.svc.FindById(ctx, uc.Uid)
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	type User struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		AboutMe  string `json:"aboutMe"`
		Birthday string `json:"birthday"`
	}
	ctx.JSON(http.StatusOK, User{
		Nickname: u.Nickname,
		Email:    u.Email,
		AboutMe:  u.AboutMe,
		Birthday: u.Birthday.Format(time.DateOnly),
	})
}

func (h *UserHandler) SendSMSLoginCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 你这边可以校验 Req
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "请输入手机号码",
		})
		return
	}
	err := h.codeSvc.Send(ctx, bizLogin, req.Phone)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送成功",
		})
	case service.ErrCodeSendTooMany:
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "短信发送太频繁，请稍后再试",
		})
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		// 补日志的
	}
}

// LoginSMS 使用短信验证码登录
func (h *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	ok, err := h.codeSvc.Verify(ctx, bizLogin, req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统异常",
		})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码不对，请重新输入",
		})
		return
	}
	dmUser, err := h.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	h.setJWTToken(ctx, dmUser.Id)
	ctx.JSON(http.StatusOK, Result{
		Msg: "登录成功",
	})
}

// 检验邮箱密码格式
func (h *UserHandler) emailPasswordFormat(
	ctx *gin.Context, email string, password string, password2 string) error {
	isEmail, err := h.emailRegexRxp.MatchString(email)
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
	isPassword, err := h.passwordRegexRxp.MatchString(password)
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
