package service_test

import (
	"MB-test/src/internal/service"
	"MB-test/src/models"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) ListOrders() ([]models.Orders, error) {
	args := m.Called()
	return args.Get(0).([]models.Orders), args.Error(1)
}

func (m *MockRepo) CreateOrder(order models.Orders) (models.Orders, error) {
	args := m.Called(order)
	return args.Get(0).(models.Orders), args.Error(1)
}

func (m *MockRepo) GetClientById(id string) (models.Client, error) {
	args := m.Called(id)
	return args.Get(0).(models.Client), args.Error(1)
}

func (m *MockRepo) FindMatchOrderToSell(order models.Orders) (models.Orders, error) {
	args := m.Called(order)
	return args.Get(0).(models.Orders), args.Error(1)
}

func (m *MockRepo) FindMatchOrderToBuy(order models.Orders) (models.Orders, error) {
	args := m.Called(order)
	return args.Get(0).(models.Orders), args.Error(1)
}

func (m *MockRepo) MakeTransactionSell(buy models.Orders, sell models.Orders) error {
	args := m.Called(buy, sell)
	return args.Error(0)
}

func (m *MockRepo) MakeTransactionBuy(buy models.Orders, sell models.Orders) error {
	args := m.Called(buy, sell)
	return args.Error(0)
}

func (m *MockRepo) GetOrderById(id string) (models.Orders, error) {
	args := m.Called(id)
	return args.Get(0).(models.Orders), args.Error(1)
}

func (m *MockRepo) UpdateStatusOrder(status int, orderId string) (models.Orders, error) {
	args := m.Called(status, orderId)
	return args.Get(0).(models.Orders), args.Error(1)
}

// Testes:
// CreateOrder:
// [x] Deve falhar caso não encontre um cliente no banco relacionado aquela order
// [x] Deve falhar se houver algum problema em conectar no banco
// [x] Deve falhar se o tipo da order for menor que 1
// [x] Deve falhar se o tipo da order for maior que 2
// [x] Deve falhar se o status da order for menor que 1
// [x] Deve falhar se o status da order for maior que 4
// [x] Deve falhar se o preco dar order em BRL for menor ou igual a 0
// [x] Deve falhar se o preco dar order em BT for menor ou igual a 0
// [x] Deve criar uma order com sucesso e retornar o id dela.
func TestCreateOrder(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := service.NewService(mockRepo)
	svc.DisableAsync = true

	client := models.Client{
		Id:         uuid.MustParse("0ee49ba6-30e6-4b8e-bfec-5bda90aa48ca"),
		BalanceBRL: 12500,
		BalanceBT:  8,
		Score:      98,
	}

	t.Run("Must create an order successfully", func(t *testing.T) {
		order := models.Orders{
			Id:            uuid.MustParse("b794a8dc-415e-435c-8a44-551cf8244e68"),
			TypeOrder:     1,
			Status:        1,
			PriceOrderBT:  100,
			PriceOrderBRL: 500,
			OwnerOrderId:  uuid.MustParse("0ee49ba6-30e6-4b8e-bfec-5bda90aa48ca"),
		}
		mockRepo.On("GetClientById", order.OwnerOrderId.String()).Return(client, nil)
		mockRepo.On("CreateOrder", mock.Anything).Return(order, nil)

		id, err := svc.CreateOrder(order)

		assert.NoError(t, err)
		assert.NotEmpty(t, id)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Should fail if a customer is not found in the database related to that order.", func(t *testing.T) {
		order := models.Orders{
			Id:            uuid.MustParse("b794a8dc-415e-435c-8a44-551cf8244e68"),
			TypeOrder:     1,
			Status:        1,
			PriceOrderBT:  100,
			PriceOrderBRL: 500,
			OwnerOrderId:  uuid.MustParse("0ee49ba6-30e6-4b8e-bfec-5bda90aa48ca"),
		}
		order.OwnerOrderId = uuid.MustParse("a7402f4d-e180-4963-bcc7-e02371a39dca")

		mockRepo.On("GetClientById", "a7402f4d-e180-4963-bcc7-e02371a39dca").Return(models.Client{}, models.ErrorNotFound)

		id, err := svc.CreateOrder(order)

		assert.Error(t, err)
		assert.Empty(t, id)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Should fail if the order type is less than 1.", func(t *testing.T) {
		order := models.Orders{
			Id:            uuid.MustParse("b794a8dc-415e-435c-8a44-551cf8244e68"),
			TypeOrder:     -5,
			Status:        1,
			PriceOrderBT:  100,
			PriceOrderBRL: 500,
			OwnerOrderId:  uuid.MustParse("0ee49ba6-30e6-4b8e-bfec-5bda90aa48ca"),
		}
		mockRepo.On("GetClientById", order.OwnerOrderId.String()).Return(client, nil)

		id, err := svc.CreateOrder(order)

		assert.Error(t, err)
		assert.Empty(t, id)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Should fail if the order type is greater than 2.", func(t *testing.T) {
		order := models.Orders{
			Id:            uuid.MustParse("b794a8dc-415e-435c-8a44-551cf8244e68"),
			TypeOrder:     7,
			Status:        1,
			PriceOrderBT:  100,
			PriceOrderBRL: 500,
			OwnerOrderId:  uuid.MustParse("0ee49ba6-30e6-4b8e-bfec-5bda90aa48ca"),
		}

		mockRepo.On("GetClientById", order.OwnerOrderId.String()).Return(client, nil)

		id, err := svc.CreateOrder(order)

		assert.Error(t, err)
		assert.Empty(t, id)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Should fail if the order status is less than 1.", func(t *testing.T) {
		order := models.Orders{
			Id:            uuid.MustParse("b794a8dc-415e-435c-8a44-551cf8244e68"),
			TypeOrder:     1,
			Status:        -1,
			PriceOrderBT:  100,
			PriceOrderBRL: 500,
			OwnerOrderId:  uuid.MustParse("0ee49ba6-30e6-4b8e-bfec-5bda90aa48ca"),
		}
		mockRepo.On("GetClientById", order.OwnerOrderId.String()).Return(client, nil)

		id, err := svc.CreateOrder(order)

		assert.Error(t, err)
		assert.Empty(t, id)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Should fail if the order status is greater than 4.", func(t *testing.T) {
		order := models.Orders{
			Id:            uuid.MustParse("b794a8dc-415e-435c-8a44-551cf8244e68"),
			TypeOrder:     1,
			Status:        7,
			PriceOrderBT:  100,
			PriceOrderBRL: 500,
			OwnerOrderId:  uuid.MustParse("0ee49ba6-30e6-4b8e-bfec-5bda90aa48ca"),
		}
		mockRepo.On("GetClientById", order.OwnerOrderId.String()).Return(client, nil)

		id, err := svc.CreateOrder(order)

		assert.Error(t, err)
		assert.Empty(t, id)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Should fail if the order price in BRL is less than or equal to 0.", func(t *testing.T) {
		order := models.Orders{
			Id:            uuid.MustParse("b794a8dc-415e-435c-8a44-551cf8244e68"),
			TypeOrder:     1,
			Status:        1,
			PriceOrderBT:  100,
			PriceOrderBRL: -5.40,
			OwnerOrderId:  uuid.MustParse("0ee49ba6-30e6-4b8e-bfec-5bda90aa48ca"),
		}
		mockRepo.On("GetClientById", order.OwnerOrderId.String()).Return(client, nil)

		id, err := svc.CreateOrder(order)

		assert.Error(t, err)
		assert.Empty(t, id)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Should fail if the order price in BT is less than or equal to 0.", func(t *testing.T) {
		order := models.Orders{
			Id:            uuid.MustParse("b794a8dc-415e-435c-8a44-551cf8244e68"),
			TypeOrder:     1,
			Status:        1,
			PriceOrderBT:  -750,
			PriceOrderBRL: 500,
			OwnerOrderId:  uuid.MustParse("0ee49ba6-30e6-4b8e-bfec-5bda90aa48ca"),
		}
		mockRepo.On("GetClientById", order.OwnerOrderId.String()).Return(client, nil)

		id, err := svc.CreateOrder(order)

		assert.Error(t, err)
		assert.Empty(t, id)

		mockRepo.AssertExpectations(t)
	})

}

// ListOrders:
// [x] Deve retorar um array com todas as orders cadastradas
func TestListOrders(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := service.NewService(mockRepo)

	orders := []models.Orders{
		{
			Id:            uuid.New(),
			TypeOrder:     1,
			Status:        1,
			PriceOrderBT:  1,
			PriceOrderBRL: 500,
			OwnerOrderId:  uuid.New(),
		},
		{
			Id:            uuid.New(),
			TypeOrder:     2,
			Status:        1,
			PriceOrderBT:  2,
			PriceOrderBRL: 6000,
			OwnerOrderId:  uuid.New(),
		},
		{
			Id:            uuid.New(),
			TypeOrder:     1,
			Status:        2,
			PriceOrderBT:  4,
			PriceOrderBRL: 2500,
			OwnerOrderId:  uuid.New(),
		},
		{
			Id:            uuid.New(),
			TypeOrder:     1,
			Status:        3,
			PriceOrderBT:  3,
			PriceOrderBRL: 1520,
			OwnerOrderId:  uuid.New(),
		},
		{
			Id:            uuid.New(),
			TypeOrder:     2,
			Status:        4,
			PriceOrderBT:  6,
			PriceOrderBRL: 652,
			OwnerOrderId:  uuid.New(),
		},
	}

	t.Run("Must return an array with all registered orders", func(t *testing.T) {
		mockRepo.On("ListOrders").Return(orders, nil)

		res, err := svc.ListOrders()

		assert.NoError(t, err)
		assert.NotEmpty(t, res)
		assert.Len(t, res, 5)

		mockRepo.AssertExpectations(t)
	})
}

// // UpdateStatusOrder:
// // [X] Deve falhar se o status da order for menor que 1
// // [X] Deve falhar se o status da order for maior que 4
// // [X] Deve falhar caso não encontre a order
// // [x] Deve falhar se a order já estiver done
// // [x] Deve falhar se a order já estiver cancela
// // [x] Deve retornar "status in effect for this order" se o status enviado for o mesmo que já se encontra na order
// // [x] Deve falhar se uma order estiver waiting e tentarem mudar o status para algo diferente de OPEN ou CANCEL
// // [x] Deve atualizar com sucesso
func TestUpdateStatusOrder(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := service.NewService(mockRepo)

	order := models.Orders{
		Id:            uuid.MustParse("b794a8dc-415e-435c-8a44-551cf8244e68"),
		TypeOrder:     1,
		Status:        1,
		PriceOrderBT:  100,
		PriceOrderBRL: 500,
		OwnerOrderId:  uuid.MustParse("0ee49ba6-30e6-4b8e-bfec-5bda90aa48ca"),
	}

	t.Run("Should fail if the order status is less than 1", func(t *testing.T) {
		res, err := svc.UpdateStatusOrder(-5, order.Id.String())

		assert.Error(t, err)
		assert.Empty(t, res)
		assert.Equal(t, "INVALID_INPUT: invalid status", err.Error())

		mockRepo.AssertExpectations(t)
	})

	t.Run("Should fail if the order status is greater than 4", func(t *testing.T) {
		res, err := svc.UpdateStatusOrder(7, order.Id.String())

		assert.Error(t, err)
		assert.Empty(t, res)
		assert.Equal(t, "INVALID_INPUT: invalid status", err.Error())

		mockRepo.AssertExpectations(t)
	})

	t.Run("Should fail if the order cannot be found", func(t *testing.T) {
		order.Id = uuid.MustParse("a7402f4d-e180-4963-bcc7-e02371a39dca")

		mockRepo.On("GetOrderById", "a7402f4d-e180-4963-bcc7-e02371a39dca").Return(models.Orders{}, errors.New("record not found"))
		res, err := svc.UpdateStatusOrder(2, order.Id.String())

		assert.Error(t, err)
		assert.Empty(t, res)
		assert.Equal(t, "NOT_FOUND: record not found", err.Error())
	})

	t.Run("Should fail if the order is already done", func(t *testing.T) {
		order.Status = 3

		mockRepo.On("GetOrderById", "b794a8dc-415e-435c-8a44-551cf8244e68").Return(order, nil)
		res, err := svc.UpdateStatusOrder(2, "b794a8dc-415e-435c-8a44-551cf8244e68")

		assert.Error(t, err)
		assert.Empty(t, res)
		assert.Equal(t, "INVALID_INPUT: invalid update, this order was done", err.Error())
	})

	t.Run("Should fail if the order is already canceled", func(t *testing.T) {
		orderT := models.Orders{
			Id:            uuid.MustParse("096338a2-8bc6-4d4c-a8e0-d395981b7031"),
			TypeOrder:     1,
			Status:        4,
			PriceOrderBT:  100,
			PriceOrderBRL: 500,
			OwnerOrderId:  uuid.MustParse("0ee49ba6-30e6-4b8e-bfec-5bda90aa48ca"),
		}

		mockRepo.On("GetOrderById", "096338a2-8bc6-4d4c-a8e0-d395981b7031").Return(orderT, nil)
		res, err := svc.UpdateStatusOrder(2, "096338a2-8bc6-4d4c-a8e0-d395981b7031")

		assert.Error(t, err)
		assert.Empty(t, res)
		assert.Equal(t, "INVALID_INPUT: invalid update, this order was cancel", err.Error())
	})

	t.Run("Should return 'status in effect for this order' if the status sent is the same as the one already on the order", func(t *testing.T) {
		orderT := models.Orders{
			Id:            uuid.MustParse("f5ca0998-a0c1-4e6a-bdc9-f70521f9154f"),
			TypeOrder:     1,
			Status:        1,
			PriceOrderBT:  100,
			PriceOrderBRL: 500,
			OwnerOrderId:  uuid.MustParse("0ee49ba6-30e6-4b8e-bfec-5bda90aa48ca"),
		}

		mockRepo.On("GetOrderById", "f5ca0998-a0c1-4e6a-bdc9-f70521f9154f").Return(orderT, nil)
		res, err := svc.UpdateStatusOrder(1, "f5ca0998-a0c1-4e6a-bdc9-f70521f9154f")

		assert.NoError(t, err)
		assert.NotEmpty(t, res)
		assert.Equal(t, "status in effect for this order", res)
	})

	t.Run("Should fail if an order is waiting and an attempt is made to change the status to something other than OPEN or CANCEL", func(t *testing.T) {
		orderT := models.Orders{
			Id:            uuid.MustParse("f9c1554a-3fde-4619-8e0e-c0b56a752ede"),
			TypeOrder:     1,
			Status:        2,
			PriceOrderBT:  100,
			PriceOrderBRL: 500,
			OwnerOrderId:  uuid.MustParse("0ee49ba6-30e6-4b8e-bfec-5bda90aa48ca"),
		}

		mockRepo.On("GetOrderById", "f9c1554a-3fde-4619-8e0e-c0b56a752ede").Return(orderT, nil)
		res, err := svc.UpdateStatusOrder(3, "f9c1554a-3fde-4619-8e0e-c0b56a752ede")

		assert.Error(t, err)
		assert.Empty(t, res)
		assert.Equal(t, "INVALID_INPUT: invalid update, An order waiting only change status to OPEN or CANCEL", err.Error())
	})

	t.Run("Should update successfully", func(t *testing.T) {
		orderT := models.Orders{
			Id:            uuid.MustParse("abe7dffa-9ecc-4d40-b7d3-a5e2aca31f28"),
			TypeOrder:     1,
			Status:        1,
			PriceOrderBT:  100,
			PriceOrderBRL: 500,
			OwnerOrderId:  uuid.MustParse("0ee49ba6-30e6-4b8e-bfec-5bda90aa48ca"),
		}

		orderF := models.Orders{
			Id:            uuid.MustParse("abe7dffa-9ecc-4d40-b7d3-a5e2aca31f28"),
			TypeOrder:     1,
			Status:        3,
			PriceOrderBT:  100,
			PriceOrderBRL: 500,
			OwnerOrderId:  uuid.MustParse("0ee49ba6-30e6-4b8e-bfec-5bda90aa48ca"),
		}

		mockRepo.On("GetOrderById", "abe7dffa-9ecc-4d40-b7d3-a5e2aca31f28").Return(orderT, nil)
		mockRepo.On("UpdateStatusOrder", 3, "abe7dffa-9ecc-4d40-b7d3-a5e2aca31f28").Return(orderF, nil)
		res, err := svc.UpdateStatusOrder(3, "abe7dffa-9ecc-4d40-b7d3-a5e2aca31f28")

		assert.NoError(t, err)
		assert.NotEmpty(t, res)
		expect := fmt.Sprintf("order %s updated", orderF.Id)
		assert.Equal(t, expect, "order abe7dffa-9ecc-4d40-b7d3-a5e2aca31f28 updated")
	})
}

// // GetClientById:
// // [x] Deve falhar caso não encontre o cliente
// // [x] Deve retornar o cliente
func TestGetClientById(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := service.NewService(mockRepo)

	t.Run("Should fail if the client cannot be found", func(t *testing.T) {
		client := models.Client{
			Id:         uuid.MustParse("bc7e77eb-12d1-4e3a-b4af-8682302dc0b4"),
			BalanceBRL: 12455,
			BalanceBT:  14,
			Score:      91,
		}
		mockRepo.On("GetClientById", client.Id.String()).Return(models.Client{}, errors.New("record not found"))
		res, err := svc.GetClientById(client.Id.String())

		assert.Error(t, err)
		assert.Empty(t, res)
		assert.Equal(t, "record not found", err.Error())
	})

	t.Run("Must return the client successfully", func(t *testing.T) {
		client := models.Client{
			Id:         uuid.MustParse("0b20b052-abd2-4da8-ac7e-5632118be457"),
			BalanceBRL: 12455,
			BalanceBT:  14,
			Score:      91,
		}
		mockRepo.On("GetClientById", client.Id.String()).Return(client, nil)
		res, err := svc.GetClientById(client.Id.String())

		expect := models.ClientDtoOutput{
			Id:         client.Id,
			BalanceBRL: client.BalanceBRL,
			BalanceBT:  client.BalanceBT,
			Score:      client.Score,
		}
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
		assert.Equal(t, expect, res)
	})
}
