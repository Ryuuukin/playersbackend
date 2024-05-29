package controllers

import (
	"github.com/Ryuuukin/ap-assignment1/initializers"
	"github.com/Ryuuukin/ap-assignment1/models"
	"github.com/gin-gonic/gin"
)

func PostsFetch(c *gin.Context) {
	var posts []models.Post
	initializers.DB.Find(&posts)

	c.JSON(200, gin.H{
		"posts": posts,
	})
}
