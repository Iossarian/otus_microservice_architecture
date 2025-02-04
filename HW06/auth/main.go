package main

import (
	"context"
	"log"
	"net"
	"net/url"
	"os"
	"strings"

	"github.com/Iossarian/otus_microservice_architecture/gateway/cmd"
	_ "github.com/lib/pq"
)

func main() {
	dsn := dsn(
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	err := cmd.Run(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
	}
}

func dsn(username, password, host, port, database string) string {
	dsn := strings.Builder{}
	dsn.WriteString("postgres://")

	if username != "" {
		dsn.WriteString(username)

		if password != "" {
			dsn.WriteString(":" + encodeIfNeeded(password))
		}

		dsn.WriteString("@")
	}

	hp := net.JoinHostPort(host, port)

	dsn.WriteString(hp + "/")

	if database != "" {
		dsn.WriteString(database)
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
