package main

import (
	"errors"
	"fmt"
	"os"
)

// 声明一个日志写入器的接口
type LogWriter interface {
	Log(data interface{}) error // 该接口具有一个Log方法
}

// 定义一个日志器
type Logger struct {
	writerList []LogWriter
}

// 注册一个日志写入器
func (l *Logger) RegisterWriter(writer LogWriter) {
	l.writerList = append(l.writerList, writer)
}

// 将一个data类型的数据写入日志
func (l *Logger) Log(data interface{}) {
	fmt.Printf("now data is %v\n", data)
	for _, writer := range l.writerList {
		fmt.Printf("now log driver is %#v \n", writer)
		err := writer.Log(data)
		if err != nil {
			panic(err)
		}

	}
}

// 构造函数创建一个日志器
func NewLogger() *Logger {
	return &Logger{}
}

/////////////////////////////////////////////////////////

/**
接下来, 创建一个文件写入器, 和 console写入器
*/

type fileWriter struct {
	file *os.File // 这里貌似是只取了os库的File组件, 作为fileWriter的file字段
}

// 文件写入器, 实现LogWriter的Log方法
func (f *fileWriter) Log(data interface{}) error {
	if f.file == nil { // 日志文件没有准备好
		return errors.New("日志文件没有准备好")
	}
	fmt.Printf("调用filewriter的Log方法, data值: %v \n", data)
	// 将数据序列化成字符串
	logStr := fmt.Sprintf("%v\n", data)

	// 将数据以字节数组放入文件中
	_, err := f.file.Write([]byte(logStr))
	return err
}

// 文件写入器需要设置一个log文件
func (f *fileWriter) SetFile(filename string) (*fileWriter, error) {
	// 如果文件已经打开了, 关了它
	if f.file != nil {
		fmt.Printf("这个文件已经被打开了")
		f.file.Close()
	}
	// 打开一个文件句柄
	createResult, err := os.Create(filename)
	fmt.Printf("===========>  createResult %#v \n", createResult)
	fmt.Printf("===========>  createResult type %T \n", createResult)
	if err != nil {
		fmt.Printf("文件打开失败%#v", err)
	}
	fmt.Printf("当前设置的文件是%#v, f.file is %#v \n", filename, f.file)
	f.file = createResult

	return f, err
}

// 文件写入器的构造函数
func NewFileWriter() *fileWriter {
	return &fileWriter{}
}

/////////////////////////////////////////////////////////

func main() {

	/**
	在程序中使用日志器一般会先通过代码创建日志器（Logger），
	为日志器添加输出设备（file Writer、console Writer等）。
	这些设备中有一部分需要一些参数设定，
	如文件日志写入器需要提供文件名（file Writer的Set File()方法）。
	*/

	// 1. 先创建logger
	logger := NewLogger()

	// 2. 创建文件写入器
	fileWriter := NewFileWriter()
	writer, _ := fileWriter.SetFile("./src/log/debug.log")

	fmt.Printf("设置好了filename, writer: %#v\n writer 的类型是%T \n", writer, writer)

	// 3. 给日志器logger注册一个文件写入器
	logger.RegisterWriter(writer)

	// 4. 输出一个日志
	logger.Log("这是我输出的一个日志")

}
