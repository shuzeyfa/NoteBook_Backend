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
	log.Println("✅ App is ready!")

	// Repositories & Usecases (keep exactly as before)
	noteRepo := &repository.MongoNoteRepository{
		Collection: infrastructure.Client.Database(infrastructure.DBName).Collection("notes"),
	}
	noteUsecase := &usecase.NoteUsecase{Repo: noteRepo}
	noteController := &controllers.NoteController{Control: noteUsecase}

	userRepo := &repository.MongoUserRepository{
		Collection: infrastructure.Client.Database(infrastructure.DBName).Collection("users"),
	}
	userUsecase := &usecase.UserUsecase{Repo: userRepo}
	userController := &controllers.UserController{Control: userUsecase}

	aiUsecase := &usecase.AIUsecase{Repo: noteRepo}
	aiController := &controllers.AIController{Control: aiUsecase}

	r := gin.Default()

	// ================ CORS CONFIGURATION ================
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",              // local dev
			"https://your-note-taker.vercel.app", // your production URL
			"https://*.vercel.app",               // all Vercel previews (very important)
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
			"Accept",
			"X-Requested-With",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 hours
	}))

	// Auth routes (public)
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
		api.POST("/notes/:id/ai", aiController.GenerateAI)
	}

	return r
}
