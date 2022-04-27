package tools

import (
	"China_Telecom_Monitor/configs"
	"io/ioutil"
	"os/exec"
)

func Cmd(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	stdout, err := cmd.StdoutPipe()
	if err != nil { //获取输出对象，可以从该对象中读取输出结果
		configs.Logger.Fatal(err)
	}
	defer stdout.Close()                // 保证关闭输出流
	if err := cmd.Start(); err != nil { // 运行命令
		configs.Logger.Fatal(err)
	}
	if opBytes, err := ioutil.ReadAll(stdout); err != nil { // 读取输出结果
		configs.Logger.Fatal(err)
	} else {
		configs.Logger.Info(string(opBytes))
	}
}
