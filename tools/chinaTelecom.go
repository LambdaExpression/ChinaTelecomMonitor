package tools

import (
	"China_Telecom_Monitor/configs"
	"China_Telecom_Monitor/models"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/golang-module/carbon/v2"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func ChinaTelecomLogin(username, password string) bool {

	if !checkLogin() {
		return false
	}

	login(username, password)
	return true
}

func checkLogin() bool {
	token := GetToken()
	if token == nil {
		return true
	}
	// 上次获取token时间
	getTokenTime := carbon.CreateFromTimestamp(token.LoginLastTime)
	// 下次获取token时间 = 上次获取token时间 + 获取token间隔时间
	nextTimeGetTokenTime := getTokenTime.AddSeconds(configs.LoginIntervalTime)
	// 比较 下次获取token时间 是否大于 现在时间
	if nextTimeGetTokenTime.Gt(carbon.Now()) {
		configs.Logger.Error(strconv.Itoa(configs.LoginIntervalTime) + " 秒内最多登录一次，下次获取token时间为" + nextTimeGetTokenTime.ToDateTimeString() + "，避免被封号")
		return false
	}
	return true
}

func login(mobile, password string) {

	t := time.Now().Format("20060102150400")
	e := fmt.Sprintf("iPhone 14 13.2.3%s%s%s%s0$$$0.", mobile, mobile, t, password)

	enc, err := encrypt(e)
	if err != nil {
		configs.Logger.Error("error", err)
		return
	}

	rb := models.Request[models.LoginRequestContent]{
		Content: models.LoginRequestContent{
			FieldData: models.LoginRequestFieldData{
				AccountType:                "",
				Authentication:             password,
				DeviceUid:                  fmt.Sprintf("3%s", mobile),
				IsChinatelecom:             "0",
				LoginAuthCipherAsymmertric: enc,
				LoginType:                  "4",
				PhoneNum:                   transPhone(mobile),
				SystemVersion:              "13.2.3",
			},
			Attach: "iPhone",
		},
		HeaderInfos: initHeaderInfos(mobile, "userLoginNormal", t),
	}

	loginUrl := "https://appgologin.189.cn:9031/login/client/userLoginNormal"

	result, err := post[models.LoginRequestContent, models.LoginData](loginUrl, rb, mobile, password, false, true)
	if err != nil {
		configs.Logger.Error("login post error", err)
		return
	}

	if result.HeaderInfos.Code != "0000" {
		configs.Logger.Error("login request header failed. " + result.HeaderInfos.Code + " " + result.HeaderInfos.Reason)
		return
	}

	if result.ResponseData.ResultCode != "0000" {
		configs.Logger.Error("login request response failed. responseData " + result.ResponseData.ResultCode + " " + result.ResponseData.ResultDesc)
		return
	}

	token := result.ResponseData.Data.LoginSuccessResult.Token
	ti := time.Now()
	SetToken(token, ti.Unix())

}

func GetQryImportantData(mobile, password string) *models.Result[models.ImportantData] {
	t := time.Now().Format("20060102150400")
	rb := models.Request[models.QryImportantDataRequestContent]{
		Content: models.QryImportantDataRequestContent{
			FieldData: models.QryImportantDataRequestFieldData{
				ProvinceCode:   "600101",
				CityCode:       "8441900",
				ShopId:         "20002",
				IsChinatelecom: "0",
				Account:        transPhone(mobile),
			},
			Attach: "test",
		},
		HeaderInfos: initHeaderInfos(mobile, "qryImportantData", t),
	}

	requestUrl := "https://appfuwu.189.cn:9021/query/qryImportantData"

	result, err := post[models.QryImportantDataRequestContent, models.ImportantData](requestUrl, rb, mobile, password, true, false)
	if err != nil {
		configs.Logger.Error("login post error", err)
		return nil
	}

	return &result
}

func GetUserFluxPackage(mobile, password string) *models.Result[models.UserFluxPackageData] {
	t := time.Now().Format("20060102150400")
	rb := models.Request[models.UserFluxPackageRequestContent]{
		Content: models.UserFluxPackageRequestContent{
			FieldData: models.UserFluxPackageRequestFieldData{
				QueryFlag:  "0",
				AccessAuth: "1",
				Account:    transPhone(mobile),
			},
			Attach: "test",
		},
		HeaderInfos: initHeaderInfos(mobile, "userFluxPackage", t),
	}

	requestUrl := "https://appfuwu.189.cn:9021/query/userFluxPackage"

	result, err := post[models.UserFluxPackageRequestContent, models.UserFluxPackageData](requestUrl, rb, mobile, password, true, false)
	if err != nil {
		configs.Logger.Error("login post error", err)
		return nil
	}

	return &result
}

func initHeaderInfos(mobile, code, t string) models.RequestHeaderInfos {
	if t == "" {
		t = time.Now().Format("20060102150400")
	}
	headerInfos := models.RequestHeaderInfos{
		ClientType:     "#" + configs.ClientVersion + "#channel50#iPhone 14 Pro#",
		Timestamp:      t,
		Code:           code,
		ShopId:         "20002",
		Source:         "110003",
		SourcePassword: "Sid98s",
		UserLoginName:  mobile,
	}
	return headerInfos
}

func toBodyStr(s any) string {
	b, err := json.Marshal(s)
	if err != nil {
		configs.Logger.Error("body json error", err)
		return ""
	}
	return string(b)
}

func post[C, D any](requestUrl string, requestBody models.Request[C], mobile, password string, autoLogin, isLogin bool) (models.Result[D], error) {
	token := ""
	firstLogin := false
	if !isLogin {
		token, firstLogin = getTokenStr(mobile, password)
		if token == "" {
			return models.Result[D]{
				HeaderInfos: models.HeaderInfos{
					Code:   "T801",
					Reason: "token 为空",
				},
			}, nil
		}
	}
	requestBody.HeaderInfos.Token = token
	body := toBodyStr(requestBody)
	payload := strings.NewReader(body)
	req, _ := http.NewRequest("POST", requestUrl, payload)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("Connection", "Keep-Alive")
	req.Header.Add("Accept-Encoding", "gzip")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		configs.Logger.Error("request url: ", requestUrl, " body: ", body, err)
		return models.Result[D]{}, err
	}
	defer response.Body.Close()
	b, err := io.ReadAll(response.Body)
	if err != nil {
		configs.Logger.Error("request url: ", requestUrl, " body: ", body, err)
		return models.Result[D]{}, err
	}
	result := string(b)
	configs.Logger.Debug("request url: ", requestUrl, " body: ", body, " result: ", result)

	var resultData models.Result[D]
	err = json.Unmarshal(b, &resultData)
	if err != nil {
		configs.Logger.Error("login body to json error. body:"+result, err)
		return models.Result[D]{}, err
	}
	// token 过期
	if resultData.HeaderInfos.Code == "X201" && autoLogin && !firstLogin {
		ChinaTelecomLogin(mobile, password)
		return post[C, D](requestUrl, requestBody, mobile, password, false, false)
	}
	if resultData.HeaderInfos.Code != "0000" || resultData.ResponseData.ResultCode != "0000" {
		configs.Logger.Error("request url: ", requestUrl, " body: ", body, " result: ", result)
	}
	return resultData, nil
}

func getTokenStr(mobile, password string) (string, bool) {
	token := GetToken()
	first := false
	if token == nil {
		login(mobile, password)
		token = GetToken()
		first = true
		if token == nil {
			return "", first
		}
	}
	return token.ChinaTelecomToken, first
}

// 转换手机号码
func transPhone(mobile string) string {
	var tMobile string
	for _, mo := range mobile {
		charCode := int(mo + 2&65535)
		tMobile += string(rune(charCode))
	}
	return tMobile
}

// 加密
func encrypt(message string) (string, error) {
	// 公钥
	//	publicKeyPEM := `-----BEGIN PUBLIC KEY-----
	//MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDBkLT15ThVgz6/NOl6s8GNPofd
	//WzWbCkWnkaAm7O2LjkM1H7dMvzkiqdxU02jamGRHLX/ZNMCXHnPcW/sDhiFCBN18
	//qFvy8g6VYb9QtroI09e176s+ZCtiv7hbin2cCTj99iUpnEloZm19lwHyo69u5UMi
	//PMpq0/XKBO8lYhN/gwIDAQAB
	//-----END PUBLIC KEY-----`
	publicKeyPEM := `-----BEGIN PUBLIC KEY-----
	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC+ugG5A8cZ3FqUKDwM57GM4io6
	JGcStivT8UdGt67PEOihLZTw3P7371+N47PrmsCpnTRzbTgcupKtUv8ImZalYk65 
	dU8rjC/ridwhw9ffW2LBwvkEnDkkKKRi2liWIItDftJVBiWOh17o6gfbPoNrWORc
	Adcbpk2L+udld5kZNwIDAQAB
	-----END PUBLIC KEY-----`

	// 解码公钥
	block, err := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return "", fmt.Errorf("Error decoding public key", err)
	}
	publicKey, err2 := x509.ParsePKIXPublicKey(block.Bytes)
	if err2 != nil {
		return "", fmt.Errorf("Error parsing public key", err2)
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("Error casting public key to RSA public key")
	}

	// 对消息进行加密
	ciphertext, err3 := rsa.EncryptPKCS1v15(rand.Reader, rsaPublicKey, []byte(message))
	if err3 != nil {
		return "", err3
	}

	// 将加密后的数据进行 Base64 编码
	encodedCiphertext := base64.StdEncoding.EncodeToString(ciphertext)

	return encodedCiphertext, nil
}
