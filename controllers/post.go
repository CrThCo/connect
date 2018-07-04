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
	userID := p.Ctx.Input.Param("userID")

	if err := json.Unmarshal(p.Ctx.Input.RequestBody, &post); err != nil {
		p.Ctx.Output.SetStatus(400)
		p.Data["json"] = err.Error()
		p.ServeJSON()
		return
	}
	post.Poster = userID
	if err := post.Insert(); err != nil {
		p.Ctx.Output.SetStatus(500)
		p.Data["json"] = err.Error()
		p.ServeJSON()
		return
	}
	p.Data["json"] = post
	p.ServeJSON()
}

// GetByUser controller method
// @Title GetByUser
// @Description Reterive user posts
// @Success 200 {int} []models.Post
// @Failure 403 body is empty
// @router /all [get]
func (p *PostController) GetByUser() {
	post := models.Post{}
	posts, err := post.GetByUser()
	if err != nil {
		p.Ctx.Output.SetStatus(500)
		p.Data["json"] = err.Error()
		p.ServeJSON()
		return
	}
	p.Data["json"] = posts
	p.ServeJSON()
}
