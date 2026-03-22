package controllers

import (
	"net/http"
	domain "taskmanagement/Domain"
	usecase "taskmanagement/Usecase"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	Control *usecase.UserUsecase
}

func (ctr *UserController) RegisterHandler(c *gin.Context) {

	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	user, err := ctr.Control.RegisterUser(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func (ctr *UserController) LoginUser(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	userCredential, err := ctr.Control.LoginUser(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userCredential)
}
