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
	JWT          JWT
	User         User
	Billing      Billing
	Order        Order
	Notification Notification
}

type HTTP struct {
	Host    string   `envconfig:"HTTP_HOST"`
	Port    int32    `envconfig:"HTTP_PORT"`
	Schemes []string `envconfig:"HTTP_SCHEMES" default:"http"`
}

type JWT struct {
	Secret string `envconfig:"JWT_SECRET"`
}

type User struct {
	BaseURL string `envconfig:"USER_BASE_URL"`
}

type Billing struct {
	BaseURL string `envconfig:"BILLING_BASE_URL"`
}

type Order struct {
	BaseURL string `envconfig:"ORDER_BASE_URL"`
}

type Notification struct {
	BaseURL string `envconfig:"NOTIFICATION_BASE_URL"`
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
