package controller

import (
	"net/http"

	"github.com/RifaldyAldy/diamond-wallet/model/dto"
	"github.com/RifaldyAldy/diamond-wallet/usecase"
	"github.com/RifaldyAldy/diamond-wallet/utils/common"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	uc usecase.UserUseCase
	rg *gin.RouterGroup
}

func (e *UserController) getHandler(c *gin.Context) {
	id := c.Param("id")

	response, err := e.uc.FindById(id)

	if err != nil {
		common.SendErrorResponse(c, http.StatusInternalServerError, "Error "+err.Error())
		return
	}

	common.SendSingleResponse(c, "OK", response)
}

func (e *UserController) createHandler(c *gin.Context) {
	var payload dto.UserRequestDto
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	payloadResponse, err := e.uc.CreateUser(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    payloadResponse,
	})
}

func (u *UserController) loginHandler(c *gin.Context) {
	var payload dto.LoginRequestDto
	if err := c.ShouldBindJSON(&payload); err != nil {
		common.SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	loginData, err := u.uc.LoginUser(payload)
	if err != nil {
		if err.Error() == "1" {
			common.SendErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
	}
	if err != nil {
		common.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	common.SendSingleResponse(c, "success", loginData)
}

func (p *UserController) Route() {
	p.rg.POST("/users/login", p.loginHandler)
	p.rg.POST("/users", p.createHandler)
	p.rg.GET("/users/:id", p.getHandler)
}

func NewUserController(uc usecase.UserUseCase, rg *gin.RouterGroup) *UserController {
	return &UserController{
		uc: uc,
		rg: rg,
	}
}