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

type ResultsController struct {
	resultsService *services.ResultsService
	usersService   *services.UsersService
}

func NewResultsController(resultsService *services.ResultsService, usersService *services.UsersService) *ResultsController {
	return &ResultsController{resultsService: resultsService, usersService: usersService}
}

func (rc ResultsController) CreateResult(ctx *gin.Context) {
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
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Error while reading create result request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var result models.Result

	err = json.Unmarshal(body, &result)

	if err != nil {
		log.Println("Error while unmarshaling creates result request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	response, responseErr := rc.resultsService.CreateResult(&result)

	if responseErr != nil {
		ctx.JSON(responseErr.Status, responseErr)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (rc ResultsController) DeleteResult(ctx *gin.Context) {
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
	resultId := ctx.Param("id")
	responseErr = rc.resultsService.DeleteResult(resultId)
	if responseErr != nil {
		ctx.JSON(responseErr.Status, responseErr)
		return
	}
	ctx.Status(http.StatusNoContent)
}
