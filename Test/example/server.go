//chatroom server
package main

import (
    "fmt"
    "net"
)

//定义函数checkError，用来错误处理
func checkError(err error)  {
    if err != nil {
        panic(err)
    }
}

//定义一个函数，专门负责接收信息的协程
func processInfo(conn net.Conn)  {
    buf := make([]byte, 1024)
    defer conn.Close()

    for {
        //读取数据
        numOfBytes, err := conn.Read(buf)
        if err != nil {
            break
        }

        //如果接收的字节数不为0，说明有消息发过来
        if numOfBytes != 0 {
            fmt.Printf("Has received message: %s", string(buf))
        }

    }
}

func main()  {
    //开启监听
    listen_socket, err := net.Listen("tcp","127.0.0.1:8888")
    checkError(err)
    defer listen_socket.Close()

    fmt.Println("chat room Server is waiting... ")

    //接收连接
    for {
        conn, err := listen_socket.Accept()
        checkError(err)

        //如果有客户端连接，则打开一个协程处理
        go processInfo(conn)

    }
}
