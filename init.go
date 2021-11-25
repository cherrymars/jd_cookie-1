package jd_cookie

import (
	"os"
	"strings"

	"github.com/cdle/sillyGirl/core"
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
}
