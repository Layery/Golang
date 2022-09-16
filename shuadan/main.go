package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/fsutil"

	//"github.com/gookit/goutil/dump"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"shuadan/utils"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gookit/goutil/strutil"
	"golang.org/x/net/publicsuffix"
)

const (
	USERNAME_C string = "gejinchenxiao"
	PASSWD_C   string = "555xiaomo"

	USERNAME_L string = "weidingyi"
	PASSWD_L   string = "LlF213344"

	UN_NOTICE_PRICE   float64 = 400
	PUSH_PLUS_TOKEN_L         = "13b306b602814b2a8f8798874a1ba49f"
	PUSH_PLUS_TOKEN_C         = "5029ee32cfcf4f409b19489237be3932"
	API_SECRET                = "6FadPxCXJLkxsVaDwkwGhyKNgNJfwwn14qNI872Y46A"
	RN                        = "\n"
)

var noticeType = map[int]string{
	1: "又刷到一个新单子, 请注意查看",
	2: "",
	3: "重新登录失败",
}

func getNoticeMsg(msgType int, noticeTitle string) string {
	if msg, ok := noticeType[msgType]; ok {
		if msgType == 1 {
			return noticeTitle
		} else {
			return msg
		}
	}
	return "未知的消息类型"
}

func getSign(time int64) string {
	timeStr := strconv.FormatInt(time, 10)
	return utils.GetMd5(timeStr + API_SECRET)
}

var ApiListMap = map[string]string{
	"login":         "http://www.haisirui.xin/index/Apprentice/getlogin.html",
	"notice_c":      "http://pushplus.hxtrip.com/send?" + PUSH_PLUS_TOKEN_C,
	"notice_l":      "http://pushplus.hxtrip.com/send?" + PUSH_PLUS_TOKEN_L,
	"notice_wechat": "https://app.lfeet.top/index.php?s=api/setting/wmsg",
}

var JobApiMap = map[string]string{
	"jd_job":  "http://www.haisirui.xin/public/index.php/index/jdapprentice/jd_obtain", // t=0.423290446257371
	"pick_up": "http://www.haisirui.xin/public/index.php/index/Apprentice/pickup_task",
	"tb_job":  "http://www.haisirui.xin/public/index.php/index/Apprentice/tbkj_obtain",
	"obtain":  "http://www.haisirui.xin/public/index.php/index/jsapprentice/obtain_js_task",
	"browser": "http://www.haisirui.xin/public/index.php/index/Browse/browse_obtain",
	"pdd_job": "http://www.haisirui.xin/public/index.php/index/pddapprentice/pdd_obtain",
}

type Data struct {
	UserCode  string `json:"user_code"`
	UserOid   string `json:"user_oid"`
	UserMid   string `json:"user_mid"`
	UserToken string `json:"user_token"`
}

type Person struct {
	Data
	UserName    string `json:"username"`
	Password    string `json:"password"`
	LoginStatus bool
	jar         interface{}
}

type LoginResult struct {
	Msg           string `json:"msg"`
	Status        int    `json:"status"`
	LoginSign     int    `json:"login_sign"`
	ApprenticeID  int    `json:"apprenticeId"`
	IdentitySign  int    `json:"identity_sign"`
	ReceivingSign int    `json:"receiving_sign"`
	UserSjcode    string `json:"user_sjcode"`
}

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

func (p *Person) AddRequestHeader(r *http.Request) *http.Request {

	// 构造自定义的头
	r.Header.Set("content-type", "application/x-www-form-urlencoded charset=utf8")
	r.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36")
	r.Header.Set("Referer", "https://www.baidu.com/")

	return r
}

func (p *Person) WriteCookie(cookies []*http.Cookie) error {
	jsonData, _ := json.Marshal(cookies)
	// 这里和php不同, 必须是0777的权限, 不得省略前置0
	_ = ioutil.WriteFile(p.GetDataPath()+"cookies", jsonData, 0777)
	return nil
}

// 自动识别当前程序运行路径
func (p *Person) GetDataPath() string {

	// 判断exe文件是否运行在temp目录下, 来判断是run, 还是build
	tempPath := os.TempDir()
	// 程序go run 运行 与 go build 运行时, 如何识别项目路径 ??
	file, _ := exec.LookPath(os.Args[0])
	execPath, _ := filepath.Abs(file)

	isCheck := strings.Contains(execPath, tempPath) // true: go run , false: go build

	var dataPath string

	if isCheck {
		_, file1, _, ok := runtime.Caller(1)
		if !ok {
			log.Fatal("system err!")
		}
		mainGoPath := filepath.Dir(file1)

		personDataPath := mainGoPath + string(os.PathSeparator) + "data" + string(os.PathSeparator) + utils.GetMd5(p.UserName)
		if !fsutil.IsDir(personDataPath) {
			_ = os.MkdirAll(personDataPath, 0777)
		}
		dataPath, _ = filepath.Abs(personDataPath)
		dataPath = dataPath + string(os.PathSeparator)
	} else {
		goos := runtime.GOOS
		binFileName := file
		if strings.Contains(goos, "window") { // windows
			binFileName = strutil.Substr(file, 0, -4)
		}
		binFilePath := execPath + "/../" + binFileName + "-data"
		dataPath, _ = filepath.Abs(binFilePath)
		dataPath = dataPath + string(os.PathSeparator) + utils.GetMd5(p.UserName) + string(os.PathSeparator)
		_ = os.MkdirAll(dataPath, 0777)
	}
	return dataPath
}

func (p *Person) Login() *Person {

	var jar, _ = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List}) // 为了根据域名, 安全的设置cookie

	p.BeforeLogin("http://www.haisirui.xin/index/apprentice")
	postParams := url.Values{
		"userName":   []string{p.UserName},
		"password":   []string{p.Password},
		"user_os":    []string{"android"},
		"user_model": []string{"chrome"},
		"user_code":  []string{p.Data.UserCode},
		"ip":         []string{""},
		"ip_data":    []string{""},
	}

	request, _ := http.NewRequest(http.MethodPost, ApiListMap["login"], strings.NewReader(postParams.Encode()))
	//proxy, _ := url.Parse("http://127.0.0.1:8888")
	client := http.Client{
		//Transport: &http.Transport{
		//	Proxy:http.ProxyURL(proxy),
		//},
	}

	client.Jar = jar

	// 拼接自定义的头
	request, _ = BuildHeader(request)
	request.Header.Set("X-Requested-With", "XMLHttpRequest")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	resp, _ := client.Do(request)
	defer resp.Body.Close()

	var loginResult LoginResult
	_ = json.NewDecoder(resp.Body).Decode(&loginResult)

	if loginResult.Status == 200 {
		p.LoginStatus = true

		// 保存cookie到本地
		_ = p.WriteCookie(resp.Cookies())

	} else {
		loginResult.Status = 9999
		msg, _ := json.Marshal(loginResult)
		panic(errors.New(string(msg)))
	}
	return p
}

func (p *Person) IsLogin() *Person {
	// 检查本地是否有cookie, 有cookie暂不登录
	localCookie, _ := ioutil.ReadFile(p.GetDataPath() + "cookies")
	if len(localCookie) > 0 {
		p.LoginStatus = true
		return p
	} else {
		p.LoginStatus = false
		return p.Login()
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

func (p *Person) BeforeLogin(url string) *Person {
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	// 加载html doc
	doc, _ := goquery.NewDocumentFromReader(res.Body)
	userCode, ok := doc.Find(".user_code").Attr("value")
	if ok {
		p.Data.UserCode = userCode
	}
	oid, ok := doc.Find(".user_oid").Attr("value")
	if ok {
		p.Data.UserOid = oid
	}
	mid, ok := doc.Find(".user_mid").Attr("value")
	if ok {
		p.Data.UserMid = mid
	}
	userToken, ok := doc.Find(".user_token").Attr("value")
	if ok {
		p.Data.UserToken = userToken
	}
	return p
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

func (p *Person) PanicReLogin() {
	p.LoginStatus = false
	err := LoginResult{
		Status: 302,
	}
	errMsg, _ := json.Marshal(err)
	panic(string(errMsg))
}

func (p *Person) runJob() {
	// 构造请求的cookie
	host := "http://www.haisirui.xin"
	urlModel, _ := url.Parse(host)
	cookies, _ := ioutil.ReadFile(p.GetDataPath() + "cookies")
	var cookieArray []*http.Cookie
	_ = json.NewDecoder(strings.NewReader(string(cookies))).Decode(&cookieArray)

	jar, _ := cookiejar.New(nil)
	jar.SetCookies(urlModel, cookieArray)

	for jobName, api := range JobApiMap {

		request, _ := http.NewRequest(http.MethodGet, api, nil)

		apiRootPath := request.URL.Scheme + "://" + request.URL.Host + "/public/index.php"

		//proxy, _ := url.Parse("http://127.0.0.1:8888")
		client := http.Client{
			//Transport: &http.Transport{
			//	Proxy:http.ProxyURL(proxy),
			//},
			CheckRedirect: func(req *http.Request, via []*http.Request) error { // 这里防止自动302, 返回一个error, 便会只请求一次
				return errors.New("first request end")
			},
		}
		client.Jar = jar
		// 拼接自定义的头
		request, _ = BuildHeader(request)
		request.Header.Set("X-Requested-With", "XMLHttpRequest")
		request.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
		request.Header.Set("Referer", "http://www.haisirui.xin/public/index.php/index/Apprentice/receive_task")

		request, _ = BuildQueryParams(request, map[string]string{
			"t": utils.GetRand("0."),
		})

		resp, _ := client.Do(request)
		defer resp.Body.Close()

		if resp.StatusCode == 302 {
			p.PanicReLogin()
		}

		body, _ := ioutil.ReadAll(resp.Body)

		if resp.StatusCode == 200 {
			jsonData := make(map[string]interface{})
			_ = json.Unmarshal(body, &jsonData)
			if status, ok := jsonData["status"]; ok {
				if status == float64(200) {
					// 找到任务, 抓取任务页, 保存
					taskUrl := apiRootPath + jsonData["task"].(string)
					request, _ = http.NewRequest(http.MethodGet, taskUrl, nil)
					resp, _ := client.Do(request)
					defer resp.Body.Close()
					body, _ := ioutil.ReadAll(resp.Body)

					date := time.Now().Format("2006-01-02 15-04")
					date += ".html"
					logName := p.GetDataPath() + jobName + "-" + date
					_ = ioutil.WriteFile(logName, body, 0777)

					priceFloat, noticeTitle, content := p.GenerateNoticeContent(jobName, logName)
					p.NoticeMe(1, taskUrl, noticeTitle, content, priceFloat)
				}
			}
		}
	}
	fmt.Println()
}

func (p *Person) DoJob() *Person {

	for {
		p.IsLogin()

		p.runJob()

		time.Sleep(time.Minute)
	}
}

// 程序无法自动登录, 发送通知后, die掉
func (p *Person) NoticeMe(msgType int, jumpUrl, noticeTitle, content string, currPrice float64) bool {

	if currPrice >= UN_NOTICE_PRICE {
		return false
	}

	currTime := time.Now().Unix()

	postMap := map[string]interface{}{
		"time":     currTime,
		"sign":     getSign(currTime), // 计算调用接口的签名
		"title":    getNoticeMsg(msgType, noticeTitle),
		"jump_url": jumpUrl,
		"content":  content,
		"who":      p.UserName,
	}

	postParams, _ := json.Marshal(postMap)

	postData := bytes.NewReader(postParams)

	response, _ := http.Post(ApiListMap["notice_wechat"], "application/json", postData)
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	log.Println("poster_succ", string(body))

	// 发送完通知后, 程序睡眠10分钟??
	time.Sleep(time.Minute * 10)
	return true
}

func (p *Person) DoJobNew(channel chan Person, group sync.WaitGroup) {

	fmt.Println(strutil.Substr(p.UserName, 0, 1) + " begin do job")

	channel <- *p
	group.Done()
}

func newPerson(userName, password string) *Person {
	return &Person{
		UserName:    userName,
		Password:    password,
		LoginStatus: false,
	}
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

func (p *Person) GenerateNoticeContent(jobType, logName string) (float64, string, string) {
	titleMap := map[string]string{
		"jd_job":  "京东任务",
		"pick_up": "快捷任务",
		"tb_job":  "淘宝任务",
		"obtain":  "obtain",
		"browser": "浏览器任务",
		"pdd_job": "拼多多",
	}

	body, _ := ioutil.ReadFile(logName)

	// 加载html doc
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(string(body)))

	price := doc.Find(".customer_order").Text()

	keyword, _ := doc.Find(".key_word_hidden").Attr("value")

	// 价格格式化为float64
	f := strings.NewReplacer(",", "")
	priceString := f.Replace(price)
	priceFloat, _ := strconv.ParseFloat(priceString, 10)

	// 如果log没内容, 不发通知
	lens := len(body)
	if lens <= 0 {
		priceFloat = float64(999)
	}

	time.Sleep(time.Second)

	_ = os.Remove(logName)

	noticeTitle := "【"+ titleMap[jobType]+"】【" + price + "】【"+ keyword +"】"

	return priceFloat, noticeTitle, "类别: " + titleMap[jobType] + RN + "价格: " + price + RN + "关键词: " + keyword
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
