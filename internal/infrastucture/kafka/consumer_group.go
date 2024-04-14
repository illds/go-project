package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"time"
)

type LoggingMessage struct {
	Method string
	URI    string
	Time   time.Time
}

type ConsumerGroup struct {
	ready chan bool
}

func NewConsumerGroup() ConsumerGroup {
	return ConsumerGroup{
		ready: make(chan bool),
	}
}

func (consumer *ConsumerGroup) Ready() <-chan bool {
	return consumer.ready
}

// Setup starts a new session, before ConsumeClaim
func (consumer *ConsumerGroup) Setup(_ sarama.ConsumerGroupSession) error {
	close(consumer.ready)

	return nil
}

// Cleanup finishes the session, after all ConsumeClaim finish
func (consumer *ConsumerGroup) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim reads while the session is running
func (consumer *ConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			lm := LoggingMessage{}
			if err := json.Unmarshal(message.Value, &lm); err != nil {
				fmt.Println("Consumer group error", err)
			}

			log.Printf("Message claimed: \"Method: %s, URI: %s, Time: %s\"", lm.Method, lm.URI, lm.Time)

			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}
