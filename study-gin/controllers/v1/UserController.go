package v1

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"reflect"
	model "study-gin/models"
)
type UserControllerStruct struct {
	BaseController
}

var UserController = new(UserControllerStruct)

func TestArrayMap(c *gin.Context) {
	// 普通数组, (数组的长度和类型是一个整体, 长度可以用...来省略代表不定长度的数组)
	// 普通数组的应用场景不如切片多, 因为他不能动态变增加/减少数组的key
	arr1 := [3]int{1, 2, 5}

	// 切片1
	slice1 := arr1[:]

	// 切片2
	slice2 := []interface{}{2, 3, 4}

	// 切片合并
	merge := append(slice2, slice1)

	// 再合并
	merge = append(merge, merge) // 这居然也可以, 有点php的味道了


	// 数组追加单元
	c.JSON(200, gin.H{
		"slice_list": map[string]interface{}{
			"slice_1_len": len(slice1),
			"slice_1_cap": cap(slice1),
			"slice_2_len": len(slice2),
			"slice_2_cap": cap(slice2),
			"merge": merge,
		},
		"test": map[string]interface{}{
			"aaa": [3]string{"a", "b", "c"},
		},
		"datalist": arr1,
		"arr1_len": len(arr1),
		"arr1_cap": cap(arr1),
	})
}

func TestMapSlice(c *gin.Context) {
	// 普通map
	map1 := make(map[string]interface{})
	map1["name"] = "llf"

	// 嵌套map的二维map
	var map2 = map[string] string {"age": "30", "addr": "shijiazhuang"}
	map1["map2"] = map2

	// 嵌套数组的二维map
	arr1 := [3]int {3, 4, 5}
	map1["arr1"] = arr1

	// 嵌套切片的二维map
	slice1 := []string{"I", "am", "slice"}
	map1["slice1"] = slice1
	map1["slice1_len"] = len(slice1)
	map1["slice1_cap"] = cap(slice1)

	// 遍历这个map
	str := ""

	for i, row := range map1 {
		dataType := fmt.Sprintf("dataType is: %v", reflect.TypeOf(row))
		str += "i = " + i + ", row = " + fmt.Sprintf("%v, ", row) + dataType + "\n"
	}

	//c.String(200, fmt.Sprint(str))

	c.JSON(200, gin.H{
		"datalist": map1,
	})
}


func (u *UserControllerStruct) GetUserList(c *gin.Context) {
	userModel := model.User{}
	list := userModel.GetUserList()
	c.JSON(200, gin.H{
		"list": list,
	})
}

func CommonGetPrams() string {
	return "test"
}

func LoginByStruct(c *gin.Context) {
	// 将请求参数绑定到struct结构体上
	// gin支持将参数自动绑定到结构体上，支持get/post请求， 同时也支持http请求的body请求体为json/xml格式的参数
	// go是强类型语言， 结构体中规定了email是字符串类型的数据， 则， 前端传参必须传字符串类型的数据， （考虑email传111111的情况）
	type User struct {
		Name string `json:"name" form:"name"`
		Email string `json:"email" form:"email" binding:"required"`
	}
	var u User
	log.Println(u)
	type tempUser struct {
		User
		Address string
	}
	var temp_user tempUser
	_ = c.ShouldBindJSON(&temp_user)
	log.Println(temp_user)
	c.JSON(200, gin.H{
		"status": 200,
		"obj": u,
	})
}



/**
	加载前端参数并赋值到struct上
 */
func LoadParamsByStruct(c *gin.Context) {
	// 声明一个接收参数的struct
	/**
		定义一个struct字段的标签, 通过json标签, 定义请求参数和form表单的字段关系
	 */
	type User struct {
		Name string `json:"name" form:"name"`
		Age int8 `json:"age" form:"age"`
	}
	
	// 初始化user struct
	user := User{}
	if c.ShouldBind(&user) == nil {
		log.Printf("log %v\n\n", user)
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data": user,
		})
	}

}


func LoginAction(c *gin.Context) {

	// 接收post过来的参数
	username := c.PostForm("username")
	password := c.PostForm("password")

	hash := md5.New()
	md5Pass, _ := hash.Write([]byte("111111"))
	pass, _ := hash.Write([]byte(password))


	if username == "llf" && md5Pass == pass {
		c.JSON(
				http.StatusOK,
				gin.H{
					"status": 200,
					"message": "login success",
				},
			)
	} else {
		c.JSON(
			http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				//"data":    ,
				"message": "账号或密码错误",
			},
		)
	}
}

func SayHelloAction(c *gin.Context) {
	name := c.Param("name")
	c.JSON(
		http.StatusOK, gin.H{
			"status": 200,
			"message": "hello " + name,
		},
	)
}
