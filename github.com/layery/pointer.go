package main

import (
	"fmt"
)

func main() {

	// 指针的概念
	// 变量在声明完后, 其实是在内存中开辟了一个空间, 同时把值塞到了这个空间里,
	// 同时, 内存中不同的空间的地址是不同的, 可以理解为门牌号, 他们是十六进制表示的
	name := "llf"

	// 用&符可以取到变量name的值的地址, 是一个十六进制的值
	// 被取地址的变量是什么类型的, 则这个地址就是*type
	ptr := &name
	fmt.Println(ptr)
	fmt.Printf("%T", ptr)

	age := 30
	agePtr := &age
	fmt.Printf("%T", agePtr)

	// 通过取到的变量内存地址, 也可获取到该变量的值
	val := *agePtr
	fmt.Println(val)

	// 考虑如下函数
	/**
	这个函数中, 首先声明了一个变量a是一个指针, 由于a根本就没有值, 所以计算机根本没有给这个变量开辟一块内存空间
	所以当再执行*a = 100 , 给地址的值重新赋值的时候, 根本找不到这个地址, 所以会报错
	**/
	//test()

	// 上述例子如何修复它的问题呢? 引入了两个关键字 new 和 make

	// 1. new 就是用来开辟一块指定类型内存空间的, 有了空间也就有了地址, 同时会给该值初始化为该类型的默认值
	// 注意:  new 只能用来开辟一些基本类型的空间, 如: int, float, string, bool
	newInt := new(int)
	fmt.Println("%T", &newInt)

	// 2. make 只能用于开辟 slicn, map, chan 类型的空间

	var b = map[string]int
	fmt.Println(b)
	// b = make(map[string]int, 10)

	// fmt.Println("%T", b)

}

func test() {
	var a *int
	*a = 100
	fmt.Println(a)
}
