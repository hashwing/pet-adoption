package controllers

import (
	"github.com/hashwing/log"
	"github.com/hashwing/pet-adoption/pkg/storage/db"
)

//LocalityController locality controller handler
type LocalityController struct {
	BaseController
}

// FindProvinces find all provinces
func (c *LocalityController) FindProvinces() {
	defer c.ServeJSON()
	provinces, err := db.FindProvinces()
	if err != nil {
		log.Error(err)
		c.SetErrMsg(500, "数据库错误")
	}

	c.SetResult(nil, provinces, 200)
}

// FindCitys find all citys
func (c *LocalityController) FindCitys() {
	defer c.ServeJSON()
	provinceID := c.GetString("province_id")
	if provinceID == "" {
		c.SetErrMsg(400, "请求参数错误")
		return
	}
	citys, err := db.FindCitysByProvinceID(provinceID)
	if err != nil {
		log.Error(err)
		c.SetErrMsg(500, "数据库错误")
	}

	c.SetResult(nil, citys, 200)
}

// FindLocalities find all localities
func (c *LocalityController) FindLocalities() {
	defer c.ServeJSON()
	cityID := c.GetString("city_id")
	if cityID == "" {
		c.SetErrMsg(400, "请求参数错误")
		return
	}
	localities, err := db.FindLocalitiesByCityID(cityID)
	if err != nil {
		log.Error(err)
		c.SetErrMsg(500, "数据库错误")
		return
	}

	c.SetResult(nil, localities, 200)
}

type provinceItem struct {
	UUID  string     `json:"province_uuid"`
	Name  string     `json:"province_name"`
	Citys []cityItem `json:"citys"`
}

type cityItem struct {
	UUID       string         `json:"city_uuid"`
	Name       string         `json:"city_name"`
	Localities []localityItem `json:"localities"`
}

type localityItem struct {
	UUID string `json:"locality_uuid"`
	Name string `json:"locality_name"`
}

func (c *LocalityController) Find() {
	defer c.ServeJSON()
	provinces, err := db.FindProvinces()
	if err != nil {
		log.Error(err)
		c.SetErrMsg(500, "数据库错误")
		return
	}

	citys, err := db.FindCitys()
	if err != nil {
		log.Error(err)
		c.SetErrMsg(500, "数据库错误")
		return
	}
	localities, err := db.FindLocalities()
	if err != nil {
		log.Error(err)
		c.SetErrMsg(500, "数据库错误")
		return
	}

	data := make([]provinceItem, 0)

	for _, p := range provinces {
		pr := provinceItem{
			UUID:  p.ID,
			Name:  p.Name,
			Citys: make([]cityItem, 0),
		}

		for _, c := range citys {
			if p.ID == c.ProvinceID {
				cr := cityItem{
					UUID:       c.ID,
					Name:       c.Name,
					Localities: make([]localityItem, 0),
				}
				for _, l := range localities {
					if l.CityID == c.ID {
						lr := localityItem{
							UUID: l.ID,
							Name: l.Name,
						}
						cr.Localities = append(cr.Localities, lr)
					}
				}
				pr.Citys = append(pr.Citys, cr)
			}

		}
		data = append(data, pr)
	}

	c.SetResult(nil, data, 200)

}
