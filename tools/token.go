package tools

import (
	"China_Telecom_Monitor/configs"
	"encoding/json"
)

type Token struct {
	ChinaTelecomToken string `json:"chinaTelecomToken"`
	LoginLastTime     int64  `json:"loginLastTime"`
}

var TokenFile = "/token.json"

func GetToken() *Token {
	token := Token{}
	jsonStr, err := ReadFile(configs.DataPath + TokenFile)
	if err != nil {
		configs.Logger.Error(err)
		return nil
	}
	if err := json.Unmarshal([]byte(jsonStr), &token); err != nil {
		configs.Logger.Error(err)
		return nil
	}
	return &token
}

func SetToken(chinaTelecomToken string, loginLastTime int64) {
	token := Token{
		ChinaTelecomToken: chinaTelecomToken,
		LoginLastTime:     loginLastTime,
	}
	jsonBytes, err := json.Marshal(token)
	if err != nil {
		configs.Logger.Error("SetToken json error", token, err)
		return
	}
	if err = WriteFile(configs.DataPath+TokenFile, string(jsonBytes)); err != nil {
		configs.Logger.Error("SetToken writeFile error", token, err)
		return
	}
	configs.Logger.Info("SetToken success. token", token)
}
