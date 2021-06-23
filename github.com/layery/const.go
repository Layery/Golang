package main

import (
	"fmt"
)

// 常量的声明, 和PHP类似
const PI = 3.1415926

// 常量也支持批量声明
// iota 是内置的常量内部的计数器, 可以理解为常量的行索引,
// 同时匿名变量也支持在常量里使用, 如果某一行的键没有给定值, 则默认沿用上一行的值
const (
	NAME_1 = iota
	_
	NAME_3
)

/**
	iota 在遇到const关键字时, 会默认重置为0, 每新增一行会累加1, 可以理解为常量的行索引
**/
func studyIota() {
	const (
		a, b = iota + 1, iota + 2
		c, d = iota + 1, iota + 2
		e, f
		g, h = iota + 1, iota + 2
	)

	fmt.Println(a, b, c, d, e, f, g, h)
}

/**
函数需要定义该函数的返回值的类型, 挨个定义
**/
func test() (int, int) {
	fmt.Println(NAME_1)
	fmt.Println(NAME_3)
	return 20, 30
}

func main() {
	// 字符串类型
	var name string = "llf"
	fmt.Println(name)

	// int类型
	var age int8 = 30
	fmt.Println(age)

	// 批量生明好多变量, 声明变量时, 不同类型的变量, go会给赋予该类型的默认值
	var (
		a string
		b int
		c bool
		d float32
		e rune
	)
	fmt.Println(a, b, c, d, e)

	// go 支持短变量声明语法, 该语法只能用在函数内来使用, 函数外,必须要使用var关键字声明变量
	address := "河北省"
	fmt.Println(address)

	// 匿名变量, 类似于占位符, 可以达到忽略某个值的功能
	var _, ttt = test()
	var ddd, _ = test()
	fmt.Println(ttt)
	fmt.Println(ddd)
	fmt.Println(PI)

	test()

	studyIota()
}
