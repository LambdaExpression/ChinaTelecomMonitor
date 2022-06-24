package tools

import (
	"China_Telecom_Monitor/configs"
	"China_Telecom_Monitor/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
	"github.com/golang-module/carbon/v2"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func ChinaTelecomLogin(username, password string) {

	if !checkLogin() {
		return
	}
	t := carbon.Now().Nanosecond()
	//先执行关闭，避免以前有残留的 docker 容器
	stopChromedp(t)
	time.Sleep(1 * time.Second)

	startChromedp(t)
	defer stopChromedp(t)

	// 循环扫描端口是否已经启动
	for i := 0; i <= configs.DockerWaitTime; i++ {
		if i == configs.DockerWaitTime {
			configs.Logger.Error("Docker container waiting to start timeout 容器等待启动超时")
			return
		}
		time.Sleep(1 * time.Second)
		err := CheckPorts(configs.DockerProt)
		if err != nil {
			break
		} else {
			continue
		}
	}

	timeout := time.Duration(configs.TimeOut) * time.Second
	cromeCtx, cromeCtxCancel := GetChromeCtx()
	ctx, cancel := context.WithTimeout(cromeCtx, timeout)
	defer cancel()
	defer cromeCtxCancel()

	listenForNetworkEvent(ctx)

	login(ctx, username, password)
}

func startChromedp(t int) {
	configs.Logger.Info("startChromedp:", t)
	Cmd(
		"docker",
		"run", "-d",
		"-p", configs.DockerProt+":9222",
		"-m", "218m",
		"--rm",
		"--name", "headless-shell-utf-8",
		"--init", "lambdaexpression/headless-shell-utf-8:95.0.4638.32",
		`--blink-settings="imagesEnabled=false"`,
		//`--user-agent="Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_1 like Mac OS X) AppleWebKit/603.1.30 (KHTML, like Gecko) Version/10.0 Mobile/14E304 Safari/602.1"`,
		//`--accept-language="zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7,zh-TW;q=0.6"`,
		"--disable-gpu",
		"--headless",
		//"--disable-web-security",
		//"--ignore-certificate-errors-spki-list",
	)
}

func stopChromedp(t int) {
	configs.Logger.Info("stopChromedp:", t)
	Cmd("docker", "stop", "headless-shell-utf-8")
}

func checkLogin() bool {
	// 上次获取cookie时间
	getCookieTime := carbon.CreateFromTimestamp(GetCookie().LoginLastTime)
	// 下次获取cookie时间 = 上次获取cookie时间 + 获取cookie间隔时间
	nextTimeGetCookieTime := getCookieTime.AddSeconds(configs.LoginIntervalTime)
	// 比较 下次获取cookie时间 是否大于 现在时间
	if nextTimeGetCookieTime.Gt(carbon.Now()) {
		configs.Logger.Error(strconv.Itoa(configs.LoginIntervalTime) + " 秒内最多登录一次，下次获取cookie时间为" + nextTimeGetCookieTime.ToDateTimeString() + "，避免被封号")
		return false
	}
	return true
}

func login(ctx context.Context, username, password string) {

	indexUrl := "https://e.189.cn/wap/index.do"
	var collectLink string
	var res string
	var bodyHtml string
	var jAgreementCheckboxHtml string

	var b1, b2, b3 []byte
	var ck []*network.Cookie

	err := chromedp.Run(ctx,
		chromedp.Emulate(device.IPhone8Plus),
		chromedp.Navigate(indexUrl),
	)

	time.Sleep(8 * time.Second)

	err = chromedp.Run(ctx,
		chromedp.CaptureScreenshot(&b1),
		chromedp.OuterHTML("body", &bodyHtml, chromedp.ByQuery),
	)

	if !strings.Contains(bodyHtml, "id=\"j-account-login\" class=\"bg-login-wrap block\"") {
		err = chromedp.Run(ctx,
			chromedp.WaitVisible(`j-sms-userName`, chromedp.ByID),
			chromedp.Click(`j-other-login-way2`, chromedp.ByID),
			chromedp.WaitVisible(`j-login-btn`, chromedp.ByID),
			chromedp.OuterHTML("j-agreement-checkbox", &jAgreementCheckboxHtml, chromedp.ByID),
		)
	}

	if err != nil {
		configs.Logger.Error("error", err)
		outLogonPng(b1, b2, b3)
		return
	}

	if !strings.Contains(jAgreementCheckboxHtml, "ag-ckbox") {
		err = chromedp.Run(ctx,
			chromedp.Click(`j-agreement-checkbox`, chromedp.ByID),
		)
		if err != nil {
			configs.Logger.Error("error", err)
			outLogonPng(b1, b2, b3)
			return
		}
	}

	err = chromedp.Run(ctx,
		chromedp.SendKeys(`j-userName`, username, chromedp.ByID),
		chromedp.SendKeys(`j-password`, password, chromedp.ByID),
		chromedp.Click(`j-login-btn`, chromedp.ByID),
		chromedp.CaptureScreenshot(&b2),
		chromedp.WaitVisible(`_userMobile`, chromedp.ByID),
		chromedp.Sleep(2*time.Second),
		chromedp.CaptureScreenshot(&b3),
		chromedp.Location(&collectLink),
		GetCookies(&ck),
		//chromedp.OuterHTML("body", &res, chromedp.ByQuery),
		//chromedp.OuterHTML("html", &res, chromedp.ByQuery),
		chromedp.OuterHTML("title", &res, chromedp.ByQuery),
	)
	if err != nil {
		configs.Logger.Error("error", err)
		outLogonPng(b1, b2, b3)
		return
	} else {
		configs.Logger.Info("success")
	}

	c, err := json.Marshal(ck)
	if err != nil {
		configs.Logger.Error("cookie error", err)
		outLogonPng(b1, b2, b3)
		return
	}
	configs.Logger.Info("cookies : ", string(c))

	cookies := make([]string, len(ck))
	for i, v := range ck {
		cookies[i] = v.Name + "=" + url.QueryEscape(v.Value)
	}
	ti := time.Now()
	SetCookie(strings.Join(cookies, ";"), ti.Unix())

	configs.Logger.Info("cookie : " + strings.Join(cookies, ";"))

	configs.Logger.Info("link : ", collectLink)

	outLogonPng(b1, b2, b3)

}

func outLogonPng(b1 []byte, b2 []byte, b3 []byte) {
	ti := time.Now()
	t := "time : " + ti.Format("2006.01.02 15:04:05")
	if b1 != nil {
		Watermark(b1, t, configs.DataPath+"/login/01.png")
		configs.Logger.Info("write login/01.png")
	}
	if b2 != nil {
		Watermark(b2, t, configs.DataPath+"/login/02.png")
		configs.Logger.Info("write login/02.png")
	}
	if b3 != nil {
		Watermark(b3, t, configs.DataPath+"/login/03.png")
		configs.Logger.Info("write login/03.png")
	}
}

func GetCookies(cookies *[]*network.Cookie) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			c, err := network.GetAllCookies().Do(ctx)
			if err != nil {
				return err
			}
			*cookies = c
			return nil
		}),
	}
}

func GetChromeCtx() (context.Context, context.CancelFunc) {

	allocOpts := chromedp.DefaultExecAllocatorOptions[:]
	allocOpts = append(allocOpts,
		chromedp.DisableGPU,
		//chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.UserAgent(`Mozilla/5.0 (Linux; Android 10; 16th Build/QKQ1.191222.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/83.0.4103.106 Mobile Safari/537.36 MzmApp/4.1.11 FlymePreApp/4.1.11 AppVersionCode/4111 MzChannel/flyme DeviceBrand/meizu DeviceModel/16s`),
		chromedp.Flag("accept-language", `zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7,zh-TW;q=0.6`),
		chromedp.Flag("headless", true),
		chromedp.Flag("ignore-certificate-errors", true), //忽略错误
		chromedp.Flag("disable-web-security", true),
		chromedp.NoFirstRun,
	)

	allocCtx, allocCtxCancel := chromedp.NewRemoteAllocator(context.Background(), "ws://127.0.0.1:"+configs.DockerProt)
	chromeCtx, chromeCtxCancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(configs.Logger.Infof),
	)

	return chromeCtx, func() {
		allocCtxCancel()
		chromeCtxCancel()
	}
}

//监听
func listenForNetworkEvent(ctx context.Context) {
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *network.GetCookiesParams:
			configs.Logger.Info("cookie", ev.Urls)
		case *network.EventRequestWillBeSent:
			req := ev.Request
			if len(req.Headers) != 0 {
				if listenUrl(req.URL) {
					configs.Logger.Infof("Request : %s %s", req.URL, req.Headers)
				}
			}
		case *network.EventResponseReceived:
			resp := ev.Response
			if len(resp.Headers) != 0 {
				// log.Printf("received headers: %s", resp.Headers)

				if listenUrl(resp.URL) {
					configs.Logger.Infof("Response : %s %s %d", resp.URL, resp.Headers, resp.Status)

					// 打印body信息
					//go func() {
					//	// print response body
					//	c := chromedp.FromContext(ctx)
					//	rbp := network.GetResponseBody(ev.RequestID)
					//	body, err := rbp.Do(cdp.WithExecutor(ctx, c.Target))
					//	if err != nil {
					//		fmt.Println(err)
					//	} else {
					//		fmt.Printf("%s %s\n", resp.URL, body)
					//	}
					//}()
				}

			}

		case *runtime.EventExceptionThrown:
			configs.Logger.Infof("Event Exception, console time > %s \n", ev.Timestamp.Time())
			configs.Logger.Infof("\tException Type > %s \n", ev.ExceptionDetails.Exception.Type.String())
			configs.Logger.Infof("\tException Description > %s \n", ev.ExceptionDetails.Exception.Description)
			configs.Logger.Infof("\tException Stacktrace Text > %s \n", ev.ExceptionDetails.Exception.ClassName)
		case *runtime.StackTrace:
			configs.Logger.Infof("Stack Trace, console type > %s \n", ev.Description)
			for _, frames := range ev.CallFrames {
				configs.Logger.Infof("Frame line # %s\n", frames.LineNumber)
			}
		case *runtime.EventConsoleAPICalled:
			configs.Logger.Infof("Event Console API Called, console type > %s call:\n", ev.Type)
			for _, arg := range ev.Args {
				configs.Logger.Infof("%s - %s\n", arg.Type, arg.Value)
			}

		}
	})
}

func listenUrl(url string) bool {
	return strings.Index(url, ".gif") == -1 &&
		strings.Index(url, "data:image") == -1 &&
		strings.Index(url, ".png") == -1 &&
		strings.Index(url, ".js") == -1 &&
		strings.Index(url, ".css") == -1 &&
		strings.Index(url, "189.c") > -1
}

func GetFlowDetail(autoLogin bool) *models.DetailRequest {
	detail := "https://e.189.cn/store/user/package_detail.do"
	body := request(detail, 0, autoLogin)
	//configs.Logger.Info("GetDetail", body)

	detailResult := models.DetailRequest{}
	if err := json.Unmarshal([]byte(body), &detailResult); err != nil {
		configs.Logger.Error(err)
		return nil
	}
	return &detailResult
}

func GetBalance(autoLogin bool) *models.BalanceNew {
	balance := "https://e.189.cn/store/user/balance_new.do"
	body := request(balance, 0, autoLogin)
	//configs.Logger.Info("GetBalance", body)

	balanceNew := models.BalanceNew{}
	if err := json.Unmarshal([]byte(body), &balanceNew); err != nil {
		configs.Logger.Error(err)
		return nil
	}
	return &balanceNew
}

func GetFlowPackage(autoLogin bool) *models.FlowPackage {
	flow := "https://e.189.cn/store/wap/flowPackage.do"
	body := request(flow, 0, autoLogin)
	//configs.Logger.Info("GetBalance", body)

	flowPackage := models.FlowPackage{}
	if err := json.Unmarshal([]byte(body), &flowPackage); err != nil {
		configs.Logger.Error(err)
		return nil
	}
	return &flowPackage
}

func request(requestUrl string, count int, autoLogin bool) string {
	setCookie()
	err := configs.Browser.Open(requestUrl)
	if err != nil {
		configs.Logger.Error(err)
		return `{"result":-100}`
	}

	body := strings.ReplaceAll(configs.Browser.Body(), `&#34;`, `"`)
	configs.Logger.Info("request ", requestUrl, " ", body)
	detailResult := map[string]interface{}{}
	if err = json.Unmarshal([]byte(body), &detailResult); err != nil {
		configs.Logger.Error(err)
		return `{"result":-101}`
	}

	if detailResult["result"] != nil && fmt.Sprint(detailResult["result"]) == "-10001" && count > 0 {
		configs.Logger.Error(errors.New("More than the number of revolutions 超过重试次数"))
		return `{"result":-102}`
	} else if detailResult["result"] != nil && fmt.Sprint(detailResult["result"]) == "-10001" && autoLogin {
		// -10001 用户未登录
		ChinaTelecomLogin(configs.Username, configs.Password)
		return request(requestUrl, count+1, autoLogin)
	} else if detailResult["result"] != nil && fmt.Sprint(detailResult["result"]) == "-10001" && !autoLogin {
		// -10001 用户未登录
		return `{"result":-1}`
	}

	return body
}

func setCookie() {
	configs.Browser.AddRequestHeader("cookie", GetCookie().ChinaTelecomCookie)
}

func uurl(requestUrl string) *url.URL {
	var scheme string
	if strings.HasPrefix(strings.ToLower(requestUrl), "https") {
		scheme = "https"
	} else {
		scheme = "http"
	}

	return &url.URL{
		Scheme: scheme,
		Host:   "e.189.cn:443",
	}

	//Scheme      string
	//Opaque      string    // encoded opaque data
	//User        *Userinfo // username and password information
	//Host        string    // host or host:port
	//Path        string    // path (relative paths may omit leading slash)
	//RawPath     string    // encoded path hint (see EscapedPath method)
	//ForceQuery  bool      // append a query ('?') even if RawQuery is empty
	//RawQuery    string    // encoded query values, without '?'
	//Fragment    string    // fragment for references, without '#'
	//RawFragment string    // encoded fragment hint (see EscapedFragment method)
}

func cookie(c *network.Cookie) *http.Cookie {
	return &http.Cookie{
		Name:     c.Name,
		Value:    c.Value,
		Path:     c.Path,
		Domain:   c.Domain,
		Expires:  time.Unix(int64(c.Expires), 0),
		MaxAge:   int(time.Unix(int64(c.Expires), 0).Unix() - time.Now().Unix()),
		Secure:   c.Secure,
		HttpOnly: c.HTTPOnly,
		SameSite: sameSite(c.SameSite),
	}
}

func sameSite(sameSite network.CookieSameSite) http.SameSite {
	switch sameSite.String() {
	case "lax":
		return http.SameSiteLaxMode
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteDefaultMode
	}
}

// 检查端口是否存在
func CheckPorts(port string) error {
	var err error

	tcpAddress, err := net.ResolveTCPAddr("tcp4", ":"+port)
	if err != nil {
		return err
	}
	listener, err := net.ListenTCP("tcp", tcpAddress)
	if err != nil {
		return err
	} else {
		listener.Close()
	}

	return nil
}
