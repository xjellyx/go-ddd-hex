package xgin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/olongfen/go-ddd-hex/internal/application"
	"github.com/olongfen/go-ddd-hex/internal/domain/vo"
	"github.com/olongfen/go-ddd-hex/lib/response"
)

type UserExecCtl struct {
	domain application.UserInterface
}

// ChangePasswd .
// @Tags User用户
// @Summary 修改密码
// @Description 参数传递新旧密码
// @Accept json
// @Produce json
// @Param oldPasswd body string true "旧密码"
// @Param newPasswd body string true "新密码"
// @Success 200  {object} response.Response
// @Failure 400  {object} response.Response "jwt验证失败"
// @Router /api/users/change-passwd [put]
func (u *UserExecCtl) ChangePasswd(c *gin.Context) {
	var (
		res       interface{}
		err       error
		oldPasswd string
		newPasswd string
	)

	defer func() {
		if err != nil {
			response.NewGinResponse(c).Fail(response.CodeFail, err).Response()
		} else {
			response.NewGinResponse(c).Success(res).Response()
		}
	}()
	id := c.GetString("id")
	oldPasswd = c.Param("oldPasswd")
	newPasswd = c.Param("newPasswd")
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

// Create .
// @Tags User用户
// @Summary 创建user记录
// @Description 参数是一个数组对象
// @Accept json
// @Produce json
// @Param [object] body []vo.UserVOForm true "表单数组"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response "jwt验证失败"
// @Router /api/v1/users/   [post]
func (u *UserExecCtl) Create(c *gin.Context) {
	var (
		res  interface{}
		form []*vo.UserVOForm
		err  error
	)
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
	data, _err := u.domain.Create(c.Request.Context(), form)
	if _err != nil {
		err = _err
		return
	}

	res = data
}
