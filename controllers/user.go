package controllers

import (
	"net/http"

	"github.com/MartinResearchSociety/connect-api/db"
	"github.com/MartinResearchSociety/connect-api/utils"
	"github.com/gin-gonic/gin"
)

// UserController struct
type UserController struct {
	*Controller
}

// Signup user
func (u *UserController) Signup(c *gin.Context) {
	user := &db.User{}
	if err := c.Bind(user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	if err := user.Insert(); err != nil {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSONP(http.StatusCreated, &user)
	return
}

// List of all users
func (u *UserController) List(c *gin.Context) {
	list, err := db.GetUserList()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSONP(http.StatusOK, &list)
}

// Info of the user
func (u *UserController) Info(c *gin.Context) {
	uid := c.Param("uid")
	user := db.User{}
	err := db.GetUser(uid, &user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSONP(http.StatusOK, &user)
}

// Update user profile info
func (u *UserController) Update(c *gin.Context) {
	uid := c.GetString("UID")
	user := &db.User{}

	if err := db.GetUser(uid, user); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := c.Bind(user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	if err := db.UpdateUser(uid, user); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"message": "updated successfully!",
	})
}

// Upload user profile picture
func (u *UserController) Upload(c *gin.Context) {
	uid := c.GetString("UID")
	file, hdrs, err := c.Request.FormFile("image")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	defer file.Close()

	filname, err := utils.SaveFile(file, hdrs, uid, "profile")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	user := &db.User{}
	if err := db.GetUser(uid, user); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	user.Image = filname
	if err := db.UpdateUser(uid, user); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSONP(http.StatusOK, gin.H{
		"message": "uploaded sucessfully",
	})
}
