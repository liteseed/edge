package commands

import (
	"github.com/urfave/cli/v2"
)

var Commands = []*cli.Command{
	Balance,
	Migrate,
	Stake,
	Start,
	Unstake,
}
