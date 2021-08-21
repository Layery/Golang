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
	Name    string
	age     int8
	Address string `json:"local, int"`    // 表示用local代替原始的Address的字段名
	sex     int8   `json:-`               // 表示对该字段不进行序列化
	// omitempty是忽略该字段所有空值的场景
	Email   string `json:email,omitempty` // 由于Go中字段首字母大写, 代表外部可访问的变量, 也可以是使用json TAG 更加定制化的生成json格式数据
}

/**
Golang中, 如果定义的标识符是首字母大写的, 那么就是对外可见的, 可以理解为PHP中类成员方法所具有的public属性,
Golang中, 没有额外的关键字, 通过首字母大写, 即可对外可见, 首字母小写即只对当前这个包内可见
*/

// 如下所示, 声明了一个student, 类型是Person, 那么该类型就只是对当前包内可见的
type student Person

func main() {

	var zhangsan = new(Person)
	fmt.Printf("%#v \n", zhangsan)

	// 创建一个小明
	xiaoming := student{
		Name:    "小明", // 由于Name字段首字母大写, 所以Name字段可以被json包读取到
		age:     30,   // 由于age字段不是首字母大写, 所以age字段不会被json包读取到
		Address: "",
		sex:     2,
		Email:   "www.baidu.com",
	}
	fmt.Printf("main包自己本身可以读取得到age字段 %#v \n\n", xiaoming)

	// 将小明格式化成json
	data, err := json.Marshal(xiaoming)
	if err != nil {
		fmt.Printf("found an err %v", err)
		return
	}
	fmt.Printf("json包读取不到age字段%#s", data)
}
