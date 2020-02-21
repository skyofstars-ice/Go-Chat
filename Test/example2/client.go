// client project main.go
package main
 
import (
	"fmt"
	"net"
)
 
var nick string = ""
 
func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Println("conn fail...")
	}
	defer conn.Close()
	fmt.Println("connect server successed \n")
 
	//给自己取一个昵称吧
	fmt.Printf("Make a nickname:")
	fmt.Scanf("%s", &nick)
	fmt.Println("hello : ", nick)
	conn.Write([]byte("nick|" + nick))
 
	go Handle(conn)
 
	var msg string
	for {
		msg = ""
		fmt.Scan(&msg)
		conn.Write([]byte("say|" + nick + "|" + msg))
		if msg == "quit" {
			conn.Write([]byte("quit|" + nick))
			break
		}
	}
 
}
 
func Handle(conn net.Conn) {
 
	for {
 
		data := make([]byte, 255)
		msg_read, err := conn.Read(data)
		if msg_read == 0 || err != nil {
			break
		}
 
		fmt.Println(string(data[0:msg_read]))
	}
