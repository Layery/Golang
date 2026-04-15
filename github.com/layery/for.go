package main

import (
	"fmt"
	"runtime"
)

var arr = [3]int{1, 2, 3}

func fixForFuncType1() {
	for i := 0; i < len(arr); i++ {
		temp := i
		msg := fmt.Sprintf("采用方式1来改造, i被赋值给temp, temp的地址: %v, temp的值: %v", &temp, temp)
		fmt.Println(msg)
	}
	fmt.Println()
}

func fixForFuncType2() {
	for i := 0; i < len(arr); i++ {
		msg := fmt.Sprintf("采用方式2来改造, i传入之前地址: %v, 值: %v, ", &i, i)
		func(i int) {
			tmp := fmt.Sprintf("i传入之后, i的地址: %v, 值: %v", &i, i)
			fmt.Println(msg + tmp)
		}(i)
	}
	fmt.Println()
}

func init() {
	fmt.Printf("当前golang版本: %v \n\n", runtime.Version())
}

func main() {

	// 声明一个数组
	var arr = [3]int{1, 2, 3}

	// 由数组声明一个切片
	slice := arr[:]

	/**********    经典for循环遍历数组  ****************/
	fmt.Printf("1. 先来看看for循环遍历数组的情况: \n")
	for i := 0; i < len(arr); i++ {
		msg := fmt.Sprintf("for循环遍历数组: %d(%p) => %d(%p) ", i, &i, arr[i], &arr[i])
		fmt.Println(msg)
	}

	/**********    for range 遍历数组  ****************/
	fmt.Println("\n\n2. 接下来看看for range 遍历数组的情况:")
	for ri, rv := range arr {
		msg := fmt.Sprintf("for range 遍历数组: %d(%p) => %d(%p)", ri, &ri, rv, &rv)
		fmt.Println(msg)
	}

	/**********    经典for循环遍历切片  ****************/
	fmt.Println("\n\n3. 接下来看看for循环遍历切片的情况:")
	for i := 0; i < len(slice); i++ {
		msg := fmt.Sprintf("for循环遍历切片: %d(%p) => %d(%p) ", i, &i, slice[i], &slice[i])
		fmt.Println(msg)
	}

	/**********    for range遍历切片  ****************/
	fmt.Println("\n\n4. 接下来看看for range遍历切片的情况:")
	for rr, vv := range slice {
		msg := fmt.Sprintf("for range 遍历切片, %d(%p) => %d(%p)", rr, &rr, vv, &vv)
		fmt.Println(msg)
	}

	/**
	从 Go 1.22 开始，每次循环都会创建新的循环变量，不再共用同一个地址,
	下边这段代码, 在1.22之前的版本输出的结果和代码文字描述能对应的上, 1.22之后, 就不再是同一个地址了, 值也不同
	*/
	for i := 0; i < len(slice); i++ {
		defer func() {
			msg := fmt.Sprintf("1.22之后, 我们不在是同一个内存地址: %p, 值每次也不一样: %v", &i, i)
			fmt.Println(msg)
		}()

	}

	/**
	如果想在1.22之后的版本实现上述代码执行结果和文字描述能对应的上, 关键点在于让defer每次捕获的i的地址都是相同的即可
	即: 把i的声明操作, 放在for循环的外层, 不过这样貌似没啥实际应用上的意义, 只是编码实现了1.22版本之前的效果
	如下所示:
	*/
	var i int
	for i = 0; i < len(slice); i++ {
		defer func() {
			msg := fmt.Sprintf("在1.22+版本中实现1.22之前版本的效果: 我们是同一个内存地址: %p, 值每次都一样: %v", &i, i)
			fmt.Println(msg)
		}()
		if i == len(slice)-1 {
			fmt.Println()
		}
	}
	fmt.Println()

	defer func() {
		msg := `
	这是在学习 Go 程序设计 中遇到的一个比较重要的一个警告。这是个 Go 语言的词法作用域规则的陷阱。
1. 1.22版本之前的for循环中, key的值始终是在同一个内存地址上重新赋值
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
	为什么会这样? 因为在1.22之前的for range中, key 和 value 一直是在各自同一个内存地址上发生的变化, 同时defer的特性
是等待函数执行完毕之后才开始执行的一个栈存储结构, 当defer内的函数开始执行的时候, 此时的i已经被
遍历完, 当i=2时, 满足i<len(arr), 所以i++后 i变为了3 所以输出的是 3 3 3
`
		fmt.Println(note + answer)
	}()

}
