package jd_cookie

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/httplib"
	"github.com/buger/jsonparser"
	"github.com/cdle/sillyGirl/core"
	"github.com/cdle/sillyGirl/develop/qinglong"
	cron "github.com/robfig/cron/v3"
)

type JdNotify struct {
	ID           string
	Pet          bool
	Fruit        bool
	DreamFactory bool
	Note         string
	PtKey        string
	AssetCron    string
	PushPlus     string
}

var cc *cron.Cron

var jdNotify = core.NewBucket("jdNotify")

func assetPush(pt_pin string) {
	jn := &JdNotify{
		ID: pt_pin,
	}
	jdNotify.First(jn)
	if jn.PushPlus != "" {
		pushpluspush("资产变动通知", GetAsset(&JdCookie{
			PtPin: pt_pin,
			PtKey: jn.PtKey,
		}), jn.PushPlus)
		return
	}
	qqGroup := jd_cookie.GetInt("qqGroup")
	if jn.PtKey != "" && pt_pin != "" {
		pt_key := jn.PtKey
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
		}
	}
}

func initNotify() {
	var ccc = map[string]cron.EntryID{}
	cc = cron.New(cron.WithSeconds())
	cc.Start()
	jdNotify.Foreach(func(_, v []byte) error {
		aa := &JdNotify{}
		json.Unmarshal(v, aa)
		if aa.AssetCron != "" {
			if rid, err := cc.AddFunc(aa.AssetCron, func() {
				assetPush(aa.ID)
			}); err == nil {
				ccc[aa.ID] = rid
			}
		}
		return nil
	})
	go func() {
		for {
			envs, _ := qinglong.GetEnvs("JD_COOKIE")
			for _, env := range envs {
				if env.Status != 0 {
					continue
				}
				pt_pin := core.FetchCookieValue(env.Value, "pt_pin")
				pt_key := core.FetchCookieValue(env.Value, "pt_key")
				if pt_pin != "" && pt_key != "" {
					jn := &JdNotify{
						ID: pt_pin,
					}
					jdNotify.First(jn)
					if jn.PtKey != pt_key {
						jn.PtKey = pt_key
						jdNotify.Create(jn)
					}
				}
			}
			time.Sleep(time.Hour)
		}
	}()
	core.AddCommand("", []core.Function{
		{
			Rules: []string{`raw ^任务通知$`},
			Cron:  jd_cookie.Get("task_Notify", "2 7,13,19 * * *"),
			Admin: true,
			Handle: func(_ core.Sender) interface{} {
				envs, _ := qinglong.GetEnvs("JD_COOKIE")
				for _, env := range envs {
					initPetTown(env.Value, nil)
					initFarm(env.Value, nil)
					dream(env.Value, nil)
				}
				return "推送完成"
			},
		},
		{
			Rules: []string{`raw ^关闭(.+)通知$`},
			Handle: func(s core.Sender) interface{} {
				class := s.Get()
				pin := pin(s.GetImType())
				uid := fmt.Sprint(s.GetUserID())
				accounts := []string{}
				pin.Foreach(func(k, v []byte) error {
					if string(v) == uid {
						accounts = append(accounts, string(k))
					}
					return nil
				})
				for i := range accounts {
					jn := &JdNotify{
						ID: accounts[i],
					}
					jdNotify.First(jn)
					if class == "京喜工厂" {
						jn.DreamFactory = true
					}
					if class == "东东农场" {
						jn.Fruit = true
					}
					if class == "东东萌宠" {
						jn.Pet = true
					}
					jdNotify.Create(jn)
				}
				return fmt.Sprintf("已为你关闭%d个账号的"+class+"通知。", len(accounts))
			},
		},
		{
			Rules: []string{`raw ^账号管理$`},
			Handle: func(s core.Sender) interface{} {
				if groupCode := jd_cookie.Get("groupCode"); !s.IsAdmin() && groupCode != "" && s.GetChatID() != 0 && !strings.Contains(groupCode, fmt.Sprint(s.GetChatID())) {
					s.Continue()
					return nil
				}
				pin := pin(s.GetImType())
				uid := fmt.Sprint(s.GetUserID())
				accounts := []string{}
				pin.Foreach(func(k, v []byte) error {
					if string(v) == uid {
						accounts = append(accounts, string(k))
					}
					return nil
				})
				num := len(accounts)
				if num == 0 {
					return "抱歉，你还没有绑定的账号呢~"
				}
				ask := fmt.Sprintf("请在20秒内从1~%d中选择你要操作的账号：\n", num)
				for i := range accounts {
					jn := &JdNotify{
						ID: accounts[i],
					}
					jdNotify.First(jn)
					note := ""
					if jn.Note != "" {
						note = jn.Note
					} else {
						note = jn.ID
					}
					ask += fmt.Sprintf("%d. %s\n", i+1, note)
				}
				s.Reply(strings.Trim(ask, "\n"))
				rt := s.Await(s, func(s core.Sender) interface{} {
					return core.Range([]int{1, num})
				}, time.Second*20)
				switch rt.(type) {
				case nil:
					return "超时，已退出会话。"
				case int:
					pt_pin := accounts[rt.(int)-1]
					jn := &JdNotify{
						ID: pt_pin,
					}
					jdNotify.First(jn)
					ask := "请在20秒内选择操作：\n1. 查询账号资产\n"

					if jn.Note == "" {
						ask += "2. 添加账户备注信息\n"
					} else {
						ask += "2. 修改账户备注信息\n"
					}
					if jn.Pet {
						ask += "3. 开启东东萌宠通知\n"
					} else {
						ask += "3. 关闭东东萌宠通知\n"
					}
					if jn.Fruit {
						ask += "4. 开启东东果园通知\n"
					} else {
						ask += "4. 关闭东东果园通知\n"
					}
					if jn.DreamFactory {
						ask += "5. 开启京喜工厂通知\n"
					} else {
						ask += "5. 关闭京喜工厂通知\n"
					}
					if jn.AssetCron == "" {
						ask += "6. 添加资产推送时间\n"
					} else {
						ask += "6. 修改资产推送时间\n"
					}
					ask += "7. 解绑当前账号\n8. 设置微信push+通知(推荐)\n9. 退出当前会话"
					s.Reply(ask)
					rt := s.Await(s, func(s core.Sender) interface{} {
						return core.Range([]int{1, 9})
					}, time.Second*20)
					switch rt.(type) {
					case nil:
						return "超时，已退出会话。"
					case int:
						switch rt.(int) {
						case 1:
							if jn.PtKey == "" {
								return "账号已过期，暂时无法查询。"
							}
							return GetAsset(&JdCookie{
								PtPin: pt_pin,
								PtKey: jn.PtKey,
							})
						case 2:
							s.Reply("请输入新的账号备注信息：")
							jn.Note = s.Await(s, nil).(string)
						case 3:
							jn.Pet = !jn.Pet
						case 4:
							jn.Fruit = !jn.Fruit
						case 5:
							jn.DreamFactory = !jn.DreamFactory
						case 6:
							s.Reply("请输入资产推送时间(格式00:00:00，对应时、分、秒):")
							rt := s.Await(s, nil).(string)
							_, err := time.ParseInLocation("2006-01-02 15:04:05", time.Now().Format("2006-01-02"+" ")+rt, time.Local)
							if err != nil {
								s.Reply("格式错误，已退出会话。")
								return nil
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
									return err
								}
							}
						case 7:
							pin.Set(pt_pin, "")
							return "解绑成功，会话结束。"
						case 8: //欢迎叼毛看内裤
							data, _ := httplib.Get("https://www.pushplus.plus/api/common/wechat/getQrcode").Bytes()
							qrCodeUrl, _ := jsonparser.GetString(data, "data", "qrCodeUrl")
							qrCode, _ := jsonparser.GetString(data, "data", "qrCode")
							if qrCodeUrl == "" {
								return "嗝屁了"
							}
							s.Reply("请在30秒内打开微信扫描二维码关注公众号：\n" + core.ToImage(qrCodeUrl))
							ck := ""
							n := time.Now()
							for {
								if n.Add(time.Second * 30).Before(time.Now()) {
									return "扫码超时。"
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
							s.Reply("扫码成功，将尝试为你推送资产信息。")
							pushpluspush("资产推送通知", GetAsset(&JdCookie{
								PtPin: jn.ID,
								PtKey: jn.PtKey,
							}), jn.PushPlus)
						case 9:
							return "已退出会话。"
						}
					}
					jdNotify.Create(jn)
					return "操作成功，会话结束。"
				}
				return nil
			},
		},
	})
}

func a叉哦叉哦(pt_pin, class, content string) {
	u := &JdNotify{
		ID: pt_pin,
	}
	jdNotify.First(u)
	if u.DreamFactory && class == "京喜工厂" {
		return
	}
	if u.Fruit && class == "东东农场" {
		return
	}
	if u.Pet && class == "东东萌宠" {
		return
	}
	if u.Note == "" {
		u.Note = u.ID
	}
	u.Note, _ = url.QueryUnescape(u.Note)
	if u.PushPlus != "" {
		pushpluspush(class+"通知", content+"\n\n通知没有用？请对登录机器人说“关闭"+class+"通知”或“账号管理”，根据提示进行关闭。", u.PushPlus)
		return
	}
	Notify(pt_pin, class+"通知("+u.Note+")：\n"+content+"\n\n通知没有用？请对我说“关闭"+class+"通知”或“账号管理”，根据提示进行关闭。")
}

func pushpluspush(title, content, token string) {
	req := httplib.Post("http://www.pushplus.plus/send")
	req.JSONBody(map[string]string{
		"token":    token,
		"title":    title,
		"content":  content,
		"template": "txt",
	})
	req.Response()
}
