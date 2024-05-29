package main

import (
	"net/http"

	"github.com/Ryuuukin/ap-assignment1/controllers"
	"github.com/Ryuuukin/ap-assignment1/initializers"
	"github.com/Ryuuukin/ap-assignment1/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVars()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowCredentials = true
	config.AllowOrigins = []string{"http://localhost:3000"} // Update with your frontend URL
	r.Use(cors.New(config))

	// Public routes
	r.POST("/signup", middlewares.RateLimitMiddleware, controllers.Signup)
	r.POST("/login", middlewares.RateLimitMiddleware, controllers.Login)

	r.GET("/emailver/:username/:verPass", middlewares.RateLimitMiddleware, controllers.EmailverGEThandler)
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to E-Player's Website!"})
	})

	// Protected routes
	protected := r.Group("/")
	protected.Use(middlewares.RequireAuth, middlewares.RateLimitMiddleware)
	{
		protected.GET("/validate", controllers.Validate)
		protected.GET("/posts", controllers.PostsFetch)

		// Admin routes
		admin := protected.Group("/admin")
		admin.Use(middlewares.RequireAdmin)
		{
			admin.GET("/users", controllers.AdminListUsers)
			admin.PUT("/users/:id", controllers.AdminUpdateUser)
			admin.DELETE("/users/:id", controllers.AdminDeleteUser)
			admin.GET("/stats", controllers.AdminStats)
			admin.POST("/send-email", controllers.AdminSendEmail)

			// Admin post routes
			admin.GET("/posts", controllers.AdminListPosts)
			admin.POST("/create-posts", controllers.AdminCreatePost)
			admin.PUT("/posts/:id", controllers.AdminUpdatePost)
			admin.DELETE("/posts/:id", controllers.AdminDeletePost)
		}
	}

	r.Run()
}
