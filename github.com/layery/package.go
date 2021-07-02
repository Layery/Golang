package main

import (
	"helper/logger"
)

/**
	存放go文件的目录, 叫做一个包, 包里的标识符如果以大写字母开头的,
	就是共享的

	当标识符是首字母大写的时候, 就可以被别的包所访问
 */
func main()  {
	fileWriter := logger.NewFileWriter()
	_, err := fileWriter.SetFile("./debug.txt")
	logger.Out(err)
	logger.NewLogger().Log("debug")
}