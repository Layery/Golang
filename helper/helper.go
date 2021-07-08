package helper

import "fmt"
import "github.com/astaxie/beego"

func LogWriter(data interface{})  {
	fmt.Printf("%v", data)
}


func main(){
	beego.Run()
}
