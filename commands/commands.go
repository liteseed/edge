package commands

import (
	"github.com/urfave/cli/v2"
)

var Bungo = []*cli.Command{
	Benchmark,
	Start,
}
