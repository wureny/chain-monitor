package main

import (
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
	/*log.SetDefault(log.NewLogger(log.NewTerminalHandlerWithLevel(os.Stderr, log.LevelInfo, true)))
	monitor := synchronizer.NewMonitor(context.Background())
	err := monitor.Start(context.Background())
	if err != nil {
		os.Exit(1)
	}*/
}
