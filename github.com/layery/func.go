package main

import (
	"fmt"
)

func add(a int, b int) (int, int) {
	return a + b, a - b
}

/**
 函数的参数支持使用...的语法, 来支持变参, 它的参数个数不是固定的,
 可以随意指定参数个数
 */
func TotalSum (valList ... int) int  {
	total := 0
	for _, val := range valList {
		total += val
	}
	return total
}



func main() {
	// go的匿名函数

	sum, _ := add(3, 4)

	fmt.Printf("%v \n", sum)

	for i := 0; i < 5; i++ {
		defer fmt.Printf("%d \n", i) // 使用了defer之后, 打印顺序变成了倒序了
		fmt.Printf("%d \n", i) // 不使用defer时, 打印顺序就是正序的
	}
	// 由于使用了...的语法, 该函数支持多个参数
	rs := TotalSum(3, 4, 5)
	fmt.Println(rs)
}
