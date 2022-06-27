package main

import (
	"China_Telecom_Monitor/configs"
	"China_Telecom_Monitor/models"
	"China_Telecom_Monitor/tools"
	"flag"
	"fmt"
	"github.com/LambdaExpression/surf/agent"
	"github.com/LambdaExpression/surf/surf"
	"github.com/golang-module/carbon/v2"
	"github.com/kataras/iris/v12"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
	"time"
)

var Version = "v1.0.6"
var GoVersion = "not set"
var GitCommit = "not set"
var BuildTime = "not set"

func main() {

	initFlag()
	if configs.PrintVersion {
		version()
		return
	}

	initLogger()

	if checkFlag() {
		return
	}

	initBrowser()
	initCron()

	// 初始化流量获取
	go cronSummary()

	initIris()
}

// 初始化配置
func initFlag() {
	flag.StringVar(&configs.Prot, "prot", "8080", "--prot 8080")
	flag.StringVar(&configs.DockerProt, "dockerProt", "9222", "--dockerProt 9222 #登录容器使用的端口")
	flag.StringVar(&configs.Username, "username", "", "--username 1xxxxxxxxxx #电信账号用户名, 必填")
	flag.StringVar(&configs.Password, "password", "", "--password xxxxx #电信账号密码, 必填")
	flag.IntVar(&configs.LoginIntervalTime, "loginIntervalTime", 43200, "--loginIntervalTime 43200 #电信登录间隔时间（防止被封号），秒")
	flag.Int64Var(&configs.TimeOut, "timeOut", 30, "--timeOut 30 #访问电信接口请求超时时间，秒")
	flag.IntVar(&configs.IntervalsTime, "intervalsTime", 180, "--intervalsTime 180 #接口防止重刷时间")
	flag.IntVar(&configs.DockerWaitTime, "dockerWaitTime", 60, "--dockerWaitTime 60 #登录容器等待启动时间")

	flag.StringVar(&configs.LogLevel, "logLevel", "info", "--logLevel info # 日志等级")
	flag.StringVar(&configs.LogEncoding, "logEncoding", "console", "--logEncoding console # 日志输出格式 console 或 json")

	flag.StringVar(&configs.DataPath, "dataPath", "./data", "--dataPath ./data # 数据日志文件保存路径")

	flag.BoolVar(&configs.Dev, "dev", false, "--dev false # 开发模式,开启后将支持以下接口： /refresh 手动更新流量，/loginLog 查看登录截图日志")
	flag.BoolVar(&configs.PrintVersion, "version", false, "--version 打印程序构建版本")

	flag.Parse()
}

func version() {
	fmt.Printf("Version: %s\n", Version)
	fmt.Printf("Go Version: %s\n", GoVersion)
	fmt.Printf("Git Commit: %s\n", GitCommit)
	fmt.Printf("Build Time: %s\n", BuildTime)
}

func checkFlag() bool {
	if configs.Username == "" {
		configs.Logger.Error("--username 电信账号不能为空")
		return true
	}
	if configs.Password == "" {
		configs.Logger.Error("--password 电信密码不能为空")
		return true
	}
	return false
}

// 初始化 zap日志框架
func initLogger() {
	level := getLevel()
	encoding := configs.LogEncoding
	// 保留两个变量，但设置成同一个文件
	stdout := configs.DataPath + "/log/stdout.log"
	stderr := configs.DataPath + "/log/stdout.log"

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,    // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 短路径编码器
	}
	atom := zap.NewAtomicLevelAt(level)
	config := zap.Config{
		Level:         atom,          // 日志级别
		Development:   true,          // 开发模式，堆栈跟踪
		Encoding:      encoding,      // 输出格式 console 或 json
		EncoderConfig: encoderConfig, // 编码器配置
		//InitialFields:    map[string]interface{}{"serviceName": "gogs-backup"}, // 初始化字段，如：添加一个服务器名称
		OutputPaths:      []string{"stdout", stdout}, // 输出到指定文件 stdout（标准输出，正常颜色） stderr（错误输出，红色）
		ErrorOutputPaths: []string{"stderr", stderr},
	}
	tools.Create(stdout)
	tools.Create(stderr)
	logger, err := config.Build()
	if err != nil {
		panic(fmt.Sprintf("logger 初始化失败: %v", err))
	}
	configs.Logger = logger.Sugar()
	configs.Logger.Info("logger 初始化成功")
}

// 获取 日志等级
func getLevel() zapcore.Level {
	levelStr := strings.TrimSpace(strings.ToLower(configs.LogLevel))
	var level zapcore.Level
	switch levelStr {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "dpanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}
	return level
}

// 初始化 浏览器
func initBrowser() {
	bow := surf.NewBrowser()
	bow.SetUserAgent(agent.Chrome())
	bow.HistoryJar().SetMax(1)
	bow.SetTimeout(time.Duration(configs.TimeOut) * time.Second)
	configs.Browser = bow
}

// 初始化 定时任务
func initCron() {
	cronApp := cron.New(cron.WithSeconds()) //精确到秒
	spec := "0 0 */2 * * ?"                 //cron表达式，每2小时一次，保障cookie有效
	_, err := cronApp.AddFunc(spec, func() {
		go cronSummary()
	})
	if err != nil {
		configs.Logger.Error(`0x000007 `, err)
		return
	}
	cronApp.Start()
}

// 定时获取流量信息
func cronSummary() {
	t := carbon.Now()
	detailRequest := tools.GetFlowDetail(true)
	balance := tools.GetBalance(true)
	flowPackage := tools.GetFlowPackage(true)
	configs.Summary = tools.ToSummary2(detailRequest, flowPackage, balance, configs.Username, t)
}

// 初始化访问接口
func initIris() {
	irisApp := iris.New()
	irisApp.Use(middleware)
	irisApp.Handle(iris.MethodGet, "/show/flow", flow)
	irisApp.Handle(iris.MethodGet, "/show/detail", packageDetail)
	irisApp.Handle(iris.MethodGet, "/show/flowPackage", flowPackage)
	if configs.Dev {
		irisApp.Handle(iris.MethodGet, "/refresh", refresh)
		irisApp.Handle(iris.MethodGet, "/image/{filename:string}", image)
		irisApp.Handle(iris.MethodGet, "/loginLog", loginLog)
	}
	err := irisApp.Run(iris.Addr(":" + configs.Prot))
	if err != nil {
		configs.Logger.Error("InitIris error", err)
	}
}

func middleware(ctx iris.Context) {
	configs.Logger.Infof("iris access Path: %s | IP: %s", ctx.Path(), ctx.RemoteAddr())
	ctx.Next()
}

var flowLastTime carbon.Carbon

// 获取流量信息接口
func flow(ctx iris.Context) {

	if carbon.Now().Lt(flowLastTime.AddSeconds(configs.IntervalsTime)) {
		ctx.JSON(iris.Map{"code": 200, "data": desensitization(configs.Summary)})
		return
	}

	flowLastTime = carbon.Now()

	t := carbon.Now()
	detailRequest := tools.GetFlowDetail(false)
	if detailRequest.Result != 0 {
		go cronSummary()
	} else {
		balance := tools.GetBalance(true)
		flowPackage := tools.GetFlowPackage(true)
		configs.Summary = tools.ToSummary2(detailRequest, flowPackage, balance, configs.Username, t)
	}
	summary := desensitization(configs.Summary)
	ctx.JSON(iris.Map{"code": 200, "data": summary})
}

var packageDetailVisitLastTime carbon.Carbon
var packageDetailDetailRequest *models.DetailRequest

func packageDetail(ctx iris.Context) {
	if carbon.Now().Lt(packageDetailVisitLastTime.AddSeconds(configs.IntervalsTime)) {
		ctx.JSON(&packageDetailDetailRequest)
		return
	}
	packageDetailDetailRequest = tools.GetFlowDetail(false)
	packageDetailVisitLastTime = carbon.Now()
	ctx.JSON(&packageDetailDetailRequest)
}

var flowPackageVisitLastTime carbon.Carbon
var flowPackageDetailRequest *models.FlowPackage

func flowPackage(ctx iris.Context) {
	if carbon.Now().Lt(flowPackageVisitLastTime.AddSeconds(configs.IntervalsTime)) {
		ctx.JSON(&flowPackageDetailRequest)
		return
	}
	flowPackageDetailRequest = tools.GetFlowPackage(false)
	flowPackageVisitLastTime = carbon.Now()
	ctx.JSON(&flowPackageDetailRequest)
}

func desensitization(summary models.Summary) models.Summary {
	if configs.Dev {
		if len(summary.Username) == 11 {
			summary.Username = summary.Username[0:3] + "****" + summary.Username[7:11]
		}
	} else {
		if len(summary.Username) > 0 {
			summary.Username = ""
		}
	}
	return summary
}

// 手动触发获取流量信息
func refresh(ctx iris.Context) {
	go cronSummary()
	ctx.JSON(iris.Map{"code": 200})
}

func image(ctx iris.Context) {
	filename := ctx.Params().GetStringDefault("filename", "01.png")
	ctx.ServeFile(configs.DataPath+"/login/"+filename, false)
}

func loginLog(ctx iris.Context) {
	ctx.HTML(`
		<html>
		  <div style="text-align: center;">
			<br/>
			<h3>登录图片1</h3>
			<br/>
			<image style="max-width:500px" src="/image/01.png"></image>
			<br/>
			<h3>登录图片2</h3>
			<br/>
			<image style="max-width:500px" src="/image/02.png"></image>
			<br/>
			<h3>登录图片3</h3>
			<br/>
			<image style="max-width:500px" src="/image/03.png"></image>
		  </div>
		</html>
	`)
}
