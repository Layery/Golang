package main

import "fmt"

var arr = [3]int{1, 2, 3}

func fixForFuncType1() {
	for i := 0; i < len(arr); i++ {
		temp := i
		msg := fmt.Sprintf("采用方式1来改造, 内存地址一样: %v, 拷贝一份i的值: %v", &i, temp)
		fmt.Println(msg)
	}
	fmt.Println()
}

func fixForFuncType2() {
	for i := 0; i < len(arr); i++ {
		msg := fmt.Sprintf("采用方式2来改造, i传入之前地址: %v, 值 : %v", &i, i)
		fmt.Println(msg)
		func(i int) {
			msg := fmt.Sprintf("采用方式2来改造, 将i传入匿名函数后地址变了(这是为啥, 因为Golang中函数的参数都是值拷贝, 虽然值一样但是其实内存地址是新拷贝出来的): %v, 值: %v", &i, i)
			fmt.Println(msg)
		}(i)
	}
	fmt.Println()
}

func main() {
	var arr = [3]int{1, 2, 3}

	defer func() {
		msg := `
	这是在学习 Go 程序设计 中遇到的一个比较重要的一个警告。这是个 Go 语言的词法作用域规则的陷阱。
1. for循环中, v的值始终是在同一个内存地址上重新赋值
2. 函数调用中, 参数始终是值拷贝, 除非显式使用值引用传参
`
		fmt.Println(msg)
	}()

	defer func() {
		fixForFuncType1() // 方式1改造

		fixForFuncType2() // 方式2改造
	}()

	defer func() {
		answer := `
	那么如何改变这个现象呢? 
1. 可以在循环体内, 对i进行重新赋值
2. 将i作为匿名函数的形参传进去 
`
		note := `
	为什么会这样? 因为for循环体中, i 一直是在同一个内存地址上发生的变化, 同时defer的特性
是等待函数执行完毕之后才开始执行的一个栈存储结构, 当defer内的函数开始执行的时候, 此时的i已经被
遍历完, 当i=2时, 满足i<len(arr), 所以i++后 i变为了3 所以输出的是 3 3 3
`
		fmt.Println(note + answer)
	}()

	for i := 0; i < len(arr); i++ {
		msg := fmt.Sprintf("对于数组, 我们是同一个内存地址: %v, 但是值却每次都不一样: %v", &arr[i], arr[i])
		fmt.Println(msg)

	}

	fmt.Println()
	var myslice = []int{1, 2, 3}

	for i := 0; i < len(myslice); i++ {
		msg := fmt.Sprintf("对于切片, 我们是同一个内存地址: %v, 但是值却每次都不一样: %v", &i, i)
		fmt.Println(msg)
	}

	fmt.Println()

	for i := 0; i < len(myslice); i++ {
		defer func() {
			msg := fmt.Sprintf("我们是同一个内存地址: %v, 但是值却每次都一样: %v", &i, i)
			fmt.Println(msg)
		}()
	}

}
