package main

import "fmt"

func main() {

	s := "layery"

	// 字符串类型转数组
	var arr = []byte(s) // 将字符串转为byte类型
	arr[0] = 'L'        // 替换数组第0个单元

	fmt.Printf("字符串通过[]byte()的方式转换数组后值为: %+v\n", arr)
	for i, v := range arr {
		fmt.Printf("当前i: %d, 当前v: %c\n", i, v)
	}

	newString := string(arr) // 将数组重新转为字符串
	fmt.Println(newString)

	// 第2种方式
	str := "hello"
	str = "c" + str[1:] // 字符串虽然不能更改, 但可以做切片操作
	fmt.Printf("%s\n", str)

}
