package service

import (
	"GOHW-1/internal/model"
)

type storage interface {
	CourierTakeOrder(orders model.Order) error
	CourierGiveOrder(orderID int) error
	ClientGiveOrder(clientID int, ordersID []string) error
	ClientGetOrders(clientID int, n int, onlyUserOrders bool) ([]model.Order, error)
	ClientRefund(clientID int, orderID int) error
	RefundList(pageNumber int) ([]model.Order, error)
	PickUpPointWrite(pickUpPoint model.PickUpPoint) error
	PickUpPointsRead() ([]model.PickUpPoint, error)
}

type Service struct {
	storage storage
}

func New(storage storage) Service {
	return Service{storage: storage}
}

// CourierTakeOrder accepts and writes order from courier into file
func (service Service) CourierTakeOrder(order model.Order) error {
	return service.storage.CourierTakeOrder(order)
}

// CourierGiveOrder deletes given order from file
func (service Service) CourierGiveOrder(orderID int) error {
	return service.storage.CourierGiveOrder(orderID)
}

// ClientGiveOrder changes given orders' boolean variables `isGiven` to true
func (service Service) ClientGiveOrder(clientID int, ordersID []string) error {
	return service.storage.ClientGiveOrder(clientID, ordersID)
}

// ClientGetOrders gets all client orders
func (service Service) ClientGetOrders(clientID int, n int, onlyUserOrders bool) ([]model.Order, error) {
	return service.storage.ClientGetOrders(clientID, n, onlyUserOrders)
}

// ClientRefund accepts refund from customer
func (service Service) ClientRefund(clientID int, orderID int) error {
	return service.storage.ClientRefund(clientID, orderID)
}

// RefundList returns slice of refunded orders
func (service Service) RefundList(pageNumber int) ([]model.Order, error) {
	return service.storage.RefundList(pageNumber)
}

// PickUpPointWrite takes a new pick-up point and adding it into file
func (service Service) PickUpPointWrite(pickUpPoint model.PickUpPoint) error {
	return service.storage.PickUpPointWrite(pickUpPoint)
}

// PickUpPointsRead gets a slice with all pick-up points
func (service Service) PickUpPointsRead() ([]model.PickUpPoint, error) {
	return service.storage.PickUpPointsRead()
}
