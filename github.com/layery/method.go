package main

import "fmt"

/**

结构体中的方法和接收者

Golang中一切皆是类型, 所以会出现int, string类型也支持定义方法的魔性语法,
这里可以类比PHP的一切皆是对象

方法: 可以理解为PHP里的类方法, 是这个类所属的
接收者: 可以理解为PHP的类本身($this 或 self)

*/

// 声明一个Person类型的结构体,
// 当一个变量的首字母是大写的, 意味着它是对外可见的
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
//    方法里有两个传参, 第二个是接收者的类型, 第一个是接收者类型的首字母简写(约定成俗)
func (p Person) Say() {
	// fmt.Printf("%s %d", p.name+" 今年"+p.age+"岁了")
	fmt.Printf("hello my name is %s \n", p.name)
}

func (p *Person) Hello(name string) {
	fmt.Printf("hi I'm " + name + "\n")
}

/**
方法的传参是指针类型的时候, 可以修改结构体内部的属性
考虑一下3种情况下使用指针传参:
	1. 需要修改接受者的值
	2. 接受者是一个比较大的对象, 拷贝它会比较消耗资源, 但是拷贝指针代价就小很多
	3. 保证一致性, 如果接收者的某一个方法使用了指针传参, 则为了保证一致性, 其他的方法也要使用指针传参
*/
func (p *Person) SetName(name string) *Person {
	p.name = name
	fmt.Println("now person's name is " + p.name)
	return p
}

func main() {

	// 利用构造函数, 获取一个person结构体
	// 相当于PHP的new Object("llf", 30)

	man := newPerson("llf", 30)
	fmt.Printf("man 的类型是一个指针 %T \n", man)
	// 接下来, 调用我们的say方法, 由于构造函数返回的是一个指针
	// 所以我们需要先取到指针的值, 再调用say方法
	(*man).Say()

	// golang里, 也可以不用显示的写出构造函数的值, 而直接用指针调用方法,
	man.Say()

	man.Hello("aaaa")

	man.SetName("weidingyi")

	man.Say()

	fmt.Printf("现在man的类型是%T, man的值是%v \n", man, man)
}
