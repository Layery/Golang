package v1

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	model "study-gin/models"
	"study-gin/utils"

	"github.com/gin-gonic/gin"
)

type IndexControllerStruct struct {
	BaseController
}

var IndexController = new(IndexControllerStruct)

func (this *IndexControllerStruct) GetDistList(c *gin.Context) {

	datalist := make(map[string]string, 0)
	c.JSON(200, gin.H{
		"datalist": datalist,
	})
}

// 练习文件读取
func (this *IndexControllerStruct) TestFileReader(c *gin.Context) {
	// 利用ioutil包来处理文件相关操作
	dir, _ := os.Getwd() // get current project path

	// 在当前路径下创建一个文件
	err := ioutil.WriteFile("./test.log", make([]byte, 0), 000)

	// create file by os pkg
	fs, err := os.Create("./os-create.log")

	// os.open打开的文件只能读， 不能写
	//fs, err := os.Open("./os-create.log")

	if err != nil {
		log.Fatal(err)
	}
	err = os.Chmod("./os-create.log", 0644)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Chown("./os-create.log", 99, 99)

	n, e := fs.WriteString("刘龙飞手动阀阿斯蒂芬")

	if e != nil {
		log.Fatal("err is :", e)
	}

	log.Printf("os create file is %v, num is: %v", fs, n)

	// read file with ioUtil pkg
	str, err := ioutil.ReadFile(dir + "/controllers/v1/UserController.go")
	if err != nil {
		log.Fatal("err: ", err)
	}

	// 字符串是一个类型为字节的切片数组, 所以可以遍历它
	strList := string(str)
	resultStr := ""
	for _, v := range strList {
		resultStr += string(v)
	}
	hostname, _ := os.Hostname()
	groups, _ := os.Getgroups()
	this.JsonReturn(200, gin.H{
		"hostname":  hostname,
		"getuid":    os.Getuid(),
		"geteuid":   os.Geteuid(),
		"getgroups": groups,
	}, c)
}

func (index *IndexControllerStruct) TestClosure(c *gin.Context) {

	f1, msg1 := closureTest()
	f1()
	log.Println("msg is ====> ", msg1)

	f2, msg2 := closureTest()
	f2()
	log.Println("msg is ====> ", msg2)

	////////////////////////////////////////

	/**
		一个普通函数, 返回值是一个匿名函数, 则这个函数可以称之为闭包
	    闭包会将普通函数作用域的变量捕捉, 然后再编译执行阶段, 会一直保存这个变量的值不被销毁
	    所以闭包会修改变量, 应谨慎使用
	*/

	///////////////////////////////////////
	var i int8
	i = 2
	log.Println("before run i value is: ", i)
	test := closureTest2(i)
	test()
	test()
	test()

	log.Println("after run i value is: ", i)

}

func closureTest2(i int8) func() {

	return func() {
		i++
		i *= 2
		log.Println("now i's value is: ", i)
	}

}

func closureTest() (func(), string) {
	var i int8
	var f = func() {
		i++
		i *= 2
		log.Printf("that i's value is : %#v <====", i)
	}

	msg := fmt.Sprintf("after run that i's value is %v \n", i)
	return f, msg
}

func IndexAction(c *gin.Context) {

	// 来试试如何操作数据库吧
	log.Println("begin get user list")

	params := utils.CommonGetParams()

	log.Printf("获取所有的get参数%#v \n", params)

	userModel := model.User{}

	list := userModel.CreateUser()

	log.Printf("===> %#v", list)

}
