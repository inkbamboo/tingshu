package scripts

import (
	"github.com/bmbstack/ripple"
	"github.com/urfave/cli/v2"
)

func Init(c *cli.Context) {
	ripple.InitConfigWithPath(c.String("env"), c.String("conf"))
	ripple.GetConfig().Set("env", c.String("env"))
}
