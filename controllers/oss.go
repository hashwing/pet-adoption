package controllers

import (
	"github.com/hashwing/log"
	"github.com/hashwing/pet-adoption/pkg/config"
	"github.com/hashwing/pet-adoption/pkg/storage/oss"
)

//OssController oss controller handler
type OssController struct {
	BaseController
}

func (c *OssController) GetToken() {
	defer c.ServeJSON()
	token := oss.GetPolicyToken(config.Cfg.OssConfig.AccessKeyId,
		config.Cfg.OssConfig.AccessKeySecret,
		config.Cfg.OssConfig.Host,
		config.Cfg.OssConfig.CallbackUrl,
		config.Cfg.OssConfig.UploadDir,
		config.Cfg.OssConfig.ExpireTime,
	)
	c.SetResult(nil, token, 200)
}

func (c *OssController) GetImg() {
	key := c.GetString("key")
	data, err := oss.Get(
		config.Cfg.OssConfig.InHost,
		config.Cfg.OssConfig.AccessKeyId,
		config.Cfg.OssConfig.AccessKeySecret,
		config.Cfg.OssConfig.UploadDir,
		config.Cfg.OssConfig.UploadDir+"/"+key,
	)
	if err != nil {
		log.Error(err)
		return
	}
	c.Ctx.ResponseWriter.Write(data)
}
