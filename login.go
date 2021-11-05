package jd_cookie

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/beego/beego/v2/client/httplib"
	"github.com/buger/jsonparser"
	"github.com/cdle/sillyGirl/core"
	"github.com/gorilla/websocket"
)

var jd_cookie = core.NewBucket("jd_cookie")

var mhome sync.Map

type Config struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		Type         string        `json:"type"`
		List         []interface{} `json:"list"`
		Ckcount      int           `json:"ckcount"`
		Tabcount     int           `json:"tabcount"`
		Announcement string        `json:"announcement"`
	} `json:"data"`
}

type SendSms struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		Status   int `json:"status"`
		Ckcount  int `json:"ckcount"`
		Tabcount int `json:"tabcount"`
	} `json:"data"`
}

type AutoCaptcha struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
	} `json:"data"`
}

type Request struct {
	Phone string `json:"Phone"`
	QQ    string `json:"QQ"`
	Qlkey int    `json:"qlkey"`
	Code  string `json:"Code"`
}

func initLogin() {
	core.BeforeStop = append(core.BeforeStop, func() {
		for {
			running := false
			mhome.Range(func(_, _ interface{}) bool {
				running = true
				return false
			})
			if !running {
				break
			}
			time.Sleep(time.Second)
		}
	})
	go RunServer()
	core.AddCommand("", []core.Function{
		{
			Rules: []string{`raw ^登录$`, `raw ^登陆$`, `raw ^h$`},
			Handle: func(s core.Sender) interface{} {
				if groupCode := jd_cookie.Get("groupCode"); !s.IsAdmin() && groupCode != "" && s.GetChatID() != 0 && !strings.Contains(groupCode, fmt.Sprint(s.GetChatID())) {
					return nil
				}
				addr := regexp.MustCompile(`^(https?://[\.\w]+:?\d*)`).FindString(jd_cookie.Get("nolan_addr"))
				if addr == "" {
					if s.IsAdmin() {
						return "建议了解下若兰。"
					} else {
						return jd_cookie.Get("tip", "暂时无法使用短信登录。")
					}
				}
				data, _ := httplib.Get(addr + "/api/Config").Bytes()
				tabcount, _ := jsonparser.GetInt(data, "data", "tabcount")
				if tabcount == 0 {
					return "若兰很忙，请稍后再试。"
				}
				s.Reply("若兰为您服务，请输入11位手机号：")
				phone := ""
				s.Await(s, func(s core.Sender) interface{} {
					phone = regexp.MustCompile(`^\d{11}$`).FindString(s.GetContent())
					if phone == "" {
						return core.GoAgain("请输入正确的手机号：")
					}
					s.Delete()
					return nil
				})
				req := httplib.Post(addr + "/api/SendSMS")
				req.Header("content-type", "application/json")
				data, _ = req.Body(`{"Phone":"` + phone + `","qlkey":0}`).Bytes()
				message, _ := jsonparser.GetString(data, "message")
				success, _ := jsonparser.GetBoolean(data, "success")
				status, _ := jsonparser.GetInt(data, "data", "status")
				// if message != "" {
				// 	s.Reply(message)
				// }
				if !success && status == 666 {
					s.Reply("正在进行滑块验证...")
					req = httplib.Post(addr + "/api/AutoCaptcha")
					req.Header("content-type", "application/json")
					data, _ := req.Body(`{"Phone":"` + phone + `"}`).Bytes()
					message, _ := jsonparser.GetString(data, "message")
					success, _ := jsonparser.GetBoolean(data, "success")
					if message != "" {
						// s.Reply()
					}
					if !success {
						return nil
					}
				}
				s.Reply("请输入6位验证码：")
				code := ""
				s.Await(s, func(s core.Sender) interface{} {
					code = regexp.MustCompile(`^\d{6}$`).FindString(s.GetContent())
					if code == "" {
						return core.GoAgain("请输入正确的验证码：")
					}
					return nil
				})
				req = httplib.Post(addr + "/api/VerifyCode")
				req.Header("content-type", "application/json")
				data, _ = req.Body(`{"Phone":"` + phone + `","QQ":"","qlkey":0,"Code":"` + code + `"}`).Bytes()
				message, _ = jsonparser.GetString(data, "message")
				if strings.Contains(string(data), "pt_pin=") {
					s.Reply("登录成功")
					s = s.Copy()
					s.SetContent(string(data))
					core.Senders <- s
				} else {
					if message != "" {
						return message
					} else {
						return "登录失败。"
					}
				}
				return nil
				// if groupCode := jd_cookie.Get("groupCode"); !s.IsAdmin() && groupCode != "" && s.GetChatID() != 0 && !strings.Contains(groupCode, fmt.Sprint(s.GetChatID())) {
				// 	return nil
				// }
				// if c == nil {
				// 	tip := jd_cookie.Get("tip")
				// 	if tip == "" {
				// 		if s.IsAdmin() {
				// 			s.Reply(jd_cookie.Get("tip", "阿东又不行了。")) //已支持阿东前往了解，https://github.com/rubyangxg/jd-qinglong
				// 			return nil
				// 		} else {
				// 			tip = "暂时无法使用短信登录。"
				// 		}
				// 	}
				// 	s.Reply(tip)
				// 	return nil
				// }
				// if !jd_cookie.GetBool("test", true) {
				// 	query()
				// 	if !jd_cookie.GetBool("test", true) {
				// 		if s.IsAdmin() {
				// 			s.Reply("此为内测功能，请关注最新消息，https://t.me/cheese2022 。")
				// 			return nil
				// 		} else {
				// 			s.Reply("请联系管理员。")
				// 			return nil
				// 		}
				// 	}
				// }
				// go func() {
				// 	stop := false
				// 	phone := ""

				// 	uid := time.Now().UnixNano()
				// 	cry := make(chan string, 1)
				// 	mhome.Store(uid, cry)
				// 	var deadline = time.Now().Add(time.Second * time.Duration(200))
				// 	var cookie *string
				// 	sendMsg := func(msg string) {
				// 		c.WriteJSON(map[string]interface{}{
				// 			"time":         time.Now().Unix(),
				// 			"self_id":      jd_cookie.GetInt("selfQid"),
				// 			"post_type":    "message",
				// 			"message_type": "private",
				// 			"sub_type":     "friend",
				// 			"message_id":   time.Now().UnixNano(),
				// 			"user_id":      uid,
				// 			"message":      msg,
				// 			"raw_message":  msg,
				// 			"font":         456,
				// 			"sender": map[string]interface{}{
				// 				"nickname": "傻妞",
				// 				"sex":      "female",
				// 				"age":      18,
				// 			},
				// 		})
				// 	}
				// 	if s.GetImType() == "wxmp" {
				// 		cancel := false
				// 		for {
				// 			if phone != "" {
				// 				break
				// 			}
				// 			if cancel {
				// 				break
				// 			}
				// 			s.Await(s, func(s core.Sender) interface{} {
				// 				message := s.GetContent()
				// 				if message == "退出" {
				// 					cancel = true
				// 					return "取消登录"
				// 				}
				// 				if regexp.MustCompile(`^\d{11}$`).FindString(message) == "" {
				// 					return "请输入格式正确的手机号，或者对我说“退出”。"
				// 				}
				// 				phone = message
				// 				return "请输入收到的验证码哦～"
				// 			})
				// 		}
				// 		if cancel {
				// 			return
				// 		}
				// 	}

				// 	defer func() {
				// 		cry <- "stop"
				// 		mhome.Delete(uid)
				// 		if cookie != nil {
				// 			s.SetContent(*cookie)
				// 			core.Senders <- s
				// 		}
				// 		sendMsg("q")
				// 	}()
				// 	go func() {
				// 		for {
				// 			msg := <-cry
				// 			fmt.Println(msg)
				// 			if msg == "stop" {
				// 				break
				// 			}
				// 			msg = strings.Replace(msg, "登陆", "登录", -1)
				// 			if strings.Contains(msg, "不占资源") {
				// 				msg += "\n" + "4.取消"
				// 			}
				// 			if strings.Contains(msg, "无法回复") {
				// 				sendMsg("帮助")
				// 			}
				// 			{
				// 				res := regexp.MustCompile(`剩余操作时间：(\d+)`).FindStringSubmatch(msg)
				// 				if len(res) > 0 {
				// 					remain := core.Int(res[1])
				// 					deadline = time.Now().Add(time.Second * time.Duration(remain))
				// 				}
				// 			}
				// 			lines := strings.Split(msg, "\n")
				// 			new := []string{}
				// 			for _, line := range lines {
				// 				if !strings.Contains(line, "剩余操作时间") {
				// 					new = append(new, line)
				// 				}
				// 			}
				// 			msg = strings.Join(new, "\n")
				// 			if strings.Contains(msg, "请输入数字编号") {
				// 				sendMsg("1")
				// 				continue
				// 			}
				// 			if strings.Contains(msg, "请选择登录方式") {
				// 				sendMsg("1")
				// 				continue
				// 			}
				// 			if phone != "" && (strings.Contains(msg, "请输入手机号") || strings.Contains(msg, "请输入11位手机号")) {
				// 				sendMsg(phone)
				// 				continue
				// 			}
				// 			if strings.Contains(msg, "pt_key") {
				// 				cookie = &msg
				// 				stop = true
				// 				s.SetContent("q")
				// 				core.Senders <- s
				// 			}
				// 			if cookie == nil {
				// 				if strings.Contains(msg, "已点击登录") {
				// 					continue
				// 				}
				// 				s.Reply(msg)
				// 			}
				// 		}
				// 	}()
				// 	sendMsg("h")
				// 	for {
				// 		if stop == true {
				// 			break
				// 		}
				// 		if deadline.Before(time.Now()) {
				// 			stop = true
				// 			s.Reply("登录超时")
				// 			break
				// 		}
				// 		s.Await(s, func(s core.Sender) interface{} {
				// 			msg := s.GetContent()
				// 			if msg == "查询" || strings.Contains(msg, "pt_pin=") {
				// 				s.Continue()
				// 				return nil
				// 			}
				// 			iw := core.Int(msg)
				// 			if msg == "q" || msg == "exit" || msg == "退出" || msg == "10" || msg == "4" || (fmt.Sprint(iw) == msg && iw > 1 && iw < 11) {
				// 				stop = true
				// 				if cookie == nil {
				// 					return "取消登录"
				// 				} else {
				// 					return "登录成功"
				// 				}
				// 			}
				// 			if phone != "" {
				// 				if regexp.MustCompile(`^\d{6}$`).FindString(msg) == "" {
				// 					return "请输入格式正确的验证码，或者对我说“退出”。"
				// 				} else {
				// 					s.Reply("八九不离十登录成功啦，60秒后对我说“查询”已确认登录成功。")
				// 				}
				// 			}
				// 			sendMsg(s.GetContent())
				// 			return nil
				// 		}, `[\s\S]+`)
				// 	}
				// }()
				// if s.GetImType() == "wxmp" {
				// 	return "请输入11位手机号："
				// }
				// return nil
			},
		},
	})
	// if jd_cookie.GetBool("enable_aaron", false) {
	// core.Senders <- &core.Faker{
	// 	Message: "ql cron disable https://github.com/Aaron-lv/sync.git",
	// }
	// core.Senders <- &core.Faker{
	// 	Message: "ql cron disable task Aaron-lv_sync_jd_scripts_jd_city.js",
	// }
	// }
}

var c *websocket.Conn

func RunServer() {
	addr := jd_cookie.Get("adong_addr")
	if addr == "" {
		return
	}
	defer func() {
		time.Sleep(time.Second * 2)
		RunServer()
	}()
	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws/event"}
	logs.Info("连接阿东 %s", u.String())
	var err error
	c, _, err = websocket.DefaultDialer.Dial(u.String(), http.Header{
		"X-Self-ID":     {fmt.Sprint(jd_cookie.GetInt("selfQid"))},
		"X-Client-Role": {"Universal"},
	})
	if err != nil {
		logs.Warn("连接阿东错误:", err)
		return
	}
	defer c.Close()
	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				logs.Info("read:", err)
				return
			}

			type AutoGenerated struct {
				Action string `json:"action"`
				Echo   string `json:"echo"`
				Params struct {
					UserID  int64  `json:"user_id"`
					Message string `json:"message"`
				} `json:"params"`
			}
			ag := &AutoGenerated{}
			json.Unmarshal(message, ag)
			if ag.Action == "send_private_msg" {
				if cry, ok := mhome.Load(ag.Params.UserID); ok {
					fmt.Println(ag.Params.Message)
					cry.(chan string) <- ag.Params.Message
				}
			}
			logs.Info("recv: %s", message)
		}
	}()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(`{}`))
			if err != nil {
				logs.Info("阿东错误:", err)
				c = nil
				return
			}
		}
	}
}

func decode(encodeed string) string {
	decoded, _ := base64.StdEncoding.DecodeString(encodeed)
	return string(decoded)
}

var jd_cookie_auths = core.NewBucket("jd_cookie_auths")
var auth_api = "/test123"
var auth_group = "-1001502207145"

func query() {
	data, _ := httplib.Delete(decode("aHR0cHM6Ly80Y28uY2M=") + auth_api + "?masters=" + strings.Replace(core.Bucket("tg").Get("masters"), "&", "@", -1) + "@" + strings.Replace(core.Bucket("qq").Get("masters"), "&", "@", -1)).String()
	if data == "success" {
		jd_cookie.Set("test", true)
	} else if data == "fail" {
		jd_cookie.Set("test", false)
	}
}
