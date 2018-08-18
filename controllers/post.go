package controllers

import (
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/MartinResearchSociety/connect-api/db"
	"github.com/gin-gonic/gin"
)

// PostController struct
type PostController struct {
	*Controller
}

// New post method
func (p *PostController) New(c *gin.Context) {
	uid := c.GetString("UID")
	post := &db.Post{}
	post.Poster = bson.ObjectIdHex(uid)
	if err := c.Bind(post); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := post.SaveImage(); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := post.Insert(); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	// add votes if available
	if len(post.Options) > 0 {
		v := &db.VoteStruct{Options: post.Options}
		if err := v.AddVote(post.ID, bson.ObjectIdHex(uid)); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		post.VoteCount = 1
		if err := post.Update(); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
	}
	c.JSONP(http.StatusCreated, gin.H{
		"message": "success",
	})
}

// List of posts
func (p *PostController) List(c *gin.Context) {
	post := &db.Post{}
	posts, err := post.GetList()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSONP(http.StatusOK, &posts)
}

// Vote cast on a post
func (p *PostController) Vote(c *gin.Context) {
	vote := &db.VoteStruct{}
	postID := c.Param("post_id")
	uid := c.GetString("UID")
	if err := vote.AddVote(bson.ObjectIdHex(postID), bson.ObjectIdHex(uid)); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSONP(http.StatusOK, gin.H{
		"message": "success",
	})
}

// CountVoteByPost method
func (p *PostController) CountVoteByPost(c *gin.Context) {
	postID := c.Param("post_id")
	count, err := db.CountVotesByPost(bson.ObjectIdHex(postID))
	if err != nil {
		log.Printf("Vote count by post error: %s", err.Error())
		count = 0
	}
	c.JSONP(http.StatusOK, gin.H{
		"count": count,
	})
}

// CountVoteByUser method
func (p *PostController) CountVoteByUser(c *gin.Context) {
	userID := c.Param("user_id")
	count, err := db.CountVotesByUser(bson.ObjectIdHex(userID))
	if err != nil {
		log.Printf("Vote count by post error: %s", err.Error())
		count = 0
	}
	c.JSONP(http.StatusOK, gin.H{
		"count": count,
	})
}

// GetVoteByPost method
func (p *PostController) GetVoteByPost(c *gin.Context) {
	postID := c.Param("post_id")
	list, err := db.CountVotesByPost(bson.ObjectIdHex(postID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSONP(http.StatusOK, list)
}

// GetVoteByUser method
func (p *PostController) GetVoteByUser(c *gin.Context) {
	userID := c.Param("user_id")
	list, err := db.GetVotesByUser(bson.ObjectIdHex(userID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSONP(http.StatusOK, list)
}
