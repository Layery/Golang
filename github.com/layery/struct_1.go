package main

import "fmt"

// 结构体
// 结构体类似于PHP的Object

type PersonModel struct {
	name, addr string // 相同类型的字段, 可以写在同一行
	age        int8
}

func getOlder(p1 PersonModel, p2 PersonModel) (int8, PersonModel) {
	if p1.age > p2.age {
		return p1.age, p1
	} else {
		return p2.age, p2
	}
}

/**
    因为结构体是值类型, 所以在传递时会占用较大的内存, 为了解决这个问题, 一般的, 在构造函数里,
	我们只返回结构体的指针即可, 这样节省内存, 提升程序效率
*/
func newPersonModel(name, addr string, age int8) *PersonModel {
	return &PersonModel{
		name: name,
		addr: addr,
		age:  age, // 这里结尾的,号不要忘记了
	}
}

func main() {
	var zhangsan PersonModel
	zhangsan.name = "张三"
	zhangsan.addr = "北京"
	zhangsan.age = 126
	fmt.Println(zhangsan)

	// 第2种声明方式, 顺序必须和声明时的顺序一致
	wangwu := PersonModel{"王五", "沙河", 20}
	fmt.Println(wangwu)
	// 第3种方式, 可以自定义顺序
	zhaoliu := PersonModel{
		age:  127,
		name: "赵六",
		addr: "龙门客栈",
	}
	fmt.Println(zhaoliu)

	// 匿名结构体, 只使用一次, 无需提前定义结构体类型, 直接使用
	var lisi struct {
		name string
		age  int8
	}
	fmt.Println(lisi)

	age, res := getOlder(zhangsan, zhaoliu)
	fmt.Printf("%v的年龄大, 他%d岁了 \n", res.name, age)

	// go的struct 类似于PHP的class, 但是go里没有构造函数的概念
	// go里通过函数去给结构体声明一个构造函数, 并约定成俗的使用new前缀命名函数
	test := newPersonModel("llf", "天津", 30)
	fmt.Println(test.name)
}
