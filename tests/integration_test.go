package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ryuuukin/ap-assignment1/controllers"
	"github.com/Ryuuukin/ap-assignment1/initializers"
	"github.com/Ryuuukin/ap-assignment1/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLoginAndAccessPosts(t *testing.T) {

	initializers.LoadEnvVars()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
	// Initialize the router
	router := gin.Default()

	// Set up the routes
	router.POST("/signup", middlewares.RateLimitMiddleware, controllers.Signup)
	router.POST("/login", middlewares.RateLimitMiddleware, controllers.Login)
	router.GET("/emailver/:username/:verPass", middlewares.RateLimitMiddleware, controllers.EmailverGEThandler)
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to E-Player's Website!"})
	})
	protected := router.Group("/")
	protected.Use(middlewares.RequireAuth, middlewares.RateLimitMiddleware)
	{
		protected.GET("/validate", controllers.Validate)
		protected.GET("/posts", controllers.PostsFetch)
	}

	// Perform login request
	loginData := map[string]string{
		"Email":    "akineshova00@gmail.com", // Use existing user's email
		"Password": "Admin",                  // Use existing user's password
	}
	loginJSON, _ := json.Marshal(loginData)

	loginReq, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginJSON))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRec := httptest.NewRecorder()
	router.ServeHTTP(loginRec, loginReq)

	// Assert login response
	assert.Equal(t, http.StatusOK, loginRec.Code)

	var loginResponse map[string]interface{}
	json.Unmarshal(loginRec.Body.Bytes(), &loginResponse)
	authToken, _ := loginResponse["Authorization"].(string)
	fmt.Println(authToken)
	assert.NotEmpty(t, authToken)

	// Mock the cookie with the token
	cookie := &http.Cookie{
		Name:  "Authorization",
		Value: authToken,
		Path:  "/",
	}

	// Create a new request and add the cookie
	postsReq, _ := http.NewRequest("GET", "/posts", nil)
	postsReq.AddCookie(cookie)
	postsRec := httptest.NewRecorder()
	router.ServeHTTP(postsRec, postsReq)

	// Assert /posts response
	assert.Equal(t, http.StatusOK, postsRec.Code)
	// Optionally, you can further validate the response body or headers if needed
}
