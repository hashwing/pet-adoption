package controllers

import (
	"encoding/json"

	"github.com/hashwing/log"
	"github.com/hashwing/pet-adoption/pkg/common"
	"github.com/hashwing/pet-adoption/pkg/storage/db"
)

// AdoptionController adoption
type AdoptionController struct {
	BaseController
}

func (c *AdoptionController) FindPetClass() {
	defer c.ServeJSON()
	pcs, err := db.FindPetClass()
	if err != nil {
		log.Error(err)
		c.SetErrMsg(500, "获取类型失败")
		return
	}
	c.SetResult(nil, pcs, 200)
}

//PublicList
func (c *AdoptionController) PublicList() {
	defer c.ServeJSON()
	cityID := c.GetString("city_id")
	localityID := c.GetString("locality_id")
	petClassID := c.GetString("petClass_id")
	key := c.GetString("key")
	page, err := c.GetInt("page")
	if err != nil {
		log.Error(err)
		c.SetErrMsg(400, "请求参数错误")
		return
	}
	size, err := c.GetInt("page_size")
	if err != nil {
		log.Error(err)
		c.SetErrMsg(400, "请求参数错误")
		return
	}
	start := size * (page - 1)
	ps, err := db.FindPetPublics(cityID, localityID, petClassID, key, start, size)
	if err != nil {
		log.Error(err)
		c.SetErrMsg(500, "数据库错误")
		return
	}
	for i, p := range ps {
		u, err := db.GetUser(p.UserID)
		if err != nil {
			log.Error(err)
			continue
		}
		ps[i].User = *u
	}
	c.SetResult(nil, ps, 200)
}

//
func (c *AdoptionController) PublicListByUser() {
	defer c.ServeJSON()
	uid := c.GetUID()
	ps, err := db.FindPetPublicsByUser(uid)
	if err != nil {
		log.Error(err)
	}
	c.SetResult(nil, ps, 200)
}

func (c *AdoptionController) CreatePublic() {
	defer c.ServeJSON()
	var pp db.PetPublic
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &pp)
	if err != nil {
		c.SetErrMsg(400, err.Error())
		return
	}
	pp.ID = common.NewUUID()
	pp.UserID = c.GetUID()
	pp.State = db.PetPublicState

	err = db.CreatePetPublics(pp)
	if err != nil {
		c.SetErrMsg(500, err.Error())
		return
	}

	c.SetResult(nil, nil, 204)
}

func (c *AdoptionController) GetPublic() {
	defer c.ServeJSON()
	uuid := c.Ctx.Input.Param(":uuid")
	p, err := db.GetPetPublic(uuid)
	if err != nil {
		c.SetErrMsg(500, err.Error())
		return
	}
	u, err := db.GetUser(p.UserID)
	if err != nil {
		log.Error(err)
	}
	p.User = *u
	c.SetResult(nil, p, 200)
}

func (c *AdoptionController) UpdatePublic() {
	defer c.ServeJSON()
	var pp db.PetPublic
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &pp)
	if err != nil {
		c.SetErrMsg(400, err.Error())
		return
	}
	pp.ID = c.Ctx.Input.Param(":uuid")
	if pp.ID == "" {
		c.SetErrMsg(400, "id 不能为空")
		return
	}
	pp.UserID = c.GetUID()

	err = db.UpdatePetPublics(pp)
	if err != nil {
		c.SetErrMsg(500, err.Error())
		return
	}

	c.SetResult(nil, nil, 204)
}

func (c *AdoptionController) DeletePublic() {
	defer c.ServeJSON()
	uuid := c.Ctx.Input.Param(":uuid")
	if uuid == "" {
		c.SetErrMsg(400, "参数错误")
		return
	}
	uid := c.GetUID()
	err := db.DelPetPublics(uid, uuid)
	if err != nil {
		log.Error(err)
		c.SetErrMsg(500, err.Error())
		return
	}
	c.SetResult(nil, nil, 204)
}

func (c *AdoptionController) ApplyListByUser() {
	defer c.ServeJSON()
	uid := c.GetUID()
	applys, err := db.FindAdoptionApplyByUserID(uid)
	if err != nil {
		c.SetErrMsg(500, err.Error())
		return
	}
	for i, ap := range applys {
		pp, err := db.GetPetPublic(ap.PetID)
		if err != nil {
			log.Error(err)
			continue
		}
		applys[i].Pet = *pp
	}
	c.SetResult(nil, applys, 200)
}

func (c *AdoptionController) ApplyListByPet() {
	defer c.ServeJSON()
	petID := c.Ctx.Input.Param(":pet_id")
	if petID == "" {
		c.SetErrMsg(400, "pet_id 不能为空")
		return
	}
	applys, err := db.FindAdoptionApplyByPetID(petID)
	if err != nil {
		log.Error(err)
		c.SetErrMsg(500, err.Error())
		return
	}

	for i, ap := range applys {
		u, err := db.GetUser(ap.UserID)
		if err != nil {
			log.Error(err)
			continue
		}
		applys[i].User = *u
	}
	c.SetResult(nil, applys, 200)
}

func (c *AdoptionController) GetApply() {
	defer c.ServeJSON()
	uuid := c.Ctx.Input.Param(":uuid")
	if uuid == "" {
		c.SetErrMsg(400, "uuid 不能为空")
		return
	}
	apply, err := db.GetAdoptionApply(uuid)
	if err != nil {
		log.Error(err)
		c.SetErrMsg(500, err.Error())
		return
	}
	pp, err := db.GetPetPublic(apply.PetID)
	if err != nil {
		log.Error(err)
	}
	apply.Pet = *pp

	c.SetResult(nil, apply, 200)
}

func (c *AdoptionController) CreateApply() {
	defer c.ServeJSON()
	uid := c.GetUID()
	petID := c.Ctx.Input.Param(":pet_id")
	var apply db.AdoptionApply
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &apply)
	if err != nil {
		c.SetErrMsg(400, err.Error())
		return
	}
	apply.PetID = petID
	if apply.PetID == "" {
		c.SetErrMsg(400, "pet_id 不能为空")
		return
	}
	apply.ID = common.NewUUID()
	apply.UserID = uid
	err = db.CreateAdoptionApply(apply)
	if err != nil {
		c.SetErrMsg(500, err.Error())
		return
	}
	c.SetResult(nil, nil, 204)
}

func (c *AdoptionController) UpdateApply() {
	defer c.ServeJSON()
	uuid := c.Ctx.Input.Param(":uuid")
	petID := c.Ctx.Input.Param(":pet_id")
	var apply db.AdoptionApply
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &apply)
	if err != nil {
		c.SetErrMsg(400, err.Error())
		return
	}
	apply.PetID = petID
	if apply.PetID == "" {
		c.SetErrMsg(400, "pet_id 不能为空")
		return
	}
	apply.ID = uuid
	err = db.UpdateAdoptionApply(apply)
	if err != nil {
		log.Error(err)
		c.SetErrMsg(500, err.Error())
		return
	}
	c.SetResult(nil, nil, 204)
}

func (c *AdoptionController) DelApply() {
	defer c.ServeJSON()
	uuid := c.Ctx.Input.Param(":uuid")
	uid := c.GetUID()
	err := db.DelAdoptionApply(uid, uuid)
	if err != nil {
		c.SetErrMsg(500, err.Error())
		return
	}
	c.SetResult(nil, nil, 204)
}
