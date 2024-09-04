package main

import (
	"context"
	"github.com/ethereum/go-ethereum/log"
	"os"
)

func main() {
	log.SetDefault(log.NewLogger(log.NewTerminalHandlerWithLevel(os.Stderr, log.LevelInfo, true)))
	app := NewCli()
	if err := app.RunContext(context.Background(), os.Args); err != nil {
		log.Error("Application failed")
		os.Exit(1)
	}
}
