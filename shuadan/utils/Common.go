package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func BuildQueryParams(req *http.Request, params map[string]string) (*http.Request, error) {
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	return req, nil
}

func BuildHeader(req *http.Request) (*http.Request, error) {
	header := map[string]string{
		"Host":                      "www.haisirui.xin",
		"Upgrade-Insecure-Requests": "1",
		"User-Agent":                "Mozilla/5.0 (Linux; Android 9; MI 8 Build/PKQ1.180729.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/72.0.3626.121 Mobile Safari/537.36",
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		"x-requested-with":          "XBrowser",
		"Accept-Encoding":           "gzip, deflate",
		"Accept-Language":           "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7",
		"Connection":                "keep-alive",
		"Origin":                    "http://www.haisirui.xin",
		"Referer":                   "http://www.haisirui.xin/index/apprentice",
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
	return req, nil
}


func GetMd5(str string) string {
	data := []byte(str)
	hash := md5.Sum(data)
	return fmt.Sprintf("%x", hash)
}

func GetRand(prefix string) string {
	res := fmt.Sprintf("%016v", rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(10000000000000000))
	if prefix != "" {
		res = prefix + res
	}
	return res
}
func Date(format, timestamp string) string {
	temp, _ := strconv.ParseInt(timestamp, 10, 64)
	timeInt64 := time.Unix(temp, 0)
	r := strings.NewReplacer(
		"Y", "2006",
		"m", "01",
		"d", "02",
		"H", "15",
		"i", "04",
		"s", "05",
	)
	format = r.Replace(format)
	return timeInt64.Format(format)
}

func P(data interface{}) {
	byteList, _ := json.Marshal(data)
	var out bytes.Buffer
	err := json.Indent(&out, byteList, "", "\t")
	if err != nil {
		return
	}
	fmt.Println(out.String())
	//
	//fmt.Printf("%#v,         type: %T \n\n", data, data)
	fmt.Println()
	fmt.Println()
	os.Exit(0)
}

func T(data interface{}) {
	fmt.Printf("%#v,         type: %T \n\n", data, data)
	fmt.Println()
	fmt.Println()
}

































