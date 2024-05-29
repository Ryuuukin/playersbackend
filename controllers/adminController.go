package controllers

import (
	"fmt"
	"net/http"

	"github.com/Ryuuukin/ap-assignment1/initializers"
	"github.com/Ryuuukin/ap-assignment1/models"
	"github.com/Ryuuukin/ap-assignment1/utils"
	"github.com/gin-gonic/gin"
)

// User CRUD
func AdminListUsers(c *gin.Context) {
	var users []models.User
	initializers.DB.Find(&users)

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

func AdminUpdateUser(c *gin.Context) {
	id := c.Param("id")

	var body struct {
		Name  string
		Email string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	var user models.User
	initializers.DB.First(&user, id)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	initializers.DB.Model(&user).Updates(models.User{
		Name:  body.Name,
		Email: body.Email,
	})

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func AdminDeleteUser(c *gin.Context) {
	id := c.Param("id")

	initializers.DB.Delete(&models.User{}, id)

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully deleted user",
	})
}

func AdminStats(c *gin.Context) {
	var userCount int64
	var postCount int64

	initializers.DB.Model(&models.User{}).Count(&userCount)
	initializers.DB.Model(&models.Post{}).Count(&postCount)

	c.JSON(http.StatusOK, gin.H{
		"user_count": userCount,
		"post_count": postCount,
	})
}

// Post CRUD
func AdminListPosts(c *gin.Context) {
	var posts []models.Post
	initializers.DB.Find(&posts)

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})
}

func AdminCreatePost(c *gin.Context) {
	var body struct {
		Title string
		Body  string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	post := models.Post{Title: body.Title, Body: body.Body}
	result := initializers.DB.Create(&post)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create post",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully created post",
		"post":    post,
	})
}

func AdminUpdatePost(c *gin.Context) {
	id := c.Param("id")

	var body struct {
		Title string
		Body  string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	var post models.Post
	initializers.DB.First(&post, id)

	if post.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid post ID",
		})
		return
	}

	initializers.DB.Model(&post).Updates(models.Post{
		Title: body.Title,
		Body:  body.Body,
	})

	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}

func AdminDeletePost(c *gin.Context) {
	id := c.Param("id")

	var post models.Post
	initializers.DB.Delete(&post, id)

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully deleted post",
	})
}

// controllers/adminController.go

func AdminSendEmail(c *gin.Context) {
	var body struct {
		Subject string `json:"subject"`
		Content string `json:"content"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	var users []models.User
	initializers.DB.Find(&users)

	c.JSON(http.StatusOK, gin.H{"message": "Emails sent successfully"})

	for _, user := range users {
		fmt.Println(user)
		err := utils.SendMailSimple(body.Subject, body.Content, []string{user.Email})
		fmt.Printf("Failed to send email to %s: %v\n", user.Email, err)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email to " + user.Email})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Emails sent successfully"})
}
