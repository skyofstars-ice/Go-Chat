package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
//发送消息
func MessageSend(conn net.Conn) {
	var input string
	//无限循环，发送消息
	for {
	   //监听输入将键入的信息记下来
		reader := bufio.NewReader(os.Stdin)
		//将记下来的消息以每行开始读
		data, _, _ := reader.ReadLine()
		//将data字符化存入到input
		input = string(data)
     //如果input为exit就退出，结束连接
		if strings.ToUpper(input) == "EXIT" {
			conn.Close()
			break
		} else {
		  //否则，往这个连接中写入字节流
			_, err := conn.Write([]byte (input))
			if err != nil {
				conn.Close()
				fmt.Println("client connect failure" + err.Error())
			}
		}
	}
}

func main() {
  //连接8080端口
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	CheckError(err)
	defer conn.Close()
   //开启协程发送消息，无限循环的
	go MessageSend(conn)
	//创建1024长度的字节流
	buf := make([]byte, 1024)
	//无限循环，读取字节流信息并打印
	for {
		_, err := conn.Read(buf)
		CheckError(err)
		fmt.Println("server receive message content: " + string(buf))
	}

	fmt.Println("client program done")
}
