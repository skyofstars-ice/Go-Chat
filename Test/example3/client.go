package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
)

func scandata(conn net.Conn) {
    for {
        //设置可以读取带空格的myscan
        myscan := bufio.NewReader(os.Stdin)
        buf, _, _ := myscan.ReadLine()
        if string(buf) == "q" {
            os.Exit(0)
        }
        conn.Write(buf)
    }
}
func main() {
    conn, _ := net.Dial("tcp", "127.0.0.1:9009")
    buf := make([]byte, 1024)
    args := os.Args
    conn.Write([]byte(args[1]))
    go scandata(conn)
    for {
        n, _ := conn.Read(buf)
        fmt.Println(string(buf[:n]))
    }
}
