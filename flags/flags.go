package flags

import "github.com/urfave/cli/v2"

const prefix = "MONITOR"

func prefixEnvVar(vari string) []string {
	return []string{prefix + "_" + vari}
}

var (
	RpcUrlFlag = &cli.StringFlag{
		Name:        "rpcUrl",
		Category:    "",
		DefaultText: "",
		FilePath:    "",
		Usage:       "",
		Required:    true,
		Hidden:      false,
		HasBeenSet:  false,
		Value:       "",
		Destination: nil,
		Aliases:     nil,
		EnvVars:     prefixEnvVar("RPC_URL"),
		TakesFile:   false,
		Action:      nil,
	}

	MaxLenOfWatch = &cli.IntFlag{
		Name:        "maxLenOfWatch",
		Category:    "",
		DefaultText: "",
		FilePath:    "",
		Usage:       "",
		Required:    true,
		Hidden:      false,
		HasBeenSet:  false,
		Value:       0,
		Destination: nil,
		Aliases:     nil,
		EnvVars:     prefixEnvVar("MAX_LEN_OF_WATCH"),
		Base:        0,
		Action:      nil,
	}

	MaxLenOfOperation = &cli.IntFlag{
		Name:        "maxLenOfOperation",
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

func init() {

}

var Flags []cli.Flag
