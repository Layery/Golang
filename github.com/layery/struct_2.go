package main

import "fmt"

/**

匿名结构体  &&  结构体的匿名字段

*/

// 1. 结构体的匿名字段, 即: 只有类型声明, 没有值, 那么如何访问结构体的字段呢?
//    答案是访问它的类型即是访问了字段
//    结构体中的字段具有唯一性, 不可重复
type Human struct {
	string
	int
}

/**

2. 嵌套结构体, 考虑如下两个结构体, Person中有city字段, 同时Address中也有city字段

*/
type Address struct {
	province string
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

func main() {
	person1 := new(Human)

	person2 := Human{
		"weidingyi",
		30,
	}

	fmt.Println(person1.int, person2.string, "\n\n")

	// 使用嵌套结构体
	person3 := Person{
		name: "layery",
		age:  30,
		sex:  1,
		Address: Address{
			province: "北京-zi",
			city:     "北京-zi",
			area:     "朝阳区-zi",
		},
	}

	fmt.Printf("%#v\n", person3)

	fmt.Printf("可以通过.的方式来访问结构体的内部字段 %#v\n", person3.Address.province)

	// 也可以直接访问结构内部的子结构体的字段
	fmt.Printf("也可以直接访问子结构体的内部字段%#v\n", person3.Address.city)

	// 当父结构体嵌套了多个结构体时, 如果子结构体中有重名的字段, 则访问的时候,
	// 必须指定访问的是哪个子结构体
	fmt.Printf("访问了Address下的city: %#v\n", person3.Address.city)
	fmt.Printf("访问了Info下的city: %#v\n", person3.Info.city)

	fmt.Printf("\n\n\n%#v\n", person3)
}
