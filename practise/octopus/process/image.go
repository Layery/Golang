package process

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func RunImage(ctx *cli.Context) error {

	// 搜索的公众号名称
	name := ctx.String("name")
	page := ctx.Int("page")

	proxy := ctx.String("proxy")

	fmt.Println(name, page, proxy)
	return nil
}
