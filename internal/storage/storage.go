package storage

import (
	"GOHW-1/internal/model"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"slices"
	"strconv"
	"sync"
	"time"
)

var availableOrdersFileName = "available_orders.json"
var refundedOrdersFileName = "refunded_orders.json"
var pickUpPointsFileName = "pick_up_points.json"

type Storage struct {
	availableOrdersFile *os.File
	refundedOrdersFile  *os.File
	pickUpPointsFile    *os.File
	rwmutex             sync.RWMutex
}

func New() (Storage, error) {
	availableOrdersFile, err := os.OpenFile(availableOrdersFileName, os.O_CREATE, 0777)
	if err != nil {
		return Storage{}, err
	}
	refundedOrdersFile, err := os.OpenFile(refundedOrdersFileName, os.O_CREATE, 0777)
	if err != nil {
		return Storage{}, err
	}
	pickUpPointsFile, err := os.OpenFile(pickUpPointsFileName, os.O_CREATE, 0777)
	if err != nil {
		return Storage{}, err
	}
	return Storage{
			availableOrdersFile: availableOrdersFile,
			refundedOrdersFile:  refundedOrdersFile,
			pickUpPointsFile:    pickUpPointsFile,
			rwmutex:             sync.RWMutex{}},
		nil
}

// CourierTakeOrder accepts and writes order from courier into file
func (s *Storage) CourierTakeOrder(order model.Order) error {
	availableOrders, err := s.GetOrders(availableOrdersFileName)
	if err != nil {
		return err
	}

	if time.Now().After(order.ExpirationDate) {
		return errors.New("order expiration date in the past")
	}
	for _, orderVal := range availableOrders {
		if orderVal.ID == order.ID {
			return errors.New("order has been already accepted")
		}
	}

	newOrder := OrderDTO{
		ID:             order.ID,
		ClientID:       order.ClientID,
		ExpirationDate: order.ExpirationDate,
		Weight:         order.Weight,
		Price:          order.Price,
		Packaging:      order.Packaging,
		IsGiven:        false,
		GivenTime:      time.Time{},
	}

	availableOrders = append(availableOrders, newOrder)
	if err = writeOrders(availableOrders, availableOrdersFileName); err != nil {
		return err
	}
	return nil
}

// CourierGiveOrder deletes given order from file
func (s *Storage) CourierGiveOrder(orderID int) error {
	availableOrders, err := s.GetOrders(availableOrdersFileName)
	if err != nil {
		return err
	}

	for ind, order := range availableOrders {
		if order.ID == orderID {
			if !order.IsGiven && time.Now().After(order.ExpirationDate) {
				availableOrders = append(availableOrders[:ind], availableOrders[ind+1:]...)
				return writeOrders(availableOrders, availableOrdersFileName)
			}
			return errors.New("the order was given or the expiration date is not over yet")
		}
	}

	refundedOrders, err := s.GetOrders(refundedOrdersFileName)
	if err != nil {
		return err
	}

	for ind, order := range refundedOrders {
		if order.ID == orderID {
			refundedOrders = append(refundedOrders[:ind], refundedOrders[ind+1:]...)
			return writeOrders(refundedOrders, refundedOrdersFileName)
		}
	}

	return errors.New("the order was not found")
}

// Write orders into file
func writeOrders(orders []OrderDTO, fileName string) error {
	rawBytes, err := json.Marshal(orders)
	if err != nil {
		return err
	}

	if err = os.WriteFile(fileName, rawBytes, 0777); err != nil {
		return err
	}
	return nil
}

// GetOrders gets slice with all orders
func (s *Storage) GetOrders(fileName string) ([]OrderDTO, error) {
	file := s.availableOrdersFile
	if fileName == availableOrdersFileName {
		file = s.availableOrdersFile
	} else if fileName == refundedOrdersFileName {
		file = s.refundedOrdersFile
	} else {
		return nil, errors.New("file not found")
	}

	reader := bufio.NewReader(file)
	rawBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var orders []OrderDTO
	if len(rawBytes) == 0 {
		return orders, nil
	}

	if err = json.Unmarshal(rawBytes, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

// ClientGiveOrder changes given orders' boolean variables `isGiven` to true
func (s *Storage) ClientGiveOrder(clientID int, ordersID []string) error {
	availableOrders, err := s.GetOrders(availableOrdersFileName)
	if err != nil {
		return err
	}

	isOrderPresent := make(map[int]bool)
	for _, order := range ordersID {
		atoi, err := strconv.Atoi(order)
		if err != nil {
			return err
		}
		isOrderPresent[atoi] = true
	}

	count := 0
	for _, order := range availableOrders {
		if isOrderPresent[order.ID] {
			if order.IsGiven {
				return errors.New("order has already been given")
			}
			if time.Now().After(order.ExpirationDate) {
				return errors.New("order expired")
			}
			if order.ClientID != clientID {
				return errors.New("order does not belong to the client")
			}
			count += 1
		}
	}

	if count < len(ordersID) {
		return errors.New("not all orders are found")
	}

	for ind, order := range availableOrders {
		if isOrderPresent[order.ID] {
			order.IsGiven = true
			order.GivenTime = time.Now()
			availableOrders = append(append(availableOrders[:ind], order), availableOrders[ind+1:]...)
		}
	}

	if err = writeOrders(availableOrders, availableOrdersFileName); err != nil {
		return err
	}
	return nil
}

// ClientGetOrders gets all client orders
func (s *Storage) ClientGetOrders(clientID int, n int, onlyUserOrders bool) ([]model.Order, error) {
	availableOrders, err := s.GetOrders(availableOrdersFileName)
	if err != nil {
		return nil, err
	}

	slices.Reverse(availableOrders)
	clientOrders := make([]model.Order, 0)
	for _, order := range availableOrders {
		if !onlyUserOrders || (order.ClientID == clientID && !order.IsGiven) {
			clientOrders = append(clientOrders, model.Order{
				ID:             order.ID,
				ClientID:       order.ClientID,
				ExpirationDate: order.ExpirationDate,
				Weight:         order.Weight,
				Price:          order.Price,
				Packaging:      order.Packaging,
			})
		}
		if n != -1 && len(clientOrders) >= n {
			break
		}
	}
	return clientOrders, nil
}

// ClientRefund accepts refund from customer
func (s *Storage) ClientRefund(clientID int, orderID int) error {
	availableOrders, err := s.GetOrders(availableOrdersFileName)
	if err != nil {
		return err
	}

	refundedOrders, err := s.GetOrders(refundedOrdersFileName)
	if err != nil {
		return err
	}

	for ind, order := range availableOrders {
		if order.ID == orderID {
			if order.ClientID == clientID {
				if order.IsGiven && time.Since(order.GivenTime) < 48*time.Hour {
					availableOrders = append(availableOrders[:ind], availableOrders[ind+1:]...)
					err = writeOrders(availableOrders, availableOrdersFileName)
					if err != nil {
						return err
					}

					refundedOrders = append(refundedOrders, order)
					err := writeOrders(refundedOrders, refundedOrdersFileName)
					if err != nil {
						return err
					}
					return nil
				}
				return errors.New("it has been more than 2 days since it was given or order was not given")
			}
			return errors.New("order does not belong to the client")
		}
	}
	return errors.New("order was not found")
}

// RefundList returns slice of refunded orders
func (s *Storage) RefundList(pageNumber int) ([]model.Order, error) {
	refundedOrders, err := s.GetOrders(refundedOrdersFileName)
	if err != nil {
		return nil, err
	}

	if len(refundedOrders) < (pageNumber-1)*10+1 {
		return nil, errors.New("page does not exists")
	}

	ordersOnPage := make([]model.Order, 0)
	for _, order := range refundedOrders[(pageNumber-1)*10 : min(pageNumber*10, len(refundedOrders))] {
		ordersOnPage = append(ordersOnPage, model.Order{
			ID:             order.ID,
			ClientID:       order.ClientID,
			ExpirationDate: order.ExpirationDate,
			Weight:         order.Weight,
			Price:          order.Price,
			Packaging:      order.Packaging,
		})
	}
	return ordersOnPage, nil
}

// PickUpPointWrite takes a new pick-up point and adding it into file
func (s *Storage) PickUpPointWrite(pickUpPoint model.PickUpPoint) error {
	pickUpPoints, err := s.PickUpPointsRead()
	if err != nil {
		return err
	}

	pickUpPoint.ID = int64(len(pickUpPoints) + 1)
	pickUpPoints = append(pickUpPoints, pickUpPoint)
	if err = s.writePickUpPoints(pickUpPoints); err != nil {
		return err
	}
	return nil
}

// writePickUpPoints writes a slice of pick-up points into file
func (s *Storage) writePickUpPoints(pickUpPoints []model.PickUpPoint) error {
	rawBytes, err := json.Marshal(pickUpPoints)
	if err != nil {
		return err
	}

	s.rwmutex.Lock()
	err = os.WriteFile(pickUpPointsFileName, rawBytes, 0777)
	s.rwmutex.Unlock()
	return err
}

// PickUpPointsRead gets a slice with all pick-up points
func (s *Storage) PickUpPointsRead() ([]model.PickUpPoint, error) {
	s.rwmutex.RLock()
	defer s.rwmutex.RUnlock()

	file, err := os.Open(pickUpPointsFileName)
	if err != nil {
		return nil, fmt.Errorf("unable to open the file")
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("unable to get the information about the file")
	}

	// Is file empty
	if fileInfo.Size() == 0 {
		return []model.PickUpPoint{}, nil
	}

	rawBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var pickUpPoints []model.PickUpPoint
	if err := json.Unmarshal(rawBytes, &pickUpPoints); err != nil {
		return nil, err
	}

	return pickUpPoints, nil
}
