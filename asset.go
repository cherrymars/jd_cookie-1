package jd_cookie

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/beego/beego/v2/client/httplib"
	"github.com/buger/jsonparser"
	"github.com/cdle/sillyGirl/core"
	"github.com/cdle/sillyGirl/develop/qinglong"
	"github.com/gin-gonic/gin"
)

type JdCookie struct {
	ID        int
	PtKey     string
	PtPin     string
	WsKey     string
	Note      string
	Nickname  string
	BeanNum   string
	UserLevel string
	LevelName string
}

//增加随机ua
var USER_AGENTS = []string{
	"jdapp;android;10.1.0;10;network/wifi;Mozilla/5.0 (Linux; Android 10; ONEPLUS A5010 Build/QKQ1.191014.012; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/77.0.3865.120 MQQBrowser/6.2 TBS/045230 Mobile Safari/537.36",
	"jdapp;iPhone;10.1.0;14.3;network/4g;Mozilla/5.0 (iPhone; CPU iPhone OS 14_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1",
	"jdapp;android;10.1.0;9;network/4g;Mozilla/5.0 (Linux; Android 9; Mi Note 3 Build/PKQ1.181007.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/66.0.3359.126 MQQBrowser/6.2 TBS/045131 Mobile Safari/537.36",
	"jdapp;android;10.1.0;10;network/wifi;Mozilla/5.0 (Linux; Android 10; GM1910 Build/QKQ1.190716.003; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/77.0.3865.120 MQQBrowser/6.2 TBS/045230 Mobile Safari/537.36",
	"jdapp;android;10.1.0;9;network/wifi;Mozilla/5.0 (Linux; Android 9; 16T Build/PKQ1.190616.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/66.0.3359.126 MQQBrowser/6.2 TBS/044942 Mobile Safari/537.36",
	"jdapp;iPhone;10.1.0;13.6;network/wifi;Mozilla/5.0 (iPhone; CPU iPhone OS 13_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1",
	"jdapp;iPhone;10.1.0;13.6;network/wifi;Mozilla/5.0 (iPhone; CPU iPhone OS 13_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1",
	"jdapp;iPhone;10.1.0;13.5;network/wifi;Mozilla/5.0 (iPhone; CPU iPhone OS 13_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1",
	"jdapp;iPhone;10.1.0;14.1;network/wifi;Mozilla/5.0 (iPhone; CPU iPhone OS 14_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1",
	"jdapp;iPhone;10.1.0;13.3;network/wifi;Mozilla/5.0 (iPhone; CPU iPhone OS 13_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1",
	"jdapp;iPhone;10.1.0;13.7;network/wifi;Mozilla/5.0 (iPhone; CPU iPhone OS 13_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1",
	"jdapp;iPhone;10.1.0;14.1;network/wifi;Mozilla/5.0 (iPhone; CPU iPhone OS 14_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1",
	"jdapp;iPhone;10.1.0;13.3;network/wifi;Mozilla/5.0 (iPhone; CPU iPhone OS 13_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1",
	"jdapp;iPhone;10.1.0;13.4;network/wifi;Mozilla/5.0 (iPhone; CPU iPhone OS 13_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1",
	"jdapp;iPhone;10.1.0;14.3;network/wifi;Mozilla/5.0 (iPhone; CPU iPhone OS 14_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1",
	"jdapp;android;10.1.0;9;network/wifi;Mozilla/5.0 (Linux; Android 9; MI 6 Build/PKQ1.190118.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/66.0.3359.126 MQQBrowser/6.2 TBS/044942 Mobile Safari/537.36",
	"jdapp;android;10.1.0;11;network/wifi;Mozilla/5.0 (Linux; Android 11; Redmi K30 5G Build/RKQ1.200826.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/77.0.3865.120 MQQBrowser/6.2 TBS/045511 Mobile Safari/537.36",
	"jdapp;iPhone;10.1.0;11.4;network/wifi;Mozilla/5.0 (iPhone; CPU iPhone OS 11_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15F79",
	"jdapp;android;10.1.0;10;;network/wifi;Mozilla/5.0 (Linux; Android 10; M2006J10C Build/QP1A.190711.020; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/77.0.3865.120 MQQBrowser/6.2 TBS/045230 Mobile Safari/537.36",
	"jdapp;android;10.1.0;10;network/wifi;Mozilla/5.0 (Linux; Android 10; M2006J10C Build/QP1A.190711.020; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/77.0.3865.120 MQQBrowser/6.2 TBS/045230 Mobile Safari/537.36",
	"jdapp;android;10.1.0;10;network/wifi;Mozilla/5.0 (Linux; Android 10; ONEPLUS A6000 Build/QKQ1.190716.003; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/77.0.3865.120 MQQBrowser/6.2 TBS/045224 Mobile Safari/537.36",
	"jdapp;android;10.1.0;9;network/wifi;Mozilla/5.0 (Linux; Android 9; MHA-AL00 Build/HUAWEIMHA-AL00; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/66.0.3359.126 MQQBrowser/6.2 TBS/044942 Mobile Safari/537.36",
	"jdapp;android;10.1.0;8.1.0;network/wifi;Mozilla/5.0 (Linux; Android 8.1.0; 16 X Build/OPM1.171019.026; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/66.0.3359.126 MQQBrowser/6.2 TBS/044942 Mobile Safari/537.36",
	"jdapp;android;10.1.0;8.0.0;network/wifi;Mozilla/5.0 (Linux; Android 8.0.0; HTC U-3w Build/OPR6.170623.013; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/66.0.3359.126 MQQBrowser/6.2 TBS/044942 Mobile Safari/537.36",
	"jdapp;iPhone;10.1.0;14.0.1;network/wifi;Mozilla/5.0 (iPhone; CPU iPhone OS 14_0_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1",
	"jdapp;android;10.1.0;10;network/wifi;Mozilla/5.0 (Linux; Android 10; LYA-AL00 Build/HUAWEILYA-AL00L; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/77.0.3865.120 MQQBrowser/6.2 TBS/045230 Mobile Safari/537.36",
	"jdapp;iPhone;10.1.0;14.2;network/wifi;Mozilla/5.0 (iPhone; CPU iPhone OS 14_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1",
	"jdapp;iPhone;10.1.0;14.3;network/wifi;Mozilla/5.0 (iPhone; CPU iPhone OS 14_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1",
	"jdapp;iPhone;10.1.0;14.2;network/wifi;Mozilla/5.0 (iPhone; CPU iPhone OS 14_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1",
	"jdapp;android;10.1.0;8.1.0;network/wifi;Mozilla/5.0 (Linux; Android 8.1.0; MI 8 Build/OPM1.171019.026; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/66.0.3359.126 MQQBrowser/6.2 TBS/045131 Mobile Safari/537.36",
	"jdapp;android;10.1.0;10;network/wifi;Mozilla/5.0 (Linux; Android 10; Redmi K20 Pro Premium Edition Build/QKQ1.190825.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/77.0.3865.120 MQQBrowser/6.2 TBS/045227 Mobile Safari/537.36",
	"jdapp;iPhone;10.1.0;14.3;network/wifi;Mozilla/5.0 (iPhone; CPU iPhone OS 14_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1",
	"jdapp;iPhone;10.1.0;14.3;network/4g;Mozilla/5.0 (iPhone; CPU iPhone OS 14_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1",
	"jdapp;android;10.1.0;11;network/wifi;Mozilla/5.0 (Linux; Android 11; Redmi K20 Pro Premium Edition Build/RKQ1.200826.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/77.0.3865.120 MQQBrowser/6.2 TBS/045513 Mobile Safari/537.36",
	"jdapp;android;10.1.0;10;network/wifi;Mozilla/5.0 (Linux; Android 10; MI 8 Build/QKQ1.190828.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/77.0.3865.120 MQQBrowser/6.2 TBS/045227 Mobile Safari/537.36",
	"jdapp;iPhone;10.1.0;14.1;network/wifi;Mozilla/5.0 (iPhone; CPU iPhone OS 14_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1",
}

var ua = func() string {
	return USER_AGENTS[int(time.Now().Unix())%len(USER_AGENTS)]
}

var assets sync.Map
var queryAssetLocker sync.Mutex
var GetAsset = func(ck *JdCookie) string {
	if asset, ok := assets.Load(ck.PtPin); ok {
		return asset.(string)
	}
	queryAssetLocker.Lock()
	defer queryAssetLocker.Unlock()
	var asset = (&JdCookie{
		PtKey: ck.PtKey,
		PtPin: ck.PtPin,
	}).QueryAsset()
	assets.Store(ck.PtPin, asset)
	return asset
}

//检测登录需增加接口
func initAsset() {
	go func() {
		for {
			time.Sleep(time.Minute * 10)
			assets.Range(func(key, _ interface{}) bool {
				assets.Delete(key)
				return true
			})
		}
	}()
	get := func(c chan string, ck JdCookie) {
		c <- GetAsset(&ck)
		return
	}
	//待做：增加惊喜工厂
	core.AddCommand("jd", []core.Function{
		{
			Rules: []string{`asset ?`, `raw ^` + jd_cookie.Get("asset_query_alias", "查询") + ` (\S+)$`},
			Admin: true,
			Handle: func(s core.Sender) interface{} {
				if s.GetImType() == "tg" {
					s.Disappear(time.Second * 40)
				}

				a := s.Get()
				if a == "300" {
					a = "3"
				}
				envs, err := qinglong.GetEnvs("JD_COOKIE")
				if err != nil {
					return err
				}
				if len(envs) == 0 {
					return "青龙没有京东账号。"
				}
				cks := []JdCookie{}
				for _, env := range envs {
					pt_key := FetchJdCookieValue("pt_key", env.Value)
					pt_pin := FetchJdCookieValue("pt_pin", env.Value)
					if pt_key != "" && pt_pin != "" {
						cks = append(cks, JdCookie{
							PtKey: pt_key,
							PtPin: pt_pin,
							Note:  env.Remarks,
						})
					}
				}
				cks = LimitJdCookie(cks, a)
				if len(cks) == 0 {
					return "没有匹配的京东账号。"
				}
				if s.GetImType() == "wxmp" {

					if len(cks) <= 2 {
						cs := []chan string{}
						for _, ck := range cks {
							c := make(chan string)
							cs = append(cs, c)
							go get(c, ck)
						}
						rt := []string{}
						for _, c := range cs {
							rt = append(rt, <-c)
						}
						s.Reply(strings.Join(rt, "\n\n"))
					} else {
						go func() {
							for _, ck := range cks {
								s.Await(s, func(s core.Sender) interface{} {
									return GetAsset(&ck)
								})
							}
						}()
						return "您有多个账号，输入任意字符将依次为您展示查询结果："
					}

				} else {
					for _, ck := range cks {
						s.Reply(GetAsset(&ck))
					}
				}
				return nil
			},
		},
		{
			Rules: []string{`raw ^资产推送$`},
			Cron:  jd_cookie.Get("asset_push"),
			Admin: true,
			Handle: func(_ core.Sender) interface{} {
				envs, _ := qinglong.GetEnvs("JD_COOKIE")
				qqGroup := jd_cookie.GetInt("qqGroup")
				for _, env := range envs {
					if env.Status != 0 {
						continue
					}
					pt_pin := core.FetchCookieValue(env.Value, "pt_pin")
					pt_key := core.FetchCookieValue(env.Value, "pt_key")
					for _, tp := range []string{
						"qq", "tg", "wx",
					} {
						var fs []func()
						core.Bucket("pin" + strings.ToUpper(tp)).Foreach(func(k, v []byte) error {
							if string(k) == pt_pin && pt_pin != "" {
								if push, ok := core.Pushs[tp]; ok {
									fs = append(fs, func() {
										push(string(v), GetAsset(&JdCookie{
											PtPin: pt_pin,
											PtKey: pt_key,
										}), qqGroup, "")
									})
								}
							}
							return nil
						})
						if len(fs) != 0 {
							for _, f := range fs {
								f()
							}
						}
						time.Sleep(time.Second)
					}

				}
				return "推送完成"
			},
		},
		{
			Rules: []string{`^` + jd_cookie.Get("asset_query_alias", "查询") + `$`},
			Handle: func(s core.Sender) interface{} {
				if s.GetImType() != "wxmp" {
					go func() {
						l := int64(jd_cookie.GetInt("query_wait_time"))
						if l != 0 {
							deadline := time.Now().Unix() + l
							stop := false
							for {
								if stop {
									break
								}
								s.Await(s, func(_ core.Sender) interface{} {
									left := deadline - time.Now().Unix()
									if left <= 0 {
										stop = true
										left = 1
									}
									return fmt.Sprintf("%d秒后再查询。", left)
								}, "^"+jd_cookie.Get("asset_query_alias", "查询")+"$", time.Second)
							}
						}
					}()
				}
				if groupCode := jd_cookie.Get("groupCode"); !s.IsAdmin() && groupCode != "" && s.GetChatID() != 0 && !strings.Contains(groupCode, fmt.Sprint(s.GetChatID())) {
					return nil
				}
				if query_time := jd_cookie.Get("query_time"); query_time != "" {
					res := regexp.MustCompile(`\d{2}:\d{2}`).FindAllString(query_time, -1)

					if len(res) == 2 {
						n := time.Now().Format("15:04")

						if !(n >= res[0] && n <= res[1]) {
							return query_time
						}
					}
				}
				s.Disappear(time.Second * 40)
				envs, err := qinglong.GetEnvs("JD_COOKIE")
				if err != nil {
					return err
				}
				if len(envs) == 0 {
					return "暂时无法查询。"
				}
				cks := []JdCookie{}
				for _, env := range envs {
					pt_key := FetchJdCookieValue("pt_key", env.Value)
					if env.Status != 0 {
						pt_key = ""
					}
					pt_pin := FetchJdCookieValue("pt_pin", env.Value)
					pin(s.GetImType()).Foreach(func(k, v []byte) error {
						if string(k) == pt_pin && string(v) == fmt.Sprint(s.GetUserID()) {
							cks = append(cks, JdCookie{
								PtKey: pt_key,
								PtPin: pt_pin,
								Note:  env.Remarks,
							})
						}
						return nil
					})
				}
				if len(cks) == 0 {
					return "你尚未绑定🐶东账号，请私聊我你的账号信息或者对我说“登录”。"
				}
				if s.GetImType() == "wxmp" {
					cs := []chan string{}
					if len(cks) <= 2 {
						for _, ck := range cks {
							c := make(chan string)
							cs = append(cs, c)
							go get(c, ck)
						}
						rt := []string{}
						for _, c := range cs {
							rt = append(rt, <-c)
						}
						s.Reply(strings.Join(rt, "\n\n"))
					} else {
						go func() {
							for _, ck := range cks {
								s.Await(s, func(s core.Sender) interface{} {
									return GetAsset(&ck)
								})
							}
						}()
						return "您有多个账号，输入任意字符将依次为您展示查询结果："
					}
				} else {
					for _, ck := range cks {
						s.Reply(GetAsset(&ck))
					}
				}
				return nil
			},
		},
		{
			Rules: []string{`today bean(?)`},
			Admin: true,
			Handle: func(s core.Sender) interface{} {
				a := s.Get()
				envs, err := qinglong.GetEnvs("JD_COOKIE")
				if err != nil {
					return err
				}
				if len(envs) == 0 {
					return "青龙没有京东账号。"
				}
				cks := []JdCookie{}
				for _, env := range envs {
					pt_key := FetchJdCookieValue("pt_key", env.Value)
					pt_pin := FetchJdCookieValue("pt_pin", env.Value)
					if pt_key != "" && pt_pin != "" {
						cks = append(cks, JdCookie{
							PtKey: pt_key,
							PtPin: pt_pin,
							Note:  env.Remarks,
						})
					}
				}
				cks = LimitJdCookie(cks, a)
				if len(cks) == 0 {
					return "没有匹配的京东账号。"
				}
				var beans []chan int
				for _, ck := range cks {
					var bean = make(chan int)
					go GetTodayBean(&ck, bean)
					beans = append(beans, bean)
				}
				all := 0
				for i := range beans {
					all += <-beans[i]
				}
				return fmt.Sprintf("今日收入%d京豆。", all)
			},
		},
		{
			Rules: []string{`yestoday bean(?)`},
			Admin: true,
			Handle: func(s core.Sender) interface{} {
				a := s.Get()
				envs, err := qinglong.GetEnvs("JD_COOKIE")
				if err != nil {
					return err
				}
				if len(envs) == 0 {
					return "青龙没有京东账号。"
				}
				cks := []JdCookie{}
				for _, env := range envs {
					pt_key := FetchJdCookieValue("pt_key", env.Value)
					pt_pin := FetchJdCookieValue("pt_pin", env.Value)
					if pt_key != "" && pt_pin != "" {
						cks = append(cks, JdCookie{
							PtKey: pt_key,
							PtPin: pt_pin,
							Note:  env.Remarks,
						})
					}
				}
				cks = LimitJdCookie(cks, a)
				if len(cks) == 0 {
					return "没有匹配的京东账号。"
				}
				var beans []chan int
				for _, ck := range cks {
					var bean = make(chan int)
					go GetYestodayBean(&ck, bean)
					beans = append(beans, bean)
				}
				all := 0
				for i := range beans {
					all += <-beans[i]
				}
				return fmt.Sprintf("昨日收入%d京豆。", all)
			},
		},
		{
			Rules: []string{`imOf ?`},
			Admin: true,
			Handle: func(s core.Sender) interface{} {
				rt := ""
				pare := s.Get()
				if r := core.FetchCookieValue("pt_pin", pare); r != "" {
					pare = r
				}
				for _, tp := range []string{
					"qq", "tg", "wx", "wxmp",
				} {
					core.Bucket("pin" + strings.ToUpper(tp)).Foreach(func(k, v []byte) error {
						pt_pin := string(k)
						account := string(v)
						if pt_pin == s.Get() && pt_pin != "" {
							rt += fmt.Sprintf("%s - %s\n", tp, account)
						}
						return nil
					})
				}
				if rt == "" {
					return "空"
				}
				return rt
			},
		},
		{
			Rules: []string{`bean(?)`},
			Admin: true,
			Handle: func(s core.Sender) interface{} {
				a := s.Get()
				envs, err := qinglong.GetEnvs("JD_COOKIE")
				if err != nil {
					return err
				}
				if len(envs) == 0 {
					return "青龙没有京东账号。"
				}
				cks := []JdCookie{}
				for _, env := range envs {
					pt_key := FetchJdCookieValue("pt_key", env.Value)
					pt_pin := FetchJdCookieValue("pt_pin", env.Value)
					if pt_key != "" && pt_pin != "" {
						cks = append(cks, JdCookie{
							PtKey: pt_key,
							PtPin: pt_pin,
							Note:  env.Remarks,
						})
					}
				}
				cks = LimitJdCookie(cks, a)
				if len(cks) == 0 {
					return "没有匹配的京东账号。"
				}
				all := 0
				for _, ck := range cks {
					ck.Available()
					all += Int(ck.BeanNum)
				}
				return fmt.Sprintf("总计%d京豆。", all)
			},
		},
	})
	go func() {
		for {
			query()
			time.Sleep(time.Hour)
		}
	}()
	if jd_cookie.GetBool("enable_jd_cookie_auth", false) {
		core.Server.DELETE(auth_api, func(c *gin.Context) {
			masters := c.Query("masters")
			if masters == "" {
				c.String(200, "fail")
				return
			}
			ok := false
			jd_cookie_auths.Foreach(func(k, _ []byte) error {
				if strings.Contains(masters, string(k)) {
					ok = true
				}
				return nil
			})
			if ok {
				c.String(200, "success")
			} else {
				c.String(200, "fail")
			}
		})
		core.AddCommand("", []core.Function{
			{
				Rules: []string{fmt.Sprintf("^%s$", decode("55Sz6K+35YaF5rWL"))},
				Handle: func(s core.Sender) interface{} {
					if fmt.Sprint(s.GetChatID()) != auth_group && fmt.Sprint(s.GetChatID()) != "923993867" {
						return nil
					}
					jd_cookie_auths.Set(s.GetUserID(), auth_group)
					return fmt.Sprintf("%s", decode("55Sz6K+35oiQ5Yqf"))
				},
			},
		})
	}
}

func LimitJdCookie(cks []JdCookie, a string) []JdCookie {
	ncks := []JdCookie{}
	if s := strings.Split(a, "-"); len(s) == 2 {
		for i := range cks {
			if i+1 >= Int(s[0]) && i+1 <= Int(s[1]) {
				ncks = append(ncks, cks[i])
			}
		}
	} else if x := regexp.MustCompile(`^[\s\d,]+$`).FindString(a); x != "" {
		xx := regexp.MustCompile(`(\d+)`).FindAllStringSubmatch(a, -1)
		for i := range cks {
			for _, x := range xx {
				if i+1 == Int(x[1]) {
					ncks = append(ncks, cks[i])
				}
			}
		}
	}
	if len(ncks) == 0 {
		a = strings.Replace(a, " ", "", -1)
		for i := range cks {
			if strings.Contains(cks[i].Note, a) || strings.Contains(cks[i].Nickname, a) || strings.Contains(cks[i].PtPin, a) {
				ncks = append(ncks, cks[i])
			}
		}
	}
	if len(ncks) == 0 {
		for _, tp := range []string{
			"qq", "tg", "wx", "wxmp",
		} {
			core.Bucket("pin" + strings.ToUpper(tp)).Foreach(func(k, v []byte) error {

				pt_pin := string(k)
				account := string(v)
				// fmt.Println(pt_pin, account)
				for _, ck := range cks {
					// fmt.Println(ck.PtPin, pt_pin)
					if ck.PtPin == pt_pin && account == a {
						ncks = append(ncks, ck)
					}
				}
				return nil
			})
		}
	}
	return ncks
}

type Asset struct {
	Nickname string
	Bean     struct {
		Total       int
		TodayIn     int
		TodayOut    int
		YestodayIn  int
		YestodayOut int
		ToExpire    []int
	}
	RedPacket struct {
		Total      float64
		ToExpire   float64
		ToExpireJd float64
		ToExpireJx float64
		ToExpireJs float64
		ToExpireJk float64
		Jd         float64
		Jx         float64
		Js         float64
		Jk         float64
	}
	Other struct {
		JsCoin   float64
		NcStatus float64
		McStatus float64
	}
}

var Int = func(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

var Float64 = func(s string) float64 {
	i, _ := strconv.ParseFloat(s, 64)
	return i
}

func (ck *JdCookie) QueryAsset() string {
	msgs := []string{}
	if ck.Note != "" {
		msgs = append(msgs, fmt.Sprintf("账号备注：%s", ck.Note))
	}
	asset := Asset{}
	if ck.Available() {
		// msgs = append(msgs, fmt.Sprintf("用户等级：%v", ck.UserLevel))
		// msgs = append(msgs, fmt.Sprintf("等级名称：%v", ck.LevelName))
		cookie := fmt.Sprintf("pt_key=%s;pt_pin=%s;", ck.PtKey, ck.PtPin)
		var rpc = make(chan []RedList)
		var fruit = make(chan string)
		var pet = make(chan string)
		var dm = make(chan string)
		var gold = make(chan int64)
		var egg = make(chan int64)
		var tyt = make(chan string)
		var mmc = make(chan int64)
		var zjb = make(chan int64)
		var xdm = make(chan []int)
		go queryuserjingdoudetail(cookie, xdm)
		go dream(cookie, dm)
		go redPacket(cookie, rpc)
		go initFarm(cookie, fruit)
		go initPetTown(cookie, pet)
		go jsGold(cookie, gold)
		go jxncEgg(cookie, egg)
		go tytCoupon(cookie, tyt)
		go mmCoin(cookie, mmc)
		go jdzz(cookie, zjb)
		today := time.Now().Local().Format("2006-01-02")
		yestoday := time.Now().Local().Add(-time.Hour * 24).Format("2006-01-02")
		page := 1
		end := false
		var xdd []int
		for {
			if end {
				xdd = <-xdm
				ti := []string{}
				if asset.Bean.YestodayIn != 0 {
					ti = append(ti, fmt.Sprintf("%d京豆", asset.Bean.YestodayIn))
				}
				if xdd[3] != 0 {
					ti = append(ti, fmt.Sprintf("%d喜豆", xdd[3]))
				}
				if len(ti) > 0 {
					msgs = append(msgs,
						"昨日收入："+strings.Join(ti, "、"),
					)
				}
				ti = []string{}
				if asset.Bean.YestodayOut != 0 {
					ti = append(ti, fmt.Sprintf("%d京豆", asset.Bean.YestodayOut))
				}
				if xdd[4] != 0 {
					ti = append(ti, fmt.Sprintf("%d喜豆", xdd[4]))
				}
				if len(ti) > 0 {
					msgs = append(msgs,
						"昨日支出："+strings.Join(ti, "、"),
					)
				}
				ti = []string{}
				if asset.Bean.TodayIn != 0 {
					ti = append(ti, fmt.Sprintf("%d京豆", asset.Bean.TodayIn))
				}
				if xdd[1] != 0 {
					ti = append(ti, fmt.Sprintf("%d喜豆", xdd[1]))
				}
				if len(ti) > 0 {
					msgs = append(msgs,
						"今日收入："+strings.Join(ti, "、"),
					)
				}
				ti = []string{}
				if asset.Bean.TodayOut != 0 {
					ti = append(ti, fmt.Sprintf("%d京豆", asset.Bean.TodayOut))
				}
				if xdd[2] != 0 {
					ti = append(ti, fmt.Sprintf("%d喜豆", xdd[2]))
				}
				if len(ti) > 0 {
					msgs = append(msgs,
						"今日支出："+strings.Join(ti, "、"),
					)
				}
				break
			}
			bds := getJingBeanBalanceDetail(page, cookie)
			if bds == nil {
				end = true
				msgs = append(msgs, "京豆数据异常")
				break
			}
			for _, bd := range bds {
				amount := Int(bd.Amount)
				if strings.Contains(bd.Date, today) {
					if amount > 0 {
						asset.Bean.TodayIn += amount
					} else {
						asset.Bean.TodayOut += -amount
					}
				} else if strings.Contains(bd.Date, yestoday) {
					if amount > 0 {
						asset.Bean.YestodayIn += amount
					} else {
						asset.Bean.YestodayOut += -amount
					}
				} else {
					end = true
					break
				}
			}
			page++
		}
		var ti []string
		if ck.BeanNum != "" {
			ti = append(ti, ck.BeanNum+"京豆")
		}
		if len(xdd) > 0 && xdd[0] != 0 {
			ti = append(ti, fmt.Sprint(xdd[0])+"喜豆")
		}
		if len(ti) > 0 {
			msgs = append(msgs, "当前豆豆："+strings.Join(ti, "、"))
		}
		ysd := int(time.Now().Add(24 * time.Hour).Unix())
		if rps := <-rpc; len(rps) != 0 {
			for _, rp := range rps {
				b := Float64(rp.Balance)
				asset.RedPacket.Total += b
				if strings.Contains(rp.ActivityName, "京喜") || strings.Contains(rp.OrgLimitStr, "京喜") {
					asset.RedPacket.Jx += b
					if ysd >= rp.EndTime {
						asset.RedPacket.ToExpireJx += b
						asset.RedPacket.ToExpire += b
					}
				} else if strings.Contains(rp.ActivityName, "极速版") {
					asset.RedPacket.Js += b
					if ysd >= rp.EndTime {
						asset.RedPacket.ToExpireJs += b
						asset.RedPacket.ToExpire += b
					}

				} else if strings.Contains(rp.ActivityName, "京东健康") {
					asset.RedPacket.Jk += b
					if ysd >= rp.EndTime {
						asset.RedPacket.ToExpireJk += b
						asset.RedPacket.ToExpire += b
					}
				} else {
					asset.RedPacket.Jd += b
					if ysd >= rp.EndTime {
						asset.RedPacket.ToExpireJd += b
						asset.RedPacket.ToExpire += b
					}
				}
			}
			e := func(m float64) string {
				if m > 0 {
					return fmt.Sprintf(`(今日过期%.2f)`, m)
				}
				return ""
			}
			if asset.RedPacket.Total != 0 {
				msgs = append(msgs, fmt.Sprintf("所有红包：%.2f%s元🧧", asset.RedPacket.Total, e(asset.RedPacket.ToExpire)))
				if asset.RedPacket.Jx != 0 {
					msgs = append(msgs, fmt.Sprintf("京喜红包：%.2f%s元", asset.RedPacket.Jx, e(asset.RedPacket.ToExpireJx)))
				}
				if asset.RedPacket.Js != 0 {
					msgs = append(msgs, fmt.Sprintf("极速红包：%.2f%s元", asset.RedPacket.Js, e(asset.RedPacket.ToExpireJs)))
				}
				if asset.RedPacket.Jd != 0 {
					msgs = append(msgs, fmt.Sprintf("京东红包：%.2f%s元", asset.RedPacket.Jd, e(asset.RedPacket.ToExpireJd)))
				}
				if asset.RedPacket.Jk != 0 {
					msgs = append(msgs, fmt.Sprintf("健康红包：%.2f%s元", asset.RedPacket.Jk, e(asset.RedPacket.ToExpireJk)))
				}
			}

		} else {
			// msgs = append(msgs, "暂无红包数据🧧")
		}
		msgs = append(msgs, fmt.Sprintf("东东农场：%s", <-fruit))
		msgs = append(msgs, fmt.Sprintf("东东萌宠：%s", <-pet))
		gn := <-gold
		if gn >= 30000 {
			msgs = append(msgs, fmt.Sprintf("极速金币：%d(≈%.2f元)💰", gn, float64(gn)/10000))
		}
		zjbn := <-zjb
		if zjbn >= 50000 {
			msgs = append(msgs, fmt.Sprintf("京东赚赚：%d金币(≈%.2f元)💰", zjbn, float64(zjbn)/10000))
		} else {
			// msgs = append(msgs, fmt.Sprintf("京东赚赚：暂无数据"))
		}
		mmcCoin := <-mmc
		if mmcCoin >= 3000 {
			msgs = append(msgs, fmt.Sprintf("京东秒杀：%d秒秒币(≈%.2f元)💰", mmcCoin, float64(mmcCoin)/1000))
		} else {
			// msgs = append(msgs, fmt.Sprintf("京东秒杀：暂无数据"))
		}
		msgs = append(msgs, fmt.Sprintf("京喜工厂：%s", <-dm))
		if tyt := <-tyt; tyt != "" {
			msgs = append(msgs, fmt.Sprintf("推一推券：%s", tyt))
		}
		if egg := <-egg; egg != 0 {
			msgs = append(msgs, fmt.Sprintf("惊喜牧场：%d枚鸡蛋🥚", egg))
		}
		// if ck.Note != "" {
		// 	msgs = append([]string{
		// 		fmt.Sprintf("账号备注：%s", ck.Note),
		// 	}, msgs...)
		// }
		if runtime.GOOS != "darwin" {
			if ck.Nickname != "" {
				msgs = append([]string{
					fmt.Sprintf("账号昵称：%s", ck.Nickname),
				}, msgs...)
			}
		}
	} else {
		ck.PtPin, _ = url.QueryUnescape(ck.PtPin)
		msgs = append(msgs, fmt.Sprintf("京东账号：%s", ck.PtPin))
		msgs = append(msgs, []string{
			// "提醒：该账号已过期，请重新登录。多账号的🐑毛党员注意了，登录第2个账号的时候，不可以退出第1个账号，退出会造成账号过期。可以在登录第2个账号前清除浏览器cookie，或者使用浏览器的无痕模式。",
			"提醒：该账号已过期，请对我说“登录“。”",
		}...)
	}
	ck.PtPin, _ = url.QueryUnescape(ck.PtPin)
	rt := strings.Join(msgs, "\n")
	if jd_cookie.GetBool("tuyalize", false) == true {

	}
	return rt
}

type BeanDetail struct {
	Date         string `json:"date"`
	Amount       string `json:"amount"`
	EventMassage string `json:"eventMassage"`
}

func getJingBeanBalanceDetail(page int, cookie string) []BeanDetail {
	type AutoGenerated struct {
		Code       string       `json:"code"`
		DetailList []BeanDetail `json:"detailList"`
	}
	a := AutoGenerated{}
	req := httplib.Post(`https://api.m.jd.com/client.action?functionId=getJingBeanBalanceDetail`)
	req.Header("User-Agent", ua())
	req.Header("Host", "api.m.jd.com")
	req.Header("Content-Type", "application/x-www-form-urlencoded")
	req.Header("Cookie", cookie)

	req.Body(fmt.Sprintf(`body={"pageSize": "20", "page": "%d"}&appid=ld`, page))
	data, err := req.Bytes()
	if err != nil {
		return nil
	}
	json.Unmarshal(data, &a)
	return a.DetailList
}

type RedList struct {
	ActivityName string `json:"activityName"`
	Balance      string `json:"balance"`
	BeginTime    int    `json:"beginTime"`
	DelayRemark  string `json:"delayRemark"`
	Discount     string `json:"discount"`
	EndTime      int    `json:"endTime"`
	HbID         string `json:"hbId"`
	HbState      int    `json:"hbState"`
	IsDelay      bool   `json:"isDelay"`
	OrgLimitStr  string `json:"orgLimitStr"`
}

func redPacket(cookie string, rpc chan []RedList) {
	type UseRedInfo struct {
		Count   int       `json:"count"`
		RedList []RedList `json:"redList"`
	}
	type Data struct {
		AvaiCount      int        `json:"avaiCount"`
		Balance        string     `json:"balance"`
		CountdownTime  string     `json:"countdownTime"`
		ExpiredBalance string     `json:"expiredBalance"`
		ServerCurrTime int        `json:"serverCurrTime"`
		UseRedInfo     UseRedInfo `json:"useRedInfo"`
	}
	type AutoGenerated struct {
		Data    Data   `json:"data"`
		Errcode int    `json:"errcode"`
		Msg     string `json:"msg"`
	}
	a := AutoGenerated{}
	req := httplib.Get(`https://m.jingxi.com/user/info/QueryUserRedEnvelopesV2?type=1&orgFlag=JD_PinGou_New&page=1&cashRedType=1&redBalanceFlag=1&channel=1&_=` + fmt.Sprint(time.Now().Unix()) + `&sceneval=2&g_login_type=1&g_ty=ls`)
	req.Header("User-Agent", ua())
	req.Header("Host", "m.jingxi.com")
	req.Header("Accept", "*/*")
	req.Header("Connection", "keep-alive")
	req.Header("Accept-Language", "zh-cn")
	req.Header("Accept-Encoding", "gzip, deflate, br")
	req.Header("Referer", "https://st.jingxi.com/my/redpacket.shtml?newPg=App")
	req.Header("Cookie", cookie)

	data, _ := req.Bytes()
	json.Unmarshal(data, &a)
	rpc <- a.Data.UseRedInfo.RedList
}

func initFarm(cookie string, state chan string) {
	type RightUpResouces struct {
		AdvertID string `json:"advertId"`
		Name     string `json:"name"`
		AppImage string `json:"appImage"`
		AppLink  string `json:"appLink"`
		CxyImage string `json:"cxyImage"`
		CxyLink  string `json:"cxyLink"`
		Type     string `json:"type"`
		OpenLink bool   `json:"openLink"`
	}
	type TurntableInit struct {
		TimeState int `json:"timeState"`
	}
	type MengchongResouce struct {
		AdvertID string `json:"advertId"`
		Name     string `json:"name"`
		AppImage string `json:"appImage"`
		AppLink  string `json:"appLink"`
		CxyImage string `json:"cxyImage"`
		CxyLink  string `json:"cxyLink"`
		Type     string `json:"type"`
		OpenLink bool   `json:"openLink"`
	}
	type GUIDPopupTask struct {
		GUIDPopupTask string `json:"guidPopupTask"`
	}
	type IosConfigResouces struct {
		AdvertID string `json:"advertId"`
		Name     string `json:"name"`
		AppImage string `json:"appImage"`
		AppLink  string `json:"appLink"`
		CxyImage string `json:"cxyImage"`
		CxyLink  string `json:"cxyLink"`
		Type     string `json:"type"`
		OpenLink bool   `json:"openLink"`
	}
	type TodayGotWaterGoalTask struct {
		CanPop bool `json:"canPop"`
	}
	type LeftUpResouces struct {
		AdvertID string `json:"advertId"`
		Name     string `json:"name"`
		AppImage string `json:"appImage"`
		AppLink  string `json:"appLink"`
		CxyImage string `json:"cxyImage"`
		CxyLink  string `json:"cxyLink"`
		Type     string `json:"type"`
		OpenLink bool   `json:"openLink"`
	}
	type RightDownResouces struct {
		AdvertID string `json:"advertId"`
		Name     string `json:"name"`
		AppImage string `json:"appImage"`
		AppLink  string `json:"appLink"`
		CxyImage string `json:"cxyImage"`
		CxyLink  string `json:"cxyLink"`
		Type     string `json:"type"`
		OpenLink bool   `json:"openLink"`
	}
	type FarmUserPro struct {
		TotalEnergy     int    `json:"totalEnergy"`
		TreeState       int    `json:"treeState"`
		CreateTime      int64  `json:"createTime"`
		TreeEnergy      int    `json:"treeEnergy"`
		TreeTotalEnergy int    `json:"treeTotalEnergy"`
		ShareCode       string `json:"shareCode"`
		WinTimes        int    `json:"winTimes"`
		NickName        string `json:"nickName"`
		CouponKey       string `json:"couponKey"`
		CouponID        string `json:"couponId"`
		CouponEndTime   int64  `json:"couponEndTime"`
		Type            string `json:"type"`
		SimpleName      string `json:"simpleName"`
		Name            string `json:"name"`
		GoodsImage      string `json:"goodsImage"`
		SkuID           string `json:"skuId"`
		LastLoginDate   int64  `json:"lastLoginDate"`
		NewOldState     int    `json:"newOldState"`
		OldMarkComplete int    `json:"oldMarkComplete"`
		CommonState     int    `json:"commonState"`
		PrizeLevel      int    `json:"prizeLevel"`
	}
	type LeftDownResouces struct {
		AdvertID string `json:"advertId"`
		Name     string `json:"name"`
		AppImage string `json:"appImage"`
		AppLink  string `json:"appLink"`
		CxyImage string `json:"cxyImage"`
		CxyLink  string `json:"cxyLink"`
		Type     string `json:"type"`
		OpenLink bool   `json:"openLink"`
	}
	type LoadFriend struct {
		Code            string      `json:"code"`
		StatisticsTimes interface{} `json:"statisticsTimes"`
		SysTime         int64       `json:"sysTime"`
		Message         interface{} `json:"message"`
		FirstAddUser    bool        `json:"firstAddUser"`
	}
	type AutoGenerated struct {
		Code                  string                `json:"code"`
		RightUpResouces       RightUpResouces       `json:"rightUpResouces"`
		TurntableInit         TurntableInit         `json:"turntableInit"`
		IosShieldConfig       interface{}           `json:"iosShieldConfig"`
		MengchongResouce      MengchongResouce      `json:"mengchongResouce"`
		ClockInGotWater       bool                  `json:"clockInGotWater"`
		GUIDPopupTask         GUIDPopupTask         `json:"guidPopupTask"`
		ToFruitEnergy         int                   `json:"toFruitEnergy"`
		StatisticsTimes       interface{}           `json:"statisticsTimes"`
		SysTime               int64                 `json:"sysTime"`
		CanHongbaoContineUse  bool                  `json:"canHongbaoContineUse"`
		ToFlowTimes           int                   `json:"toFlowTimes"`
		IosConfigResouces     IosConfigResouces     `json:"iosConfigResouces"`
		TodayGotWaterGoalTask TodayGotWaterGoalTask `json:"todayGotWaterGoalTask"`
		LeftUpResouces        LeftUpResouces        `json:"leftUpResouces"`
		MinSupportAPPVersion  string                `json:"minSupportAPPVersion"`
		LowFreqStatus         int                   `json:"lowFreqStatus"`
		FunCollectionHasLimit bool                  `json:"funCollectionHasLimit"`
		Message               interface{}           `json:"message"`
		TreeState             int                   `json:"treeState"`
		RightDownResouces     RightDownResouces     `json:"rightDownResouces"`
		IconFirstPurchaseInit bool                  `json:"iconFirstPurchaseInit"`
		ToFlowEnergy          int                   `json:"toFlowEnergy"`
		FarmUserPro           FarmUserPro           `json:"farmUserPro"`
		RetainPopupLimit      int                   `json:"retainPopupLimit"`
		ToBeginEnergy         int                   `json:"toBeginEnergy"`
		LeftDownResouces      LeftDownResouces      `json:"leftDownResouces"`
		EnableSign            bool                  `json:"enableSign"`
		LoadFriend            LoadFriend            `json:"loadFriend"`
		HadCompleteXgTask     bool                  `json:"hadCompleteXgTask"`
		OldUserIntervalTimes  []int                 `json:"oldUserIntervalTimes"`
		ToFruitTimes          int                   `json:"toFruitTimes"`
		OldUserSendWater      []string              `json:"oldUserSendWater"`
	}
	a := AutoGenerated{}
	req := httplib.Post(`https://api.m.jd.com/client.action?functionId=initForFarm`)
	req.Header("accept", "*/*")
	req.Header("accept-encoding", "gzip, deflate, br")
	req.Header("accept-language", "zh-CN,zh;q=0.9")
	req.Header("cache-control", "no-cache")
	req.Header("cookie", cookie)
	req.Header("origin", "https://home.m.jd.com")
	req.Header("pragma", "no-cache")
	req.Header("referer", "https://home.m.jd.com/myJd/newhome.action")
	req.Header("sec-fetch-dest", "empty")
	req.Header("sec-fetch-mode", "cors")
	req.Header("sec-fetch-site", "same-site")
	req.Header("User-Agent", ua())
	req.Header("Content-Type", "application/x-www-form-urlencoded")

	req.Body(`body={"version":4}&appid=wh5&clientVersion=9.1.0`)
	data, _ := req.Bytes()
	json.Unmarshal(data, &a)
	pt_pin := core.FetchCookieValue("pt_pin", cookie)
	rt := a.FarmUserPro.Name
	not := ""
	if rt == "" {
		rt = "数据异常"
	} else {
		if a.TreeState == 2 || a.TreeState == 3 {
			rt += "已可领取⏰"
			not = rt
		} else if a.TreeState == 1 {
			rt += fmt.Sprintf("种植中，进度%.2f%%🍒", 100*float64(a.FarmUserPro.TreeEnergy)/float64(a.FarmUserPro.TreeTotalEnergy))
		} else if a.TreeState == 0 {
			rt = "您忘了种植新的水果⏰"
			not = rt
		}
	}
	if state != nil {
		state <- rt
	} else if not != "" {
		a叉哦叉哦(pt_pin, "东东农场", not)
	}
}

func initPetTown(cookie string, state chan string) {
	type ResourceList struct {
		AdvertID string `json:"advertId"`
		ImageURL string `json:"imageUrl"`
		Link     string `json:"link"`
		ShopID   string `json:"shopId"`
	}
	type PetPlaceInfoList struct {
		Place  int `json:"place"`
		Energy int `json:"energy"`
	}
	type PetInfo struct {
		AdvertID     string `json:"advertId"`
		NickName     string `json:"nickName"`
		IconURL      string `json:"iconUrl"`
		ClickIconURL string `json:"clickIconUrl"`
		FeedGifURL   string `json:"feedGifUrl"`
		HomePetImage string `json:"homePetImage"`
		CrossBallURL string `json:"crossBallUrl"`
		RunURL       string `json:"runUrl"`
		TickleURL    string `json:"tickleUrl"`
	}
	type GoodsInfo struct {
		GoodsName        string `json:"goodsName"`
		GoodsURL         string `json:"goodsUrl"`
		GoodsID          string `json:"goodsId"`
		ExchangeMedalNum int    `json:"exchangeMedalNum"`
		ActivityID       string `json:"activityId"`
		ActivityIds      string `json:"activityIds"`
	}
	type Result struct {
		ShareCode              string             `json:"shareCode"`
		HisHbFlag              bool               `json:"hisHbFlag"`
		MasterHelpPeoples      []interface{}      `json:"masterHelpPeoples"`
		HelpSwitchOn           bool               `json:"helpSwitchOn"`
		UserStatus             int                `json:"userStatus"`
		TotalEnergy            int                `json:"totalEnergy"`
		MasterInvitePeoples    []interface{}      `json:"masterInvitePeoples"`
		ShareTo                string             `json:"shareTo"`
		PetSportStatus         int                `json:"petSportStatus"`
		UserImage              string             `json:"userImage"`
		MasterHelpReward       int                `json:"masterHelpReward"`
		ShowHongBaoExchangePop bool               `json:"showHongBaoExchangePop"`
		ShowNeedCollectPop     bool               `json:"showNeedCollectPop"`
		PetSportReward         string             `json:"petSportReward"`
		NewhandBubble          bool               `json:"newhandBubble"`
		ResourceList           []ResourceList     `json:"resourceList"`
		ProjectBubble          bool               `json:"projectBubble"`
		MasterInvitePop        bool               `json:"masterInvitePop"`
		MasterInviteReward     int                `json:"masterInviteReward"`
		MedalNum               int                `json:"medalNum"`
		MasterHelpPop          bool               `json:"masterHelpPop"`
		MeetDays               int                `json:"meetDays"`
		PetPlaceInfoList       []PetPlaceInfoList `json:"petPlaceInfoList"`
		MedalPercent           float64            `json:"medalPercent"`
		CharitableSwitchOn     bool               `json:"charitableSwitchOn"`
		PetInfo                PetInfo            `json:"petInfo"`
		NeedCollectEnergy      int                `json:"needCollectEnergy"`
		FoodAmount             int                `json:"foodAmount"`
		InviteCode             string             `json:"inviteCode"`
		RulesURL               string             `json:"rulesUrl"`
		PetStatus              int                `json:"petStatus"`
		GoodsInfo              GoodsInfo          `json:"goodsInfo"`
	}
	type AutoGenerated struct {
		Code       string `json:"code"`
		ResultCode string `json:"resultCode"`
		Message    string `json:"message"`
		Result     Result `json:"result"`
	}
	a := AutoGenerated{}
	req := httplib.Post(`https://api.m.jd.com/client.action?functionId=initPetTown`)
	req.Header("Host", "api.m.jd.com")
	req.Header("User-Agent", ua())
	req.Header("cookie", cookie)
	req.Header("Content-Type", "application/x-www-form-urlencoded")

	req.Body(`body={}&appid=wh5&loginWQBiz=pet-town&clientVersion=9.0.4`)
	data, _ := req.Bytes()
	json.Unmarshal(data, &a)
	rt := ""
	pt_pin := core.FetchCookieValue("pt_pin", cookie)
	not := ""
	if a.Code == "0" && a.ResultCode == "0" && a.Message == "success" {
		if a.Result.UserStatus == 0 {
			rt = "请手动开启活动⏰"
			not = rt

		} else if a.Result.GoodsInfo.GoodsName == "" {
			rt = "你忘了选购新的商品⏰"
			not = rt

		} else if a.Result.PetStatus == 5 {
			rt = a.Result.GoodsInfo.GoodsName + "已可领取⏰"
			not = rt

		} else if a.Result.PetStatus == 6 {
			rt = a.Result.GoodsInfo.GoodsName + "未继续领养新的物品⏰"
			not = rt
		} else {
			rt = a.Result.GoodsInfo.GoodsName + fmt.Sprintf("领养中，进度%.2f%%，勋章%d/%d🐶", a.Result.MedalPercent, a.Result.MedalNum, a.Result.GoodsInfo.ExchangeMedalNum)
		}
	} else {
		rt = "数据异常"
	}
	if state != nil {
		state <- rt
	} else if not != "" {
		a叉哦叉哦(pt_pin, "东东萌宠", not)
	}
}

func jsGold(cookie string, state chan int64) { //

	type BalanceVO struct {
		CashBalance       string `json:"cashBalance"`
		EstimatedAmount   string `json:"estimatedAmount"`
		ExchangeGold      string `json:"exchangeGold"`
		FormatGoldBalance string `json:"formatGoldBalance"`
		GoldBalance       int    `json:"goldBalance"`
	}
	type Gears struct {
		Amount         string `json:"amount"`
		ExchangeAmount string `json:"exchangeAmount"`
		Order          int    `json:"order"`
		Status         int    `json:"status"`
		Type           int    `json:"type"`
	}
	type Data struct {
		Advertise      string    `json:"advertise"`
		BalanceVO      BalanceVO `json:"balanceVO"`
		Gears          []Gears   `json:"gears"`
		IsGetCoupon    bool      `json:"isGetCoupon"`
		IsGetCouponEid bool      `json:"isGetCouponEid"`
		IsLogin        bool      `json:"isLogin"`
		NewPeople      bool      `json:"newPeople"`
	}
	type AutoGenerated struct {
		Code      int    `json:"code"`
		Data      Data   `json:"data"`
		IsSuccess bool   `json:"isSuccess"`
		Message   string `json:"message"`
		RequestID string `json:"requestId"`
	}
	a := AutoGenerated{}
	req := httplib.Post(`https://api.m.jd.com?functionId=MyAssetsService.execute&appid=market-task-h5`)
	req.Header("Accept", "application/json, text/plain, */*")
	req.Header("Accept-Encoding", "gzip, deflate, br")
	req.Header("Cookie", cookie)
	req.Header("Content-Type", "application/x-www-form-urlencoded")
	req.Header("Origin", "https://gold.jd.com")
	req.Header("Host", "api.m.jd.com")
	req.Header("Connection", "keep-alive")
	req.Header("User-Agent", ua())
	req.Header("Referer", "https://gold.jd.com/")

	req.Body(`functionId=MyAssetsService.execute&body={"method":"goldShopPage","data":{"channel":1}}&_t=` + fmt.Sprint(time.Now().Unix()) + `&appid=market-task-h5;`)
	data, _ := req.Bytes()
	json.Unmarshal(data, &a)
	if state != nil {
		state <- int64(a.Data.BalanceVO.GoldBalance)
	}
}

func jxncEgg(cookie string, state chan int64) {
	req := httplib.Get("https://m.jingxi.com/jxmc/queryservice/GetHomePageInfo?channel=7&sceneid=1001&activeid=null&activekey=null&isgift=1&isquerypicksite=1&_stk=activeid%2Cactivekey%2Cchannel%2Cisgift%2Cisquerypicksite%2Csceneid&_ste=1&h5st=20210818211830955%3B4408816258824161%3B10028%3Btk01w8db21b2130ny2eg0siAPpNQgBqjGzYfuG6IP7Z%2BAOB40BiqLQ%2Blglfi540AB%2FaQrTduHbnk61ngEeKn813gFeRD%3Bd9a0b833bf99a29ed726cbffa07ba955cc27d1ff7d2d55552878fc18fc667929&_=1629292710957&sceneval=2&g_login_type=1&g_ty=ls")
	req.Header("User-Agent", ua())
	req.Header("Host", "m.jingxi.com")
	req.Header("Accept", "*/*")
	req.Header("Connection", "keep-alive")
	req.Header("Accept-Language", "zh-cn")
	req.Header("Accept-Encoding", "gzip, deflate, br")
	req.Header("Referer", "https://st.jingxi.com/pingou/jxmc/index.html?nativeConfig=%7B%22immersion%22%3A1%2C%22toColor%22%3A%22%23e62e0f%22%7D&;__mcwvt=sjcp&ptag=7155.9.95")

	req.Header("Cookie", cookie)
	data, _ := req.Bytes()

	egg, _ := jsonparser.GetInt(data, "data", "eggcnt")
	state <- egg
}

func tytCoupon(cookie string, state chan string) {

	type DiscountInfo struct {
		High string        `json:"high"`
		Info []interface{} `json:"info"`
	}
	type ExtInfo struct {
		Num5              string `json:"5"`
		Num12             string `json:"12"`
		Num16             string `json:"16"`
		Num21             string `json:"21"`
		Num52             string `json:"52"`
		Num54             string `json:"54"`
		Num74             string `json:"74"`
		BusinessLabel     string `json:"business_label"`
		LimitOrganization string `json:"limit_organization"`
		UserLabel         string `json:"user_label"`
	}
	type Useable struct {
		AreaDesc         string        `json:"areaDesc"`
		AreaType         int           `json:"areaType"`
		Batchid          string        `json:"batchid"`
		BeanNumForPerson int           `json:"beanNumForPerson"`
		BeanNumForPlat   int           `json:"beanNumForPlat"`
		BeginTime        string        `json:"beginTime"`
		CanBeSell        bool          `json:"canBeSell"`
		CanBeShare       bool          `json:"canBeShare"`
		CompleteTime     string        `json:"completeTime"`
		CouponKind       int           `json:"couponKind"`
		CouponStyle      int           `json:"couponStyle"`
		CouponTitle      string        `json:"couponTitle"`
		Couponid         string        `json:"couponid"`
		Coupontype       int           `json:"coupontype"`
		CreateTime       string        `json:"createTime"`
		Discount         string        `json:"discount"`
		DiscountInfo     DiscountInfo  `json:"discountInfo"`
		EndTime          string        `json:"endTime"`
		ExpireType       int           `json:"expireType"`
		ExtInfo          ExtInfo       `json:"extInfo"`
		HourCoupon       int           `json:"hourCoupon"`
		IsOverlay        int           `json:"isOverlay"`
		LimitStr         string        `json:"limitStr"`
		LinkStr          string        `json:"linkStr"`
		OperateTime      string        `json:"operateTime"`
		OrderID          string        `json:"orderId"`
		OverlayDesc      string        `json:"overlayDesc"`
		PassKey          string        `json:"passKey"`
		Pin              string        `json:"pin"`
		PlatFormInfo     string        `json:"platFormInfo"`
		Platform         int           `json:"platform"`
		PlatformDetails  []interface{} `json:"platformDetails"`
		PwdKey           string        `json:"pwdKey"`
		Quota            string        `json:"quota"`
		SellID           string        `json:"sellId"`
		ShareID          string        `json:"shareId"`
		ShopID           string        `json:"shopId"`
		ShopName         string        `json:"shopName"`
		State            int           `json:"state"`
		UseTime          string        `json:"useTime"`
		VenderID         string        `json:"venderId"`
	}
	type Coupon struct {
		Curtimestamp           int       `json:"curtimestamp"`
		ExpiredCount           int       `json:"expired_count"`
		IsHideBaiTiaoInJxWxapp int       `json:"isHideBaiTiaoInJxWxapp"`
		IsHideMailInWxapp      int       `json:"isHideMailInWxapp"`
		Useable                []Useable `json:"useable"`
		UseableCount           int       `json:"useable_count"`
		UsedCount              int       `json:"used_count"`
	}
	type AutoGenerated struct {
		Coupon    Coupon `json:"coupon"`
		ErrMsg    string `json:"errMsg"`
		ErrorCode int    `json:"errorCode"`
		HasNext   int    `json:"hasNext"`
		Jdpin     string `json:"jdpin"`
		State     int    `json:"state"`
		Uin       string `json:"uin"`
	}
	a := AutoGenerated{}
	req := httplib.Get(`https://m.jingxi.com/activeapi/queryjdcouponlistwithfinance?state=1&wxadd=1&filterswitch=1&_=1629296270692&sceneval=2&g_login_type=1&callback=jsonpCBKB&g_ty=ls`)
	req.Header("Accept", "*/*")
	req.Header("Connection", "keep-alive")
	req.Header("Accept-Encoding", "gzip, deflate, br")
	req.Header("Cookie", cookie)
	req.Header("Content-Type", "application/x-www-form-urlencoded")
	req.Header("Host", "m.jingxi.com")

	req.Header("User-Agent", ua())
	req.Header("Referer", "https://st.jingxi.com/my/coupon/jx.shtml?sceneval=2&ptag=7155.1.18")

	data, _ := req.Bytes()
	res := regexp.MustCompile(`jsonpCBKB[(](.*)\s+[)];}catch`).FindSubmatch(data)
	rt := ""
	if len(res) > 0 {
		json.Unmarshal(res[1], &a)
		num := 0
		toexp := 0
		tm := int(time.Now().Unix() * 1000)
		for _, cp := range a.Coupon.Useable {
			if strings.Contains(cp.CouponTitle, "推推5.01") {
				num++
				if Int(cp.EndTime) < tm {
					toexp++
				}
			}
		}
		if num == 0 {
			rt = ""
		} else {
			rt = fmt.Sprintf("%d张5元优惠券", num)
			if toexp > 0 {
				rt += fmt.Sprintf("(今天将过期%d张)⏰", toexp)
			} else {
				rt += "🎰"
			}
		}
	}
	state <- rt
}

func mmCoin(cookie string, state chan int64) {
	req := httplib.Post(`https://api.m.jd.com/client.action`)
	req.Header("Host", "api.m.jd.com")
	req.Header("Accept", "application/json, text/plain, */*")
	req.Header("Origin", "https://h5.m.jd.com")

	req.Header("User-Agent", ua())
	req.Header("cookie", cookie)
	req.Header("Content-Type", "application/x-www-form-urlencoded")

	req.Body(`uuid=3245ad3d16ab2153c69f9ca91cd2e931b06a3bb8&clientVersion=10.1.0&client=wh5&osVersion=&area=&networkType=wifi&functionId=homePageV2&body=%7B%7D&appid=SecKill2020`)
	data, _ := req.Bytes()
	mmc, _ := jsonparser.GetInt(data, "result", "assignment", "assignmentPoints")
	state <- mmc
}

func jdzz(cookie string, state chan int64) { //
	req := httplib.Get(`https://api.m.jd.com/client.action?functionId=interactTaskIndex&body={}&client=wh5&clientVersion=9.1.0`)
	req.Header("Host", "api.m.jd.com")
	req.Header("Accept-Language", "zh-cn")
	req.Header("Accept-Encoding", "gzip, deflate, br")
	req.Header("Referer", "http://wq.jd.com/wxapp/pages/hd-interaction/index/index")
	req.Header("User-Agent", ua())
	req.Header("cookie", cookie)
	req.Header("Content-Type", "application/json")

	data, _ := req.Bytes()
	mmc, _ := jsonparser.GetString(data, "data", "totalNum")
	state <- int64(Int(mmc))
}

func (ck *JdCookie) Available() bool {
	if ck.PtKey == "" {
		return false
	}
	cookie := "pt_key=" + ck.PtKey + ";pt_pin=" + ck.PtPin + ";"
	if ck == nil {
		return true
	}
	req := httplib.Get("https://me-api.jd.com/user_new/info/GetJDUserInfoUnion")
	req.Header("Cookie", cookie)
	req.Header("Accept", "*/*")
	req.Header("Accept-Language", "zh-cn,")
	req.Header("Connection", "keep-alive,")
	req.Header("Referer", "https://home.m.jd.com/myJd/newhome.action?sceneval=2&ufc=&")
	req.Header("Host", "me-api.jd.com")
	req.Header("User-Agent", ua())

	data, err := req.Bytes()
	if err != nil {
		return av2(ck)
	}
	ui := &UserInfoResult{}
	if nil != json.Unmarshal(data, ui) {
		return av2(ck)
	}
	switch ui.Retcode {
	// case "1001": //ck.BeanNum
	// 	if ui.Msg == "not login" {
	// 		return false
	// 	}
	case "0":
		realPin := url.QueryEscape(ui.Data.UserInfo.BaseInfo.CurPin)
		if realPin != ck.PtPin {
			if realPin == "" {
				return av2(ck)
			} else {
				ck.PtPin = realPin
			}
		}
		if ui.Data.UserInfo.BaseInfo.Nickname != ck.Nickname || ui.Data.AssetInfo.BeanNum != ck.BeanNum || ui.Data.UserInfo.BaseInfo.UserLevel != ck.UserLevel || ui.Data.UserInfo.BaseInfo.LevelName != ck.LevelName {
			ck.UserLevel = ui.Data.UserInfo.BaseInfo.UserLevel
			ck.LevelName = ui.Data.UserInfo.BaseInfo.LevelName
			ck.Nickname = ui.Data.UserInfo.BaseInfo.Nickname
			ck.BeanNum = ui.Data.AssetInfo.BeanNum
		}
		return true
	}
	return av2(ck)
}

func av2(ck *JdCookie) bool {
	req := httplib.Get(`https://m.jingxi.com/user/info/GetJDUserBaseInfo?_=1629334995401&sceneval=2&g_login_type=1&g_ty=ls`)
	req.Header("User-Agent", ua())
	req.Header("Host", "m.jingxi.com")
	req.Header("Accept", "*/*")
	req.Header("Connection", "keep-alive")
	req.Header("Accept-Language", "zh-cn")
	req.Header("Accept-Encoding", "gzip, deflate, br")
	req.Header("Referer", "https://st.jingxi.com/my/userinfo.html?&ptag=7205.12.4")
	req.Header("Cookie", "pt_key="+ck.PtKey+";pt_pin="+ck.PtPin+";")

	data, err := req.Bytes()
	if err != nil {
		return true
	}
	ck.Nickname, _ = jsonparser.GetString(data, "nickname")
	return !strings.Contains(string(data), "login")
}

func av3(ck *JdCookie) bool {
	req := httplib.Get(`https://wq.jd.com/user_new/info/GetJDUserInfoUnion?sceneval=2`)
	req.Header("User-Agent", ua())
	req.Header("Host", "wq.jd.com")
	req.Header("Accept", "*/*")
	req.Header("Connection", "keep-alive")
	req.Header("Accept-Language", "zh-cn")
	req.Header("Accept-Encoding", "gzip, deflate, br")
	req.Header("Referer", "https://home.m.jd.com/myJd/newhome.action?sceneval=2&ufc=&")

	req.Header("Cookie", "pt_key="+ck.PtKey+";pt_pin="+ck.PtPin+";")
	data, err := req.Bytes()
	if err != nil {
		return av2(ck)
	}
	ck.Nickname, _ = jsonparser.GetString(data, "data", "userInfo", "baseInfo", "nickname")
	ck.BeanNum, _ = jsonparser.GetString(data, "data", "assetInfo", "beanNum")
	if ck.Nickname != "" {
		return true
	} else {
		return av2(ck)
	}
}

type UserInfoResult struct {
	Data struct {
		JdVvipCocoonInfo struct {
			JdVvipCocoon struct {
				DisplayType   int    `json:"displayType"`
				HitTypeList   []int  `json:"hitTypeList"`
				Link          string `json:"link"`
				Price         string `json:"price"`
				Qualification int    `json:"qualification"`
				SellingPoints string `json:"sellingPoints"`
			} `json:"JdVvipCocoon"`
			JdVvipCocoonStatus string `json:"JdVvipCocoonStatus"`
		} `json:"JdVvipCocoonInfo"`
		JdVvipInfo struct {
			JdVvipStatus string `json:"jdVvipStatus"`
		} `json:"JdVvipInfo"`
		AssetInfo struct {
			AccountBalance string `json:"accountBalance"`
			BaitiaoInfo    struct {
				AvailableLimit     string `json:"availableLimit"`
				BaiTiaoStatus      string `json:"baiTiaoStatus"`
				Bill               string `json:"bill"`
				BillOverStatus     string `json:"billOverStatus"`
				Outstanding7Amount string `json:"outstanding7Amount"`
				OverDueAmount      string `json:"overDueAmount"`
				OverDueCount       string `json:"overDueCount"`
				UnpaidForAll       string `json:"unpaidForAll"`
				UnpaidForMonth     string `json:"unpaidForMonth"`
			} `json:"baitiaoInfo"`
			BeanNum    string `json:"beanNum"`
			CouponNum  string `json:"couponNum"`
			CouponRed  string `json:"couponRed"`
			RedBalance string `json:"redBalance"`
		} `json:"assetInfo"`
		FavInfo struct {
			FavDpNum    string `json:"favDpNum"`
			FavGoodsNum string `json:"favGoodsNum"`
			FavShopNum  string `json:"favShopNum"`
			FootNum     string `json:"footNum"`
			IsGoodsRed  string `json:"isGoodsRed"`
			IsShopRed   string `json:"isShopRed"`
		} `json:"favInfo"`
		GrowHelperCoupon struct {
			AddDays     int     `json:"addDays"`
			BatchID     int     `json:"batchId"`
			CouponKind  int     `json:"couponKind"`
			CouponModel int     `json:"couponModel"`
			CouponStyle int     `json:"couponStyle"`
			CouponType  int     `json:"couponType"`
			Discount    float64 `json:"discount"`
			LimitType   int     `json:"limitType"`
			MsgType     int     `json:"msgType"`
			Quota       float64 `json:"quota"`
			RoleID      int     `json:"roleId"`
			State       int     `json:"state"`
			Status      int     `json:"status"`
		} `json:"growHelperCoupon"`
		KplInfo struct {
			KplInfoStatus string `json:"kplInfoStatus"`
			Mopenbp17     string `json:"mopenbp17"`
			Mopenbp22     string `json:"mopenbp22"`
		} `json:"kplInfo"`
		OrderInfo struct {
			CommentCount     string        `json:"commentCount"`
			Logistics        []interface{} `json:"logistics"`
			OrderCountStatus string        `json:"orderCountStatus"`
			ReceiveCount     string        `json:"receiveCount"`
			WaitPayCount     string        `json:"waitPayCount"`
		} `json:"orderInfo"`
		PlusPromotion struct {
			Status int `json:"status"`
		} `json:"plusPromotion"`
		UserInfo struct {
			BaseInfo struct {
				AccountType    string `json:"accountType"`
				BaseInfoStatus string `json:"baseInfoStatus"`
				CurPin         string `json:"curPin"`
				DefinePin      string `json:"definePin"`
				HeadImageURL   string `json:"headImageUrl"`
				LevelName      string `json:"levelName"`
				Nickname       string `json:"nickname"`
				Pinlist        string `json:"pinlist"`
				UserLevel      string `json:"userLevel"`
			} `json:"baseInfo"`
			IsHideNavi     string `json:"isHideNavi"`
			IsHomeWhite    string `json:"isHomeWhite"`
			IsJTH          string `json:"isJTH"`
			IsKaiPu        string `json:"isKaiPu"`
			IsPlusVip      string `json:"isPlusVip"`
			IsQQFans       string `json:"isQQFans"`
			IsRealNameAuth string `json:"isRealNameAuth"`
			IsWxFans       string `json:"isWxFans"`
			Jvalue         string `json:"jvalue"`
			OrderFlag      string `json:"orderFlag"`
			PlusInfo       struct {
			} `json:"plusInfo"`
			XbScore string `json:"xbScore"`
		} `json:"userInfo"`
		UserLifeCycle struct {
			IdentityID      string `json:"identityId"`
			LifeCycleStatus string `json:"lifeCycleStatus"`
			TrackID         string `json:"trackId"`
		} `json:"userLifeCycle"`
	} `json:"data"`
	Msg       string `json:"msg"`
	Retcode   string `json:"retcode"`
	Timestamp int64  `json:"timestamp"`
}

func FetchJdCookieValue(ps ...string) string {
	var key, cookies string
	if len(ps) == 2 {
		if len(ps[0]) > len(ps[1]) {
			key, cookies = ps[1], ps[0]
		} else {
			key, cookies = ps[0], ps[1]
		}
	}
	match := regexp.MustCompile(key + `=([^;]*);{0,1}`).FindStringSubmatch(cookies)
	if len(match) == 2 {
		return match[1]
	} else {
		return ""
	}
}

func GetTodayBean(ck *JdCookie, state chan int) {
	cookie := fmt.Sprintf("pt_key=%s;pt_pin=%s;", ck.PtKey, ck.PtPin)
	today := time.Now().Local().Format("2006-01-02")
	page := 1
	end := false
	in := 0
	defer func() {
		state <- in
	}()
	for {
		if end {
			return
		}
		bds := getJingBeanBalanceDetail(page, cookie)
		if bds == nil {
			break
		}
		for _, bd := range bds {
			amount := Int(bd.Amount)
			if strings.Contains(bd.Date, today) {
				if amount > 0 {
					in += amount
				} else {

				}
			} else {
				end = true
				break
			}
		}
		page++
	}
	return
}

func GetYestodayBean(ck *JdCookie, state chan int) {
	cookie := fmt.Sprintf("pt_key=%s;pt_pin=%s;", ck.PtKey, ck.PtPin)
	today := time.Now().Local().Format("2006-01-02")
	yestoday := time.Now().Local().Add(-time.Hour * 24).Format("2006-01-02")
	page := 1
	end := false
	in := 0
	defer func() {
		state <- in
	}()
	for {
		if end {
			return
		}
		bds := getJingBeanBalanceDetail(page, cookie)
		if bds == nil {
			break
		}
		for _, bd := range bds {
			amount := Int(bd.Amount)
			if strings.Contains(bd.Date, yestoday) {
				if amount > 0 {
					in += amount
				} else {

				}
			} else if strings.Contains(bd.Date, today) {

			} else {
				end = true
				break
			}
		}
		page++
	}
	return
}

type XBeanDetail struct {
	Amount      int    `json:"amount"`
	Createdate  string `json:"createdate"`
	Visibleinfo string `json:"visibleinfo"`
}

//欢迎叼毛来抄代码
func queryuserjingdoudetail(cookie string, e下水道 chan []int) {
	type AutoGenerated struct {
		Detail []XBeanDetail `json:"detail"`
		Ret    int           `json:"ret"`
		Retmsg string        `json:"retmsg"`
	}
	a := AutoGenerated{}
	req := httplib.Get(`https://m.jingxi.com/activeapi/queryuserjingdoudetail?pagesize=10&type=16&sceneval=2`)
	req.Header("User-Agent", "jdpingou;"+ua())
	req.Header("Host", "m.jingxi.com")
	req.Header("Accept", "*/*")
	req.Header("Connection", "keep-alive")
	req.Header("Accept-Language", "zh-cn")
	req.Header("Referer", "https://st.jingxi.com/sns/202105/31/xi_bean/index.html?nativeConfig=%7B%22immersion%22%3A%221%22%2C%20%22layoutUnderNavi%22%3A%220%22%2C%20%22showTitle%22%3A%221%22%2C%20%22toColor%22%3A%20%22%23f84b3b%22%7D")
	req.Header("Accept-Encoding", "gzip, deflate, br")
	req.Header("Cookie", cookie)

	data, err := req.Bytes()
	if err != nil {
		return
	}
	json.Unmarshal(data, &a)
	e叼毛 := []int{0, 0, 0, 0, 0}
	today := time.Now().Local().Format("2006/01/02")
	yestoday := time.Now().Local().Add(-time.Hour * 24).Format("2006/01/02")
	for _, v := range a.Detail {
		e叼毛[0] += v.Amount
		if strings.Contains(v.Createdate, today) {
			if v.Amount > 0 {
				e叼毛[1] += v.Amount
			} else {
				e叼毛[2] += -v.Amount
			}
		} else if strings.Contains(v.Createdate, yestoday) {
			if v.Amount > 0 {
				e叼毛[3] += v.Amount
			} else {
				e叼毛[4] += -v.Amount
			}
		}
	}
	e下水道 <- e叼毛 //叼毛去下水道
}

//欢迎叼毛来抄代码
func dream(cookie string, state chan string) {
	type AssistCondition struct {
		AssistConditionMsg    string `json:"assistConditionMsg"`
		AssistNumCurrent      int    `json:"assistNumCurrent"`
		AssistNumLimit        int    `json:"assistNumLimit"`
		AssistNumMax          int    `json:"assistNumMax"`
		AssistRemindKey       string `json:"assistRemindKey"`
		AssistRemindUser      string `json:"assistRemindUser"`
		CommodityAppLimitFlag int    `json:"commodityAppLimitFlag"`
		FactoryStatus         int    `json:"factoryStatus"`
		HireNumLimit          int    `json:"hireNumLimit"`
		ReAssistFlag          int    `json:"reAssistFlag"`
		SharePin              string `json:"sharePin"`
		SharePinHeadImage     string `json:"sharePinHeadImage"`
	}
	type AssistMaterialTuanCondition struct {
		AssistAppFlag    int           `json:"assistAppFlag"`
		AssistSelfFlag   int           `json:"assistSelfFlag"`
		CommodityList    []interface{} `json:"commodityList"`
		LimitTime        int           `json:"limitTime"`
		MaterialName     string        `json:"materialName"`
		MaterialPicture  string        `json:"materialPicture"`
		MaterialStatus   int           `json:"materialStatus"`
		OutOfStockFlag   int           `json:"outOfStockFlag"`
		RemindMsg        string        `json:"remindMsg"`
		SharePin         string        `json:"sharePin"`
		SharePinNickname string        `json:"sharePinNickname"`
		StartTime        int           `json:"startTime"`
		TuanID           string        `json:"tuanId"`
	}
	type DeviceList struct {
		CreateTime  int `json:"createTime"`
		DeviceDimID int `json:"deviceDimId"`
		DeviceID    int `json:"deviceId"`
		FactoryID   int `json:"factoryId"`
		UpdateTime  int `json:"updateTime"`
	}
	type FactoryList struct {
		CreateTime int    `json:"createTime"`
		FactoryID  int    `json:"factoryId"`
		Name       string `json:"name"`
		UpdateTime int    `json:"updateTime"`
	}
	type NewFactoryFlower struct {
		FactoryFlowerSendFlag int `json:"factoryFlowerSendFlag"`
		SendElectric          int `json:"sendElectric"`
	}
	type PickSiteInfo struct {
		Address         string `json:"address"`
		CityID          int    `json:"cityId"`
		CityName        string `json:"cityName"`
		CountryID       int    `json:"countryId"`
		CountryName     string `json:"countryName"`
		DcID            int    `json:"dcId"`
		ProvinceID      int    `json:"provinceId"`
		ProvinceName    string `json:"provinceName"`
		Sid             int    `json:"sid"`
		SiteID          string `json:"siteId"`
		SiteName        string `json:"siteName"`
		SiteURL         string `json:"siteUrl"`
		ToastChangeSite bool   `json:"toastChangeSite"`
		TownID          int    `json:"townId"`
		TownName        string `json:"townName"`
		Weight          int    `json:"weight"`
	}
	type ProductionList struct {
		BeginTime        int   `json:"beginTime"`
		CommodityDimID   int   `json:"commodityDimId"`
		CreateTime       int   `json:"createTime"`
		DataMark         int   `json:"dataMark"`
		DeviceID         int   `json:"deviceId"`
		EndTime          int   `json:"endTime"`
		ExchangeStatus   int   `json:"exchangeStatus"`
		FactoryID        int   `json:"factoryId"`
		InvestedElectric int   `json:"investedElectric"`
		NeedElectric     int   `json:"needElectric"`
		ProductionID     int64 `json:"productionId"`
		Status           int   `json:"status"`
		UpdateTime       int   `json:"updateTime"`
	}
	type ProductionStage struct {
		IsReachEnd                 int    `json:"isReachEnd"`
		ProductionStageAwardStatus int    `json:"productionStageAwardStatus"`
		ProductionStageProgress    string `json:"productionStageProgress"`
	}
	type Speciality struct {
		FactoryFlowerQualification int `json:"factoryFlowerQualification"`
		FactoryFlowerStatus        int `json:"factoryFlowerStatus"`
		SkinQualification          int `json:"skinQualification"`
		SkinStatus                 int `json:"skinStatus"`
	}
	type User struct {
		CreateTime                int    `json:"createTime"`
		CurrentLevel              int    `json:"currentLevel"`
		DataMark                  int    `json:"dataMark"`
		DeviceID                  string `json:"deviceId"`
		Electric                  int    `json:"electric"`
		EncryptPin                string `json:"encryptPin"`
		HeadImage                 string `json:"headImage"`
		HongBaoValue              string `json:"hongBaoValue"`
		IsJXNewUser               int    `json:"isJXNewUser"`
		IsProductSpecialCommodity int    `json:"isProductSpecialCommodity"`
		MosaicPin                 string `json:"mosaicPin"`
		NewPlayerWelfareFlag      int    `json:"newPlayerWelfareFlag"`
		NextLevelPercent          int    `json:"nextLevelPercent"`
		Nickname                  string `json:"nickname"`
		NpcStep                   int    `json:"npcStep"`
		Pin                       string `json:"pin"`
		ShareQywx                 string `json:"shareQywx"`
		UpdateTime                int    `json:"updateTime"`
		UserIdentity              string `json:"userIdentity"`
		Xid                       string `json:"xid"`
		Zone                      string `json:"zone"`
	}
	type UserAttrExtInfo struct {
		Electric              int `json:"electric"`
		InvestElectricLimDays int `json:"investElectricLimDays"`
		LastProduceInvestTime int `json:"lastProduceInvestTime"`
		ProductLimFlag        int `json:"productLimFlag"`
		RewardType            int `json:"rewardType"`
		UserType              int `json:"userType"`
	}
	type Data struct {
		AssistCondition             AssistCondition             `json:"assistCondition"`
		AssistMaterialTuanCondition AssistMaterialTuanCondition `json:"assistMaterialTuanCondition"`
		DeviceList                  []DeviceList                `json:"deviceList"`
		FactoryList                 []FactoryList               `json:"factoryList"`
		NeedSelectPickSite          int                         `json:"needSelectPickSite"`
		NewFactoryFlower            NewFactoryFlower            `json:"newFactoryFlower"`
		PickSiteInfo                PickSiteInfo                `json:"pickSiteInfo"`
		ProductionList              []ProductionList            `json:"productionList"`
		ProductionStage             ProductionStage             `json:"productionStage"`
		Speciality                  Speciality                  `json:"speciality"`
		SystemVersion               string                      `json:"systemVersion"`
		User                        User                        `json:"user"`
		UserAttrExtInfo             UserAttrExtInfo             `json:"userAttrExtInfo"`
	}
	type AutoGenerated struct {
		Data    Data   `json:"data"`
		Msg     string `json:"msg"`
		NowTime int    `json:"nowTime"`
		Ret     int    `json:"ret"`
	}
	url := "https://m.jingxi.com/dreamfactory/userinfo/GetUserInfo?zone=dream_factory&pin=&sharePin=&shareType=&materialTuanPin=&materialTuanId=&needPickSiteInfo=1&source=&_time=1637631683565&_ts=1637631683565&timeStamp=&_stk=_time,_ts,materialTuanId,materialTuanPin,needPickSiteInfo,pin,sharePin,shareType,source,timeStamp,zone&_ste=1&_=1637631683575&sceneval=2&g_login_type=1&g_ty=ls"

	req := httplib.Get(url)
	req.Header("Host", "m.jingxi.com")
	req.Header("Accept", "*/*")
	req.Header("Connection", "keep-alive")
	req.Header("Cookie", cookie)
	req.Header("User-Agent", "jdpingou;"+ua())
	req.Header("Accept-Language", "zh-cn")
	req.Header("Referer", "https://st.jingxi.com/pingou/dream_factory/index.html?ptag=7155.9.46")
	req.Header("Accept-Encoding", "gzip, deflate, br")

	data, _ := req.Bytes()
	a := &AutoGenerated{}
	json.Unmarshal(data, a)

	desc := ""
	not := true
	if state != nil {
		not = false
	}
	if len(a.Data.ProductionList) > 0 && len(a.Data.FactoryList) > 0 {
		var production = a.Data.ProductionList[0]
		if production.InvestedElectric >= production.NeedElectric {
			if production.ExchangeStatus == 1 {
				desc = "可以兑换商品了"
			}
			if production.ExchangeStatus == 3 {
				desc = "商品兑换已超时，请选择新商品进行制造"
			}
			// await exchangeProNotify()
		} else {
			not = false
			desc = fmt.Sprintf(`预计最快还需%d天生产完毕`, (production.NeedElectric-production.InvestedElectric)/(2*60*60*24))

		}
	} else {
		if len(a.Data.FactoryList) == 0 {
			desc = "请手动开启活动"
		} else if len(a.Data.ProductionList) == 0 {
			desc = "请手动选购商品进行生产"
		}
	}
	desc += "🏭"
	if state != nil {
		state <- desc
	}
	if not {
		a叉哦叉哦(core.FetchCookieValue("pt_pin", cookie), "京喜工厂", desc)
	}
}
