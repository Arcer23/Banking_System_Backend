package routes 

import (
	"github.com/gin-gonic/gin"
	"jwt-banking/controllers"
)

func SetUpRoutes(router *gin.Engine) {
	router.GET("/user", controllers.Hello());
	router.POST("/register", controllers.RegisterUser());
	router.POST("/login", controllers.LoginUser());
}