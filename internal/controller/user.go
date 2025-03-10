package controller

import (
	"net/http"

	"github.com/YoubaImkf/go-auth-api/internal/service"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// @Summary      Get all users
// @Description  Get a list of all users
// @Tags         user
// @Produce      json
// @Success      200  {array}  model.User
// @Router       /users [get]
func (c *UserController) GetAllUsers(ctx *gin.Context) {
	users, err := c.userService.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, users)
}
