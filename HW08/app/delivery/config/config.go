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
	HTTP     HTTP
	Postgres Postgres
}

type HTTP struct {
	Host    string   `envconfig:"HTTP_HOST"`
	Port    int32    `envconfig:"HTTP_PORT"`
	Schemes []string `envconfig:"HTTP_SCHEMES" default:"http"`
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
