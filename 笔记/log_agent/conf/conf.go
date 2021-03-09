package conf

import (
	"gopkg.in/ini.v1"
	"os"
)

// fileExist 检查配置文件是否存在
func fileExist(path string) bool{
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}


// Loadconf 加载并解析.ini配置文件
func Loadconf(path string,Section string, confMap interface{}) error {

	fileExist(path)

	// 加载配置文件
	cfg, err := ini.Load(path)
	if err != nil {
		panic(err)
	}

	// 返回一个配置文件对象
	err = cfg.Section(Section).MapTo(confMap)
	if err != nil{
		panic(err)
	}

	return nil
}