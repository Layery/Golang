package main

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"
)

type myError struct {
	Msg, Line, FilePath string
}

func (e *myError) Error() string {
	return fmt.Sprintf("%s, %s, $s \n", e.Msg, e.FilePath, e.Line)
}


func handleErrorByType() {
	defer func() {
		if err := recover(); err != nil {
			if err, ok := err.(myError); ok {
				fmt.Printf("程序出错了: %s, 第%s行, 文件地址:%s \n", err.Msg, err.Line, err.FilePath)
			}
		}
	}()




	_, file, line, _ := runtime.Caller(0)

	panic(myError{"参数校验不通过", strconv.Itoa(line), file})
}


func handleErrorByWrap() {
	fd, err := os.Open("./log.log")


	if err != nil {

		err2 := fmt.Errorf("这里是包裹的第2层错误: [%w]", err)
		err3 := fmt.Errorf("这里是包裹的第3层错误: [%w]", err2)
		err33 := fmt.Errorf("这里是包裹的第33层错误: [%w]", err3)

		fmt.Printf("当前最新的错误: %s \n\n", err33)



		err4 := errors.Unwrap(err33)

		fmt.Println(err4)

		err5 := errors.Unwrap(err4)

		DD(err5)
	}



	_ = fd



}


// 使用自定义的类型承载错误信息,
// 优点: 能记录更多的错误上下文信息
// 缺点: 错误类型需要被公开, 以能够被别的地方调用, 可能会泄露源码信息
func main () {


	//go handleErrorByType()

	go handleErrorByWrap()


	time.Sleep(time.Second * 5)
	select {}
}