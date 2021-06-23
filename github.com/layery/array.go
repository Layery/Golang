package main

import (
	"fmt"
	f "fmt"
)

// 数组相关内容
func main() {
	// 数组的长度, 和数组的类型, 在声明数组的时候必须指定, 这是数组的两个必要条件
	// 初始化数组时, 编译器会默认赋给指定类型的默认值
	// 这种方式变量名后边可以没有=号
	var class [2]int
	fmt.Println(class)

	// 2. 也可以定义时就初始化自定义的值
	// 这种方式需要加上=号
	var class2 = [4]string{"北京", "上海", "天津", "重庆"}
	fmt.Println(class2)

	// 3. 也可通过编译器自行获取数组长度
	var class3 = [...]int{2, 3, 4, 6}
	fmt.Println(len(class3))

	// 4. 还可以通过索引值的方式来初始化数组, 类似于PHP的下标数组,但是其中未指定的键, 将会默认有个类型的初始值, 这点跟PHP不一样
	var class4 = [...]string{2: "java", 14: "php"}
	var class5 = [...]int{2: 2, 14: 14}
	var class6 = [...]bool{2: true, 14: true}
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
	for i := 0; i < len(class4); i++ {
		fmt.Println(i, "---->"+class4[i])
	}

	// for range 方式 类似于PHP的foreach遍历
	// 这里边也可以使用匿名变量, 来忽略一个变量
	for i, v := range class6 {
		fmt.Println(i, v)
	}
	for _, v := range class6 {
		fmt.Println(v)
	}

	// go 的二维数组真恶心 变量名后边先跟一个一级数组的长度, 再跟一个二级数组的长度
	var multiArray = [3][2]int{
		{111, 222},
		{111, 222},
		{111, 222},
	}
	fmt.Println(multiArray)
	// 如果想使用...的方式, 写成如下两种方式
	// 1. 外层...简写
	var multiArray2 = [...][]int{
		{111, 2222},
		{111, 2222},
	}
	fmt.Println(multiArray2)
	// 2. 两层都不写... todo: 这种居然也可以
	var multiArray3 = [][]int{
		{111, 2222},
		{111, 2222},
	}
	fmt.Println(multiArray3)
	// for range 嵌套 for 遍历多维数组
	for i, row := range multiArray3 {
		fmt.Println(i, row)
		for j := 0; j < len(row); j++ {
			fmt.Println(row[j])
		}
	}

	// 练习题, 找出值相加为8的键的组合
	var test = []int{1, 2, 3, 4, 5, 6, 7, 8}
	f.Println(test)
	var tmp [4]int
	fmt.Println(tmp)
	// for i, v := range test {
	// 	a := [2]int{
	// 		i, v
	// 	}
	// 	append(tmp, a)
	// }
	// fmt.Println(tmp)

}
