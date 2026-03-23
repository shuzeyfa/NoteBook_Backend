package controllers

import (
	"errors"
	"net/http"
	domain "taskmanagement/Domain"
	usecase "taskmanagement/Usecase"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NoteController struct {
	Control *usecase.NoteUsecase
}

func getUserIDFromContext(c *gin.Context) (primitive.ObjectID, error) {
	id, exists := c.Get("user_id")
	if !exists {
		return primitive.NilObjectID, errors.New("user not authenticated")
	}
	return id.(primitive.ObjectID), nil
}

func (ctr *NoteController) GetAllNote(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	notes, err := ctr.Control.GetAllNote(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, notes)
}

func (ctr *NoteController) CreateNote(c *gin.Context) {
	var note domain.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	result, err := ctr.Control.CreateNote(note, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (ctr *NoteController) UpdateNote(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	var note domain.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	note.ID = id

	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	result, err := ctr.Control.UpdateNote(note, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (ctr *NoteController) DeleteNote(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	err = ctr.Control.DeleteNote(id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "note deleted successfully"})
}
