package main

import (
	"fmt"
)

/**

Golang里的继承的实现方式

*/

type Animal struct {
	name string // 动物都有名字
	feet int8   // 动物都有腿
}

// 声明一个构造函数
func newAnimal(name string, feet int8) *Animal  {
	return &Animal{
		name: name,
		feet: feet,
	}
}


// 定义一个Animal的方法, 动物都会跑
func (a *Animal) sayRun ()  {
	fmt.Printf("我是%v, 我用%d条腿来跑\n\n", a.name, a.feet)
}
// 声明一个设置name的方法
func (a *Animal) setName(name string) *Animal {
	a.name = name
	return a
}
// 声明一个设置腿的方法
func (a *Animal) setLegs(legs int8) *Animal  {
	a.feet = legs
	return a
}



////////////////////////////////////////////////////////////

/**
	接下来定义一个狗子, 同时拥有结构体匿名字段Animal, 以达到继承目的
 */
type DogModel struct {
	Animal
}

/**
	定义一个兔子, 同时拥有结构体匿名字段Animal, 以达到继承目的
 */
type Rabbit struct {
	Animal
}

func main() {
 	//dog := newAnimal("狗子", 4)
 	// new一只狗子
	dog := new(DogModel)
	fmt.Printf("%#v \n", dog)
	dog.setName("狗子").setLegs(4).sayRun()

	// new一只兔子
	rabbit := new(Rabbit)
	rabbit.setName("兔子").setLegs(2).sayRun()

	// new一只喵
	cat := new(Animal)
	fmt.Printf("%#v\n\n", cat)
	cat.setName("喵").setLegs(4).sayRun()
	fmt.Printf("%#v \n", *cat) // *号用来去指针的值




}
