package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"skill-order/server/dao"
)

func skillOrder() {
	api := "http://localhost:9000/v2/skill"

	resp, err := http.Get(api)
	defer resp.Body.Close()

	if err != nil {
		log.Fatal("当前请求失败, err: ", err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("%#v \n\n", string(body))
}

func testSkillOrder(num int) {
	fmt.Printf("开启%d个协程 \n", num)
	for i := 0; i < num; i++ {
		go func() {
			skillOrder()
		}()
	}
	select {}
}

func main() {
	// 并发200个请求, 下单同一个商品
	act := flag.String("act", "eg: drop:初始化数据", "执行的动作")
	num := flag.Int("num", 300, "设置启动协程的数量")

	flag.Parse()

	switch *act {
	case "drop":
		nt := make(chan string)
		go func(ch chan string) {
			dao.InitMysql()
			dao.InitGoodsData()
			ch <- "数据初始化完毕"
		}(nt)
		fmt.Println(<-nt)
		return
	default:
		testSkillOrder(*num)
		return
	}

}
