package jd_cookie

func initEnEn() {
	// core.AddCommand("jd", []core.Function{

	// 	{
	// 		Rules: []string{"eueu ?"},
	// 		Admin: true,
	// 		Handle: func(s core.Sender) interface{} {
	// 			envs, err := qinglong.GetEnvs("JD_WSCK")
	// 			if err != nil {
	// 				return err
	// 			}
	// 			yes := false
	// 			for _, env := range envs {
	// 				if strings.Contains(env.Value, s.Get()) {
	// 					yes = true
	// 					pin := core.FetchCookieValue("pin", env.Value)
	// 					pt_key, err := getKey(env.Value)
	// 					if err != nil {
	// 						return err
	// 					}
	// 					s.Reply(fmt.Sprintf("pt_key=%s;pt_pin=%s;", pt_key, pin))
	// 				}
	// 			}
	// 			if !yes {
	// 				return "找不到转换目标"
	// 			}
	// 			return nil
	// 		},
	// 	},
	// })

}
