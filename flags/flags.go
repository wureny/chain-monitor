package flags

import "github.com/urfave/cli/v2"

var (
	RpcUrlFlag = &cli.StringFlag{
		Name:        "",
		Category:    "",
		DefaultText: "",
		FilePath:    "",
		Usage:       "",
		Required:    false,
		Hidden:      false,
		HasBeenSet:  false,
		Value:       "",
		Destination: nil,
		Aliases:     nil,
		EnvVars:     nil,
		TakesFile:   false,
		Action:      nil,
	}

	MaxLenOfWatch = &cli.IntFlag{
		Name:        "",
		Category:    "",
		DefaultText: "",
		FilePath:    "",
		Usage:       "",
		Required:    false,
		Hidden:      false,
		HasBeenSet:  false,
		Value:       0,
		Destination: nil,
		Aliases:     nil,
		EnvVars:     nil,
		Base:        0,
		Action:      nil,
	}
)
