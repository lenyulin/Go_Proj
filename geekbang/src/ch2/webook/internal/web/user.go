package web

import (
	"ch2/webook/internal/domain"
	"ch2/webook/internal/service"
	"encoding/json"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	emailRegexPattern = `^[A-Za-z0-9\u4e00-\u9fa5]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`
	pwdRegexPattern   = `^(?=.*[0-9])(?=.*[A-Z])(?=.*[a-z])(?=.*[!@#$%^&*,\._])[0-9a-zA-Z!@#$%^&*,\\._]{8,12}$`
)

type UserHandler struct {
	emailRexExp    *regexp.Regexp
	passwordRexExp *regexp.Regexp
	svc            *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		emailRexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRexExp: regexp.MustCompile(pwdRegexPattern, regexp.None),
		svc:            svc,
	}
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	server.POST("/users/signup", h.SignUp)
	server.POST("/users/login", h.Login)
	server.POST("/users/edit", h.Edit)
	server.GET("/users/profile", h.Profile)
}

func (h *UserHandler) SignUp(ctx *gin.Context) {
	type Signup struct {
		Email           string `json:"email"`
		Password        string `json:"pwd"`
		ConfirmPassword string `json:"confirmPwd"`
	}
	var req Signup
	if err := ctx.Bind(&req); err != nil {
		return
	}
	isEmail, err := h.emailRexExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "System timeout")
		return
	}
	if !isEmail {
		ctx.String(http.StatusOK, "Email format error")
		return
	}
	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "Password and ConfirmPassword not equal")
		return
	}
	isPwd, err := h.passwordRexExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "System timeout")
		return
	}
	if !isPwd {
		ctx.String(http.StatusOK, "Password format error")
		return
	}
	err = h.svc.Signup(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	switch err {
	case nil:
		ctx.String(http.StatusOK, "Register Succeed")
	case service.ErrduplicateEmail:
		ctx.String(http.StatusOK, "Email exist")
	default:
		ctx.String(http.StatusOK, "System Failure")
	}
}

func (h *UserHandler) Login(ctx *gin.Context) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	u, err := h.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		sess := sessions.Default(ctx)
		sess.Set("userId", u.Id)
		sess.Options(sessions.Options{
			MaxAge: 900,
		})
		err = sess.Save()
		if err != nil {
			ctx.String(http.StatusOK, "System Error")
		}
		ctx.String(http.StatusOK, "Login Succeed")
	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "User Or Password Error")
	default:
		ctx.String(http.StatusOK, "System Error")
	}
}
func (h *UserHandler) Edit(ctx *gin.Context) {
	var user domain.User
	if err := ctx.Bind(&user); err != nil {
		ctx.String(http.StatusOK, "System Error")
		return
	}
	if len(user.NickName) == 0 || len(user.NickName) > 30 {
		ctx.String(http.StatusOK, "NickName cannot be less than 1 or longer than 30")
		return
	}
	if len(user.Introduction) == 0 || len(user.Introduction) > 100 {
		ctx.String(http.StatusOK, "Introduction cannot be less than 1 or longer than 100")
		return
	}
	err := h.svc.Edit(ctx, user)
	if err != nil {
		ctx.String(http.StatusOK, "System Error")
		return
	}
	ctx.String(http.StatusOK, "Edit Succeed")
}
func (h *UserHandler) Profile(ctx *gin.Context) {
	type Req struct {
		Email string `json:"email"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "System Error")
		return
	}
	req.Email = "12345678@qq.com"
	if req.Email == "" {
		ctx.String(http.StatusOK, "Email cannot be empty")
		return
	}
	u, err := h.svc.Profile(ctx, req.Email)
	if err != nil {
		return
	}
	us, err := json.Marshal(u)
	if err != nil {
		return
	}
	ctx.String(http.StatusOK, string(us))
}
