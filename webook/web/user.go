package web

import (
	"errors"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	s                *gin.Engine
	emailRegexRxp    *regexp.Regexp
	passwordRegexRxp *regexp.Regexp
}

func NewUserHandler(server *gin.Engine) *UserHandler {
	const (
		emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	)
	return &UserHandler{
		s:                server,
		emailRegexRxp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegexRxp: regexp.MustCompile(passwordRegexPattern, regexp.None),
	}
}

func (h *UserHandler) RegisterRoutes() {
	h.s.POST("/users/signup", h.SignUp)
	h.s.POST("/users/login", h.Login)
	h.s.POST("/users/edit", h.Edit)
	h.s.POST("/users/profile", h.Profile)

}

func (h *UserHandler) SignUp(context *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req SignUpReq
	if err := context.ShouldBindJSON(&req); err != nil {
		context.String(http.StatusOK, "系统错误")
		return
	}
	err := h.emailPasswordFormat(context, req.Email, req.Password, req.ConfirmPassword)
	if err != nil {
		return
	}

	context.String(http.StatusOK, "注册成功")

	//context.String(http.StatusOK, "Hello, this is signup")
}

func (h *UserHandler) Login(context *gin.Context) {
	//TODO: implement
}

func (h *UserHandler) Edit(context *gin.Context) {
	//TODO: implement
}

func (h *UserHandler) Profile(context *gin.Context) {
	//TODO: implement
}

// 检验邮箱密码格式
func (h *UserHandler) emailPasswordFormat(
	context *gin.Context, email string, password string, password2 string) error {
	isEmail, err := h.emailRegexRxp.MatchString(email)
	if err != nil {
		context.String(http.StatusOK, "系统错误")
		return errors.New("系统错误")
	}
	if !isEmail {
		context.String(http.StatusOK, "邮箱格式错误")
		return errors.New("邮箱格式错误")
	}
	if password != password2 {
		context.String(http.StatusOK, "两次密码不一致")
		return errors.New("两次密码不一致")
	}
	isPassword, err := h.passwordRegexRxp.MatchString(password)
	if err != nil {
		context.String(http.StatusOK, "系统错误")
		return errors.New("系统错误")
	}
	if !isPassword {
		context.String(http.StatusOK, "密码格式错误")
		return errors.New("密码格式错误")
	}
	return nil
}
