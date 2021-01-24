// 简易的GO_TCP代理程序
// 启动方法有两种
//1、直接执行 ./go_forwardport 本地监听地址:端口    目标服务器地址:端口
//2、在当前目录下配置config.ini
//格式如下：
//			0.0.0.0:1234	1.1.1.1:5678
//			0.0.0.0:6666	1.1.1.1:7642
//
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"
)

var (
	filepath string = "./config.ini"
	arge []string
)

func tcp_handle(tcpConn net.Conn,dst_addr string){
	remote_tcp,err:=net.Dial("tcp",dst_addr) //连接目标服务器
	if err!=nil{
		fmt.Println(err)
		return
	}
	go io.Copy(remote_tcp,tcpConn)
	go io.Copy(tcpConn,remote_tcp)
}
func tcp_listen(local_addr string, dst_addr string){
	ln,err:=net.Listen("tcp",local_addr)
	if err!=nil{
		fmt.Println("tcp_listen:",err)
		return
	}
	defer ln.Close()
	for {
		tcp_Conn,err:=ln.Accept() //接受tcp客户端连接，并返回新的套接字进行通信
		if err!=nil{
			fmt.Println("Accept:",err)
			return
		}
		go tcp_handle(tcp_Conn,dst_addr)   //创建新的协程进行转发
	}
}
func cal(arge []string)  {
	var (
		local_addr string
		dst_addr string
	)
	for i := 1 ; i < len(arge) ; i++{
		if i%2 == 1 {
			local_addr = arge[i]
		}
		if i%2 == 0 {
			dst_addr = arge[i]
			fmt.Println(local_addr + "->" + dst_addr)
			go tcp_listen(local_addr, dst_addr)
		}
	}
	fmt.Println("端口转发正在运行，请勿关闭本程序")
	for{
		time.Sleep(time.Duration(2)*time.Second)
	}
}
func main(){

	content ,err :=ioutil.ReadFile(filepath)
	if err !=nil {
		arge = os.Args
		if len(arge) < 2 {
			fmt.Println("参数格式：   本地监听地址:端口    目标服务器地址:端口")
			return
		}
	}else{
		arge = strings.Fields("0 " + string(content))
	}
	cal(arge)

}

