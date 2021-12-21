package jd_cookie

import (
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"github.com/cdle/sillyGirl/core"
	"golang.org/x/net/proxy"
)

func init() {
	if !core.Bucket("qinglong").GetBool("enable_qinglong", true) {
		return
	}
	data, _ := os.ReadFile("dev.go")
	if !strings.Contains(string(data), "jd_cookie") && !jd_cookie.GetBool("enable_jd_cookie") {
		return
	}
	initAsset()
	initCheck()
	initEnEn()
	initEnv()
	initHelp()
	initLogin()
	initSubmit()
	initTyt()
	initNotify()
	buildHttpTransportWithProxy()
}

var Transport *http.Transport

func buildHttpTransportWithProxy() {
	addr := jd_cookie.Get("http_proxy")
	if strings.Contains(addr, "http://") {
		if addr != "" {
			u, err := url.Parse(addr)
			if err != nil {
				logs.Warn("can't connect to the http proxy:", err)
				return
			}
			Transport = &http.Transport{Proxy: http.ProxyURL(u)}
		}
	}
	if strings.Contains(addr, "sock5://") || strings.Contains(addr, "socks5://") {
		addr = strings.Replace(addr, "sock5://", "", -1)
		addr = strings.Replace(addr, "socks5://", "", -1)
		var auth *proxy.Auth
		v := strings.Split(addr, "@")
		if len(v) == 3 {
			auth = &proxy.Auth{
				User:     v[1],
				Password: v[2],
			}
			addr = v[0]
		}
		dialer, err := proxy.SOCKS5("tcp", addr, auth, proxy.Direct)
		if err != nil {
			logs.Warn("can't connect to the sock5 proxy:", err)
			return
		}
		Transport = &http.Transport{
			Dial: dialer.Dial,
		}
	}
	if Transport != nil {
		logs.Info("芝士夹心饼已启用。")
	}
}
