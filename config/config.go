package config

import (
	"github.com.wureny/chain-monitor/flags"
	"github.com/urfave/cli/v2"
)

type Config struct {
	RpcUrl            string
	MaxLenOfOperation int64
	MaxLenOfWatch     int64
}

func LoadConfig(ctx *cli.Context) Config {
	cfg := NewConfig(ctx)
	//TODO 一些检查
	return cfg
}

func NewConfig(ctx *cli.Context) Config {
	return Config{
		MaxLenOfOperation: ctx.Int64(flags.MaxLenOfOperationFlag.Name),
		MaxLenOfWatch:     ctx.Int64(flags.MaxLenOfWatchFlag.Name),
		RpcUrl:            ctx.String(flags.RpcUrlFlag.Name),
	}
}
