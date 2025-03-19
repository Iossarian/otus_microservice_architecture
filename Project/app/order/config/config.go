package config

import (
	"net"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	HTTP         HTTP
	Orchestrator Orchestrator
	Warehouse    Warehouse
	Billing      Billing
	Delivery     Delivery
	Notification Notification
	JWT          JWT
	Kafka        Kafka
	Postgres     Postgres
}

type HTTP struct {
	Host    string   `envconfig:"HTTP_HOST"`
	Port    int32    `envconfig:"HTTP_PORT"`
	Schemes []string `envconfig:"HTTP_SCHEMES" default:"http"`
}

type Orchestrator struct {
	BaseURL string `envconfig:"ORCHESTRATOR_BASE_URL"`
}

type Warehouse struct {
	BaseURL string `envconfig:"WAREHOUSE_BASE_URL"`
}

type Billing struct {
	BaseURL string `envconfig:"BILLING_BASE_URL"`
}

type Delivery struct {
	BaseURL string `envconfig:"DELIVERY_BASE_URL"`
}

type Notification struct {
	BaseURL string `envconfig:"NOTIFICATION_BASE_URL"`
}

type JWT struct {
	Secret string `envconfig:"JWT_SECRET"`
}

type Kafka struct {
	Broker            string `envconfig:"KAFKA_BROKER" default:"kafka:29092"`
	OrderCreatedTopic string `envconfig:"KAFKA_ORDER_CREATED_TOPIC" default:"order.created"`
}

func Load() (Config, error) {
	cnf := Config{}

	if err := godotenv.Load(".env"); err != nil && !errors.Is(err, os.ErrNotExist) {
		return cnf, errors.Wrap(err, "read .env file")
	}

	if err := envconfig.Process("", &cnf); err != nil {
		return cnf, errors.Wrap(err, "read environment")
	}

	return cnf, nil
}

func (c *Config) HTTPAddr() string {
	return net.JoinHostPort(c.HTTP.Host, strconv.Itoa(int(c.HTTP.Port)))
}
