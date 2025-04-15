package controllers

import (
	"fmt"
	"jwt-banking/database"
	"jwt-banking/models"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Hello() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	}
}

func RegisterUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		//parsing and validating the input
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//checking the email exists

		var existingUser models.User

		if err := database.DB.Where("email =  ?", user.Email).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
			return
		}

		//hashing the password

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash the password"})
			return
		}
		//creating the user
		user = models.User{
			Name:     user.Name,
			Email:    user.Email,
			Password: hashedPassword,
		}

		//inserting the user to the database
		if err := database.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create a User "})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User Registered Successfully", "email": user.Email, "name": user.Name})
		
	}
}

func LoginUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.User
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request body",
			})
			return
		}

		var user models.User
		if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
			fmt.Println("User not found:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
			return
		}

		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
			return
		}

		// jwt claims
		claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": strconv.Itoa(int(user.ID)),
			"exp": time.Now().Add(24 * time.Hour).Unix(),
		})

		secret := os.Getenv("JWT_SECRET") // Use environment variable for the secret
		token, err := claims.SignedString([]byte(secret))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not provide the token"})
			return
		}

		// Set cookie with more security options
		c.SetCookie("jwt", token, 3600*24, "/", "localhost", true, true)

		c.JSON(http.StatusAccepted, gin.H{"message": "Login was successful"})
	}
}
