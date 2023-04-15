package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gookit/goutil/dump"

	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"shuadan/utils"
	"strings"
	"sync"
	"time"
	"github.com/gookit/goutil/strutil"
)




func (p *Person) AddRequestHeader(r *http.Request) *http.Request {

	// 构造自定义的头
	r.Header.Set("content-type", "application/x-www-form-urlencoded charset=utf8")
	r.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36")
	r.Header.Set("Referer", "https://www.baidu.com/")

	return r
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
		utils.P(number)
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

func (p Person) LoginLocal() {
	defer func() {
		catch := recover()
		if catch != nil {
			log.Fatal(catch)
		}
	}()

	MyPostForm() // 发送post请求, body参数格式化为url_encode类型

	MyPostWithJsonData() // 发送post数据, 采用application/json 的头类型

	//MyPostWithCookie() // 发送post请求, 携带cookie

	//MyHeadRequest()
}



func (p *Person) DoJobNew(channel chan Person, group sync.WaitGroup) {

	fmt.Println(strutil.Substr(p.UserName, 0, 1) + " begin do job")

	channel <- *p
	group.Done()
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

		dump.P(string(buf[i]))

	}

	return num, err
}

func testChan() {

	rs := os.Getpid()

	fmt.Println("testChan ppid:", os.Getppid())
	fmt.Println("testChan pid:", rs)

	go func() {
		fmt.Println("go func ppid:", os.Getppid())
		fmt.Println("go func pid:", os.Getpid())

		go func() {
			fmt.Println("go func ppid:", os.Getppid())
			fmt.Println("go func pid:", os.Getpid())
		}()
	}()

}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func newUser(u, p string) *User {

	return &User{
		Username: u,
		Password: p,
	}
}

func (u *User) Login() error {
	postParams := url.Values{
		"username": []string{u.Username},
		"password": []string{u.Password},
	}
	req, _ := http.NewRequest(http.MethodPost, "http://yixinyuan.cc/a.php/index/login", strings.NewReader(postParams.Encode()))
	req.Header.Add("content-type", "application/json")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	//proxy, _ := url.Parse("http://127.0.0.1:8888")
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return errors.New("error")
		},
		//Transport: &http.Transport{
		//	Proxy:http.ProxyURL(proxy),
		//},
	}

	resp, err := client.Do(req)

	defer resp.Body.Close()

	if err != nil {
		panic(err)
	}

	dataPath := "E:/www/haisirui/" + u.Username + string(os.PathSeparator)

	_ = os.MkdirAll(dataPath, 0777)

	data, _ := json.Marshal(resp.Cookies())
	_ = ioutil.WriteFile(dataPath+"cookie", data, 0777)

	return nil
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
	//ioReaderStudy()

	ioReaderBufioStudy()
}

func main() {

	go func() {
		defer func() {
			catch := recover()
			if catch != nil {
				var err LoginResult
				person := newPerson(USERNAME_C, PASSWD_C)

				s := catch.(string)
				_ = json.NewDecoder(strings.NewReader(s)).Decode(&err)
				if err.Status == 302 { // 5kqljb892qqlcmgga0cbuilcc7
					_ = ioutil.WriteFile(person.GetDataPath()+"cookies", nil, 0777)
					person.DoJob()
				}

				if err.Status == 9999 { // 登录失败, 通知我
					person.NoticeMe(3, "https://www.baidu.com", "", "", 0)
				}
			}
		}()

		p := newPerson(USERNAME_C, PASSWD_C)

		_ = p.DoJob()

	}()

	time.Sleep(time.Second * 3)

	go func() {
		defer func() {
			catch := recover()
			if catch != nil {
				var err LoginResult
				person := newPerson(USERNAME_L, PASSWD_L)

				s := catch.(string)

				_ = json.NewDecoder(strings.NewReader(s)).Decode(&err)
				if err.Status == 302 { // 5kqljb892qqlcmgga0cbuilcc7
					_ = ioutil.WriteFile(person.GetDataPath()+"cookies", nil, 0777)
					person.DoJob()
				}

				if err.Status == 9999 { // 登录失败, 通知我
					person.NoticeMe(3, "https://www.baidu.com", "", "", 0)
				}
			}
		}()

		p := newPerson(USERNAME_L, PASSWD_L)

		_ = p.DoJob()

	}()

	select {}
}
