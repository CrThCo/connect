package controllers

import "github.com/gin-gonic/gin"

// Controller type
type Controller struct {
	Error gin.H       `json:"error"`
	Data  interface{} `json:"data"`
}
