package router

import (
	"log"
	"taskmanagement/Delivery/controllers"
	infrastructure "taskmanagement/Infrastructure"
	repository "taskmanagement/Repository"
	usecase "taskmanagement/Usecase"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	infrastructure.ConnectDB()
	log.Println("✅ App is ready!")

	// Note
	noteRepo := &repository.MongoNoteRepository{
		Collection: infrastructure.Client.Database(infrastructure.DBName).Collection("notes"),
	}
	noteUsecase := &usecase.NoteUsecase{Repo: noteRepo}
	noteController := &controllers.NoteController{Control: noteUsecase}

	// User
	userRepo := &repository.MongoUserRepository{
		Collection: infrastructure.Client.Database(infrastructure.DBName).Collection("users"),
	}
	userUsecase := &usecase.UserUsecase{Repo: userRepo}
	userController := &controllers.UserController{Control: userUsecase}

	r := gin.Default()

	r.POST("/register", userController.RegisterHandler)
	r.POST("/login", userController.LoginUser)

	// Protected routes
	api := r.Group("/")
	api.Use(infrastructure.AuthMiddleware())
	{
		api.GET("/notes", noteController.GetAllNote)
		api.POST("/notes", noteController.CreateNote)
		api.PUT("/notes/:id", noteController.UpdateNote)
		api.DELETE("/notes/:id", noteController.DeleteNote)
	}

	return r
}
