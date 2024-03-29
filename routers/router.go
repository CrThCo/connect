// @APIVersion 1.0.0
// @Title Connect API Console
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/MartinResearchSociety/connect/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
				&controllers.LoginController{},
			),
		),
		beego.NSNamespace("/post",
			beego.NSInclude(
				&controllers.PostController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
