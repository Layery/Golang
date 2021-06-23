package main

import "fmt"

// 结构体

// 结构体类似于PHP的Object

type PersonModel struct {
	name, addr string
	age        int8
}

func main() {
	var zhangsan PersonModel

	fmt.Printf("%T", zhangsan)
}
