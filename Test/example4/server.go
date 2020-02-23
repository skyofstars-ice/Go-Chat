package main

import (
	"fmt"
	"net"
	"strings"
)
//创建基于链接的哈希map，方便知道链接对应的端口号
var onlineConns = make(map[string]net.Conn) 
//创建一个频道，长度1000.用来存储信息
var messageQueue = make(chan string, 1000)
//创建一个返回布尔值的频道，用来应对没有信息的情况
var quitChan = make(chan bool)
//检查错误，通用方法
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
//将消息存进频道中
func ProcessInfo(conn net.Conn){
	buf := make([]byte , 1024) 创建个1024字节的字节流
	defer conn.Close() //程序运行完毕时关闭链接
    //循环体
	for{
	   //从信息流读出消息
		numOfBytes,err:= conn.Read(buf) 
		if err != nil{
			panic(err)
		}
     //如果字节不等于零，就把消息字符化发到消息频道中
	   if numOfBytes != 0{
			message := string(buf[0:numOfBytes])
			messageQueue <- message
		}
	}
}
//消费消息
func ConsumeMessage(){
	for{
	   //select是配合频道使用的，选择进哪个通信操作
		select{
			case message:=<-messageQueue:
				//对消息进行解析
				doProcessMessage(message)
			case <- quitChan:
			   //空消息，什么都不做~
				break
			}
	}
}
//对消息进行解析
func doProcessMessage(message string){
    //将消息以“#”分开，如 127.0.0.1:3389#你好
    //就把这个消息分成两段，一段为127.0.0.1:3389，一段为你好
    //这里面没有考虑消息体里面包括#好的情况
	contents := strings.Split(message,"#")
	if len(contents) >1 {
		addr := contents[0]
		sendMessage := contents[1]
     //将消息体格式化
		addr =strings.Trim(addr, " ")
		//判断哈希map中这个地址存不存在，存在就往这个链接里写消息体
		if conn,ok := onlineConns[addr]; ok{
			_,err := conn.Write([]byte(sendMessage))
			if err != nil{
			  //不存在就写发送失败
				fmt.Println("online conns send failure!!!!")
			}
		}
	}
}
//主函数
func main(){
  //监听127.0.0.1:8080端口。
	listen_socket , err := net.Listen("tcp","127.0.0.1:8080")
	CheckError(err)
	//程序退出后监听关闭
	defer listen_socket.Close()
   //增加用户体验，让用户知道程序在运行中
	fmt.Println("server is waitting...")
  //开启协程来消费消息
	go ConsumeMessage()
  //无限循环，来监听消息并处理消息
	for{
		conn,err := listen_socket.Accept()
		CheckError(err)
     //打印出连接过来的远端地址
		addr := fmt.Sprintf("%s", conn.RemoteAddr())
		//将这个连接并存储到hashmap
		onlineConns[addr] = conn
     //循环hashmap，打印连接上来的连接
		for i := range onlineConns{
			fmt.Println(i)
		}
		//运行协程，将连接存到频道中。程序不停，此过程不断
		go ProcessInfo(conn)
	}
}
