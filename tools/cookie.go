package tools

import (
	"China_Telecom_Monitor/configs"
	"encoding/json"
)

type Cookie struct {
	ChinaTelecomCookie string `json:"chinaTelecomCookie"`
	LoginLastTime      int64  `json:"loginLastTime"`
}

var CookieFile = "/cookie.json"

func GetCookie() *Cookie {
	cookie := Cookie{}
	jsonStr, err := ReadFile(configs.DataPath + CookieFile)
	if err != nil {
		configs.Logger.Error(err)
		return &cookie
	}
	if err := json.Unmarshal([]byte(jsonStr), &cookie); err != nil {
		configs.Logger.Error(err)
		return nil
	}
	return &cookie
}

func SetCookie(chinaTelecomCookie string, loginLastTime int64) {
	cookie := Cookie{
		ChinaTelecomCookie: chinaTelecomCookie,
		LoginLastTime:      loginLastTime,
	}
	jsonBytes, err := json.Marshal(cookie)
	if err != nil {
		configs.Logger.Error("SetCookie json error", cookie, err)
		return
	}
	if err = WriteFile(configs.DataPath+CookieFile, string(jsonBytes)); err != nil {
		configs.Logger.Error("SetCookie writeFile error", cookie, err)
		return
	}
}
