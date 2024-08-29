package router

import (
	"github.com/subashshakya/SFSS/controllers"
	"github.com/subashshakya/SFSS/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	fileRoutes := router.Group("/files")
	{
		fileRoutes.GET("/fetch_all/:id", controllers.GetUserFiles)
		fileRoutes.PATCH("/update", controllers.UpdateSecureFile)
		fileRoutes.POST("/create", controllers.MakeSecureFile)
		fileRoutes.DELETE("/delete", controllers.DeleteFile)
		fileRoutes.GET("/:id", controllers.GetSecureFileByID)
	}

	secretRoutes := router.Group("/secret")
	{
		secretRoutes.POST("/create", controllers.CreateSuperSecret)
		secretRoutes.GET("/:id", controllers.ReadSuperSecret)
		secretRoutes.PATCH("/update", controllers.UpdatedSuperSecret)
		secretRoutes.DELETE("/delete", controllers.DeleteSuperSecret)
		secretRoutes.GET("/fetch_all/:id", controllers.GetSuperSecretsForUser)
	}

	sharingRoutes := router.Group("/sharing")
	{
		sharingRoutes.Use(middlewares.CheckInvalidToken())
		sharingRoutes.POST("/secure_file", controllers.ShareSecureFile)
		sharingRoutes.POST("/super_secret", controllers.ShareSuperSecret)
		sharingRoutes.GET("/files/:id", controllers.GetFileSharedOfAUser)
		sharingRoutes.GET("/secrets/:id", controllers.GetSecretSharedOfAUser)
	}

	userRoutes := router.Group("/user")
	{
		userRoutes.POST("/sign_up", controllers.UserSignUp)
		userRoutes.POST("/sign_in", controllers.UserSignIn)
		userRoutes.GET("/:id", controllers.GetUser)
		userRoutes.PATCH("/update", controllers.UpdateUser)
		userRoutes.DELETE("/delete/:id", controllers.DeleteUser)
	}
}
