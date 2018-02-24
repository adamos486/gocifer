package controllers

import (
	"fmt"
	"net/http"

	"cobb.com/server/events/models"
	"github.com/gin-gonic/gin"
	"database/sql"
	"io/ioutil"
	"bytes"
)

type EventsService interface {
	AddCannedRowToEventsDB() (interface{}, sql.Result, error)
	AddNewEvent(name string, description string) (*models.EventRow, sql.Result, error)
	GetAllEvents() (*[]models.GetAllEventsResponseStruct, error)
}

var service EventsService

func NewEventsApiClient(eventsService EventsService) {
	service = eventsService
}

func AddEvent(ctx *gin.Context) {
	var bodyBytes []byte
	var data map[string]string

	if ctx.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(ctx.Request.Body)
	}
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	err := ctx.BindJSON(&data)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Post body must be in JSON",
		})
		return
	}
	name, nameOk := data["name"]
	description, _ := data["description"]
	if !nameOk {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "missing fields and values",
		})
		ctx.Done()
		return
	}

	created, _, err := service.AddNewEvent(name, description)
	if err != nil {
		handleServerError("Error adding new event", err, ctx)
		return
	}
	if created == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "Event created but couldn't fetch",
		})
		ctx.Done()
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"created": created,
		})
	}
}

func GetEvents(ctx *gin.Context) {
	eventRows, err := service.GetAllEvents()
	if err != nil {
		handleServerError("Error getting all events", err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"list": eventRows,
	})
	ctx.Done()
}

func handleServerError(caller string, err error, ctx *gin.Context) {
	fmt.Println(caller, ": -> ", err.Error())
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"code":    http.StatusInternalServerError,
		"message": err.Error(),
	})
	ctx.Done()
}
