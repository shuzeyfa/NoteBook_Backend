package router

import (
	"log"
	"taskmanagement/Delivery/controllers"
	infrastructure "taskmanagement/Infrastructure"
	repository "taskmanagement/Repository"
	usecase "taskmanagement/Usecase"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	infrastructure.ConnectDB()
	log.Println("App is ready!")

	userRepo := &repository.MongoUserRepository{
		Collection: infrastructure.Client.Database(infrastructure.DBName).Collection("user"),
	}

	UserUsecase := &usecase.UserUsecase{
		Repo: userRepo,
	}

	UserController := &controllers.UserController{
		Control: UserUsecase,
	}

	r := gin.Default()

	r.Use(cors.Default())

	r.POST("/register", UserController.RegisterHandler)
	r.POST("/login", UserController.LoginUser)

	return r
}
