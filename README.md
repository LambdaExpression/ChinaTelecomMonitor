# ChinaTelecomMonitor 
**中国电信 手机话费、流量、语音通话监控**

本工具是部署在服务器(或x86软路由等设备) 使用 docker [lambdaexpression/headless-shell-utf-8](https://hub.docker.com/r/lambdaexpression/headless-shell-utf-8) 进行模拟浏览器登录获取cookie，按需定时获取电信手机话费、流量、语音通话使用情况，通过接口返回数据。
可配合 Scriptables 插件 [ChinaTelecomPanel](https://lambdaexpression.github.io/ScriptablesComponent/ChinaTelecomPanel/) 一起使用


### 版本更新

*(注：应用如果在正常期间突然无法使用，大概率是电信登录页面的页面元素发生改变，导致登录逻辑无法正常执行。在保证自身使用的是最新版本后，任然无法解决可通过 [Issuess](https://github.com/LambdaExpression/ChinaTelecomMonitor/issues) 进行反馈。电信的这种行为，个人怀疑这是电信的防爬虫策略，现阶段应用也只能发现一种就通过版本迭代兼容一种)*

**v1.0.6.1**

update:
 - 优化 网页截图

**v1.0.6**

update:
 - 优化 电信检查环境时间逻辑 [#9](https://github.com/LambdaExpression/ChinaTelecomMonitor/issues/9)

**v1.0.5**

add:
 - 添加 /show/flowPackage 接口 [#8](https://github.com/LambdaExpression/ChinaTelecomMonitor/issues/8)

update:
 - 更新 通用流量和专用流量统计规则 [#8](https://github.com/LambdaExpression/ChinaTelecomMonitor/issues/8)

**v1.0.4**

add:
 - 添加 -version 版本打印

update:

 - 优化 自适应电信登录页面登录动作 [#5](https://github.com/LambdaExpression/ChinaTelecomMonitor/issues/5)

**v1.0.3**

add:
 - 添加 /show/detail 接口 [#4](https://github.com/LambdaExpression/ChinaTelecomMonitor/issues/4)

**v1.0.2**

add:
 - 添加自动判断是否勾选用户隐私协议 [#3](https://github.com/LambdaExpression/ChinaTelecomMonitor/issues/3)

**v1.0.1**

update：
 - 修复新版本电信登录页面默认勾选隐私协议，导致登录失败问题 [#1](https://github.com/LambdaExpression/ChinaTelecomMonitor/issues/1)

### 1.准备

- 1.准备一个可正常登录[电信](https://e.189.cn/wap/index.do)账号密码
- 2.安装 docker (可自行查询谷哥或度娘)
- 3.执行 `docker pull lambdaexpression/headless-shell-utf-8:95.0.4638.32`，下载 [lambdaexpression/headless-shell-utf-8](https://hub.docker.com/r/lambdaexpression/headless-shell-utf-8) 容器到本地
- 4.下载本应用 `wget https://github.com/LambdaExpression/ChinaTelecomMonitor/releases/download/v1.0.6.1/China_Telecom_Monitor_amd64`
- 5.应用授权 `chmod +x ./China_Telecom_Monitor_amd64`

### 2.启动应用

```
$ ./China_Telecom_Monitor_amd64 --prot 8081 --dockerProt 9222 --username '电信账号' --password '电信密码'
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
    	--dev false # 开发模式,开启后将支持以下接口： /refresh 手动更新流量，/loginLog 查看登录截图日志
  -dockerProt string
    	--dockerProt 9222 #登录容器使用的端口 (default "9222")
  -dockerWaitTime int
    	--dockerWaitTime 60 #登录容器等待启动时间 (default 60)
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
curl http://127.0.0.1:8081/show/detail
```

这些接口主要提供给有二次开发需求的用户使用。接口的数据是从电信的 https://e.189.cn/store/user/package_detail.do 接口获取而来，没有进行数据二次处理，完全是原始数据输出。[#4](https://github.com/LambdaExpression/ChinaTelecomMonitor/issues/4)


```
curl http://127.0.0.1:8081/show/flowPackage
```

接口数据来源自 https://e.189.cn/store/wap/flowPackage.do ，数据同样不进行二次处理，原样输出 [#8](https://github.com/LambdaExpression/ChinaTelecomMonitor/issues/8#issuecomment-1165158557)

**服务后台运行**

我通过 issues ，发现有部分用户，并不会让服务在后台不挂断运行。我在这里只提出其中一种方案，使用 Linux的 `nohup` 命令。具体的还请大家在网上自行学习。

下面是结合 `nohup` 后的启动命令，这样就能保证大家退出服务器后，服务仍然在运行
```
$ nohup ./China_Telecom_Monitor_amd64 --prot 8081 --dockerProt 9222 --username '电信账号' --password '电信密码' >/dev/null &
```

**异常情况**

如果上述第3步，返回数据有问题或无法访问。
- 1. 先查看一下同级目录下的 `./data/login/01.png`、`./data/login/02.png` 和 `./data/login/03.png` 这三张截图，是否有存在账号密码错误、登录要求输入验证码等情况。处理“登录要求输入验证码”，需要先手动访问 [https://e.189.cn/wap/index.do](https://e.189.cn/wap/index.do) 进行一次正常登陆，正常情况下就能恢复了。
- 2. 如果上述截图无法找到问题，可以查看一下 `./data/log/stdout.log` 是否存在异常信息（注意：首次启动时提示 `open ./data/cookie.json: no such file or directory` 是正常情况）

