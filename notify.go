package jd_cookie

import (
	"fmt"
	"strings"
	"time"

	"github.com/cdle/sillyGirl/core"
	"github.com/cdle/sillyGirl/develop/qinglong"
)

type JdNotify struct {
	ID string
	// Nickname     string
	Pet          bool
	Fruit        bool
	DreamFactory bool
}

var jdNotify = core.Bucket("jdNotify")

var notTip = "\n\n通知没有用？请对我说“账号管理”，根据提示进行关闭。"

func initNotify() {
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
			Rules: []string{`raw ^账号管理$`},
			Handle: func(s core.Sender) interface{} {
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
					ask += fmt.Sprintf("%d. %s\n", i+1, accounts[i])
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

					if jn.Pet {
						ask += "2. 开启东东萌宠通知\n"
					} else {
						ask += "2. 关闭东东萌宠通知\n"
					}
					if jn.Fruit {
						ask += "3. 开启东东果园通知\n"
					} else {
						ask += "3. 关闭东东果园通知\n"
					}

					if jn.DreamFactory {
						ask += "4. 开启京喜工厂通知\n"
					} else {
						ask += "4. 关闭京喜工厂通知\n"
					}
					ask += "5. 解绑当前账号\n6. 退出当前会话"
					s.Reply(ask)
					rt := s.Await(s, func(s core.Sender) interface{} {
						return core.Range([]int{1, 6})
					}, time.Second*20)
					switch rt.(type) {
					case nil:
						return "超时，已退出会话。"
					case int:
						switch rt.(int) {
						case 1:
							return "请使用“查询”命令进行资产查询。"
						case 2:
							jn.Pet = !jn.Pet
						case 3:
							jn.Fruit = !jn.Fruit
						case 4:
							jn.DreamFactory = !jn.DreamFactory
						case 5:
							pin.Set(pt_pin, "")
							return "解绑成功。"
						case 6:
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
