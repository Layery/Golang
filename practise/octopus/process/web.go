package process

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"octopus/output"
	"octopus/types"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/urfave/cli/v2"
)

const (
	BaseURL        = "https://t66y.com/thread0806.php"
	MaxConcurrency = 5
)

var handlerMap = map[string]types.CrawlerHandler{
	"caoliu": crawlCaoliu,
}

var proxyAddr string

var (
	globalClient     *http.Client
	globalClientOnce sync.Once
)

func init() {
	// 强制 net 包优先返回 IPv4 地址
	//net.DefaultResolver.PreferGo = true
}

func buildHTTPClient(proxyAddr string) *http.Client {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	if proxyAddr != "" {
		proxyURL, err := url.Parse(proxyAddr)
		if err != nil {
			//panic(errors.New("parse proxyAddr failed"))
		}
		client.Transport = &http.Transport{
			DisableKeepAlives: true,
			Proxy:             http.ProxyURL(proxyURL),
			TLSNextProto:      make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
			// TLS 配置
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
			// 超时设置
			ResponseHeaderTimeout: 30 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}
	}
	return client
}

// GetClient 在 Controller 内部调用
func GetClient(proxyAddr string) *http.Client {
	defer func() {
		if err := recover(); err != nil {
			log.Println("GetClient panic:", err)
		}
	}()
	if ok, err := regexp.Match("https?://", []byte(proxyAddr)); !ok || err != nil {
		log.Println("proxy url incorrect: ", proxyAddr)
		return nil
	}
	globalClientOnce.Do(func() {
		globalClient = buildHTTPClient(proxyAddr)
	})
	return globalClient
}

func RunWeb(ctx *cli.Context) error {
	defer log.Println("End.")

	// 根据name分发到不同的方法, 一等函数, 高阶函数?
	name := ctx.String("name")
	proxyAddr = ctx.String("proxy")

	// 分发逻辑
	handler, exists := handlerMap[name]
	if !exists {
		log.Printf("不支持的站点名称：%s ", name)
		return nil
	}

	defer func() {
		if panicVal := recover(); panicVal != nil {
			// 判断panic抛出的是不是context.Context
			if c, ok := panicVal.(context.Context); ok {
				log.Printf("错误上下文：%v", c.Err())
			} else {
				log.Printf("panic: %v \n", panicVal)
				// 打印完整的堆栈信息
				//debug.PrintStack()
			}
			//log.Fatal("爬虫执行异常，程序退出")
		}
	}()
	log.Println("开始执行")
	return handler(ctx)
}

func crawlCaoliu(cliContext *cli.Context) error {
	category := cliContext.String("category")
	page := cliContext.Int("page")
	proxyAddr = cliContext.String("proxy")
	ext := filepath.Ext(cliContext.String("output"))             // .html
	name := strings.TrimSuffix(cliContext.String("output"), ext) // ./output/dist/web
	finalOutput := fmt.Sprintf("%s.%d%s", name, page, ext)
	file, _ := output.NewFileWriter(finalOutput)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()

	baseUrl := BaseURL + "?fid=" + category + "&page=" + strconv.Itoa(page)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseUrl, nil)
	client := GetClient(proxyAddr)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("获取帖子列表失败: %v\n", err)
		return err
	}
	defer resp.Body.Close()
	// 4. 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil || body == nil {
		fmt.Printf("❌ 读取响应失败: %v\n", err)
		return err
	}
	fmt.Printf("✅ 获取帖子列表成功\n")
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		fmt.Printf("❌ 解析文档失败: %v\n", err)
		return err
	}
	// 找到所有帖子链接（语法类似 jQuery）
	wg := sync.WaitGroup{}
	semaphore := make(chan struct{}, MaxConcurrency) // 限制并发数为 5
	doc.Find(".tr3.t_one.tac").Each(func(i int, dom *goquery.Selection) {
		wg.Add(1)
		go func(index int, selection *goquery.Selection) {
			defer wg.Done()
			semaphore <- struct{}{}        // 获取信号量
			defer func() { <-semaphore }() // 释放信号量
			_ = crawlCaoliuItem(selection, file)
		}(i, dom)
	})
	wg.Wait()
	_ = file.Close()
	_ = OpenBrowser(finalOutput)
	return nil
}

func crawlCaoliuItem(dom *goquery.Selection, file *output.FileWriter) error {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("craw item panic: %v \n", err)
			//stack := debug.Stack()
			//log.Printf("📚 item堆栈跟踪:\n%s", string(stack))
			return
		}
	}()
	a := dom.Find("a").Eq(0)
	href, _ := a.Attr("href")
	urlObj, _ := url.Parse(BaseURL)
	href = urlObj.Scheme + "://" + urlObj.Host + href

	log.Println("开始处理:" + href)
	client := GetClient(proxyAddr)

	req, err := http.NewRequest("GET", href, nil)
	if err != nil {
		log.Printf("build req err: %v\n", err) // 单个item获取失败的情况下, 直接return
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("item send req err: %v\n", err) // 单个item获取失败的情况下, 直接return
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	// 解析出来内容, 生成html文件
	jsCode, _ := doc.Find("#conttpc a").Eq(1).Attr("onclick")

	re := regexp.MustCompile("https?://.*?'")
	href = re.FindString(jsCode)
	href = strings.TrimRight(href, "'")
	if href == "" {
		log.Printf("href is null, title: %s\n", a.Text())
		return nil
	}
	_ = file.Write(a.Text(), href)
	return nil
}
