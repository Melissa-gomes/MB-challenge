package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ClientDtoOutput struct {
	Id         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id,omitempty"`
	BalanceBRL float64   `json:"balance_brl"`
	BalanceBT  float64   `json:"balance_bt"`
	Score      int       `json:"score,omitempty"`
}

type OrderDtoOutput struct {
	Id            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id,omitempty"`
	OwnerOrderId  uuid.UUID `gorm:"type:uuid;not null" json:"owner_order_id"` // Chave estrangeira
	PriceOrderBRL float64   `json:"price_order_brl"`
	PriceOrderBT  float64   `json:"price_order_bt"`
	TypeOrder     string    `json:"type_order"`
	Status        string    `json:"status,omitempty"`
}

type Client struct {
	Id         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id,omitempty"`
	BalanceBRL float64   `json:"balance_brl"`
	BalanceBT  float64   `json:"balance_bt"`
	Score      int       `json:"score,omitempty"`
	CreatedAt  time.Time `json:"created_at" gorm:"default:now()`

	Orders []Orders `gorm:"foreignKey:OwnerOrderId" json:"orders,omitempty"`
}

type Orders struct {
	Id            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id,omitempty"`
	OwnerOrderId  uuid.UUID `gorm:"type:uuid;not null" json:"owner_order_id"` // Chave estrangeira
	PriceOrderBRL float64   `json:"price_order_brl"`
	PriceOrderBT  float64   `json:"price_order_bt"`
	TypeOrder     int       `json:"type_order"`
	Status        int       `json:"status,omitempty"`
	CreatedAt     time.Time `json:"created_at" gorm:"default:now()"`

	Client Client `gorm:"foreignKey:OwnerOrderId;references:Id" json:"client"` // Relacionamento
}

const (
	OPEN = iota + 1
	WAITING
	DONE
	CANCEL
)

const (
	BUY = iota + 1
	SELL
)

func TranslateStatus(status int) string {
	switch status {
	case 1:
		return "OPEN"
	case 2:
		return "WAITING"
	case 3:
		return "DONE"
	case 4:
		return "CANCEL"
	default:
		return "status not found"
	}
}

func TranslateTypeOrder(status int) string {
	switch status {
	case 1:
		return "BUY"
	case 2:
		return "SELL"
	default:
		return "type_order not found"
	}
}

type Error struct {
	Message    string
	Kind       ErrorKind
	StatusCode int
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Kind, e.Message)
}

type ErrorKind string

func NewError(kind ErrorKind, message string, statusCode int) Error {
	return Error{Message: message, Kind: kind, StatusCode: statusCode}
}

const (
	ErrorKindNotFound     ErrorKind = "NOT_FOUND"
	ErrorKindInvalidInput ErrorKind = "INVALID_INPUT"
	ErrorKindBadRequest   ErrorKind = "BAD_REQUEST"
	ErrorKindForbidden    ErrorKind = "FORBIDDEN"
	ErrorKindInternal     ErrorKind = "INTERNAL_SERVER_ERROR"
	ErrorKindDatabase     ErrorKind = "DATABASE_ERROR"
)

const (
	ErrorMessageNotFound     string = "record not found"
	ErrorMessageInvalidInput string = "invalid input:"
	ErrorMessageBadRequest   string = "bad request"
	ErrorMessageForbidden    string = "forbidden"
	ErrorMessageInternal     string = "internal server error"
)

const (
	StatusCodeNotFound     int = 404
	StatusCodeInvalidInput int = 422
	StatusCodeBadRequest   int = 400
	StatusCodeForbidden    int = 403
	StatusCodeInternal     int = 500
)

var (
	ErrorNotFound                  = NewError(ErrorKindNotFound, ErrorMessageNotFound, StatusCodeNotFound)
	ErrorInvalidTypeOrder          = NewError(ErrorKindInvalidInput, "invalid type_order", StatusCodeInvalidInput)
	ErrorInvalidStatus             = NewError(ErrorKindInvalidInput, "invalid status", StatusCodeInvalidInput)
	ErrorInvalidPriceOrder         = NewError(ErrorKindInvalidInput, "It is not allowed to create orders with a price less than or equal to 0", StatusCodeInvalidInput)
	ErrorInvalidUpdateOrderDone    = NewError(ErrorKindInvalidInput, "invalid update, this order was done", StatusCodeInvalidInput)
	ErrorInvalidUpdateOrderCancel  = NewError(ErrorKindInvalidInput, "invalid update, this order was cancel", StatusCodeInvalidInput)
	ErrorInvalidUpdateOrderWaiting = NewError(ErrorKindInvalidInput, "invalid update, An order waiting only change status to OPEN or CANCEL", StatusCodeInvalidInput)
	ErrorInsufficientBalance       = NewError(ErrorKindInvalidInput, "insufficient balance", StatusCodeInvalidInput)
)
