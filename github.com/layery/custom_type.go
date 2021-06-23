package main

import "fmt"

// 自定义类型 和类型别名的实例

// 1. 使用type关键字, 来声明一个自己的类型

type MyInt int

func main() {
	var i MyInt
	fmt.Printf("%T", i)
}

// 2. 类型的别名 (实际上是给string类型加了个别名)
//    类似于php里引用某个对象的时候, use Model as MyModel
type MySting = string

// 3. 那么 自定义类型的应用场景有哪些呢?
