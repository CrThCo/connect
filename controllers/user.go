package controllers

import (
	"encoding/json"

	"github.com/MartinResearchSociety/connect/models"

	"github.com/astaxie/beego"
)

// UserController Operations about Users
type UserController struct {
	beego.Controller
}

// Post method for endpoing signup
// @Title CreateUser
// @Description create users
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.User.Id
// @Failure 403 body is empty
// @router /signup [post]
func (u *UserController) Post() {
	var user models.User
	if err := json.Unmarshal(u.Ctx.Input.RequestBody, &user); err != nil {
		u.Ctx.Output.SetStatus(400)
		u.Data["json"] = err.Error()
		u.ServeJSON()
		return
	}
	err := user.Insert()
	if err != nil {
		u.Ctx.Output.SetStatus(500)
		u.Data["json"] = err.Error()
		u.ServeJSON()
		return
	}
	u.Ctx.Output.SetStatus(200)
	u.Data["json"] = user
	u.ServeJSON()
}

// GetAll method for endpoint /
// @Title GetAll
// @Description get all Users
// @Success 200 {object} models.User
// @router / [get]
func (u *UserController) GetAll() {
	users := models.GetAllUsers()
	u.Data["json"] = users
	u.ServeJSON()
}

// @Title Get
// @Description get user by uid
// @Param	uid		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.User
// @Failure 403 :uid is empty
// @router /:uid [get]
func (u *UserController) Get() {
	uid := u.GetString(":uid")
	if uid == "" {
		u.Ctx.Output.SetStatus(400)
		u.Data["json"] = "Bad Request"
		u.ServeJSON()
		return
	}
	user, err := models.GetUser(uid)
	if err != nil {
		u.Ctx.Output.SetStatus(500)
		u.Data["json"] = err.Error()
		u.ServeJSON()
		return
	}
	u.Data["json"] = user
	u.ServeJSON()
}

// @Title Update
// @Description update the user
// @Param	uid		path 	string	true		"The uid you want to update"
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {object} models.User
// @Failure 403 :uid is not int
// @router /:uid [put]
func (u *UserController) Put() {
	uid := u.GetString(":uid")
	if uid == "" {
		u.Ctx.Output.SetStatus(400)
		u.Data["json"] = "Bad Request"
		u.ServeJSON()
		return
	}
	var user models.User
	if err := json.Unmarshal(u.Ctx.Input.RequestBody, &user); err != nil {
		u.Ctx.Output.SetStatus(500)
		u.Data["json"] = err.Error()
		u.ServeJSON()
		return
	}

	uu, err := models.UpdateUser(uid, &user)
	if err != nil {
		u.Ctx.Output.SetStatus(500)
		u.Data["json"] = err.Error()
		u.ServeJSON()
		return
	}

	u.Data["json"] = uu
	u.ServeJSON()
}


// @Title Upload
// @Description upload user image
// @router /upload [post]
func (u *UserController) Upload() {
	f, h, _ := u.GetFile("file")  
	path := beego.AppConfig.String("FileStoragePath") + h.Filename   
	f.Close()                          
	if err := u.SaveToFile("file", path); err != nil {
		u.Ctx.Output.SetStatus(500)
		u.Data["json"] = err.Error()
		u.ServeJSON()
		return
	} 
}

// @Title DownloadFile
// @router /download [get]
func (u *UserController) DownloadFile() {
	filename := u.GetString("file")
	path := beego.AppConfig.String("FileStoragePath") + filename
	u.Ctx.Output.Download(path)
}
