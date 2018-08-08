package controllers

import (
	"gopkg.in/mgo.v2/bson"
	"encoding/json"
	"log"

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
	if err := post.SaveImage(); err != nil {
		log.Println(err.Error())
		p.Ctx.Output.SetStatus(500)
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

	// add votes if available
	if len(post.Options) > 0 {
		vs := &models.VoteStruct{Options: post.Options}
		if err := vs.AddVote(post.ID, bson.ObjectIdHex(userID)); err != nil {
			p.Ctx.Output.SetStatus(500)
			p.Data["json"] = err.Error()
			p.ServeJSON()
			return
		}
	}
	p.Data["json"] = post
	p.ServeJSON()
}

// GetByUser controller method
// @Title GetByUser
// @Description Retrieve user posts
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


// Vote controller method
// @Title Vote
// @Description Reterive user posts
// @Param body body models.VoteStruct true "body for post content"
// @Success 200 {int} models.Vote.Id
// @Failure 403 body is empty
// @router /:id/vote [post]
func (p *PostController) Vote(postId string, vote *models.VoteStruct) {
	if !bson.IsObjectIdHex(postId) {
		log.Printf("Post id is invalid: id=%v\n", postId)
	}
	voterId := p.Ctx.Input.Param("userID")
	if err := vote.AddVote(bson.ObjectIdHex(postId), bson.ObjectIdHex(voterId)); err != nil {
		p.Ctx.Output.SetStatus(500)
		p.Data["json"] = err.Error()
		p.ServeJSON()
		return
	}
	p.Data["json"] = "success"
	p.ServeJSON()
}
