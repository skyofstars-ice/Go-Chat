package main

import (
    "fmt"
    "net"
    "strings"
)

type info struct {
    conn net.Conn
    name string
}

var ch_all chan string = make(chan string)
var ch_one chan string = make(chan string)
var ch_who chan string = make(chan string)
var infos map[string]info = make(map[string]info)

func handle(conn net.Conn) {
    defer conn.Close()

    buf := make([]byte, 100)
    n, _ := conn.Read(buf)
    name := string(buf[:n])
    var oneinfo info
    oneinfo.conn = conn
    oneinfo.name = name
    addr := conn.RemoteAddr().String()
    infos[addr] = oneinfo
    msg := name + "进入聊天室"
    ch_all <- msg

    for {
        n, _ := conn.Read(buf)
        if n == 0 {
            fmt.Printf("%s断开连接\n", addr)
            msg = name + "离开聊天室"
            delete(infos, addr)
            ch_all <- msg
            fmt.Println(msg)
            return
        }
        if string(buf[:n])[0] == '@' {
            sli := strings.Fields(string(buf[1:n])) //隔开
            who := sli[0]
            msg = strings.Join(sli[1:], "") //后边的再拼接回去
            ch_who <- who
            ch_one <- name + "->me : " + msg //单发
            continue
        }
        ch_all <- name + "->all : " + string(buf[:n]) //群发
    }
}

func sendone() {
    for {
        who := <-ch_who
        msg := <-ch_one
        for _, v := range infos {
            if v.name == who {
                v.conn.Write([]byte(msg))
                break
            }
        }
    }
}
func sendall() {
    for {
        msg := <-ch_all
        for _, val := range infos {
            val.conn.Write([]byte(msg))
        }
    }
}
func main() {
    listener, _ := net.Listen("tcp", ":9009")
    defer listener.Close()
    go sendall()
    go sendone()
    for {
        conn, _ := listener.Accept()
        fmt.Printf("%s建立连接\n", conn.RemoteAddr().String())
        go handle(conn)
    }
}
