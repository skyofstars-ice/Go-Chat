// chat room client
package main

import (
    "fmt"
    "net"
)

func checkError(err error)  {
    if err != nil {
        panic(err)
    }
}

func main()  {
    //连接服务端的信息
    conn, err := net.Dial("tcp","127.0.0.1:8888")
    checkError(err)
    defer conn.Close()

    //向服务端发送信息
    conn.Write([]byte("Now connection.success!"))

    fmt.Println("Has sent the message!")
}
