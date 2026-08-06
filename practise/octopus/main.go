package main

import (
	"log"
	"octopus/process"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	// 创建
	app := &cli.App{
		Name:            "小章鱼",
		Usage:           "不可以👀😍",
		HideHelpCommand: true,
		UsageText:       "octopus.exe [--proxy HOST:PORT] web -n NAME [-c CATEGORY] [-p PAGE]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "proxy",
				Aliases:  []string{"p"},
				Value:    "",
				Usage:    "设置全局代理",
				Required: false,
			},
		},
		Commands: []*cli.Command{
			{
				Name: "web",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Aliases:  []string{"n"},
						Name:     "name",
						Required: true,
						Usage:    "待抓的站点名称",
					},
					&cli.StringFlag{
						Aliases: []string{"p"},
						Name:    "page",
						Value:   "1",
						Usage:   "paginate, default 1st page.",
					},
					&cli.StringFlag{
						Aliases: []string{"c"},
						Name:    "category",
						Value:   "22",
						Usage:   "指定待抓取的分区板块", // 7: 技术讨论 8: 新时代的我们 16: 达盖尔的旗帜 22: 在线小电影
					},
					&cli.StringFlag{
						Aliases: []string{"o"},
						Name:    "output",
						Value:   "./output/dist/web.html",
					},
				},
				Usage: "抓网页",
				Action: func(ctx *cli.Context) error {
					process.RunWeb(ctx)
					return nil
				},
			},
			{
				Name: "image",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Aliases:  []string{"n"},
						Name:     "name",
						Required: true,
						Usage:    "待抓的站点名称",
					},
					&cli.StringFlag{
						Aliases: []string{"p"},
						Name:    "page",
						Value:   "1",
						Usage:   "paginate, default 1st page.",
					},
				},
				Usage: "抓图片",
				Action: func(ctx *cli.Context) error {
					return process.RunImage(ctx)
				},
			},
			{
				Name: "debug",
				Action: func(ctx *cli.Context) error {
					return process.RunDebug(ctx)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
