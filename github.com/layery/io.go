package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
)

func DD(data ...interface{}) { // 声明可变参数函数时，需要在参数列表的最后一个参数类型之前加上省略符号"..."，这表示该函数会接收任意数量的该类型参数。
	d, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	var out bytes.Buffer
	err = json.Indent(&out, d, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out.String())
}

// 从标准输入读入数据, 并直接写入到文件中
func WriteFileFromBufio() {
	fd, err := os.OpenFile("./io.log", os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
	}

	defer fd.Close()

	reader := bufio.NewReader(os.Stdin)

	p := make([]byte, 3)
	for {
		n, err := reader.Read(p)
		if err != nil {
			if err == io.EOF {
				log.Fatal("read over")
				break
			}
			log.Fatal(err)
		}

		// 将读到的内容, 实时写入文件中去
		if n > 0 {
			_, writeErr := fd.WriteString(string(p[0:n]))
			if writeErr != nil {
				log.Fatal(writeErr)
			}

		}
		fmt.Println(string(p[0:n]))
	}
}

func readDir(path string) {
	rs, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for key, value := range rs {
		_ = key
		if value.IsDir() {
			readDir(path + string(os.PathSeparator) + value.Name())
		}

		log.Println(value.Name())
	}

}

func MyPostForm() {
	api := "http://debug.cc/debug.php"
	client := http.Client{}
	resp, err := client.PostForm(api, url.Values{"name": []string{
		"wdy",
	}})
	if err != nil {
		panic(errors.New("body"))
	}

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func MyPostWithJsonData() {
	api := "http://micro.cc:880/site/test"

	// 只要是想自定义头, 就得使用NewRequest方法
	urlModel, _ := url.Parse(api)

	// 声明一个client, 给client 设置上cookiejar, 他会自己携带上cookie
	client := http.Client{}
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	client.Jar = jar

	// 采用application/json格式的数据, 发送post数据
	type postData struct {
		Name string `json:"name"`
		Age  int8   `json:"age"`
	}
	jsonPostData := postData{
		"llf",
		20,
	}
	post, err := json.Marshal(jsonPostData)

	req, err := http.NewRequest(http.MethodPost, urlModel.String(), bytes.NewReader(post))
	if err != nil {
		panic(err)
	}

	// 构造自定义的header
	req.Header.Set("content-type", "application/json charset=utf8")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36")
	req.Header.Set("Referer", "https://www.baidu.com/")
	req.Header.Set("Host", "https://sohu.cn")

	// 发送请求
	resp, _ := client.Do(req) // 第一次post请求
	defer resp.Body.Close()

	var response struct {
		Code int                    `json:"code"`
		Msg  string                 `json:"msg"`
		Data map[string]interface{} `json:"data"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&response)

	// 检测map中的key是否存在
	if number, ok := response.Data["test_num"]; !ok {
		fmt.Println(number)
	} else {
		// utils.P(number)
	}

	//
	//
	//resp, _ = client.Do(req) // 第二次post请求
	//defer resp.Body.Close()
	//
	//body, _ = ioutil.ReadAll(resp.Body)
	//
	//utils.P(string(body))

}

func LoginLocal() {
	defer func() {
		catch := recover()
		if catch != nil {
			log.Fatal(catch)
		}
	}()

	MyPostForm() // 发送post请求, body参数格式化为url_encode类型

	MyPostWithJsonData() // 发送post数据, 采用application/json 的头类型
}

type MyReader struct {
	// 资源
	reader io.Reader
	// 当前读取到的位置
	cur int
}

// 创建一个实例
func newMyReader(reader io.Reader) *MyReader {
	return &MyReader{reader: reader}
}

func (m *MyReader) Read(p []byte) (int, error) {
	num, err := m.reader.Read(p)

	// 给每个字节加个括号包裹起来
	buf := make([][]byte, num) // 声明一个等同于当前读到的长度的切片

	for i := 0; i < num; i++ {
		charStr := string(p[i])
		charStr = "[" + charStr + "]"
		ss := []byte(charStr)
		buf[i] = ss

		log.Println(string(buf[i]))

	}

	return num, err
}

func ioReaderStudy() {
	str := "abcdefghjkmnopq"
	p := make([]byte, 4)

	// 使用一个 Reader 作为另一个 Reader 的实现是一种常见的用法.
	// 这样做可以让一个 Reader 重用另一个 Reader 的逻辑,
	// 下面展示通过更新 alphaReader 以接受 io.Reader 作为其来源.

	reader_str := strings.NewReader(str)

	reader_byte := bytes.NewReader([]byte{'A', 'B', 'C', 'D', 'E', 'F', 'G'})

	reader_file, _ := os.Open("./config.toml")

	response, _ := http.Get("https://www.baidu.com")
	defer response.Body.Close()

	reader_net := response.Body

	reader_stdin := os.Stdin

	_ = reader_str
	_ = reader_byte
	_ = reader_file
	_ = reader_net

	myreader := newMyReader(reader_stdin)

	for {
		n, err := myreader.Read(p)
		if err != nil {
			if err == io.EOF {
				fmt.Println("读完了")
				break
			}
			fmt.Println("myreader 读取错误")
			os.Exit(0)
		}

		// 将读取到的长度

		fmt.Println("读取到的长度", string(p[:n]))
	}
}

func futuMianShiTi() {

	type T rune

	type myT T
	var (
		t  T
		p  *T
		i1 interface{} = t
		i2 interface{} = p
	)

	//var t1 T

	//var t2 myT

	//fmt.Printf("%p ==> %#v \n", &t1, t1)
	//fmt.Printf("%p ==> %#v \n", &t2, t2)

	fmt.Printf("%#v === %#v\n", i1 == t, i1 == nil)
	fmt.Printf("%#v === %#v\n", i2 == p, i2 == nil)
	fmt.Printf("\n%p ==> %v\n %p ==> %v\n %p ==> %v\n %p ==> %v\n", &t, t, &p, p, &i1, i1, &i2, i2)
}

func ioReaderBufioStudy() {
	res, err := os.Open("./config.toml")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Close()

	var p = make([]byte, 100)

	reader := bufio.NewScanner(res)

	_ = reader.Scan()

	_ = p

	for reader.Scan() {

		log.Println(reader.Text())
	}

}

func Debug() {
	ioReaderStudy()

	ioReaderBufioStudy()
}

func main() {

	go WriteFileFromBufio()

	go readDir(".")

	select {}

}
