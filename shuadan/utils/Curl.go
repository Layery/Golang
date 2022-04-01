package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

type Curl struct {
	Client *http.Client
	Url string
	Method string
	Body string
}

func NewCurl() *Curl {
	return new(Curl)
}


type RequestOption func() (*Curl, error)

func (c *Curl) Instance() RequestOption {
	return func() (*Curl, error) {
		return c, nil
	}
}


func (opt RequestOption) SetUrl(url string) RequestOption {
	return func() (c *Curl, e error) {
		c, err := opt()
		fmt.Printf("----> %#v \n", c)
		if err != nil {
			return c, err
		}
		c.Url = url
		return c, err
	}
}

func (c Curl) SetMethod(method string) *Curl {
	c.Method = method
	return &c
}

func (c Curl) SetBody(dada map[string]interface{}) Curl {
	postParams := ""
	for k, v := range dada {
		postParams += fmt.Sprintf("%s=%v&", k, v)
	}
	postParams = strings.TrimRight(postParams, "&")
	c.Body = postParams
	return c
}

func (c Curl) SetHeader()  {
	
}

func (c Curl) Send() {
	jar, _ := cookiejar.New(nil) // 为了根据域名, 安全的设置cookie
	client := http.Client{}
	client.Jar = jar

	request, _ := http.NewRequest(c.Method, c.Url, strings.NewReader(c.Body))
	request.Header.Set("X-Requested-With", "XMLHttpRequest")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")



	resp, _ := client.Do(request)
	body, _ := ioutil.ReadAll(resp.Body)

	P(string(body))




}


func init() {
}



































