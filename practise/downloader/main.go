package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"runtime"
)

func main() {
	// 默认的并发数
	currCyn := runtime.NumCPU()

	app := cli.App{
		Name: "downloader",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "url",
				Aliases:  []string{"u"},
				Usage:    "`URL` 即将下载的文件地址",
				Value:    "https://apache.claz.org/zookeeper/zookeeper-3.7.0/apache-zookeeper-3.7.0-bin.tar.gz",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Output `filename`",
			},
			&cli.IntFlag{
				Name:    "concurrency",
				Aliases: []string{"n"},
				Value:   currCyn,
				Usage:   "最大并发线程数, 默认为4",
			},
		},
		Action: func(context *cli.Context) error {
			url := context.String("url")
			currCyn := context.Int("concurrency")
			filename := context.String("filename")
			return NewDownloader(currCyn).Download(url, filename)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
