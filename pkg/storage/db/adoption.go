package db

import (
	"errors"
	"time"

	"github.com/hashwing/log"
)

type Province struct {
	ID   string `json:"uuid" xorm:"'uuid'"`
	Name string `json:"name" xorm:"name"`
}

func AddProvince(p Province) error {
	var old Province
	isExit, err := MysqlDB.Where("name=?", p.Name).Get(&old)
	if isExit {
		return nil
	}
	if err != nil {
		return err
	}
	_, err = MysqlDB.Insert(p)
	return err
}

func FindProvinces() ([]Province, error) {
	var provinces []Province
	err := MysqlDB.Find(&provinces)
	return provinces, err
}

type City struct {
	ID         string `json:"uuid" xorm:"'uuid'"`
	ProvinceID string `json:"province_id" xorm:"province_id"`
	Name       string `json:"name" xorm:"name"`
}

func AddCity(c City) error {
	var old City
	isExit, err := MysqlDB.Where("name=? and province_id=?", c.Name, c.ProvinceID).Get(&old)
	if isExit {
		return nil
	}
	if err != nil {
		return err
	}
	_, err = MysqlDB.Insert(c)
	return err
}
func FindCitys() ([]City, error) {
	var citys []City
	err := MysqlDB.Find(&citys)
	return citys, err
}

func FindCitysByProvinceID(provinceID string) ([]City, error) {
	var citys []City
	err := MysqlDB.Where("province_id=?", provinceID).Find(&citys)
	return citys, err
}

// Locality
type Locality struct {
	ID     string `json:"uuid" xorm:"'uuid'"`
	CityID string `json:"city_id" xorm:"city_id"`
	Name   string `json:"name" xorm:"name"`
}

func AddLocality(l Locality) error {
	var old Locality
	isExit, err := MysqlDB.Where("name=? and city_id=?", l.Name, l.CityID).Get(&old)
	if isExit {
		return nil
	}
	if err != nil {
		return err
	}
	_, err = MysqlDB.Insert(l)
	return err
}

func FindLocalities() ([]Locality, error) {
	var localities []Locality
	err := MysqlDB.Find(&localities)
	return localities, err
}

func FindLocalitiesByCityID(cityID string) ([]Locality, error) {
	var localities []Locality
	err := MysqlDB.Where("city_id=?", cityID).Find(&localities)
	return localities, err
}

type PetClass struct {
	ID   string `json:"uuid" xorm:"'uuid'"`
	Name string `json:"name" xorm:"name"`
}

func AddPetClass(p PetClass) error {
	_, err := MysqlDB.Insert(p)
	return err
}

func FindPetClass() ([]PetClass, error) {
	var petClasses []PetClass
	err := MysqlDB.Find(&petClasses)
	return petClasses, err
}

const (
	PetPublicIn  = 1
	PetPublicOut = 0
)

// PetPublic
type PetPublic struct {
	ID               string    `json:"uuid" xorm:"'uuid'"`
	Title            string    `json:"title" xorm:"title"`
	UserID           string    `json:"user_id" xorm:"user_id"`
	LocalityID       string    `json:"locality_id" xorm:"locality_id"`
	LocalityName     string    `json:"locality_name" xorm:"locality_name"`
	Free             bool      `json:"free" xorm:"free"`
	PetName          string    `json:"pet_name" xorm:"pet_name"`
	PetClassID       string    `json:"pet_classid" xorm:"pet_classid"`
	PetVariety       string    `json:"pet_variety" xorm:"pet_variety"`
	PetAge           int       `json:"pet_age" xorm:"pet_age"`
	PetSex           int       `json:"pet_sex" xorm:"pet_sex"`
	PetDisposition   string    `json:"pet_disposition" xorm:"pet_disposition"`
	PetVaccine       int       `json:"pet_vaccine" xorm:"pet_vaccine"`
	PetSterilization int       `json:"pet_sterilization" xorm:"pet_sterilization"`
	State            bool      `json:"pet_state" xorm:"pet_state"`
	PetDescription   string    `json:"pet_description" xorm:"Text pet_description"`
	AdopteReq        string    `json:"adoption_request" xorm:"varchar(2048) adoption_request"`
	PetImages        []string  `json:"pet_images" xorm:"varchar(1024) pet_images"`
	Created          time.Time `json:"created" xorm:"created"`
	Updated          time.Time `json:"updated" xorm:"updated"`
}

//FindPetPublics
func FindPetPublics(cityID, localityID, petClassID, key string, start, count int) ([]PetPublic, error) {
	var adps []PetPublic
	expr := ""
	if cityID != "" {
		var ls []Locality
		err := MysqlDB.Where("city_id=?", cityID).Find(&ls)
		if err != nil {
			return adps, err
		}
		if len(ls) == 0 {
			return adps, errors.New("can't found city_id" + cityID)
		}

		for i, l := range ls {
			expr += "locality_id='" + l.ID + "'"
			if i < len(ls)-1 {
				expr += " or "
			}
		}
		log.Debug(expr)

	}

	if localityID != "" {
		expr = "locality_id='" + localityID + "'"
	}

	if key != "" {
		if expr != "" {
			expr += " "
		}
		expr += "title like %%" + key + "%%"
	}

	if expr != "" {
		err := MysqlDB.Where(expr).Limit(count, start).Find(&adps)
		if err != nil {
			return adps, err
		}
		return adps, nil
	}
	err := MysqlDB.Limit(count, start).Find(&adps)
	if err != nil {
		return adps, err
	}
	return adps, nil
}

func FindPetPublicsByUser(uid string) ([]PetPublic, error) {
	var adps []PetPublic
	err := MysqlDB.Where("user_id=?", uid).Find(&adps)
	return adps, err
}

func CreatePetPublics(adp PetPublic) error {
	_, err := MysqlDB.Insert(adp)
	return err
}

func UpdatePetPublics(adp PetPublic) error {
	_, err := MysqlDB.Where("uuid=? and user_id=?", adp.ID, adp.UserID).Update(adp)
	return err
}

func DelPetPublics(user_id, uuid string) error {
	var adp PetPublic
	s := MysqlDB.NewSession()
	_, err := s.Where("uuid=? and user_id=?", uuid, user_id).Delete(adp)
	if err != nil {
		s.Rollback()
		return err
	}
	apply := AdoptionApply{
		State:  ApplyDel,
		Remark: "宠物发布已经移除",
	}
	_, err = s.Where("pet_id=?", uuid).Update(&apply)
	if err != nil {
		s.Rollback()
		return err
	}
	err = s.Commit()
	if err != nil {
		s.Rollback()
		return err
	}
	return nil
}

const (
	ApplyWait = 1
	ApplyPass = iota
	ApplyFail
	ApplyDel
)

type AdoptionApply struct {
	ID     string      `json:"uuid" xorm:"'uuid'"`
	UserID string      `json:"user_id" xorm:"user_id"`
	PetID  string      `json:"pet_id" xorm:"pet_id"`
	State  int         `json:"state" xorm:"state"`
	Remark string      `json:"remark" xorm:"remark"`
	Infos  interface{} `json:"infos" xorm:"json infos"`
}

func FindAdoptionApplyByPetID(petID, uid string) ([]AdoptionApply, error) {
	var adapplys []AdoptionApply
	err := MysqlDB.Where("pet_id=?,user_id=?", petID, uid).Find(&adapplys)
	return adapplys, err
}

func FindAdoptionApplyByUserID(userID string) ([]AdoptionApply, error) {
	var adapplys []AdoptionApply
	err := MysqlDB.Where("user_id=?", userID).Find(&adapplys)
	return adapplys, err
}

func CreateAdoptionApply(adapply AdoptionApply) error {
	_, err := MysqlDB.Insert(adapply)
	return err
}

func UpdateAdoptionApply(adapply AdoptionApply) error {
	_, err := MysqlDB.Where("uuid=? user_id=?", adapply.ID, adapply.UserID).Update(adapply)
	return err
}

func DelAdoptionApply(user_id, uuid string) error {
	var adapply AdoptionApply
	_, err := MysqlDB.Where("uuid=? and user_id=?", uuid, user_id).Delete(adapply)
	return err
}
