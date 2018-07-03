package controllers

import (
	"encoding/json"

	"github.com/MartinResearchSociety/connect/models"
	"github.com/astaxie/beego"
)

type PostController struct {
	beego.Controller
}

// @Title CreatPost
// @Description create new post
// @Param body body models.Post true "body for post content"
// @Success 200 {int} models.Post.Id
// @Failure 403 body is empty
// @router /save [post]
func (p *PostController) NewPost() {
	var post models.Post

	json.Unmarshal(p.Ctx.Input.RequestBody, &post)
	err := post.Insert()
	if err != nil {
		p.Abort("400")
		return
	}
	p.Data["json"] = post
	p.ServeJSON()
}
