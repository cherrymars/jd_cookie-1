package jd_cookie

import (
	"fmt"
	"time"

	"github.com/beego/beego/v2/adapter/httplib"
	"github.com/cdle/sillyGirl/core"
	"github.com/cdle/sillyGirl/develop/qinglong"
)

func init() {
	core.AddCommand("jd", []core.Function{
		{
			Rules: []string{"travel help"},
			Admin: true,
			Handle: func(_ core.Sender) interface{} {
				envs, err := qinglong.GetEnvs("JD_COOKIE")
				if err != nil {
					return err
				}
				inviteId := "ZXASTT0225KkcRBgdpFaDJU6hx6QDJwFjRWn6u7zB55awQ"
				inviteIds := []string{}
				for _, env := range envs {
					req := httplib.Post("https://api.m.jd.com/client.action?functionId=travel_getHomeData")
					req.Header("Host", "api.m.jd.com")
					req.Header("Cookie", env.Value)
					req.Header("Content-Type", "application/x-www-form-urlencoded")
					req.Header("Origin", "https://wbbny.m.jd.com")
					req.Header("Accept-Encoding", "gzip, deflate, br")
					req.Header("Connection", "keep-alive")
					req.Header("Accept", "application/json, text/plain, */*")
					req.Body(`functionId=travel_getHomeData&body={"inviteId":"` + inviteId + `"}&client=wh5&clientVersion=1.0.0`)
					data, _ := req.String()
					fmt.Println(data)
					// msg, _ := jsonparser.GetString(data, "msg")
					// inviteId, _ := jsonparser.GetString(data, "data", "result", "homeMainInfo", "inviteId")
					// bizMsg, _ := jsonparser.GetString(data, "data", "bizMsg")
					if inviteId != "" {
						inviteIds = append(inviteIds, inviteId)
					}
					time.Sleep(time.Second)
				}

				return nil
			},
		},
	})

}
