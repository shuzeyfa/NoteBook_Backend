package controllers

import (
	"net/http"
	domain "taskmanagement/Domain"
	usecase "taskmanagement/Usecase"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AIController struct {
	Control *usecase.AIUsecase
}

func (ctr *AIController) GenerateAI(c *gin.Context) {
	noteID, _ := primitive.ObjectIDFromHex(c.Param("id"))

	var req domain.GenerateAIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	response, err := ctr.Control.GenerateResponse(noteID, userID, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.GenerateAIResponse{Response: response})
}
