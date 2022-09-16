package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

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

































