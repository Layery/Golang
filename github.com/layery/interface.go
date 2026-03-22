package main

import (
	"fmt"
)

type People interface {
	Speak(string) string
}

type Stduent struct{}

func (stu *Stduent) Speak(think string) (talk string) {
	if think == "love" {
		talk = "You are a good boy"
	} else {
		talk = "hi"
	}
	return
}

func main() {
	/**
		go当中的多态, 由于peo定义的变量类型是一个接口People, student虽然实现了People接口, 但是本质上他们俩的
	类型还是不一样的, People是接口, 但是Student是结构体, 所以无法直接赋值, 只要把Student的引用赋值给peo就好了
	因为People接口的零值是nil, Student结构体的指针也是nil
	*/
	//var peo People = Stduent{}
	var peo People = &Stduent{}
	think := "love"
	fmt.Println(peo.Speak(think))
}
