package controller

import (
	"GOHW-1/internal/infrastucture/kafka"
	"context"
	"github.com/IBM/sarama"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func ConsumerGroup(brokers []string, topicName string) {
	keepRunning := true
	log.Println("Starting a new Sarama consumer...")

	config := sarama.NewConfig()
	config.Version = sarama.MaxVersion

	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Group.ResetInvalidOffsets = true

	config.Consumer.Group.Heartbeat.Interval = 3 * time.Second
	config.Consumer.Group.Session.Timeout = time.Minute
	config.Consumer.Group.Rebalance.Timeout = time.Minute

	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRoundRobin}

	consumer := kafka.NewConsumerGroup()
	group := "route"

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(brokers, group, config)
	if err != nil {
		log.Panicf("error creating consumer group client: %v", err)
	}
	defer client.Close()

	consumptionIsPaused := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := client.Consume(ctx, []string{topicName}, &consumer); err != nil {
				log.Panicf("error from consumer: %v", err)
			}
			// Check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
		}
	}()

	<-consumer.Ready()
	log.Println("Sarama consumer up and running!")

	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			log.Println("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			log.Println("terminating: via signal")
			keepRunning = false
		case <-sigusr1:
			toggleConsumptionFlow(client, &consumptionIsPaused)
		}
	}

	cancel()
	wg.Wait()
}

func toggleConsumptionFlow(client sarama.ConsumerGroup, isPaused *bool) {
	if *isPaused {
		client.ResumeAll()
		log.Println("Resuming consumption")
	} else {
		client.PauseAll()
		log.Println("Pausing consumption")
	}

	*isPaused = !*isPaused
}
