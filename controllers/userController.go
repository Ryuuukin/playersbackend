package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/Ryuuukin/ap-assignment1/initializers"
	"github.com/Ryuuukin/ap-assignment1/logging"
	"github.com/Ryuuukin/ap-assignment1/models"
	"github.com/Ryuuukin/ap-assignment1/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	var body struct {
		Name     string
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		logging.LogError("Signup", "Failed to read body", nil)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		logging.LogError("Signup", "Failed to hash the password", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash the password",
		})

		return
	}

	// Generate email verification hash
	rand.Seed(time.Now().UnixNano())
	var alphaNumRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	emailVerRandRune := make([]rune, 64)
	for i := range emailVerRandRune {
		emailVerRandRune[i] = alphaNumRunes[rand.Intn(len(alphaNumRunes))]
	}
	emailVerPassword := string(emailVerRandRune)
	emailVerPWhash, err := bcrypt.GenerateFromPassword([]byte(emailVerPassword), bcrypt.DefaultCost)
	if err != nil {
		logging.LogError("Signup", "Failed to generate verification hash", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate verification hash",
		})
		return
	}

	user := models.User{Name: body.Name, Email: body.Email, Password: string(hash), Verhash: string(emailVerPWhash)}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		logging.LogError("Signup", "Failed to create User", result.Error)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create User",
		})

		return
	}

	// Send verification email
	var domName string = "http://localhost:8080"

	subject := "Email Verification"
	emailBody := domName + `/emailver/` + user.Name + `/` + emailVerPassword + " click to verify email"

	err = utils.SendMailSimple(subject, emailBody, []string{user.Email}) // Assuming you have a SendEmail function
	if err != nil {
		logging.LogError("Signup", "Failed to send verification email", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to send verification email",
		})
		return
	}

	logging.LogUserCreation(user.Name, "N/A")
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully registered. Confirmation email sent.",
		"user":    user,
	})
}

func Login(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		logging.LogError("Login", "Failed to read body", nil)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		logging.LogError("Login", "Invalid Email or password", nil)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Email or password",
		})

		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		logging.LogError("Login", "Invalid Email or password", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Email or password",
		})

		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		logging.LogError("Login", "Failed to create token", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid to create token",
		})

		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "/", "localhost", true, true)

	logging.LogUserLogin(user.Email)
	c.JSON(http.StatusOK, gin.H{
		"Authorization": tokenString,
		"message":       "Successfully logged in",
	})
}

func EmailverGEThandler(c *gin.Context) {
	var u models.User
	u.Name = c.Param("username")
	linkVerPass := c.Param("verPass")
	err := initializers.DB.Where("Name = ?", u.Name).First(&u).Error
	if err != nil {
		fmt.Println("error selecting ver_hash in db by Username, err:", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid verification link"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Verhash), []byte(linkVerPass))
	if err == nil {
		err = initializers.DB.Model(&u).Update("is_active", true).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Please try the email confirmation link again",
			})
			return
		}
		c.Redirect(http.StatusFound, "http://localhost:3000/email-confirmed")
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid verification link",
		})
	}
}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")

	logging.LogValidation(user.(models.User).Email)
	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}
