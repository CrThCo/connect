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
		post.VoteCount = 1
		if err := post.Update()	; err != nil {
			log.Printf("Error updating user: %v", err)
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
func (p *PostController) Vote(postID string, vote *models.VoteStruct) {
	if !bson.IsObjectIdHex(postID) {
		log.Printf("Post id is invalid: id=%v\n", postID)
	}
	voterID := p.Ctx.Input.Param("userID")
	if err := vote.AddVote(bson.ObjectIdHex(postID), bson.ObjectIdHex(voterID)); err != nil {
		p.Ctx.Output.SetStatus(500)
		p.Data["json"] = err.Error()
		p.ServeJSON()
		return
	} 
	p.Data["json"] = "success"
	p.ServeJSON()
}

// Vote controller method
// @Title Vote
// @Description Retrieve vote counts by filtering either postID or userID
// @Success 200 {int} count
// @router /vote/count [get]
func (p *PostController) VoteCount(postID, userID *string) {
	var count int
	if postID != nil {
		count, _ = models.CountVotesByPost(bson.ObjectIdHex(*postID))
	} else if userID != nil {
		count, _ = models.CountVotesByUser(bson.ObjectIdHex(*userID))
	}
	
	p.Data["json"] = count
	p.ServeJSON()
}

// Vote controller method
// @Title Vote
// @Description Retrieve vote counts by filtering either postID or userID
// @Success 200 {int} count
// @router /vote/get [get]
func (p *PostController) GetVotesBy(postID, userID *string) {
	var res []bson.M
	if postID != nil {
		res, _ = models.GetVotesByPost(bson.ObjectIdHex(*postID))
	} else {
		voterID := p.Ctx.Input.Param("userID")
		res, _ = models.GetVotesByUser(bson.ObjectIdHex(voterID))
	}
	
	p.Data["json"] = res
	p.ServeJSON()
}

