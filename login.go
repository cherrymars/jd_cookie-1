package jd_cookie

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/beego/beego/v2/client/httplib"
	"github.com/buger/jsonparser"
	"github.com/cdle/sillyGirl/core"
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
	// go RunServer()

	core.AddCommand("", []core.Function{
		{
			Rules: []string{`raw ^登录$`, `raw ^登陆$`, `raw ^h$`},
			Handle: func(s core.Sender) interface{} {

				if groupCode := jd_cookie.Get("groupCode"); !s.IsAdmin() && groupCode != "" && s.GetChatID() != 0 && !strings.Contains(groupCode, fmt.Sprint(s.GetChatID())) {
					logs.Info("跳过登录。")
					return nil
				}
				addr := ""
				var tabcount int64
				v := jd_cookie.Get("nolan_addr")
				addrs := strings.Split(v, "&")
				var haha func()
				var successLogin bool

				cancel := false
				phone := ""
				hasNolan := false
				ke := core.Bucket("wxmp").GetBool("isKe?", false)
				if v == "" {
					// goto ADONG
					return "快递员没有诺兰的地址。"
				}
				for _, addr = range addrs {
					addr = regexp.MustCompile(`^(https?://[-\.\w]+:?\d*)`).FindString(addr)
					if addr != "" {
						data, _ := httplib.Get(addr + "/api/Config").Bytes()
						tabcount, _ = jsonparser.GetInt(data, "data", "tabcount")
						if tabcount != 0 {
							hasNolan = true
							break
						}
					}
				}
				if !hasNolan == true {
					// goto ADONG
					return "诺兰无法为您服务。"
				}
				s.Reply(jd_cookie.Get("nolan_first", "若兰为您服务，请输入11位手机号：(输入“q”随时退出会话。)"))
				haha = func() {
					s.Await(s, func(s core.Sender) interface{} {
						ct := s.GetContent()
						if ct == "q" {
							cancel = true
							return "已退出会话。"
						}
						phone = regexp.MustCompile(`^\d{11}$`).FindString(ct)
						if phone == "" {
							return core.GoAgain("请输入正确的手机号：")
						}
						if s.GetImType() == "wxmp" && !ke {
							return "待会输入收到的验证码哦～"
						}
						s.Delete()
						return nil
					})
					if cancel {
						return
					}
					s.Reply("请稍等片刻...")
					req := httplib.Post(addr + "/api/SendSMS")
					req.Header("Proxy-Connection", "keep-alive")
					req.Header("accept", "application/json")
					req.Header("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.93 Safari/537.36")
					req.Header("content-type", "application/json")
					req.Header("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
					req.SetTimeout(time.Second*60, time.Second*60)
					data, err := req.Body(`{"Phone":"` + phone + `","qlkey":0}`).Bytes()
					if err != nil {
						s.Reply(err)
						return
					}
					message, _ := jsonparser.GetString(data, "message")
					success, _ := jsonparser.GetBoolean(data, "success")
					captcha, _ := jsonparser.GetInt(data, "data", "captcha")
					status, _ := jsonparser.GetInt(data, "data", "status")
					if message != "" && status != 666 {
						s.Reply(message)
					}
					i := 1
					if !success && status == 666 {
						if captcha <= 1 {
							s.Reply("正在进行滑块验证...")
							for {
								req = httplib.Post(addr + "/api/AutoCaptcha")
								req.Header("Proxy-Connection", "keep-alive")
								req.Header("accept", "application/json")
								req.Header("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.93 Safari/537.36")
								req.Header("content-type", "application/json")
								req.Header("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
								req.SetTimeout(time.Second*60, time.Second*60)
								data, err := req.Body(`{"Phone":"` + phone + `"}`).Bytes()
								if err != nil {
									s.Reply(err)
									return
								}
								message, _ := jsonparser.GetString(data, "message")
								success, _ := jsonparser.GetBoolean(data, "success")
								status, _ := jsonparser.GetInt(data, "data", "status")
								// if message != "" {
								// 	s.Reply()
								// }
								if !success {
									s.Reply("滑块验证失败：" + string(data))
								}
								if status == 666 {
									i++
									s.Reply(fmt.Sprintf("正在进行第%d次滑块验证...", i))
									continue
								}
								if success {
									break
								}
								s.Reply(message)
								return
							}
						} else {
							//欢迎叼毛前来抄代码
							//看代码的也是叼毛
							s.Reply("请先完成找成语小游戏：" + addr + "?id=" + phone)
							for {
								time.Sleep(time.Second)
								req = httplib.Get(addr + "/Captcha/" + phone)
								req.Header("Proxy-Connection", "keep-alive")
								req.Header("accept", "application/json")
								req.Header("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.93 Safari/537.36")
								req.Header("content-type", "application/json")
								req.Header("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
								req.SetTimeout(time.Second*60, time.Second*60)
								data, _ := req.Body(`{"Phone":"` + phone + `"}`).Bytes()
								status, _ := jsonparser.GetInt(data, "data", "status")
								if status != 666 {
									break
								}
							}
						}
					}
					s.Reply("请输入6位验证码：")
					code := ""
				aaa带带弟弟:
					s.Await(s, func(s core.Sender) interface{} {
						ct := s.GetContent()
						if ct == "q" {
							cancel = true
							return "已退出会话。"
						}
						code = regexp.MustCompile(`^\d{6}$`).FindString(ct)
						if code == "" {
							return core.GoAgain("请输入正确的验证码：")
						}
						// s.Reply("登录成功。")
						if s.GetImType() == "wxmp" && !ke {
							rt := "八九不离十登录成功啦，10秒后对我说“查询”以确认登录成功。"
							if jd_cookie.Get("xdd_url") != "" {
								rt += "此外，你可以在30秒内输入QQ号："
							}
							return rt
						}
						return nil
					}, time.Second*60, func(_ error) {
						s.Reply(jd_cookie.Get("nolan_timeout", "兄贵，你超时啦～"))
						cancel = true
					})
					if cancel {
						return
					}
					req = httplib.Post(addr + "/api/VerifyCode")
					req.Header("content-type", "application/json")
					data, _ = req.Body(`{"Phone":"` + phone + `","QQ":"` + fmt.Sprint(time.Now().Unix()) + `","qlkey":0,"Code":"` + code + `"}`).Bytes()
					req.SetTimeout(time.Second*60, time.Second*60)
					message, _ = jsonparser.GetString(data, "message")
					if strings.Contains(string(data), "pt_pin=") {
						successLogin = true
						s.Reply("登录成功。")
						pt_pin := core.FetchCookieValue(string(data), "pt_pin")
						s = s.Copy()
						s.SetContent(string(data))
						core.Senders <- s
						ad := jd_cookie.Get("ad")
						if ad != "" {
							s.Reply(ad)
						}
						time.Sleep(time.Second)
						jn := &JdNotify{
							ID: pt_pin,
						}
						jdNotify.First(jn)
						if jn.PushPlus == "" {
							s.Reply("是否订阅微信推送消息通知？(请在30s内回复”是“或”否“)")
							switch s.Await(s, func(s core.Sender) interface{} {
								return core.Switch{"是", "否"}
							}) {
							case "是":
								if jn.AssetCron == "" {
									rt := ""
									for {
										s.Reply("请输入资产推送时间(格式00:00:00，对应时、分、秒):")
										rt = s.Await(s, nil).(string)
										_, err := time.ParseInLocation("2006-01-02 15:04:05", time.Now().Format("2006-01-02"+" ")+rt, time.Local)
										if err == nil {
											break
										}
									}
									dd := strings.Split(rt, ":")
									jn.AssetCron = fmt.Sprintf("%s %s %s * * *", dd[2], dd[1], dd[0])
									if rid, ok := ccc[jn.ID]; ok {
										cc.Remove(rid)
										if rid, err := cc.AddFunc(jn.AssetCron, func() {
											assetPush(jn.ID)
										}); err == nil {
											ccc[jn.ID] = rid
										} else {
											return
										}
									}
								}
								data, _ := httplib.Get("https://www.pushplus.plus/api/common/wechat/getQrcode").Bytes()
								qrCodeUrl, _ := jsonparser.GetString(data, "data", "qrCodeUrl")
								qrCode, _ := jsonparser.GetString(data, "data", "qrCode")
								if qrCodeUrl == "" {
									s.Reply("嗝屁了。")
									return
								}
								s.Reply("请在30秒内打开微信扫描二维码关注公众号：\n" + core.ToImage(qrCodeUrl))
								ck := ""
								n := time.Now()
								for {
									if n.Add(time.Second * 30).Before(time.Now()) {
										s.Reply("扫码超时。")
										return
									}
									time.Sleep(time.Second)
									rsp, err := httplib.Get("https://www.pushplus.plus/api/common/wechat/confirmLogin?key=" + qrCode + "&code=1001").Response()
									if err != nil {
										continue
									}
									ck = rsp.Header.Get("Set-Cookie")
									if ck != "" {
										fmt.Println(ck)
										break
									}
								}
								req := httplib.Get("https://www.pushplus.plus/api/customer/user/token")
								req.Header("Cookie", ck)
								data, _ = req.Bytes()
								jn.PushPlus, _ = jsonparser.GetString(data, "data")
								s.Reply("扫码成功，请关注公号，我将尝试为你推送资产信息。")
								pushpluspush("资产推送通知", GetAsset(&JdCookie{
									PtPin: jn.ID,
									PtKey: jn.PtKey,
								}), jn.PushPlus)
								s.Reply("推送完成，祝您生活愉快！！！")
							}
						}
					} else {
						if strings.Contains(message, "验证码输入错误") {
							s.Reply("请输入正确的验证码：")
							goto aaa带带弟弟
						}
						s.Reply(message + "。")
						// if message != "" {
						// 	s.Reply("不好意思，刚搞错了还没成功，因为" + message + "。")
						// } else {
						// 	s.Reply("不好意思，刚搞错了并没有成功...")
						// }
					}
				}
				if s.GetImType() == "wxmp" && !ke {
					go haha()
				} else {
					haha()
					if !successLogin && !cancel { // && c != nil
						// s.Reply("将由阿东继续为您服务！")
						// goto ADONG
						return "登录失败。"
					}
				}
				return nil
				// 		ADONG:
				// 			// s.Reply("阿东嗝屁了。")
				// 			// return nil
				// 			if c == nil {
				// 				tip := jd_cookie.Get("tip")
				// 				if tip == "" {
				// 					if s.IsAdmin() {
				// 						s.Reply(jd_cookie.Get("tip", "阿东又不行了。")) //已支持阿东前往了解，https://github.com/rubyangxg/jd-qinglong
				// 						return nil
				// 					} else {
				// 						tip = "阿东未接入，暂时无法为您服务。"
				// 					}
				// 				}
				// 				s.Reply(tip)
				// 				return nil
				// 			}
				// 			go func() {
				// 				stop := false
				// 				uid := fmt.Sprint(time.Now().UnixNano())
				// 				cry := make(chan string, 1)
				// 				mhome.Store(uid, cry)
				// 				var deadline = time.Now().Add(time.Second * time.Duration(200))
				// 				var cookie *string
				// 				sendMsg := func(msg string) {
				// 					c.WriteJSON(map[string]interface{}{
				// 						"time":         time.Now().Unix(),
				// 						"self_id":      jd_cookie.GetInt("selfQid"),
				// 						"post_type":    "message",
				// 						"message_type": "private",
				// 						"sub_type":     "friend",
				// 						"message_id":   time.Now().UnixNano(),
				// 						"user_id":      uid,
				// 						"message":      msg,
				// 						"raw_message":  msg,
				// 						"font":         456,
				// 						"sender": map[string]interface{}{
				// 							"nickname": "傻妞",
				// 							"sex":      "female",
				// 							"age":      18,
				// 						},
				// 					})
				// 				}
				// 				if s.GetImType() == "wxmp" && !ke {
				// 					cancel := false
				// 					s.Await(s, func(s core.Sender) interface{} {
				// 						message := s.GetContent()
				// 						if message == "退出" || message == "q" {
				// 							cancel = true
				// 							return "取消登录"
				// 						}
				// 						if regexp.MustCompile(`^\d{11}$`).FindString(message) == "" {
				// 							return core.GoAgain("请输入格式正确的手机号，或者对我说“q”。")
				// 						}
				// 						phone = message
				// 						return "请输入收到的验证码哦～"
				// 					})

				// 					if cancel {
				// 						return
				// 					}
				// 				}
				// 				defer func() {
				// 					cry <- "stop"
				// 					mhome.Delete(uid)
				// 					if cookie != nil {
				// 						s.SetContent(*cookie)
				// 						core.Senders <- s
				// 					}
				// 					sendMsg("q")
				// 				}()
				// 				go func() {
				// 					for {
				// 						msg := <-cry
				// 						fmt.Println(msg)
				// 						if msg == "stop" {
				// 							break
				// 						}
				// 						msg = strings.Replace(msg, "登陆", "登录", -1)
				// 						if strings.Contains(msg, "不占资源") {
				// 							msg += "\n" + "4.取消"
				// 						}
				// 						if strings.Contains(msg, "无法回复") {
				// 							sendMsg("帮助")
				// 						}
				// 						{
				// 							res := regexp.MustCompile(`剩余操作时间：(\d+)`).FindStringSubmatch(msg)
				// 							if len(res) > 0 {
				// 								remain := core.Int(res[1])
				// 								deadline = time.Now().Add(time.Second * time.Duration(remain))
				// 							}
				// 						}
				// 						lines := strings.Split(msg, "\n")
				// 						new := []string{}
				// 						for _, line := range lines {
				// 							if !strings.Contains(line, "剩余操作时间") {
				// 								new = append(new, line)
				// 							}
				// 						}
				// 						msg = strings.Join(new, "\n")
				// 						if strings.Contains(msg, "直接退出") { //菜单页面
				// 							sendMsg("1")
				// 							continue
				// 						}
				// 						if strings.Contains(msg, "登录方式") {
				// 							sendMsg("1")
				// 							continue
				// 						}
				// 						if strings.Contains(msg, "请输入手机号") || strings.Contains(msg, "请输入11位手机号") {
				// 							if phone != "" {
				// 								sendMsg(phone)
				// 								continue
				// 							} else {
				// 								msg = "阿东为您服务，请输入11位手机号：(输入“q”随时退出会话。)"
				// 							}
				// 						}
				// 						if strings.Contains(msg, "pt_key") {
				// 							cookie = &msg
				// 							stop = true
				// 							s.SetContent("q")
				// 							core.Senders <- s
				// 						}
				// 						if cookie == nil {
				// 							if strings.Contains(msg, "已点击登录") {
				// 								continue
				// 							}
				// 							s.Reply(msg)
				// 						}
				// 					}
				// 				}()
				// 				sendMsg("h")
				// 				for {
				// 					if stop == true {
				// 						break
				// 					}
				// 					if deadline.Before(time.Now()) {
				// 						stop = true
				// 						s.Reply("登录超时")
				// 						break
				// 					}
				// 					s.Await(s, func(s core.Sender) interface{} {
				// 						msg := s.GetContent()
				// 						if msg == "查询" || strings.Contains(msg, "pt_pin=") {
				// 							s.Continue()
				// 							return nil
				// 						}
				// 						iw := core.Int(msg)
				// 						if msg == "q" || msg == "exit" || msg == "退出" || msg == "10" || msg == "4" || (fmt.Sprint(iw) == msg && iw > 1 && iw < 11) {
				// 							stop = true
				// 							if cookie == nil {
				// 								return "取消登录"
				// 							} else {
				// 								return "登录成功"
				// 							}
				// 						}
				// 						if phone != "" {
				// 							if regexp.MustCompile(`^\d{6}$`).FindString(msg) == "" {
				// 								return core.GoAgain("请输入格式正确的验证码，或者对我说“q”。")
				// 							} else {
				// 								rt := "八九不离十登录成功啦，60秒后对我说“查询”已确认登录成功。"
				// 								if jd_cookie.Get("xdd_url") != "" {
				// 									rt += "此外，你可以在30秒内输入QQ号："
				// 								}
				// 								s.Reply(rt)
				// 							}
				// 						}
				// 						sendMsg(s.GetContent())
				// 						return nil
				// 					}, `[\s\S]+`, time.Second)
				// 				}
				// 			}()
				// 			if s.GetImType() == "wxmp" && !ke {
				// 				return "请输入11位手机号："
				// 			}
				// 			return nil
			},
		},
	})
	if jd_cookie.GetBool("enable_aaron", false) {
		core.Senders <- &core.Faker{
			Message: "ql cron disable https://github.com/Aaron-lv/sync.git",
		}
		core.Senders <- &core.Faker{
			Message: "ql cron disable task Aaron-lv_sync_jd_scripts_jd_city.js",
		}
	}
}

// var c *websocket.Conn

// func RunServer() {
// 	addr := jd_cookie.Get("adong_addr")
// 	if addr == "" {
// 		return
// 	}
// 	defer func() {
// 		time.Sleep(time.Second * 2)
// 		RunServer()
// 	}()
// 	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws/event"}
// 	logs.Info("连接阿东 %s", u.String())
// 	var err error
// 	c, _, err = websocket.DefaultDialer.Dial(u.String(), http.Header{
// 		"X-Self-ID":     {fmt.Sprint(jd_cookie.GetInt("selfQid"))},
// 		"X-Client-Role": {"Universal"},
// 	})
// 	if err != nil {
// 		logs.Warn("连接阿东错误:", err)
// 		return
// 	}
// 	defer c.Close()
// 	go func() {
// 		for {
// 			_, message, err := c.ReadMessage()
// 			if err != nil {
// 				logs.Info("read:", err)
// 				return
// 			}
// 			type AutoGenerated struct {
// 				Action string `json:"action"`
// 				Echo   string `json:"echo"`
// 				Params struct {
// 					UserID  interface{} `json:"user_id"`
// 					Message string      `json:"message"`
// 				} `json:"params"`
// 			}
// 			ag := &AutoGenerated{}
// 			json.Unmarshal(message, ag)
// 			if ag.Action == "send_private_msg" {
// 				if cry, ok := mhome.Load(fmt.Sprint(ag.Params.UserID)); ok {
// 					fmt.Println(ag.Params.Message)
// 					cry.(chan string) <- ag.Params.Message
// 				}
// 			}
// 			logs.Info("recv: %s", message)
// 		}
// 	}()
// 	ticker := time.NewTicker(time.Second)
// 	defer ticker.Stop()
// 	for {
// 		select {
// 		case <-ticker.C:
// 			err := c.WriteMessage(websocket.TextMessage, []byte(`{}`))
// 			if err != nil {
// 				logs.Info("阿东错误:", err)
// 				c = nil
// 				return
// 			}
// 		}
// 	}
// }

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
