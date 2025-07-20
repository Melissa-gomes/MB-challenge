package main

import (
	"MB-test/src/configs"
	"MB-test/src/internal/controller"
	"MB-test/src/internal/repository"
	"MB-test/src/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {

	env := configs.LoadEnv()
	db := configs.NewDatabase(env)
	configs.MigrateDb(db)
	configs.Seeders(db)

	repo := repository.NewRepository(db)
	svc := service.NewService(repo)
	ctl := controller.NewController(svc)

	router := gin.New()
	router.POST("/orders", ctl.CreateOrder)
	router.PATCH("/orders/:orderId/status/:status", ctl.UpdateStatusOrder)
	router.GET("/orders", ctl.ListOrders)
	router.GET("/client/:id", ctl.GetClientById)

	router.Run()
}
