package conf

import (
	"fmt"
	"gopkg.in/ini.v1"
	"strings"
	"todo_list/model"
)

var (
	AppMode    string
	HttpPort   string
	Db         string
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassWord string
	DbName     string
)

func Init() {
	file, err := ini.Load("./conf/config.ini")
	if err != nil {
		fmt.Println("配置文件读取错误，请检查文件路径")
	}
	LoadServer(file)
	LoadMysql(file)
	path := strings.Join([]string{DbUser, ":", DbPassWord, "@tcp(", DbHost, ":", DbPort, ")/", DbName, "?charset=utf8mb4"}, "")
	model.Database(path)
}

func LoadServer(file *ini.File) {
	AppMode = file.Section("server").Key("AppMode").MustString("debug")
	HttpPort = file.Section("server").Key("HttpPort").MustString(":8080")
}

func LoadMysql(file *ini.File) {
	Db = file.Section("mysql").Key("Db").MustString("mysql")
	DbHost = file.Section("mysql").Key("DbHost").MustString("localhost")
	DbPort = file.Section("mysql").Key("DbPort").MustString("3306")
	DbUser = file.Section("mysql").Key("DbUser").MustString("root")
	DbPassWord = file.Section("mysql").Key("DbPassWord").MustString("password")
	DbName = file.Section("mysql").Key("DbName").MustString("go")
}
