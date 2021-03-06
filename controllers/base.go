package controllers

import (
	"errors"

	"github.com/astaxie/beego"
)

// ErrorMsg the struct of error message
type ErrorMsg struct {
	Code    int    `json:"err_code"`
	Message string `json:"message"`
}

// BaseController the base controller
type BaseController struct {
	beego.Controller
}

// SetErrMsg return err msg to http
func (c *BaseController) SetErrMsg(code int, msg string) {
	errMsg := ErrorMsg{
		Code:    code,
		Message: msg,
	}
	c.SetResult(errors.New("err"), errMsg, code)
}

// SetResult return json to http
func (c *BaseController) SetResult(err error, result interface{}, errcode int, key ...string) {
	c.Ctx.Output.Status = errcode
	if err != nil {
		c.Data["json"] = result
		return
	}
	if result == nil && (len(key) == 0 || key[0] == "") {
		return
	}

	if len(key) == 0 || key[0] == "" {
		c.Data["json"] = result
	} else {
		c.Data["json"] = map[string]interface{}{key[0]: result}
	}
}

// get user id
func (c *BaseController) GetUID() string {
	uid := c.Ctx.Input.GetData("uid")
	return uid.(string)
}
