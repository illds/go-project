package controller

import (
	"GOHW-1/internal/model"
	"GOHW-1/internal/service"
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type OrderController struct {
	service *service.Service
	wg      *sync.WaitGroup
	ctx     context.Context
}

type ReadRequest struct {
	ResponseChan chan []model.PickUpPoint
}

type WriteRequest struct {
	PickUpPoint  model.PickUpPoint
	ResponseChan chan error
}

func NewOrderController(svc *service.Service, wg *sync.WaitGroup, ctx context.Context) *OrderController {
	return &OrderController{
		service: svc,
		wg:      wg,
		ctx:     ctx,
	}
}

func (controller *OrderController) CourierTakeCommand(orderID int, clientID int, availableTime time.Duration, weight float64, price float64, packaging string) error {
	if orderID <= 0 {
		return fmt.Errorf("order ID is not given or incorrect")
	}
	if clientID <= 0 {
		return fmt.Errorf("client ID is not given or incorrect")
	}
	if availableTime <= 0 {
		return fmt.Errorf("available time is not given or incorrect")
	}

	rules := model.GetPackagingRules()
	rule, ok := rules[packaging]

	if !ok {
		return fmt.Errorf("invalid packaging type")
	}

	if rule.MaxWeight > 0 && weight > rule.MaxWeight {
		return fmt.Errorf("order weight exceeds limit for %s", packaging)
	}

	price += rule.ExtraCost

	if err := controller.service.CourierTakeOrder(model.Order{
		ID:             orderID,
		ClientID:       clientID,
		ExpirationDate: time.Now().Add(availableTime),
		Weight:         weight,
		Price:          price,
		Packaging:      packaging,
	}); err != nil {
		return fmt.Errorf("failed to take order from courier: %w", err)
	}

	fmt.Println("Courier's order accepted successfully!")
	return nil
}

func (controller *OrderController) CourierReturnCommand(orderID int) error {
	if orderID <= 0 {
		return fmt.Errorf("order ID is not given or incorrect")
	}

	if err := controller.service.CourierGiveOrder(orderID); err != nil {
		return fmt.Errorf("failed to return order to courier: %w", err)
	}

	fmt.Println("The order was returned to the courier successfully!")
	return nil
}

func (controller *OrderController) ClientGiveCommand(clientID int, ordersID string) error {
	if clientID <= 0 {
		return fmt.Errorf("client ID is not given or incorrect")
	}
	if ordersID == "" {
		return fmt.Errorf("orders ID are not given")
	}
	ordersIDSlice := strings.Split(ordersID, ",")

	if err := controller.service.ClientGiveOrder(clientID, ordersIDSlice); err != nil {
		return fmt.Errorf("failed to give order to client: %w", err)
	}

	fmt.Println("Orders were given to the client successfully!")
	return nil
}

func (controller *OrderController) ClientOrdersCommand(clientID int, N int, onlyUserOrders bool) error {
	if clientID <= 0 {
		return fmt.Errorf("client ID is not given or incorrect")
	}
	if N < -1 || N == 0 {
		return fmt.Errorf("n value is incorrect")
	}

	orders, err := controller.service.ClientGetOrders(clientID, N, onlyUserOrders)
	if err != nil {
		return fmt.Errorf("failed to get orders for client: %w", err)
	}

	if len(orders) == 0 {
		fmt.Println("List of orders is empty")
		return nil
	}
	for _, order := range orders {
		fmt.Printf("%+v\n", order)
	}
	return nil
}

func (controller *OrderController) ClientRefundCommand(orderID int, clientID int) error {
	if orderID <= 0 {
		return fmt.Errorf("order ID is not given or incorrect")
	}
	if clientID <= 0 {
		return fmt.Errorf("client ID is not given or incorrect")
	}

	if err := controller.service.ClientRefund(clientID, orderID); err != nil {
		return fmt.Errorf("failed to refund order: %w", err)
	}

	fmt.Println("Order refunded successfully!")
	return nil
}

func (controller *OrderController) RefundListCommand(pageNumber int) error {
	if pageNumber <= 0 {
		return fmt.Errorf("Page number is  incorrect")
	}
	orders, err := controller.service.RefundList(pageNumber)
	if err != nil {
		return fmt.Errorf("failed to get refund list: %w", err)
	}

	if len(orders) == 0 {
		fmt.Println("List of orders is empty")
		return nil
	}
	fmt.Printf("\t\t\tPage number: %d\n", pageNumber)
	for _, order := range orders {
		fmt.Printf("%+v\n", order)
	}
	fmt.Printf("\t\t\tPage number: %d\n", pageNumber)
	return nil
}

func (controller *OrderController) InteractiveCommand() {
	readRequests := make(chan ReadRequest)
	writeRequests := make(chan WriteRequest)

	// Goroutine for reading
	controller.wg.Add(1)
	go func(ctx context.Context) {
		defer controller.wg.Done()
		for {
			select {
			case req := <-readRequests:
				pickUpPoints, err := controller.service.PickUpPointsRead()
				if err != nil {
					fmt.Println(fmt.Errorf("failed to read pick-up points: %w", err))
					close(req.ResponseChan)
					continue
				}
				req.ResponseChan <- pickUpPoints
			case <-ctx.Done():
				return
			}
		}
	}(controller.ctx)

	// Goroutine for writing
	controller.wg.Add(1)
	go func(ctx context.Context) {
		defer controller.wg.Done()
		for {
			select {
			case req := <-writeRequests:
				err := controller.service.PickUpPointWrite(req.PickUpPoint)
				req.ResponseChan <- err
			case <-ctx.Done():
				return
			}
		}
	}(controller.ctx)

	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(fmt.Errorf("failed to read input: %w", err))
			continue
		}
		input = strings.TrimSpace(input)

		parts := strings.SplitN(input, " ", 2)
		command := parts[0]

		switch command {
		case "read":
			responseChan := make(chan []model.PickUpPoint)
			readRequests <- ReadRequest{ResponseChan: responseChan}
			pickUpPoints := <-responseChan

			for _, p := range pickUpPoints {
				fmt.Printf("%+v\n", p)
			}
		case "write":
			if len(parts) < 2 {
				fmt.Println("write command requires arguments")
				continue
			}
			argParts := strings.Split(parts[1], ",")
			if len(argParts) != 3 {
				fmt.Println("expected 3 arguments for write command: Name, Address, Contact info")
				continue
			}
			responseChan := make(chan error)
			writeRequests <- WriteRequest{
				PickUpPoint: model.PickUpPoint{
					Name:    strings.TrimSpace(argParts[0]),
					Address: strings.TrimSpace(argParts[1]),
					Contact: strings.TrimSpace(argParts[2]),
				},
				ResponseChan: responseChan,
			}
			if err := <-responseChan; err != nil {
				fmt.Println(fmt.Errorf("failed to get response: %w", err))
			} else {
				fmt.Println("Pick-up point has been written successfully!")
			}
		case "exit":
			controller.wg.Wait()
			return
		default:
			fmt.Println("unknown command. List of commands: \"read, write\"")
		}
	}
}

func (controller *OrderController) HelpCommand() {
	fmt.Print("COMMANDS" +
		"\n  cr-take\n" +
		"\tAccepts and writes an order from the courier into a file\n" +
		"\tRequired flags: -cid, -oid, -at\n\n" +
		"\tCannot accept order twice or with a negative available time\n" +
		"\tExample of using: `cr-take -cid=1 -oid=1 -at=48h`\n" +
		"\n  cr-return\n" +
		"\tReturns an order to the courier (deletes the order from file)\n" +
		"\tRequired flag: -oid\n\n" +
		"\tCan return orders only with an expiration date in the past and those\n" +
		"\tthat have not been given\n" +
		"\tExample of using: `cr-return -oid=1`\n" +
		"\n  cl-give\n" +
		"\tGives one or more orders to the client\n" +
		"\tRequired flags: -cid, -oids\n\n" +
		"\tOnly orders accepted from the courier (present in the file) with expiration\n" +
		"\tdates earlier than the current date can be given. All order IDs must \n" +
		"\tbelong to the same client.\n" +
		"\tExample of using: `cl-give -cid=1 -oids=1,3,7`\n" +
		"\n  cl-orders\n" +
		"\tGet orders\n" +
		"\tRequired flag: -cid\n" +
		"\tOptional flags: -n, -ouo\n\n" +
		"\tExample of using: `cl-orders -cid=1 [-n=5] [-ouo]`\n" +
		"\n  cl-refund\n" +
		"\tAccepts a return from the client (move an order from available orders to refunded orders)\n" +
		"\tRequired flags: -cid, -oid\n\n" +
		"\tOrders can be returned within 2 days from the given date.\n" +
		"\tExample of using: `cl-refund -cid=1 -oid=1\n" +
		"\n  refund-list\n" +
		"\tProvides the list of refunds in a paginated manner.\n" +
		"\tOptional flag: -p\n\n" +
		"\tExample of using: `refund-list [-p=3]`\n" +
		"\n  interactive\n" +
		"\tLaunches an interactive mode that has two commands: write and read\n\n" +
		"\tExample of using: `write Pick-up point #1, Tomorrow Avenue, +78005553535`\n\n" +
		"FLAGS" +
		"\n  -cid int\n" +
		"\tClient ID\n" +
		"\n  -at duration\n" +
		"\tAvailable time for picking up the order (e.g. -at=48h)\n" +
		"\n  -n int\n" +
		"\tHow many last orders the user will receive (default -1)\n" +
		"\n  -oid int\n" +
		"\tOrder ID\n" +
		"\n  -oids string\n" +
		"\tComma-separated slice of ordersID (e.g. -oids=1,3,7)\n" +
		"\n  -ouo\n" +
		"\tGet only user orders\n" +
		"\n  -p int\n" +
		"\tPage number starting with 1 (default 1)")
}
