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
func TotalSum(valList ...int) int {
	total := 0
	for _, val := range valList {
		total += val
	}
	return total
}

/**
函数可以作为值, 作为类型, 来传递,
go中, 可以使用type关键字来定义函数类型的一个类型, 它的类型就是所拥有的的
相同的函数
**/

// 声明一个函数类型
type FuncTypeInt func(int) bool

// 是否是奇数
func isOdd(integer int) bool {
	if integer%2 == 0 {
		return false
	}
	return true
}

// 是否是偶数
func isEven(integer int) bool {
	if integer%2 == 0 {
		return true
	}
	return false
}

// 将声明的函数类型, 作为传参
func filter(slice []int, f FuncTypeInt) []int {
	var result []int
	for _, v := range slice {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}

func main() {
	// go的匿名函数

	sum, _ := add(3, 4)

	fmt.Printf("%v \n", sum)

	// defer 语句会在程序都执行完后开始执行
	for i := 0; i < 5; i++ {
		defer fmt.Printf("使用了defer之后, 打印顺序变成了倒序了: %d \n", i)
		fmt.Printf("不使用defer时, 打印顺序就是正序的: %d \n", i)
	}
	// 由于使用了...的语法, 该函数支持多个参数
	rs := TotalSum(3, 4, 5)
	fmt.Println(rs)

	var array = [11]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	slice := array[:5]
	fmt.Printf("type is %T, <====> value is %v \n", slice, slice)

	// 将函数, 作为参数传到filter函数中, 则 传入的函数参数将会在内部执行, 类似于PHP的array_map函数的用法
	odd := filter(slice, isOdd)

	fmt.Printf("切片内的值是: %v, 其中的奇数: %#v\n", slice, odd)
	fmt.Printf("切片内的值是: %v, 其中的偶数: %#v\n", slice, filter(slice, isEven))

}
