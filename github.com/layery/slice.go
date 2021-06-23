package main

import (
	"fmt"
	"reflect"
)

func ArraySearch(needle interface{}, hystack interface{}) (index int) {
	index = -1
	switch reflect.TypeOf(hystack).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(hystack)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(needle, s.Index(i).Interface()) == true {
				index = i
				return
			}
		}
	}
	return
}

func main() {
	// 切片的定义, 类似于数组, 只是不需要定义数组的长度了

	var slice1 []string
	var slice2 []int
	slice3 := []bool{}

	fmt.Println(slice1, "---", slice2, "---", slice3)

	// 基于数组定义一个切片, 这类似于PHP的数组函数array_slice
	arr := [5]int{1, 2, 3, 4, 5}
	slice4array := arr[2:] // 采用:号语法时, :号左边为切片的起始位置, 右边为截至位置

	slice4array1 := arr[1:2]

	fmt.Println(slice4array1)
	fmt.Println(slice4array)

	// 基于数组做出来的切片, 同样可以再次被切
	// slice通过array[i:j]来获取，其中i是数组的开始位置，j是结束位置，但不包含array[j]，它的长度是j-i。
	// slice 默认从0开始slice[:3] 等价于 slice[0:3]
	// slice[3:] 等价于 slice[3: len(array)]
	// slice[:] 等价于 slice[0: len(array)]       <====  这个写法好特么奇葩
	slice4slice := slice4array[:2]
	fmt.Println("debug", slice4slice)

	slice5 := arr[:]
	fmt.Println("从一个数组里面直接获取slice, 奇葩的:写法", slice5)

	// slice 是一个引用类型, 他总是执行一个底层的array,

	bool_slice := []int{2, 3}
	fmt.Printf("%T", bool_slice)

	// len 函数 和 cap函数
	// len 表示切片从起始位置到截止位置的长度
	// cap : 容量, 表示切片从起始位置, 到终点位置的长度

	var slice6 = make([]int, 3, 5)

	fmt.Printf("len=%d cap=%d value=%v \n", len(slice6), cap(slice6), slice6)

	// 声明一个多维切片 , 类似于声明一个多维数组
	multiSlice := [][]int{{1, 2}, {111, 222, 333}}
	fmt.Println(multiSlice)

	// 遍历数组的方式, 同样适用于遍历切片
	for i := 0; i < len(multiSlice); i++ {
		fmt.Printf("i=%d row=%v \n", i, multiSlice[i])
	}

	fmt.Print("\n\n")

	// append函数可以为切片增加单元
	currSlice := make([]int, 0, 0)
	for i, v := range multiSlice {
		fmt.Println(v)
		currSlice := append(v, i)
		fmt.Printf("do append slice %v", currSlice)
	}
	fmt.Println(currSlice)
	// append也可以直接将一个切片追加到另一个切片中
	nowSlice := append(currSlice, slice4slice...)

	fmt.Println("append也可以直接将一个切片追加到另一个切片中, 这个...的语法也好奇葩(代表将元素从切片中拆分出来, 然后全部追加到指定slince)", nowSlice)

	// 切片的删除
	// 想删除某个目标, 就把目标的索引之前的切出来, 同时追加上它之后的切片, 这样就变相的删除了指定的切片 <== 好变态啊赶脚, PHP果然是最好的语言
	newSlice := []string{"北京", "上海", "天津", "重庆"}
	del := append(newSlice[:2], newSlice[3:]...)
	fmt.Println(del)

	// 元素类型为map的切片, make函数此时只完成了切片的初始化, map里的它并没有初始化, 还是nil
	var mapSlice = make([]map[string]int, 8, 8)
	// 还需要完成内部map的初始化
	mapSlice[0] = make(map[string]int, 8)
	mapSlice[0]["zhangsan"] = 30
	fmt.Println(mapSlice)

	// 元素类型为切片的map

	var sliceOfMap = make(map[int][]int, 3)
	sliceOfMap[0] = []int{1, 2, 3}
	fmt.Printf("len=%d cap=%T val=%v \n\n", len(sliceOfMap), sliceOfMap, sliceOfMap)
	v, true := sliceOfMap[0]
	if true {
		fmt.Printf("ok %v", v)
	} else {
		fmt.Printf("no ok %v", v)
	}

}
