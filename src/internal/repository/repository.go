package repository

import (
	"MB-test/src/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

// ORDERS
func (r Repository) CreateOrder(order models.Orders) (models.Orders, error) {
	if result := r.DB.Create(&order); result.Error != nil {
		return models.Orders{}, models.NewError(models.ErrorKindDatabase, "database error: "+result.Error.Error(), models.StatusCodeInternal)
	}
	return order, nil
}

func (r Repository) ListOrders() ([]models.Orders, error) {
	orders := []models.Orders{}
	if result := r.DB.Find(&orders); result.Error != nil {
		return orders, models.NewError(models.ErrorKindDatabase, "database error: "+result.Error.Error(), models.StatusCodeInternal)
	}

	return orders, nil
}

func (r Repository) GetOrderById(id string) (models.Orders, error) {
	order := models.Orders{}

	if result := r.DB.Find(&order).Where("id = ?", id); result.Error != nil {
		return order, errors.New(result.Error.Error())
	}

	return order, nil
}

func (r Repository) FindMatchOrderToSell(order models.Orders) (models.Orders, error) {
	orderMatch := models.Orders{}

	result := r.DB.Where("price_order_bt = ?", order.PriceOrderBT).
		Where("price_order_brl >= ?", order.PriceOrderBRL).
		Where("status = 1").
		Where("type_order = 1").
		Order("price_order_brl DESC").
		First(&orderMatch)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return orderMatch, models.ErrorNotFound
	}

	if result.Error != nil {
		return orderMatch, models.NewError(models.ErrorKindDatabase, "database error: "+result.Error.Error(), models.StatusCodeInternal)
	}

	return orderMatch, nil
}

func (r Repository) FindMatchOrderToBuy(order models.Orders) (models.Orders, error) {
	orderMatch := models.Orders{}

	result := r.DB.Where("price_order_bt = ?", order.PriceOrderBT).
		Where("price_order_brl <= ?", order.PriceOrderBRL).
		Where("status = 1").
		Where("type_order = 2").
		Order("price_order_brl ASC").
		First(&orderMatch)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return orderMatch, models.ErrorNotFound
	}

	if result.Error != nil {
		return orderMatch, models.NewError(models.ErrorKindDatabase, "database error: "+result.Error.Error(), models.StatusCodeInternal)
	}

	return orderMatch, nil
}

func (r Repository) MakeTransactionBuy(buyOrder, sellOrder models.Orders) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	clientBuyer := models.Client{}
	clientSeller := models.Client{}

	if err := tx.Where("id = ?", buyOrder.OwnerOrderId).First(&clientBuyer).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("err to found client: %w", err)
	}

	if err := tx.Where("id = ?", sellOrder.OwnerOrderId).First(&clientSeller).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("err to found client: %w", err)
	}

	if clientBuyer.BalanceBRL < sellOrder.PriceOrderBRL || clientSeller.BalanceBT < sellOrder.PriceOrderBT {
		tx.Rollback()
		return fmt.Errorf("customer with insufficient balance for this transaction")
	}

	if err := tx.Model(&clientBuyer).
		Updates(map[string]interface{}{
			"balance_brl": gorm.Expr("balance_brl - ?", sellOrder.PriceOrderBRL),
			"balance_bt":  gorm.Expr("balance_bt + ?", sellOrder.PriceOrderBT),
		}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("erro in transaction: %w", err)
	}

	if err := tx.Model(&clientSeller).
		Updates(map[string]interface{}{
			"balance_brl": gorm.Expr("balance_brl + ?", sellOrder.PriceOrderBRL),
			"balance_bt":  gorm.Expr("balance_bt - ?", sellOrder.PriceOrderBT),
		}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("erro in transaction: %w", err)
	}

	if err := tx.Model(&buyOrder).Update("status", 3).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("erro to update orders status: %w", err)
	}

	if err := tx.Model(&sellOrder).Update("status", 3).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("erro to update orders status: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("err in commit: %w", err)
	}

	return nil
}

func (r Repository) MakeTransactionSell(buyOrder, sellOrder models.Orders) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	clientBuyer := models.Client{}
	clientSeller := models.Client{}

	if err := tx.Where("id = ?", buyOrder.OwnerOrderId).First(&clientBuyer).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("err to found client: %w", err)
	}

	if err := tx.Where("id = ?", sellOrder.OwnerOrderId).First(&clientSeller).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("err to found client: %w", err)
	}

	if clientBuyer.BalanceBRL < sellOrder.PriceOrderBRL || clientSeller.BalanceBT < sellOrder.PriceOrderBT {
		tx.Rollback()
		return fmt.Errorf("customer with insufficient balance for this transaction")
	}

	if err := tx.Model(&clientBuyer).
		Updates(map[string]interface{}{
			"balance_brl": gorm.Expr("balance_brl - ?", buyOrder.PriceOrderBRL),
			"balance_bt":  gorm.Expr("balance_bt + ?", buyOrder.PriceOrderBT),
		}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("erro in transaction: %w", err)
	}

	if err := tx.Model(&clientSeller).
		Updates(map[string]interface{}{
			"balance_brl": gorm.Expr("balance_brl + ?", buyOrder.PriceOrderBRL),
			"balance_bt":  gorm.Expr("balance_bt - ?", buyOrder.PriceOrderBT),
		}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("erro in transaction: %w", err)
	}

	if err := tx.Model(&buyOrder).Update("status", 3).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("erro to update orders status: %w", err)
	}

	if err := tx.Model(&sellOrder).Update("status", 3).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("erro to update orders status: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("err in commit: %w", err)
	}

	return nil
}

func (r Repository) UpdateStatusOrder(status int, orderId string) (models.Orders, error) {
	order := models.Orders{}
	if result := r.DB.Model(&order).Where("id = ?", orderId).Update("status", status); result.Error != nil {
		return models.Orders{}, models.NewError(models.ErrorKindDatabase, "database error: "+result.Error.Error(), models.StatusCodeInternal)
	}
	return order, nil
}

// CLIENTS
func (r Repository) GetClientById(id string) (models.Client, error) {
	var client models.Client

	result := r.DB.Where("id = ?", id).First(&client)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return client, models.ErrorNotFound
	}

	if result.Error != nil {
		return client, models.NewError(models.ErrorKindDatabase, "database error: "+result.Error.Error(), models.StatusCodeInternal)
	}

	return client, nil
}
