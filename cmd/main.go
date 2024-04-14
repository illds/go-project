package main

import (
	"GOHW-1/internal/configuration"
	"GOHW-1/internal/controller"
	"GOHW-1/internal/db"
	"GOHW-1/internal/infrastucture/kafka"
	"GOHW-1/internal/service"
	"GOHW-1/internal/storage"
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup

	go func() {
		<-signals
		log.Println("\nReceived shutdown signal, exiting...")

		cancel()
		wg.Wait()
		os.Exit(0)
	}()

	// Storage initialization
	strg, err := storage.New()
	if err != nil {
		log.Fatalf("cannot connect to storage: %v", err)
	}

	// Service initialization
	svc := service.New(&strg)

	// Database initialization
	dbCredentials := configuration.NewDBCredentials()
	dbCredentials.SetEnv()
	database, err := db.NewDb(ctx, dbCredentials)
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}
	defer database.GetPool(ctx).Close()

	// Kafka producer initialization
	topicName := configuration.GetTopicName()
	brokers := configuration.GetBrokers()
	kafkaProducer, err := kafka.NewProducer(*brokers)
	if err != nil {
		log.Fatalf("cannot connect to kafka: %v", err)
	}
	defer kafkaProducer.Close()

	// Kafka consumer (group) initialization in a go-routine
	go controller.ConsumerGroup(*brokers, *topicName)

	// Order controller initialization
	orderController := controller.NewOrderController(&svc, &wg, ctx)

	// Pick-up point controller initialization
	sender := controller.NewKafkaSender(kafkaProducer, *topicName)
	pickUpPointController := controller.NewPickUpPointController(database, sender)

	config := configuration.DefaultConfig()

	if len(os.Args) < 2 {
		log.Fatal("you need to write a command")
	}

	switch os.Args[1] {
	case "http":
		pickUpPointController.StartHTTPServer()
	case "cr-take":
		crTake := config.CrTake.FlagSet
		if err := crTake.Parse(os.Args[2:]); err != nil {
			log.Fatalf("failed to parse command line flags for cr-take: %v", err)
		}
		if err := orderController.CourierTakeCommand(*config.CrTake.OrderID, *config.CrTake.ClientID, *config.CrTake.AvailableTime,
			*config.CrTake.Weight, *config.CrTake.Price, *config.CrTake.Packaging); err != nil {
			log.Fatal(err)
		}
	case "cr-return":
		crReturn := config.CrReturn.FlagSet
		if err := crReturn.Parse(os.Args[2:]); err != nil {
			log.Fatalf("failed to parse command line flags for cr-return: %v", err)
		}
		if err := orderController.CourierReturnCommand(*config.CrReturn.OrderID); err != nil {
			log.Fatal(err)
		}
	case "cl-give":
		clGive := config.ClGive.FlagSet
		if err := clGive.Parse(os.Args[2:]); err != nil {
			log.Fatalf("failed to parse command line flags for cl-give: %v", err)
		}
		if err := orderController.ClientGiveCommand(*config.ClGive.ClientID, *config.ClGive.OrdersID); err != nil {
			log.Fatal(err)
		}
	case "cl-orders":
		clOrders := config.ClOrders.FlagSet
		if err := clOrders.Parse(os.Args[2:]); err != nil {
			log.Fatalf("failed to parse command line flags for cl-orders: %v", err)
		}
		if err := orderController.ClientOrdersCommand(*config.ClOrders.ClientID, *config.ClOrders.N, *config.ClOrders.OnlyUserOrders); err != nil {
			log.Fatal(err)
		}
	case "cl-refund":
		clRefund := config.ClRefund.FlagSet
		if err := clRefund.Parse(os.Args[2:]); err != nil {
			log.Fatalf("failed to parse command line flags for cl-refund: %v", err)
		}
		if err := orderController.ClientRefundCommand(*config.ClRefund.ClientID, *config.ClRefund.OrderID); err != nil {
			log.Fatal(err)
		}
	case "refund-list":
		refundList := config.RefundList.FlagSet
		if err := refundList.Parse(os.Args[2:]); err != nil {
			log.Fatalf("failed to parse command line flags for refund-list: %v", err)
		}
		if err := orderController.RefundListCommand(*config.RefundList.PageNumber); err != nil {
			log.Fatal(err)
		}
	case "interactive":
		orderController.InteractiveCommand()
	case "help":
		orderController.HelpCommand()
	default:
		log.Fatal("unknown command. Try to use `help` command")
	}
}
