package main

import (
	"encoding/json"
	"fmt"
)

/**
Golang 中结构体字段的可见性 && json序列化
*/

// 声明一个结构体
type Person struct {
	name string
	age  int8
}

/**
Golang中, 如果定义的标识符是首字母大写的, 那么就是对外可见的, 可以理解为PHP中类成员方法所具有的public属性,
Golang中, 没有额外的关键字, 通过首字母大写, 即可对外可见, 首字母小写即只对当前这个包内可见
*/

// 如下所示, 声明了一个student, 类型是Person, 那么该类型就只是对当前包内可见的
type student Person

func main() {

	var zhangsan = new(Person)
	fmt.Printf("%#v", zhangsan)

	// 创建一个小明
	xiaoming := student{
		name: "小明",
		age:  30,
	}
	fmt.Printf("%#v \n\n", xiaoming)

	// 将小明格式化成json
	data, err := json.Marshal(xiaoming)
	if err != nil {
		fmt.Printf("found an err %v", err)
		return
	}
	fmt.Printf("json is %#v", data)

}
