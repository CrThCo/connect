package controllers

import (
	"encoding/json"

	"github.com/MartinResearchSociety/connect/models"

	"github.com/astaxie/beego"
)

// Operations about Users
type UserController struct {
	beego.Controller
}

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
	if uid != "" {
		var user models.User
		json.Unmarshal(u.Ctx.Input.RequestBody, &user)
		uu, err := models.UpdateUser(uid, &user)
		if err != nil {
			u.Data["json"] = err.Error()
		} else {
			u.Data["json"] = uu
		}
	}
	u.ServeJSON()
}
