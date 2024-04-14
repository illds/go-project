//go:build integration

package tests

import (
	"GOHW-1/internal/configuration"
	"GOHW-1/internal/controller"
	"GOHW-1/internal/db"
	"GOHW-1/internal/infrastucture/kafka"
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

type TDB struct {
	DB db.Database
}

type LoggingMessage struct {
	Method string
	URI    string
	Time   time.Time
}

var (
	tdb                   *TDB
	pickUpPointController *controller.PickUpPointController
)

func init() {
	// DB initialization
	dbCredentials := configuration.NewDBCredentials()
	dbCredentials.SetEnv()
	dbCredentials.SetDBname(os.Getenv("POSTGRES_DB_TEST"))
	tdb = newFromEnv(dbCredentials)

	ctx := context.Background()
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
	//defer kafkaProducer.Close()

	sender := controller.NewKafkaSender(kafkaProducer, *topicName)
	pickUpPointController = controller.NewPickUpPointController(database, sender)

}

func newFromEnv(dbCredentials *configuration.DBCredentials) *TDB {
	newDb, err := db.NewDb(context.Background(), dbCredentials)
	if err != nil {
		panic(err)
	}
	return &TDB{DB: *newDb}
}

func getKafkaMessage(t *testing.T) LoggingMessage {
	brokers := configuration.GetBrokers()
	consumer, err := sarama.NewConsumer(*brokers, nil)
	if err != nil {
		panic("Failed to start Sarama consumer: " + err.Error())
	}

	topicName := configuration.GetTopicName()
	partitionConsumer, err := consumer.ConsumePartition(*topicName, 0, sarama.OffsetNewest)
	if err != nil {
		panic("Failed to start partition consumer: " + err.Error())
	}

	msg := <-partitionConsumer.Messages()

	lm := LoggingMessage{}
	if err := json.Unmarshal(msg.Value, &lm); err != nil {
		fmt.Println("Consumer group error", err)
	}
	return lm
}

func dummyHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Handler executed"))
	})
}
