# ChinaTelecomMonitor 



**中国电信 手机话费、流量、语音通话监控**

本工具是部署在服务器(或x86软路由等设备) 使用 接口模拟登录获取cookie，按需定时获取电信手机话费、流量、语音通话使用情况，通过接口返回数据。
可配合 Scriptables 插件 [ChinaTelecomPanel](https://lambdaexpression.github.io/ScriptablesComponent/ChinaTelecomPanel/) 一起使用


### 版本更新

**v2.0.1**

fix:
- 修复 通用流量显示为总流量问题

**v2.0**

update:
- 更新 电信请求接口。从容器化模拟登录，更换为纯接口调用登录
- 调整 日志打印规则
- 新增 /show/qryImportantData、/show/userFluxPackage 接口（需要启动应用时开启dev模式才能访问）
- 删除 /show/detail、/show/flowPackage 接口
- 删除 dockerProt、dockerWaitTime 启动参数

### 1.准备

- 1.准备一个可正常登录 电信APP 的账号密码（**注意：电信APP 的密码 和以前的h5登录密码不是同一个**）
- 2.下载本应用 `wget https://github.com/LambdaExpression/ChinaTelecomMonitor/releases/download/v2.0/China_Telecom_Monitor_amd64`
- 3.应用授权 `chmod +x ./China_Telecom_Monitor_amd64`

### 2.启动应用

```
$ ./China_Telecom_Monitor_amd64 --prot 8081 --username '电信账号' --password '电信密码'
```

### 3.测试访问

```shell
curl http://127.0.0.1:8081/show/flow
{"code":200,"data":{"id":0,"username":"","use":12276406,"total":167045874,"generalUse":12276406,"generalTotal":83159794,"specialUse":0,"specialTotal":83886080,"balance":7036,"voiceUsage":0,"voiceAmount":500,"createTime":"2022-04-26 15:37:47"}}
```
**接口参数说明**

```
{
    "code":200,
    "data":{
        "id":0,                              // 保留字段
        "username":"",                       // 手机号，默认为空，dev 模式下脱敏显示
        "use":12276406,                      // 总流量使用量，单位kb
        "total":167045874,                   // 总流量总量，单位kb
        "generalUse":12276406,               // 通用流量使用量，单位kb
        "generalTotal":83159794,             // 通用流量总量，单位kb
        "specialUse":0,                      // 专用流量使用量，单位kb
        "specialTotal":83886080,             // 专用流量总量，单位kb
        "balance":7036,                      // 话费余额，单位分
        "voiceUsage":0,                      // 通话语音使用量，单位分钟
        "voiceAmount":500,                   // 通话语音总量，单位分钟
        "createTime":"2022-04-26 15:37:47"   // 获取数据时间
        "items": [                           // 具体流量使用情况
                     {
                        "name": "国内流量",
                        "use": 12276406,
                        "total": 83159794
                     },
                     {
                        "name": "定向流量",
                        "use": 0,
                        "total": 41943040
                     },
                     {
                        "name": "闲时流量",
                        "use": 0,
                        "total": 41943040
                     }
         ]
    }
}
```

### 补充

**应用支持参数**

除了账号、密码是必填的，其他参数都可以保持默认，应用会为其设置默认值

```shell
$ ./China_Telecom_Monitor_amd64 -h
Usage of ./China_Telecom_Monitor_amd64:
  -dataPath string
    	--dataPath ./data # 数据日志文件保存路径 (default "./data")
  -dev
    	--dev false # 开发模式,开启后将支持以下接口： /refresh 手动更新流量 和 /show/qryImportantData /show/userFluxPackage 这里两个电信接口
  -intervalsTime int
    	--intervalsTime 180 #接口防止重刷时间 (default 180)
  -logEncoding string
    	--logEncoding console # 日志输出格式 console 或 json (default "console")
  -logLevel string
    	--logLevel info # 日志等级 (default "info")
  -loginIntervalTime int
    	--loginIntervalTime 43200 #电信登录间隔时间（防止被封号），秒 (default 43200)
  -password string
    	--password xxxxx #电信账号密码, 必填
  -prot string
    	--prot 8080 (default "8080")
  -timeOut int
    	--timeOut 30 #访问电信接口请求超时时间，秒 (default 30)
  -username string
    	--username 1xxxxxxxxxx #电信账号用户名, 必填
  -version
    	--version 打印程序构建版本

```

**额外接口说明**


```
// 数据来源 https://appfuwu.189.cn:9021/query/qryImportantData
curl http://127.0.0.1:8081/show/qryImportantData

// 数据来源 https://appfuwu.189.cn:9021/query/userFluxPackage
curl http://127.0.0.1:8081/show/userFluxPackage

// 数据来源 https://e.189.cn/store/user/package_detail.do
curl http://127.0.0.1:8081/show/detail （已无法使用）

// 数据来源 https://e.189.cn/store/wap/flowPackage.do
curl http://127.0.0.1:8081/show/flowPackage （已无法使用）
```

这些接口主要提供给有二次开发需求的用户使用。没有进行数据二次处理，完全是原始数据输出。

**服务后台运行**

我通过 issues ，发现有部分用户，并不会让服务在后台不挂断运行。我在这里只提出其中一种方案，使用 Linux的 `nohup` 命令。具体的还请大家在网上自行学习。

下面是结合 `nohup` 后的启动命令，这样就能保证大家退出服务器后，服务仍然在运行
```
$ nohup ./China_Telecom_Monitor_amd64 --prot 8081 --username '电信账号' --password '电信密码' >/dev/null &
```

**异常情况**


```json
{"headerInfos":{"code":"0000","reason":"操作成功"},"responseData":{"resultCode":"1000","resultDesc":"请求失败","attach":"","data":null}}
```
遇到“请求失败”，请确保在 **手机电信APP** 内测试一下登录账号密码是否正确


### 最后感谢 [boxjs](https://github.com/gsons/boxjs) 项目开源提供的电信接口
