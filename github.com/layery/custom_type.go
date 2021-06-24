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

/**

	当时的困惑, 现在已有了初步的释疑 , golang中, 不仅仅可以为struct结构体
	定义自定义的方法, 也可以为int, string, bool, 等基本类型去定义自定义方法,
	但是仅限于非系统内置的基本类型, 也就是说, 我们可以去定义自己的int类型,
	然后就可以为我们的MyInt类型去定义自定义的方法了

**/
