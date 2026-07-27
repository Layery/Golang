package types

import "github.com/urfave/cli/v2"

// 统一处理器签名，所有抓取函数必须遵守这个入参
type CrawlerHandler func(ctx *cli.Context) error
