package main

import (
	"context"
	"os"

	"delivery/cmd"
	"delivery/config"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		panic(err)
	}

	exitCode := 0

	ctx := context.Background()

	err = cmd.Run(ctx, conf)
	if err != nil {
		exitCode = 1
	}

	os.Exit(exitCode)
}
