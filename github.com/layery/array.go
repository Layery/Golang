package main

import (
	"fmt"
	"reflect"

	"github.com/layery/utils"
)

// 数组相关内容
func main() {
	// 数组的长度, 和数组的类型, 在声明数组的时候必须指定, 这是数组的两个必要条件
	// 初始化数组时, 编译器会默认赋给指定类型的默认值
	// 这种方式变量名后边可以没有=号
	var class []int
	utils.ConsoleLog(class)

	// 2. 也可以定义时就初始化自定义的值
	// 这种方式需要加上=号
	var class2 = [4]string{"北京", "上海", "天津", "重庆"}
	TypeOf := reflect.TypeOf(class2)
	fmt.Printf("class2的变量类型是: %s, kind: %s \n", TypeOf, TypeOf.Kind())

	// 3. 也可通过编译器自行获取数组长度
	var class3 = [...]int{2, 3, 4, 6}
	TypeOf = reflect.TypeOf(class3)
	fmt.Printf("class3的变量类型是: %s, kind: %s \n", TypeOf, TypeOf.Kind())

	// 4. 还可以通过索引值的方式来初始化数组, 类似于PHP的下标数组,但是其中未指定的键, 将会默认有个类型的初始值, 这点跟PHP不一样
	var class4 = [...]string{2: "java", 14: "php"}
	var class5 = [...]int{2: 2, 8: 14}
	var class6 = [...]bool{2: true, 8: true}
	fmt.Println(class4)
	fmt.Println(class5)
	fmt.Println(class6)

	// 5. 数组内不支持像PHP的任意类型, 指定是什么类型的数组, 数组里的值就都是什么类型的
	//  下边的方式定义数组会报错
	/**
		var class7 = [...]string{0: "aaaa", 1: 222}
		fmt.Println(class7)
	**/

	// 6. 数组的遍历, 可以使用for 或者 for range
	fmt.Println("采用for循环的方式遍历变量, 如下所示👇")
	for i := 0; i < len(class4); i++ {
		fmt.Println(i, "---->"+class4[i])
	}
	fmt.Println("采用for循环的方式遍历变量, 如上所示👆")

	fmt.Println()
	fmt.Println()

	// for range 方式 类似于PHP的foreach遍历
	// 这里边也可以使用匿名变量, 来忽略一个变量
	fmt.Println("采用for range循环的方式遍历变量, 如下所示👇")
	for i, v := range class6 {
		fmt.Printf("当前i: %v, i的指针%p ===> 当前v: %v, v的指针: %p \n", i, &i, v, &v)
	}

	fmt.Println()

	for _, v := range class6 {
		vv := v
		fmt.Printf("当前v: %v, v的指针: %p \n", vv, &vv)
	}

	fmt.Println("采用for range循环的方式遍历变量, 如上所示👆")

	fmt.Println()

	// go 的二维数组真恶心 变量名后边先跟一个一级数组的长度, 再跟一个二级数组的长度
	var multiArray1 = [3][2]int{
		{111, 222},
		{111, 222},
		{111, 222},
	}
	fmt.Printf("multiArray1的值: %v \n\n", multiArray1)

	// 如果想使用...的方式, 写成如下两种方式
	// 1. 外层...简写
	var multiArray2 = [...][]byte{{1, 2, 9}, {3, 1}}
	fmt.Println(multiArray2)
	fmt.Printf("multiArray2的类型 %v, kind: %v \n\n", reflect.TypeOf(multiArray2), reflect.TypeOf(multiArray2).Kind())

	// 2. 两层都不写... todo: 这种居然也可以, 这种虽然可以,
	// 但是没有了[...]自动推导, 这其实是个切片类型了
	var multiArray3 = [][]int{
		{111, 2222},
		{111, 2222},
	}
	fmt.Println(multiArray3)
	fmt.Printf("multiArray3 %v, kind: %v \n\n\n", reflect.TypeOf(multiArray3), reflect.TypeOf(multiArray3).Kind())

	// for range 嵌套 for 遍历多维数组
	for i, row := range multiArray2 {
		fmt.Println(i, row)
		for j := 0; j < len(row); j++ {
			fmt.Println(row[j])
		}
	}

	fmt.Printf("\n\n\n 一道练习题: \n")

	// 练习题, 找出值相加为8的键的组合
	var test = []int{1, 2, 3, 4, 5, 6, 7, 8}

	cnt := len(test)
	fmt.Println("cnt: ", cnt)
	var list [][]int
	for i := 0; i < cnt; i++ {
		for j := 0; j < cnt; j++ {
			if j <= i {
				continue
			}
			if test[i]+test[j] == 8 {
				temp := []int{i, j}
				list = append(list, temp)
			}
		}
	}

	fmt.Println(list)

}
