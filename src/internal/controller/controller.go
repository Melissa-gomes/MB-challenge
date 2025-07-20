package controller

import (
	"MB-test/src/internal/contracts"
	"MB-test/src/models"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	Service contracts.OperationsServiceHandler
}

func NewController(service contracts.OperationsServiceHandler) Controller {
	return Controller{Service: service}
}

func (c Controller) CreateOrder(ctx *gin.Context) {
	var order models.Orders
	err := ctx.ShouldBindJSON(&order)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	res, err := c.Service.CreateOrder(order)
	if err != nil {
		var appErr models.Error
		if errors.As(err, &appErr) {
			ctx.JSON(appErr.StatusCode, gin.H{
				"error":  appErr.Message,
				"kind":   appErr.Kind,
				"status": appErr.StatusCode,
			})
		}
	} else {
		ctx.JSON(http.StatusCreated, gin.H{
			"data": res,
		})
	}
}

func (c Controller) ListOrders(ctx *gin.Context) {
	res, err := c.Service.ListOrders()
	if err != nil {
		var appErr models.Error
		if errors.As(err, &appErr) {
			ctx.JSON(appErr.StatusCode, gin.H{
				"error":  appErr.Message,
				"kind":   appErr.Kind,
				"status": appErr.StatusCode,
			})
		}
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"data": res,
		})
	}
}

func (c Controller) UpdateStatusOrder(ctx *gin.Context) {
	orderId := ctx.Param("orderId")
	s := ctx.Param("status")
	status, err := strconv.Atoi(s)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	res, err := c.Service.UpdateStatusOrder(status, orderId)
	if err != nil {
		var appErr models.Error
		if errors.As(err, &appErr) {
			ctx.JSON(appErr.StatusCode, gin.H{
				"error":  appErr.Message,
				"kind":   appErr.Kind,
				"status": appErr.StatusCode,
			})
		}
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"data": res,
		})
	}
}

func (c Controller) GetClientById(ctx *gin.Context) {
	id := ctx.Param("id")

	res, err := c.Service.GetClientById(id)
	if err != nil {
		var appErr models.Error
		if errors.As(err, &appErr) {
			ctx.JSON(appErr.StatusCode, gin.H{
				"error":  appErr.Message,
				"kind":   appErr.Kind,
				"status": appErr.StatusCode,
			})
		}
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"data": res,
		})
	}
}
