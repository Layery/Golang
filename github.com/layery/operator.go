package main

import "fmt"

func main() {
	/**
	Go里的运算操作符, 和PHP的基本类似, 有以下几点需要注意的
	  	1. 不同类型间的变量不能直接比较, 并且必须为基本数值类型(例如: int, int32)

	位移运算有左移位, 右移位两种运算,
		1. 左移位时 低位补0
		2. 右移位是 地位被舍弃

	*/

	// 异或运算 (具体的计算方式也是先转成2进制,再计算, 每一位上, 如果不相同则为1, 相同的则为0, 例如: 5^1=4)
	a, b := yihuo()
	fmt.Printf("a: %#v --- b: %#v \n\n", a, b)

	weiyi()

	ltgt()

	// 常量子表达式的顺序有可能影响到最终的估值结果
	constExpression()




}

func constExpression() {
	/**
	 	这里居然会输出的是2.2 真是奇了个葩 ...
	 */
	var x = 1.2 + 1/2
	fmt.Printf("%T %v \n", x, x)

	xx := 3/2*0.1
	fmt.Println(xx)
	fmt.Printf("%T %v \n", xx, xx)
}

// ltgt函数, 意思为小于(<), 大于(>)号 ^_^!
func ltgt() {
	/**
		<, <=, >, >= 这几个运算符中, 两个值的类型一定要相同, 并且只能是(整型, 浮点型, 字符串型) 才可参与计算
	 */

	str1 := "layery"
	str2 := "Layery"

	if str2 < str1 {
		fmt.Println("aaaa")
	} else {
		fmt.Println("bbbb")
	}

}

func weiyi() {
	/**
		貌似它是将十进制数转成二进制后, 再开始移位运算, 最终将结果再转为十进制数
		3(10) = 011(2)
		011 << 2 结果是 01100
		01100(2) = 12(10)

	 */
	a := 3
	b := 4

	a = a << 2 // a 左移2位

	b = b >> 1 // b 右移5位

	fmt.Printf("%v == %v \n\n", a, b)
}

// 异或运算, 不占用第三个变量的情况下, 交换两个变量的值
func yihuo() (int, int) {
	var a, b int
	a = 2
	b = 3
	a = a ^ b
	fmt.Println(a)
	b = a ^ b
	fmt.Println(b)
	a = a ^ b
	fmt.Println(a)
	return a, b
}