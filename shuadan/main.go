package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gookit/goutil/strutil"
	"golang.org/x/net/publicsuffix"
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
	"strings"
	"time"
)

const (
	USERNAME  string = "gejinchenxiao"
	PASSWD    string = "555xiaomo"
	PUSH_PLUS_TOKEN_L = "13b306b602814b2a8f8798874a1ba49f"
	PUSH_PLUS_TOKEN_C = "5029ee32cfcf4f409b19489237be3932"
	RN = "\r\n"
)

var jar, _ = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List}) // 为了根据域名, 安全的设置cookie

var ApiListMap = map[string]string{
	"login": "http://www.haisirui.xin/index/Apprentice/getlogin.html",
	"notice_c": "http://pushplus.hxtrip.com/send?"+ PUSH_PLUS_TOKEN_C,
	"notice_l": "http://pushplus.hxtrip.com/send?"+ PUSH_PLUS_TOKEN_L,
}

var JobApiMap = map[string]string{
	"jd_job": "http://www.haisirui.xin/public/index.php/index/jdapprentice/jd_obtain", // t=0.423290446257371
	"pick_up": "http://www.haisirui.xin/public/index.php/index/Apprentice/pickup_task",
	"tb_job": "http://www.haisirui.xin/public/index.php/index/Apprentice/tbkj_obtain",
	"obtain": "http://www.haisirui.xin/public/index.php/index/jsapprentice/obtain_js_task",
	"browser": "http://www.haisirui.xin/public/index.php/index/Browse/browse_obtain",
	"pdd_job": "http://www.haisirui.xin/public/index.php/index/pddapprentice/pdd_obtain",
}

type Data struct {
	UserCode string `json:"user_code"`
	UserOid string `json:"user_oid"`
	UserMid string `json:"user_mid"`
	UserToken string `json:"user_token"`
	//LoginResult LoginResult `json:"login_result"`
}

type Person struct {
	client *http.Client
	Data
	UserName    string `json:"username"`
	Password    string `json:"password"`
	LoginStatus bool
}

type LoginResult struct {
	Msg string `json:"msg"`
	Status int `json:"status"`
	LoginSign int `json:"login_sign"`
	ApprenticeID int `json:"apprenticeId"`
	IdentitySign int `json:"identity_sign"`
	ReceivingSign int `json:"receiving_sign"`
	UserSjcode string `json:"user_sjcode"`
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
		"Host":             "www.haisirui.xin",
		"Upgrade-Insecure-Requests": "1",
		"User-Agent":       "Mozilla/5.0 (Linux; Android 9; MI 8 Build/PKQ1.180729.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/72.0.3626.121 Mobile Safari/537.36",
		"Accept":           "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		"x-requested-with": "XBrowser",
		"Accept-Encoding":  "gzip, deflate",
		"Accept-Language":  "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7",
		"Connection":       "keep-alive",
		"Origin":           "http://www.haisirui.xin",
		"Referer":          "http://www.haisirui.xin/index/apprentice",
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
	return req, nil
}

func (p Person) AddRequestHeader(r *http.Request) *http.Request {

	// 构造自定义的头
	r.Header.Set("content-type", "application/x-www-form-urlencoded charset=utf8")
	r.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36")
	r.Header.Set("Referer", "https://www.baidu.com/")

	return r
}

func (p Person) WriteCookie(cookies []*http.Cookie) error {
	jsonData, _ := json.Marshal(cookies)
	// 这里和php不同, 必须是0777的权限, 不得省略前置0
	_ = ioutil.WriteFile(p.GetDataPath()+"cookies", jsonData, 0777)
	return nil
}

// 自动识别当前程序运行路径
func (p Person) GetDataPath() string {

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
		_ = os.Mkdir(mainGoPath + "/data", 0777)
		dataPath, _ = filepath.Abs(mainGoPath + "/data/")
		dataPath = dataPath + string(os.PathSeparator)
	} else {
		goos := runtime.GOOS
		binFileName := file
		if strings.Contains(goos, "window") { // windows
		    binFileName = strutil.Substr(file, 0, -4)
		}
		binFilePath := execPath + "/../" + binFileName + "-data"
		dataPath, _ = filepath.Abs(binFilePath)
		dataPath = dataPath + string(os.PathSeparator)
		_ = os.Mkdir(dataPath, 0777)
	}
	return dataPath
}

func (p *Person) Login() *Person {

	p.BeforeLogin("http://www.haisirui.xin/index/apprentice")

	postParams := url.Values{
		"userName":   []string{USERNAME},
		"password":   []string{PASSWD},
		"user_os":    []string{"android"},
		"user_model": []string{"chrome"},
		"user_code":  []string{p.Data.UserCode},
		"ip": 		  []string{""},
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
		_ = p.WriteCookie(jar.Cookies(resp.Request.URL))

		log.Println("登录成功")
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

func (p *Person)doJingDongJob() {
	log.Println("begin do job")
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
			"t" : utils.GetRand("0."),
		})

		resp, _ := client.Do(request)
		defer resp.Body.Close()

		if resp.StatusCode == 302 {
			p.PanicReLogin()
		}

		body, _ := ioutil.ReadAll(resp.Body)
		log.Println(string(body))
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

					//todo 获取该任务的标题, 价格, 佣金 等关键字段,
					p.NoticeMe("又刷到一个单子, 请登录查看", p.GenerateNoticeContent(jobName, logName))
				}
			}
		}
	}
	fmt.Println()
}

func (p *Person)DoJob() *Person{
	for {
		// 获取当前时间, 22:00 -- 08:00阶段, 按照小时来睡

		p.IsLogin()

		// 构造请求的cookie
		host := "http://www.haisirui.xin"
		urlModel, _ := url.Parse(host)
		cookies, _:= ioutil.ReadFile(p.GetDataPath() + "cookies")
		var cookieArray []*http.Cookie
		_ = json.NewDecoder(strings.NewReader(string(cookies))).Decode(&cookieArray)
		jar.SetCookies(urlModel, cookieArray)


		p.doJingDongJob()

		time.Sleep(time.Minute)
	}
}

// 程序无法自动登录, 发送通知后, die掉
func (p *Person)NoticeMe(msg, content string) {
	postMap := map[string]interface{}{
		"token": PUSH_PLUS_TOKEN_L,
		"title": msg,
		"content": content,
	}
	postParams, _ := json.Marshal(postMap)

	postData := bytes.NewReader(postParams)

	response, err_me := http.Post(ApiListMap["notice_l"], "application/json", postData)
	defer response.Body.Close()
	if err_me != nil {
		log.Println("poster_err", err_me)
	}

	time.Sleep(time.Second * 2)

	postMap = map[string]interface{}{
		"token": PUSH_PLUS_TOKEN_C,
		"title": msg,
		"content": content,
	}
	postParams, _ = json.Marshal(postMap)
	postData = bytes.NewReader(postParams)
	response, _ = http.Post(ApiListMap["notice_c"], "application/json", postData)
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	log.Println("poster_succ", string(body))


	// 发送完通知后, 程序睡眠10分钟??
	time.Sleep(time.Minute * 10)
}

func IsCouldRun() bool{
	//runBeginHour := 8
	//runEndHour := 22

	now := time.Now()
	//if now.Hour() >= runBeginHour && now.Hour() < runEndHour {
	//	return true
	//} else {
	//	return false
	//}

	if now.Minute() >= 15 && now.Minute() < 16 {
		return true
	}

	if now.Minute() >= 17 {
		return true
	}
	return false
}

func (p *Person)Debug ()  {
	// http://micro.cc/site/debug

	resp, _ := http.Get("http://micro.cc/site/debug")

	date := time.Now().Format("2006-01-02 15-04")
	date += ".html"
	logName := p.GetDataPath() + "debug" + "-" + date
	body, _ := ioutil.ReadAll(resp.Body)
	_ = ioutil.WriteFile(logName, body, 0777)

	rs := p.GenerateNoticeContent("browser", logName)
	p.NoticeMe("又刷到一个单子, 请登录查看", rs)

}


func (p *Person) GenerateNoticeContent(jobType, logName string) string {
	titleMap := map[string]string{
		"jd_job": "京东任务",
		"pick_up": "快捷任务",
		"tb_job": "淘宝任务",
		"obtain": "obtain",
		"browser": "浏览器任务",
		"pdd_job": "拼多多",
	}

	body, _ := ioutil.ReadFile(logName)

	// 加载html doc
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(string(body)))

	price := doc.Find(".customer_order").Text()

	keyword, _ := doc.Find(".key_word_hidden").Attr("value")

	return "类别: " + titleMap[jobType] + RN + "价格: " + price + RN + "关键词: " + keyword
}


func main() {

	defer func() {
		catch := recover()
		if catch != nil {
			var err LoginResult
			var person Person
			s := catch.(string)
			_ = json.NewDecoder(strings.NewReader(s)).Decode(&err)
			if err.Status == 302 { // 5kqljb892qqlcmgga0cbuilcc7
				log.Println("登录已失效, 重新登录")

				_ = ioutil.WriteFile(person.GetDataPath() + "cookies", []byte{}, 0777)
				person.DoJob()
			}

			if err.Status == 9999 { // 登录失败, 通知我
				log.Println("登录失败: ",catch)
				person.NoticeMe("重新登录失败", "")
			}
		}
	}()


	var p Person

	_ = p.DoJob()

	//p.Debug()


	//go func() {
	//	person.Debug()
	//}()
	//
	//for {
	//	runtime.Gosched()
	//}
}
