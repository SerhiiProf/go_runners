package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"runners/models"
	"runners/services"

	"github.com/gin-gonic/gin"
)

const ROLE_ADMIN = "admin"
const ROLE_RUNNER = "runner"

type RunnersController struct {
	runnersService *services.RunnersService
	usersService   *services.UsersService
}

func NewRunnersController(runnersService *services.RunnersService, usersService *services.UsersService) *RunnersController {
	return &RunnersController{runnersService: runnersService, usersService: usersService}
}

func (rc RunnersController) CreateRunner(ctx *gin.Context) {
	accessToken := ctx.Request.Header.Get("Token")
	auth, responseErr := rc.usersService.AuthorizeUser(accessToken, []string{ROLE_ADMIN})
	if responseErr != nil {
		ctx.JSON(responseErr.Status, responseErr)
		return
	}

	if !auth {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	var body, err = io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Error while reading create runner request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var runner models.Runner
	err = json.Unmarshal(body, &runner)
	if err != nil {
		log.Println("Error while unmarshaling create runner request body ", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	response, responseErr := rc.runnersService.CreateRunner(&runner)
	if responseErr != nil {
		ctx.AbortWithError(responseErr.Status, responseErr)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (rc RunnersController) UpdateRunner(ctx *gin.Context) {
	accessToken := ctx.Request.Header.Get("Token")
	auth, responseErr := rc.usersService.AuthorizeUser(accessToken, []string{ROLE_ADMIN})
	if responseErr != nil {
		ctx.JSON(responseErr.Status, responseErr)
		return
	}

	if !auth {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Error while reading update runner request body")
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var runner models.Runner
	err = json.Unmarshal(body, &runner)
	if err != nil {
		log.Println("Error while unmarshalling update runner request body")
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	responseErr = rc.runnersService.UpdateRunner(&runner)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (rc RunnersController) DeleteRunner(ctx *gin.Context) {
	accessToken := ctx.Request.Header.Get("Token")
	auth, responseErr := rc.usersService.AuthorizeUser(accessToken, []string{ROLE_ADMIN})
	if responseErr != nil {
		ctx.JSON(responseErr.Status, responseErr)
		return
	}

	if !auth {
		ctx.Status(http.StatusUnauthorized)
		return
	}
	runnerId := ctx.Param("id")
	responseErr = rc.runnersService.DeleteRunner(runnerId)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (rc RunnersController) GetRunner(ctx *gin.Context) {
	accessToken := ctx.Request.Header.Get("Token")
	auth, responseErr := rc.usersService.AuthorizeUser(accessToken, []string{ROLE_ADMIN, ROLE_RUNNER})
	if responseErr != nil {
		ctx.JSON(responseErr.Status, responseErr)
		return
	}

	if !auth {
		ctx.Status(http.StatusUnauthorized)
		return
	}
	runnerId := ctx.Param("id")
	runner, responseErr := rc.runnersService.GetRunner(runnerId)
	if responseErr != nil {
		ctx.JSON(responseErr.Status, responseErr)
		return
	}
	ctx.JSON(http.StatusOK, runner)
}

func (rc RunnersController) GetRunnersBatch(ctx *gin.Context) {
	accessToken := ctx.Request.Header.Get("Token")
	auth, responseErr := rc.usersService.AuthorizeUser(accessToken, []string{ROLE_ADMIN, ROLE_RUNNER})
	if responseErr != nil {
		ctx.JSON(responseErr.Status, responseErr)
		return
	}

	if !auth {
		ctx.Status(http.StatusUnauthorized)
		return
	}
	country := ctx.Query("country")
	year := ctx.Query("year")
	runners, responseErr := rc.runnersService.GetRunnersBatch(country, year)
	if responseErr != nil {
		ctx.JSON(responseErr.Status, responseErr)
		return
	}
	ctx.JSON(http.StatusOK, runners)

}
