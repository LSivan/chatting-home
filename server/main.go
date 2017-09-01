package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
)

var conn_map = make(map[string]net.Conn)
var nickname_map = make(map[string]string)

func main() {

	listener, err := net.Listen("tcp", "localhost:1798")
	if err != nil {
		log.Fatalf("net.Listen err----------->%v\n", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go headleConn(conn)
	}
}

func headleConn(conn net.Conn) {

	addr := conn.RemoteAddr().String()
	fmt.Printf("connect successfully : %s\n", addr)

	conn.Write([]byte(fmt.Sprintf("Hello!Welcome to chatting home.Plz enter your nickname(Press Enter to finish):")))

	scanner := bufio.NewScanner(conn)
	isNickname := true
	for scanner.Scan() {
		// 第一次输入，保存昵称
		if isNickname {
			nickname_map[addr] = scanner.Text()
			// 保存会话
			_, ok := conn_map[addr]
			if !ok {
				conn_map[addr] = conn
			}
			conn.Write([]byte("Hello " + scanner.Text() +  ",chat with others and enjoy it.\n"))
			// 遍历map发送欢迎消息
			for _addr, connect := range conn_map {
				if _addr == addr {
					// 不发自己
					continue
				}
				mustWrite(connect,[]byte( "Welcome " + scanner.Text() + " to join in the chatting home! "),_addr)
				mustWrite(connect,[]byte( strconv.Itoa(len(conn_map))+" person at the chatting home this time!\n"),_addr)
			}
			isNickname = false
		} else {
			// 检测到输入,先在发话的窗口输出
			conn.Write([]byte(nickname_map[addr] + ":" + scanner.Text() + "\n"))
			// 遍历map发送消息
			for _addr, connect := range conn_map {
				if _addr == addr {
					// 不发自己
					continue
				}
				mustWrite(connect,[]byte(nickname_map[addr] + ":" + scanner.Text() + "\n"),_addr)
			}
		}

	}
	conn.Close()
}

func mustWrite(conn net.Conn,word []byte,addr string){
	i, err := conn.Write(word)
	if err != nil {
		fmt.Printf("err-----------%v,i----------->%d", err, i)
		delete(conn_map, addr)
	}
}