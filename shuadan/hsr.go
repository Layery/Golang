package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gookit/goutil/dump"

	//"github.com/gookit/goutil/dump"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"shuadan/utils"
	"strconv"
	"strings"
	"time"
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

var dataDir string


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

const (
	ERR_CODE_LOGIN_FALSE int = 10000 + iota
	ERR_CODE_COOKIE_INVALID
	ERR_CODE_UNKNOW
)



type ErrorResp struct {
	*Person
	Code int
	Msg string
}


// 自动识别当前程序运行路径
func (p *Person) GetDataPath() string {
	distPath := dataDir + utils.GetMd5(p.UserName) + string(os.PathSeparator)
	_ = os.MkdirAll(distPath, 0777)
	return distPath
}

func (p *Person) IsLogin() *Person {
	// 检查本地是否有cookie, 有cookie暂不登录
	localCookie, _ := ioutil.ReadFile(p.GetDataPath() + "cookies")
	if len(localCookie) > 0 {
		log.Println("存在cookie直接登录")
		p.LoginStatus = true
		return p
	} else {
		log.Println("不存在cookie")
		p.LoginStatus = false
		return p.Login()
	}
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
		//CheckRedirect: func(req *http.Request, via []*http.Request) error { // 这里防止自动302, 返回一个error, 便会只请求一次
		//	return errors.New("first request end")
		//},
		//Transport: &http.Transport{
		//	Proxy:http.ProxyURL(proxy),
		//},
	}

	client.Jar = jar

	// 拼接自定义的头
	request, _ = utils.BuildHeader(request)
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
		panic(ErrorResp{
			Person: p,
			Code: ERR_CODE_LOGIN_FALSE,
			Msg: "登录失败",
		})
	}
	return p
}

func (p *Person) WriteCookie(cookies []*http.Cookie) error {
	jsonData, _ := json.Marshal(cookies)
	// 这里和php不同, 必须是0777的权限, 不得省略前置0
	err := ioutil.WriteFile(p.GetDataPath()+"cookies", jsonData, 0777)
	if err != nil {
		panic(err)
	}
	return err
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
		request, _ = utils.BuildHeader(request)
		request.Header.Set("X-Requested-With", "XMLHttpRequest")
		request.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
		request.Header.Set("Referer", "http://www.haisirui.xin/public/index.php/index/Apprentice/receive_task")

		request, _ = utils.BuildQueryParams(request, map[string]string{
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
				log.Println(jsonData)
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

func (p *Person) PanicReLogin() {
	p.LoginStatus = false
	_ = os.Remove(p.GetDataPath() + "cookies")
	panic(ErrorResp{
		Person: p,
		Code: ERR_CODE_COOKIE_INVALID,
		Msg: "cookie失效重登陆",
	})
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
		//"who":      p.UserName,
		"who":      "weidingyi",
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

	//_ = os.Remove(logName)

	noticeTitle := "【"+ titleMap[jobType]+"】【" + price + "】【"+ keyword +"】"

	return priceFloat, noticeTitle, "类别: " + titleMap[jobType] + RN + "价格: " + price + RN + "关键词: " + keyword
}


func (p *Person) DoJob() *Person {

	for {
		p.IsLogin()

		p.runJob()
		log.Println("ping...")
		time.Sleep(time.Second * 10)
	}
}


func newPerson(userName, password string) *Person {
	return &Person{
		UserName:    userName,
		Password:    password,
		LoginStatus: false,
	}
}

func run () {

	defer func() {
		if err := recover(); err != nil {
			dump.P(err)
			errRet := err.(ErrorResp)

			switch errRet.Code {
			case ERR_CODE_LOGIN_FALSE:
				log.Fatal(errRet.UserName, errRet.Msg, "检查账号密码")
			case ERR_CODE_COOKIE_INVALID:
				log.Println(errRet.Msg)
				//_ = os.Remove(errRet.Person.GetDataPath() + "cookies")
				//time.Sleep(time.Second)
				errRet.Person.DoJob()
			default:
				log.Println("未知错误")
				log.Fatal(errRet)
			}

		}
	}()
	p := newPerson(USERNAME_L, PASSWD_L)
	//p := newPerson(USERNAME_C, PASSWD_C)

	p.IsLogin()

	p.runJob()
	log.Println("ping...")
}

func init() {
	ptr := flag.String("path", ".", "set storage path")
	flag.Parse()
	dataDir = *ptr
	if dataDir != "." {
		dataDir, _ = filepath.Abs(dataDir)
	}
	dataDir = dataDir + string(os.PathSeparator) + "hsr-data" + string(os.PathSeparator)
}

func main() {

	go run()

	select {}
}
