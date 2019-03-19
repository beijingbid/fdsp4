package main

import (
	"encoding/json"
	_ "fdsp4/routers"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

var sep = '\x02'
var sep_str = string(sep)

func initLogger() (err error) {

	config := make(map[string]interface{})
	config["filename"] = beego.AppConfig.String("log_path")

	// map 转 json
	configStr, err := json.Marshal(config)
	if err != nil {
		fmt.Println("initLogger failed, marshal err:", err)
		return
	}
	// log 的配置
	beego.SetLogger(logs.AdapterFile, string(configStr))
	// log打印文件名和行数
	beego.SetLogFuncCall(true)
	fmt.Println(string(configStr))
	return
}
func main() {
	beego.Run()
	initLogger()
}
