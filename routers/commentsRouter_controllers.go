package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/MartinResearchSociety/connect/controllers:LoginController"] = append(beego.GlobalControllerRouter["github.com/MartinResearchSociety/connect/controllers:LoginController"],
		beego.ControllerComments{
			Method: "Auth",
			Router: `/login`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/MartinResearchSociety/connect/controllers:LoginController"] = append(beego.GlobalControllerRouter["github.com/MartinResearchSociety/connect/controllers:LoginController"],
		beego.ControllerComments{
			Method: "Refresh",
			Router: `/token/refresh`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/MartinResearchSociety/connect/controllers:PostController"] = append(beego.GlobalControllerRouter["github.com/MartinResearchSociety/connect/controllers:PostController"],
		beego.ControllerComments{
			Method: "GetByUser",
			Router: `/all`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/MartinResearchSociety/connect/controllers:PostController"] = append(beego.GlobalControllerRouter["github.com/MartinResearchSociety/connect/controllers:PostController"],
		beego.ControllerComments{
			Method: "NewPost",
			Router: `/save`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/MartinResearchSociety/connect/controllers:PostController"] = append(beego.GlobalControllerRouter["github.com/MartinResearchSociety/connect/controllers:PostController"],
		beego.ControllerComments{
			Method: "Vote",
			Router: `/vote`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(
				param.New("voterId"),
				param.New("postId"),
				param.New("vote"),
			),
			Params: nil})

	beego.GlobalControllerRouter["github.com/MartinResearchSociety/connect/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/MartinResearchSociety/connect/controllers:UserController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/MartinResearchSociety/connect/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/MartinResearchSociety/connect/controllers:UserController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/:uid`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/MartinResearchSociety/connect/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/MartinResearchSociety/connect/controllers:UserController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:uid`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/MartinResearchSociety/connect/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/MartinResearchSociety/connect/controllers:UserController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/signup`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

}
