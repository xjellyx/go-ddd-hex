package handler

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/olongfen/go-ddd-hex/internal/application"
	"github.com/olongfen/go-ddd-hex/internal/domain/vo"
	"github.com/olongfen/go-ddd-hex/lib/response"
)

type userExecCtl struct {
	domain application.UserServiceInterface
}

// ChangePasswd .
// @Tags User用户
// @Summary 修改密码
// @Description 参数传递新旧密码
// @Accept json
// @Produce json
// @Param old_passwd body string true "旧密码"
// @Param new_passwd body string true "新密码"
// @Success 200  {object} response.Response
// @Failure 400  {object} response.Response "jwt验证失败"
// @Router /api/users/change-passwd [put]
func (u *userExecCtl) ChangePasswd(ctx context.Context) {
	var (
		res       interface{}
		err       error
		oldPasswd string
		newPasswd string
	)

	c := ctx.(*gin.Context)
	defer func() {
		if err != nil {
			response.NewGinResponse(c).Fail(response.CodeFail, err).Response()
		} else {
			response.NewGinResponse(c).Success(res).Response()
		}
	}()
	id := c.GetString("id")
	oldPasswd = c.Param("old_passwd")
	newPasswd = c.Param("new_passwd")
	if len(oldPasswd) == 0 {
		err = errors.New("old password invalid")
		return
	}
	if len(newPasswd) == 0 {
		err = errors.New("new password invalid")
		return
	}
	if err = u.domain.ChangePassword(c.Request.Context(), id, oldPasswd, newPasswd); err != nil {
		return
	}

}

// Register .
// @Tags User用户
// @Summary 创建user记录
// @Description 参数是一个数组对象
// @Accept json
// @Produce json
// @Param [object] body vo.RegisterForm true "表单数组"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response "jwt验证失败"
// @Router /api/v1/users/   [post]
func (u *userExecCtl) Register(ctx context.Context) {
	var (
		res  interface{}
		form vo.RegisterForm
		err  error
	)
	c := ctx.(*gin.Context)
	defer func() {
		if err != nil {
			response.NewGinResponse(c).Fail(response.CodeFail, err).Response()
		} else {
			response.NewGinResponse(c).Success(res).Response()
		}
	}()
	if err = c.ShouldBind(&form); err != nil {
		return
	}
	_err := u.domain.Register(c.Request.Context(), form)
	if _err != nil {
		err = _err
		return
	}
}
