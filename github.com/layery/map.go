package main

import (
	"fmt"
)

func main() {
	// 通过这种方式创建的map是一个nil, 无法使用, 想要使用还必须要初始化他
	var map1 map[string]interface{}

	fmt.Printf("map 刚刚声明之后的值是一个nil %#v \n", map1)

	map1 = make(map[string]interface{})

	map1["name"] = "llf"

	map1["list"] = make(map[string]interface{})

	fmt.Printf("map is: %#v", map1)
}
