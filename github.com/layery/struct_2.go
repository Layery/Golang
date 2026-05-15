package main

import (
	"fmt"
	"sync"

	"github.com/gookit/goutil/dump"
)

/**

匿名结构体  &&  结构体的匿名字段

匿名结构体:
	1. 不定义类型名称, 直接使用结构体赋值,可以把有关联的全局变量组合在一起使用, 或者局部的匿名结构体, 只作为临时一次性使用
	2. 这个感觉不是很常用啊
结构体的匿名字段:
	1. 顾名思义, 结构体中的某个字段, 不进行字段命名, 直接声明类型, 常见的用法如: 多个goroutine的共享变量
	2. 互斥锁的场景比较常用
*/

// Human
//  1. 结构体的匿名字段, 即: 只有类型声明, 没有值, 那么如何访问结构体的字段呢?
//     答案是访问它的类型即是访问了字段
//     结构体中的字段具有唯一性, 不可重复
type Human struct {
	string
	int
}

// Address 2. 嵌套结构体, 考虑如下两个结构体, Person中有city字段, 同时Address中也有city字段
type Address struct {
	Province string
	city     string
	area     string
}

type Info struct {
	city string
}

type Person struct {
	name    string
	age     int8
	sex     int8
	Address // 在这里, 我将Address结构体嵌套到了Person里
	Info
}

// anonymityStruct
func anonymityStruct() {
	var person struct {
		name string
		age  int8
	}
	person.name = "llf"
	person.age = 33
	fmt.Println(person)
}

// anonymityColumn
func anonymityColumn() {
	person1 := new(Human)

	fmt.Println(person1)
	// 使用嵌套结构体
	person3 := Person{
		name: "layery",
		age:  30,
		sex:  1,
		Address: Address{
			Province: "北京-Province",
			city:     "北京-city",
			area:     "朝阳区-area",
		},
	}

	dump.P(person3)

	fmt.Printf("可以通过.的方式来访问结构体的内部字段 %#v\n", person3.Address.Province)

	// 也可以直接访问结构内部的子结构体的字段
	// 此时用的go1.17的版本, 貌似不能直接访问子结构体的字段
	fmt.Printf("此时用的go1.17的版本, 貌似不能直接访问子结构体的字段%#v\n", person3.Info.city)

	// 当父结构体嵌套了多个结构体时, 如果子结构体中有重名的字段, 则访问的时候,
	// // 必须指定访问的是哪个子结构体
	fmt.Printf("访问了Address下的city: %#v\n", person3.Address.city)
	fmt.Printf("访问了Info下的city: %#v\n", person3.Info.city)
}

func anonymityLock() {
	type counter struct {
		sync.Mutex // 将sync.Mutex作为匿名字段, 可以很方便的写加锁释放锁的逻辑
		num        int
	}
	cnt := counter{num: 0}
	var wg sync.WaitGroup
	for i := 1; i <= 100; i++ { // go1.22之后,for循环每次都是迭代新的变量i
		wg.Go(func() {
			cnt.Lock()
			defer cnt.Unlock()
			cnt.num += i
		})
	}
	wg.Wait()
	dump.P(cnt)
}

func main() {

	anonymityStruct() // 匿名结构体

	anonymityColumn() // 结构体匿名字段

	anonymityLock() // 结构体匿名字段的常见用法
}
