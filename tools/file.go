package tools

import (
	"China_Telecom_Monitor/configs"
	"errors"
	"os"
	"path/filepath"
)

// 判断文件是否存在
func IsExist(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// 创建文件及其文件夹
func Create(path string) (*os.File, error) {
	//if IsExist(path) {
	//	return nil,errors.New("文件已经存在")
	//}
	dir := filepath.Dir(path)
	if !IsExist(dir) {
		e := os.MkdirAll(dir, 0755)
		if e != nil {
			configs.Logger.Error("创建文件夹失败, Error : ", e)
			return nil, errors.New("创建文件夹失败")
		}
	}

	return os.Create(path)
}

func ReadFile(file string) (string, error) {
	bs, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func WriteFile(file string, str string) error {
	return os.WriteFile(file, []byte(str), os.FileMode(int(0644)))
}
