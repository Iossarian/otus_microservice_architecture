package config

import (
	"net"
	"net/url"
	"strings"
)

type Postgres struct {
	User     string `envconfig:"POSTGRES_USER"`
	Password string `envconfig:"POSTGRES_PASSWORD"`
	Host     string `envconfig:"POSTGRES_HOST"`
	Port     string `envconfig:"POSTGRES_PORT"`
	DB       string `envconfig:"POSTGRES_DB"`
}

func (c *Config) PostgresDSN() string {
	dsn := strings.Builder{}
	dsn.WriteString("postgres://")

	if c.Postgres.User != "" {
		dsn.WriteString(c.Postgres.User)

		if c.Postgres.Password != "" {
			dsn.WriteString(":" + encodeIfNeeded(c.Postgres.Password))
		}

		dsn.WriteString("@")
	}

	hp := net.JoinHostPort(c.Postgres.Host, c.Postgres.Port)

	dsn.WriteString(hp + "/")

	if c.Postgres.DB != "" {
		dsn.WriteString(c.Postgres.DB)
	}

	dsn.WriteString("?sslmode=disable")

	return dsn.String()
}

func encodeIfNeeded(s string) string {
	encoded := url.QueryEscape(s)
	if s != encoded {
		return encoded
	}

	return s
}
