package main

import (
	"fmt"
	"reflect"

	"github.com/gookit/goutil/dump"
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

func testSlice1() {
	// 采用字面量的方式定义一个切片, 会自动推断出切片的长度和容量,
	// 虽然你看不到底层数组, 但是它底层的确就是一个len和cap都是5的匿名数组
	var slice1 = []int{1, 2, 3, 4, 5}
	fmt.Printf("slice1 len: %v, cap: %v \n", len(slice1), cap(slice1))

	// 采用make函数定义一个切片, 需要指定切片的长度和容量,
	var slice2 = make([]int, 5, 10)
	fmt.Printf("slice2 len: %v, cap: %v \n", len(slice2), cap(slice2))

	// 如果不指定容量, 那么切片的长度和容量是一样的
	var slice3 = make([]int, 5)
	fmt.Printf("slice3 len: %v, cap: %v \n", len(slice3), cap(slice3))

	// 切片底层的那个指针，指向的就是切片当前视窗里的第一个元素在内存中的地址
	fmt.Printf("slice1 的指针: %p, slice1的1号元素的指针: %p \n", slice1, &slice1[0])

	// 当切片发生截取时
	slice4 := slice1[2:3]
	fmt.Printf("slice4 len: %v, cap: %v \n", len(slice4), cap(slice4))

}

// demoSliceDefinition 演示切片定义与基础操作
// 参数: 无
// 输出: 无
func demoSliceDefinition() {
	// 切片的定义, 类似于数组, 只是不需要定义数组的长度了

	var slice1 []string
	var slice2 []int
	slice3 := []bool{true, false, true}

	fmt.Printf("%#v, %#v, %#v, \n\n", slice1, slice2, slice3)

	// 基于数组定义一个切片, 这类似于PHP的数组函数array_slice
	arr := [5]int{1, 2, 3, 4, 5}
	slice4array := arr[2:] // 采用:号语法时, :号左边为切片的起始位置, 右边为截至位置

	slice4array1 := arr[:0]

	fmt.Println(slice4array)
	fmt.Printf("===> %+v \n", slice4array1)

	// 基于数组做出来的切片, 同样可以再次被切
	// slice通过array[i:j]来获取，其中i是数组的开始位置，j是结束位置，但不包含array[j]，它的长度是j-i。
	// slice 默认从0开始slice[:3] 等价于 slice[0:3]
	// slice[3:] 等价于 slice[3: len(array)]
	// slice[:] 等价于 slice[0: len(array)]       <====  这个写法好特么奇葩

	slice4slice := slice4array[:2]
	fmt.Printf("sliec4array: %v, sliec4slice: %v \n", slice4array, slice4slice)

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
}

// demoMultiDimensionSlice 演示多维切片与遍历
// 参数: 无
// 输出: 无
func demoMultiDimensionSlice() {
	// 声明一个多维切片 , 类似于声明一个多维数组
	multiSlice := [][]int{{1, 2}, {111, 222, 333}}
	fmt.Println(multiSlice)

	// 遍历数组的方式, 同样适用于遍历切片
	for i := 0; i < len(multiSlice); i++ {
		fmt.Printf("i=%d row=%v \n", i, multiSlice[i])
	}

	fmt.Print("\n\n")
}

// demoSliceAppend 演示append函数操作切片
// 参数: multiSlice [][]int - 输入的多维切片，用于遍历并追加元素
// 输出: 无
func demoSliceAppend(multiSlice [][]int) {
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
}

// demoSliceDelete 演示切片的删除操作
// 参数: 无
// 输出: 无
func demoSliceDelete() {
	// 切片的删除
	// 想删除某个目标, 就把目标的索引之前的切出来, 同时追加上它之后的切片, 这样就变相的删除了指定的切片 <== 好变态啊赶脚, PHP果然是最好的语言
	newSlice := []string{"北京", "上海", "天津", "重庆"}
	del := append(newSlice[:2], newSlice[3:]...)
	fmt.Println(del)
}

// demoMapSlice 演示元素类型为map的切片和元素类型为切片的map
// 参数: 无
// 输出: 无
func demoMapSlice() {
	// 元素类型为map的切片, make函数此时只完成了切片的初始化, map里的它并没有初始化, 还是nil
	var mapSlice = make([]map[string]int, 8, 8)

	// 还需要完成内部map的初始化
	mapSlice[0] = make(map[string]int, 8)
	mapSlice[0]["zhangsan"] = 30
	dump.Print(mapSlice)

	// 元素类型为切片的map

	var sliceOfMap = make(map[int][]int, 3)
	sliceOfMap[0] = []int{1, 2, 3}
	fmt.Printf("len=%d type=%T val=%v \n\n", len(sliceOfMap), sliceOfMap, sliceOfMap)
	v, true := sliceOfMap[0]
	if true {
		fmt.Printf("ok %v", v)
	} else {
		fmt.Printf("no ok %v", v)
	}

	fmt.Println()
	fmt.Println()
}

// demoSliceCompare 演示切片比较与面试题
// 参数: 无
// 输出: 无
func demoSliceCompare() {
	/**
		切片之间是不能比较的，我们不能使用==操作符来判断两个切片是否含有全部相等元素。
	 	切片唯一合法的比较操作是和nil比较。
		一个nil值的切片并没有底层数组，一个nil值的切片的长度和容量都是0。  <----这段话有点绕口
	 	但是我们不能说一个长度和容量都是0的切片一定是nil，例如下面的示例:
	*/
	var myslice1 []int
	myslice2 := []int{}
	myslice3 := make([]int, 0)

	fmt.Printf("%#v \n", myslice1)

	fmt.Printf("%#v \n", myslice2)
	fmt.Printf("%#v \n", myslice3)

	/*********************   以下是一道面试题  ************************/
	var a = make([]string, 5) // 这里已经初始化切片了, 并且用字符串的零值来填充这个切片了
	for i := 0; i < 10; i++ {
		a = append(a, fmt.Sprintf("%v", i)) // 这一步再扩容, 等于是在5个空字符串的基础上, 拼接了0-9
	}
	fmt.Printf("%#v \n\n", a) // 输出: "     0123456789", 注意! 最开始输出的是几个空字符串
	/*********************   以上是一道面试题  ************************/
}

// demoSliceGrowth 演示切片的扩容规则
// 参数: 无
// 输出: 无
func demoSliceGrowth() {
	/**
	当前golang版本是1.18, 当前版本下, 切片的扩容规则是:
		1. 当所需容量大于oldcap的2倍时, 直接把所需容量赋值给newcap,
		2. 当所需容量小于oldcap的2倍时:
			2.1: 如果oldcap小256, 则视为小切片, 则新容量就是2倍的旧容量
			2.2: 如果oldcap大于等于256, 视为大切片, 则新容量是旧容量的1.25倍
	*/

	slice7 := []int{10, 20, 30, 40}

	newSlice7 := append(slice7, 50)

	fmt.Printf("slice7, len: %v, cap: %v \n", len(slice7), cap(slice7))
	fmt.Printf("newSlice7, len: %v, cap: %v \n", len(newSlice7), cap(newSlice7))
}

// demoSliceCopy 演示切片的深拷贝与浅拷贝
// 参数: 无
// 输出: 无
func demoSliceCopy() {
	/*****************    切片的深拷贝, 浅拷贝  ********************/
	fmt.Println()
	fmt.Println()

	array := [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	s1 := array[0:3:5]
	s2 := array[0:4:8]

	s1[2] = 100

	for s1key, s1val := range s1 {
		fmt.Printf("%v(%v) => %v(%v) \n", s1key, &s1key, s1val, &s1val)
	}
	fmt.Println()
	for s2key, s2val := range s2 {
		fmt.Printf("%v(%v) => %v(%v) \n", s2key, &s2key, s2val, &s2val)
	}

	_ = s2
}

func main() {

	// 切片的底层是一个数组, 无论这个数组你是否能直观的看到, 它的底层就是一个数组
	// 切片本身的数据结构只有3个字段, 分别是: 指向底层数组的指针, 切片的长度, 切片的容量
	testSlice1()
	return

	// 调用切片定义与基础操作演示
	demoSliceDefinition()

	// 调用多维切片与遍历演示
	demoMultiDimensionSlice()

	// 定义多维切片供append演示使用
	multiSlice := [][]int{{1, 2}, {111, 222, 333}}
	// 调用append函数操作切片演示
	demoSliceAppend(multiSlice)

	// 调用切片删除操作演示
	demoSliceDelete()

	// 调用map切片演示
	demoMapSlice()

	// 调用切片比较与面试题演示
	demoSliceCompare()

	// 调用切片扩容规则演示
	demoSliceGrowth()

	// 调用切片的深拷贝与浅拷贝演示
	demoSliceCopy()
}
