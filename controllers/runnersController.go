package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	//"runners-postgresql/models"
	//"runners-postgresql/services"
	"runners/models"
	"runners/services"

	"github.com/gin-gonic/gin"
)

type RunnersController struct {
	runnersService *services.RunnersService
}

func NewRunnersController(runnersService *services.RunnersService) *RunnersController {
	return &RunnersController{
		runnersService: runnersService,
	}
}

func (rh RunnersController) CreateRunner(ctx *gin.Context) {
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
	response, responseErr := rh.runnersService.CreateRunner(&runner)
	if responseErr != nil {
		ctx.AbortWithError(responseErr.Status, responseErr)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (rh RunnersController) UpdateRunner(ctx *gin.Context) {
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
	responseErr := rh.runnersService.UpdateRunner(&runner)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (rh RunnersController) DeleteRunner(ctx *gin.Context) {
	runnerId := ctx.Param("id")
	responseErr := rh.runnersService.DeleteRunner(runnerId)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (rh RunnersController) GetRunner(ctx *gin.Context) {
	runnerId := ctx.Param("id")
	runner, responseErr := rh.runnersService.GetRunner(runnerId)
	if responseErr != nil {
		ctx.JSON(responseErr.Status, responseErr)
		return
	}
	ctx.JSON(http.StatusOK, runner)
}

func (rh RunnersController) GetRunnersBatch(ctx *gin.Context) {
	country := ctx.Query("country")
	year := ctx.Query("year")
	runners, responseErr := rh.runnersService.GetRunnersBatch(country, year)
	if responseErr != nil {
		ctx.JSON(responseErr.Status, responseErr)
		return
	}
	ctx.JSON(http.StatusOK, runners)

}
