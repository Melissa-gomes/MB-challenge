package service

import (
	"MB-test/src/internal/contracts"
	"MB-test/src/models"
	"errors"
	"fmt"
	"reflect"

	"github.com/google/uuid"
)

type Service struct {
	Repo contracts.OperationsRepositoryHandle
}

func NewService(repo contracts.OperationsRepositoryHandle) *Service {
	return &Service{Repo: repo}
}

func (s Service) ListOrders() ([]models.OrderDtoOutput, error) {
	orders, err := s.Repo.ListOrders()
	if err != nil {
		return []models.OrderDtoOutput{}, err
	}

	var result []models.OrderDtoOutput
	for _, o := range orders {
		result = append(result, models.OrderDtoOutput{
			Id:            o.Id,
			OwnerOrderId:  o.OwnerOrderId,
			PriceOrderBRL: o.PriceOrderBRL,
			PriceOrderBT:  o.PriceOrderBT,
			TypeOrder:     models.TranslateStatus(o.Status),
			Status:        models.TranslateTypeOrder(o.TypeOrder),
		})
	}

	return result, nil
}

func (s Service) CreateOrder(order models.Orders) (string, error) {
	owner, err := s.Repo.GetClientById(order.OwnerOrderId.String())
	if err != nil {
		return "", err
	}
	if reflect.DeepEqual(owner, models.Client{}) {
		return "", models.ErrorNotFound
	}

	if order.TypeOrder < 1 || order.TypeOrder > 2 {
		return "", models.ErrorInvalidTypeOrder
	}

	if order.Status < 1 || order.Status > 4 {
		return "", models.ErrorInvalidStatus
	}

	if order.PriceOrderBRL <= 0 || order.PriceOrderBT <= 0 {
		return "", models.ErrorInvalidPriceOrder
	}

	order.Id = uuid.New()

	res, err := s.Repo.CreateOrder(order)
	if err != nil {
		return "", err
	}

	err = s.FindMatchOrder(order)
	if err != nil {
		if errors.Is(err, models.ErrorNotFound) {
			return res.Id.String(), nil
		}
	}

	return res.Id.String(), nil
}

func (s Service) FindMatchOrder(orderToMatch models.Orders) error {
	if orderToMatch.TypeOrder == models.SELL {
		orderMatched, err := s.Repo.FindMatchOrderToSell(orderToMatch)
		if err != nil {
			if errors.Is(err, models.ErrorNotFound) {
				return nil
			}
			return err
		}
		return s.Repo.MakeTransactionSell(orderMatched, orderToMatch)

	} else {
		orderMatched, err := s.Repo.FindMatchOrderToBuy(orderToMatch)
		if err != nil {
			if errors.Is(err, models.ErrorNotFound) {
				return nil
			}
			return err
		}
		return s.Repo.MakeTransactionBuy(orderToMatch, orderMatched)
	}
}

func (s Service) UpdateStatusOrder(status int, orderId string) (string, error) {
	if status < 1 || status > 4 {
		return "", models.ErrorInvalidStatus
	}

	order, err := s.Repo.GetOrderById(orderId)
	if err != nil {
		return "", models.ErrorNotFound
	}

	if order.Status == models.DONE {
		return "", models.ErrorInvalidUpdateOrderDone
	}

	if order.Status == models.CANCEL {
		return "", models.ErrorInvalidUpdateOrderCancel
	}

	if order.Status == status {
		return "status in effect for this order", nil
	}

	if order.Status == models.WAITING && status != models.OPEN && status != models.CANCEL {
		return "", models.ErrorInvalidUpdateOrderWaiting
	}

	_, err = s.Repo.UpdateStatusOrder(status, orderId)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("order %s updated", orderId), nil
}

func (s Service) GetClientById(id string) (models.ClientDtoOutput, error) {
	client, err := s.Repo.GetClientById(id)
	if err != nil {
		return models.ClientDtoOutput{}, err
	}

	return models.ClientDtoOutput{
		Id:         client.Id,
		BalanceBRL: client.BalanceBRL,
		BalanceBT:  client.BalanceBT,
		Score:      client.Score,
	}, nil
}
