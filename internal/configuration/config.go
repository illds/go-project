package configuration

import (
	"flag"
	"os"
	"time"
)

type AppConfig struct {
	CrTake     CrTakeConfig
	CrReturn   CrReturnConfig
	ClGive     ClGiveConfig
	ClOrders   ClOrdersConfig
	ClRefund   ClRefundConfig
	RefundList RefundListConfig
}

type CrTakeConfig struct {
	FlagSet       flag.FlagSet
	OrderID       *int
	ClientID      *int
	AvailableTime *time.Duration
	Weight        *float64
	Price         *float64
	Packaging     *string
}

type CrReturnConfig struct {
	FlagSet flag.FlagSet
	OrderID *int
}

type ClGiveConfig struct {
	FlagSet  flag.FlagSet
	ClientID *int
	OrdersID *string
}

type ClOrdersConfig struct {
	FlagSet        flag.FlagSet
	ClientID       *int
	N              *int
	OnlyUserOrders *bool
}

type ClRefundConfig struct {
	FlagSet  flag.FlagSet
	OrderID  *int
	ClientID *int
}

type RefundListConfig struct {
	FlagSet    flag.FlagSet
	PageNumber *int
}

type DBCredentials struct {
	Host     string
	Port     string
	User     string
	Password string
	DBname   string
}

func NewDBCredentials() *DBCredentials {
	return &DBCredentials{}
}

func (dbCredentials *DBCredentials) SetEnv() {
	dbCredentials.Host = os.Getenv("POSTGRES_HOST")
	dbCredentials.Port = os.Getenv("POSTGRES_PORT")
	dbCredentials.User = os.Getenv("POSTGRES_USER")
	dbCredentials.Password = os.Getenv("POSTGRES_PASSWORD")
	dbCredentials.DBname = os.Getenv("POSTGRES_DB")
}

func (dbCredentials *DBCredentials) SetCredits(host string, port string, user string, password string, dbname string) {
	dbCredentials.Host = host
	dbCredentials.Port = port
	dbCredentials.User = user
	dbCredentials.Password = password
	dbCredentials.DBname = dbname
}

func (dbCredentials *DBCredentials) SetDBname(dbname string) {
	dbCredentials.DBname = dbname
}

func GetBrokers() *[]string {
	brokers := []string{
		"127.0.0.1:9091",
		"127.0.0.1:9092",
		"127.0.0.1:9093",
	}
	return &brokers
}

func GetTopicName() *string {
	topicName := "logs"
	return &topicName
}

func DefaultConfig() AppConfig {
	crTake := flag.NewFlagSet("cr-take", flag.ExitOnError)
	crTakeOrderID := crTake.Int("oid", 0, "Order ID")
	crTakeClientID := crTake.Int("cid", 0, "Client ID")
	crTakeAvailableTime := crTake.Duration("at", 0, "Available time for picking up the order (e.g. -at=48h)")
	crTakeWeight := crTake.Float64("w", 0, "Weight")
	crTakePrice := crTake.Float64("pr", 0, "Price")
	crTakePackaging := crTake.String("pg", "", "Packaging")

	crReturn := flag.NewFlagSet("cr-return", flag.ExitOnError)
	crReturnOrderID := crReturn.Int("oid", 0, "Order ID")

	clGive := flag.NewFlagSet("cl-give", flag.ExitOnError)
	clGiveClientID := clGive.Int("cid", 0, "Client ID")
	clGiveOrdersID := clGive.String("oids", "", "Comma-separated slice of ordersID (e.g. -oids=1,3,7)")

	clOrders := flag.NewFlagSet("cl-orders", flag.ExitOnError)
	clOrdersClientID := clOrders.Int("cid", 0, "Client ID")
	clOrdersN := clOrders.Int("n", -1, "How many last orders the user will receive")
	clOrdersOnlyUserOrders := clOrders.Bool("ouo", false, "Get only user orders")

	clRefund := flag.NewFlagSet("cl-refund", flag.ExitOnError)
	clRefundOrderID := clRefund.Int("oid", 0, "Order ID")
	clRefundClientID := clRefund.Int("cid", 0, "Client ID")

	refundList := flag.NewFlagSet("refund-list", flag.ExitOnError)
	refundListPageNumber := refundList.Int("p", 1, "Page number starting with 1")

	return AppConfig{
		CrTake: CrTakeConfig{FlagSet: *crTake, OrderID: crTakeOrderID, ClientID: crTakeClientID, AvailableTime: crTakeAvailableTime,
			Weight: crTakeWeight, Price: crTakePrice, Packaging: crTakePackaging},
		CrReturn:   CrReturnConfig{FlagSet: *crReturn, OrderID: crReturnOrderID},
		ClGive:     ClGiveConfig{FlagSet: *clGive, ClientID: clGiveClientID, OrdersID: clGiveOrdersID},
		ClOrders:   ClOrdersConfig{FlagSet: *clOrders, ClientID: clOrdersClientID, N: clOrdersN, OnlyUserOrders: clOrdersOnlyUserOrders},
		ClRefund:   ClRefundConfig{FlagSet: *clRefund, OrderID: clRefundOrderID, ClientID: clRefundClientID},
		RefundList: RefundListConfig{FlagSet: *crTake, PageNumber: refundListPageNumber},
	}
}
