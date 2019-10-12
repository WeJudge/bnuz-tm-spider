package main

import (
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)


const (
	MaxIdleConns int = 100
	MaxIdleConnsPerHost int = 100
	IdleConnTimeout int = 90
	UserAgent string = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.90 Safari/537.36"
)

type UserLoginedInfo struct {
	User struct {
		Id string `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
		DepartmentId string `json:"departmentId"`
		PhoneNumber string `json:"phoneNumber"`
	} `json:"user"`
}

func noRedirect(req *http.Request, via []*http.Request) error {
	location := req.URL.String()
	if strings.Index(location, "error") > -1 {
		return errors.New("Don't redirect!")
	}
	return nil
}

func buildClient (cookies *cookiejar.Jar, redirectCallback func(*http.Request, []*http.Request)error) *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   20 * time.Second,
				KeepAlive: 20 * time.Second,
			}).DialContext,
			MaxIdleConns:        MaxIdleConns,
			MaxIdleConnsPerHost: MaxIdleConnsPerHost,
			IdleConnTimeout:     time.Duration(IdleConnTimeout) * time.Second,
		},
		Timeout: 10 * time.Second,
		Jar: cookies,
	}
	if redirectCallback != nil {
		client.CheckRedirect = redirectCallback
	}
	return client
}


func initLogin (cookies *cookiejar.Jar) string {
	req, err := http.NewRequest("GET", "http://tm.bnuz.edu.cn/login", nil)
	if err != nil { return "" }
	req.Header.Set("User-Agent", UserAgent)

	client := buildClient(cookies, nil)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil { return "" }

	xsrftoken := ""
	url, err := url.Parse("http://tm.bnuz.edu.cn/uaa/login")
	if err != nil { return "" }

	for _, c := range cookies.Cookies(url) {
		if c.Name == "UAA-XSRF-TOKEN" {
			xsrftoken = c.Value
			break
		}
	}

	return xsrftoken
}

func postLogin(cookies *cookiejar.Jar, xsrftoken string, username string, password string) bool {
	req, err := http.NewRequest(
		"POST",
		"http://tm.bnuz.edu.cn/uaa/login",
		strings.NewReader("username=" + username + "&password=" + password),
	)
	if err != nil { return false }

	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://tm.bnuz.edu.cn/ui/login")
	req.Header.Set("X-UAA-XSRF-TOKEN", xsrftoken)

	client := buildClient(cookies, noRedirect)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	location := resp.Header.Get("Location")
	if strings.Index(location, "error") > -1 {
		return false
	}
	return true
}

func getUserInfo(cookies *cookiejar.Jar) *UserLoginedInfo {
	req, err := http.NewRequest("GET", "http://tm.bnuz.edu.cn/api/user", nil)
	if err != nil { return nil }
	req.Header.Set("User-Agent", UserAgent)

	client := buildClient(cookies, nil)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil { return nil }

	var userInfo UserLoginedInfo

	ret := JSONStringToObject(string(body), &userInfo)
	if !ret {
		return nil
	}

	return &userInfo
}

//func main() {
//	cookies, err := cookiejar.New(nil)
//	if err != nil {
//		return
//	}
//
//	xsrftoken := initLogin(cookies)
//	ret := postLogin(cookies, xsrftoken, "", "")
//
//	if ret {
//		userInfo := getUserInfo(cookies)
//		if userInfo != nil {
//			fmt.Println(ObjectToJSONStringFormatted(userInfo))
//		} else {
//			fmt.Println("解析失败")
//		}
//	} else {
//		fmt.Println("登录失败")
//	}
//}
