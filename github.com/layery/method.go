package main

import "fmt"

/**

结构体中的方法和接收者

方法: 可以理解为PHP里的类方法, 是这个类所属的

*/

// 声明一个Person类型的结构体
type Person struct {
	name string
	age  int8
}

// 1. 创建一个person类型的构造函数
func newPerson(name string, age int8) *Person {
	return &Person{
		name: name,
		age:  age,
	}
}

// 2. 为刚刚创建的Person类型的结构体, 定义一个方法
func (p Person) Say() {
	// fmt.Printf("%s %d", p.name+" 今年"+p.age+"岁了")
	fmt.Printf("%v", p)
}

func main() {

	// 利用构造函数, 获取一个person结构体
	// 相当于PHP的new Object("llf", 30)

	man := newPerson("llf", 30)
	fmt.Printf("man 的类型是一个指针 %T", man)
	// 接下来, 调用我们的say方法

}
