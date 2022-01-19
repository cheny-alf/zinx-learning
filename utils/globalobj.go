package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/ziface"
)

type GlobalObj struct {
	/*
		Server
	*/
	TcpServer ziface.IServer //当前Zinx的全局对象
	Host      string         //ip地址
	TcpPort   int            //端口号
	Name      string         //当前服务器名称

	/*
		Zinx
	*/
	Version        string //版本号
	MaxConn        int    //最大连接数
	MaxPackageSize uint32 //数据包最大值

}

/**
定义一个全局的对外对象
*/
var GlobalObject *GlobalObj

//初始化obj
func init() {
	GlobalObject = &GlobalObj{
		Host:           "0.0.0.0",
		TcpPort:        8999,
		Name:           "ZinxServerApp",
		Version:        "V0.4",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}
	GlobalObject.Reload()
}

//加载zinx文件的方法
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}
