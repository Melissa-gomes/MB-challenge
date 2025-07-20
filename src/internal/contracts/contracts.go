package contracts

import (
	"MB-test/src/models"
)

type OperationsServiceHandler interface {
	CreateOrder(order models.Orders) (string, error)
	ListOrders() ([]models.OrderDtoOutput, error)
	GetClientById(id string) (models.ClientDtoOutput, error)
	UpdateStatusOrder(status int, orderId string) (string, error)
}

type OperationsRepositoryHandle interface {
	CreateOrder(order models.Orders) (models.Orders, error)
	UpdateStatusOrder(status int, orderId string) (models.Orders, error)
	GetClientById(id string) (models.Client, error)
	ListOrders() ([]models.Orders, error)
	GetOrderById(id string) (models.Orders, error)
	FindMatchOrderToSell(order models.Orders) (models.Orders, error)
	FindMatchOrderToBuy(order models.Orders) (models.Orders, error)
	MakeTransactionBuy(buyOrder, sellOrder models.Orders) error
	MakeTransactionSell(buyOrder, sellOrder models.Orders) error
}
