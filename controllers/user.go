package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego/validation"

	"github.com/hashwing/pet-adoption/pkg/auth"
	"github.com/hashwing/pet-adoption/pkg/common"
	"github.com/hashwing/pet-adoption/pkg/config"
	"github.com/hashwing/pet-adoption/pkg/storage/db"
	"github.com/hashwing/pet-adoption/pkg/wx"
)

type UserController struct {
	BaseController
}

func (c *UserController) Put() {
	defer c.ServeJSON()
	var user db.User
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &user)
	if err != nil {
		c.SetErrMsg(400, err.Error())
		return
	}
	user.ID = c.GetUID()
	err = db.UpdateUser(user)
	if err != nil {
		c.SetErrMsg(500, err.Error())
		return
	}
	c.SetResult(nil, nil, 204)
}

func (c *UserController) WxLogin() {
	defer c.ServeJSON()
	jsCode := c.GetString("js_code")

	valid := validation.Validation{}
	valid.Required(jsCode, "name")
	if valid.HasErrors() {
		c.SetErrMsg(400, valid.Errors[0].Message)
		return
	}

	s, err := wx.Code2Session(config.Cfg.WxConfig.AppID, config.Cfg.WxConfig.Secret, jsCode)
	if err != nil {
		c.SetErrMsg(401, err.Error())
		return
	}

	user, isExist, err := db.UserExistByOpenID(s.OpenID)
	if err != nil {
		c.SetErrMsg(500, err.Error())
		return
	}
	if isExist {
		token := auth.GetToken(user.ID)
		c.SetResult(nil, token, 200, "token")
		return
	}

	userID := common.NewUUID()
	userinfo := db.User{
		ID:     userID,
		OpenID: s.OpenID,
	}
	err = db.AddUser(userinfo)
	if err != nil {
		c.SetErrMsg(500, err.Error())
		return
	}
	token := auth.GetToken(userID)
	c.SetResult(nil, token, 200, "token")
}
