package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/hashwing/pet-adoption/controllers"
	"github.com/hashwing/pet-adoption/pkg/auth"
)

func init() {
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins: true,
		AllowMethods:    []string{"PUT", "POST", "GET", "DELETE", "OPTIONS"},
	}))

	authApi := beego.NewNamespace("/pet/auth",
		beego.NSRouter("wx", &controllers.UserController{}, "get:WxLogin"),
	)

	staticApi := beego.NewNamespace("/pet/static",
		beego.NSRouter("img/:key", &controllers.OssController{}, "get:GetImg;delete:DelImg"),
	)

	api := beego.NewNamespace("/pet/api",
		beego.NSRouter("v1/oss", &controllers.OssController{}, "get:GetToken"),
		beego.NSRouter("v1/user", &controllers.UserController{}, "put:Put"),
		beego.NSRouter("v1/province", &controllers.LocalityController{}, "get:FindProvinces"),
		beego.NSRouter("v1/province/:province_id/city", &controllers.LocalityController{}, "get:FindCitys"),
		beego.NSRouter("v1/province/:province_id/city/:city_id/locality", &controllers.LocalityController{}, "get:FindLocalities"),
		beego.NSRouter("v1/locality", &controllers.LocalityController{}, "get:Find"),
		beego.NSRouter("v1/pet/class", &controllers.AdoptionController{}, "get:FindPetClass"),
		beego.NSRouter("v1/adoption", &controllers.AdoptionController{}, "get:PublicList"),
		beego.NSRouter("v1/adoption/pet", &controllers.AdoptionController{}, "get:PublicListByUser;post:CreatePublic"),
		beego.NSRouter("v1/adoption/pet/:uuid", &controllers.AdoptionController{}, "put:UpdatePublic;delete:DeletePublic"),
		beego.NSRouter("v1/adoption/pet/:uuid", &controllers.AdoptionController{}, "get:GetPublic;put:UpdatePublic;delete:DeletePublic"),
		beego.NSRouter("v1/adoption/pet/:pet_id/application", &controllers.AdoptionController{}, "get:ApplyListByPet;post:CreateApply"),
		beego.NSRouter("v1/adoption/pet/:pet_id/application/:uuid", &controllers.AdoptionController{}, "put:UpdateApply;delete:DelApply"),
		beego.NSRouter("v1/adoption/application", &controllers.AdoptionController{}, "get:ApplyListByUser"),
		beego.NSRouter("v1/adoption/application/:uuid", &controllers.AdoptionController{}, "get:GetApply"),
	)

	beego.AddNamespace(authApi, staticApi, api)
	beego.InsertFilter("pet/api/*", beego.BeforeRouter, auth.JwtAuthFilter)
}
