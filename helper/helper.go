package helper

import (
	"fmt"

	"github.com/astaxie/beego"
)

func LogWriter(data interface{}) {
	fmt.Printf("%v", data)
}

func main() {
	beego.Run()
}
