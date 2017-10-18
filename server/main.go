package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
)

var connMaps = make(map[string]net.Conn)
var nicknameMaps = make(map[string]string)

func main() {
	log.Printf("The server is listening on 127.0.0.1:1798\n")
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
			nickname := scanner.Text()
			if nickname == "" {
				nickname = "user from "+addr
			}
			nicknameMaps[addr] = nickname
			// 保存会话
			_, ok := connMaps[addr]
			if !ok {
				connMaps[addr] = conn
			}
			conn.Write([]byte("Hello " + nickname  +  ",chat with others and enjoy it.\n"))
			conn.Write([]byte(strconv.Itoa(len(connMaps))+" person at the chatting home this time!\n"))
			// 遍历map发送欢迎消息
			withMapTraversal(addr,func(_addr string,connect net.Conn){
				mustWrite(connect,[]byte( "\nWelcome " + scanner.Text() + " to join in the chatting home! "),_addr)
				mustWrite(connect,[]byte( strconv.Itoa(len(connMaps))+" person at the chatting home this time!\n\n"),_addr)
			})
			isNickname = false
		} else {
			// 检测到输入,先在发话的窗口输出
			conn.Write([]byte(nicknameMaps[addr] + "(you said):" + scanner.Text() + "\n"))
			// 遍历map发送消息
			withMapTraversal(addr,func(_addr string,connect net.Conn){
				mustWrite(connect,[]byte(nicknameMaps[addr] + ":" + scanner.Text() + "\n"),_addr)
			})
		}

	}
	conn.Close()
}

func withMapTraversal(addr string,fn func(string,net.Conn)){
	for _addr, connect := range connMaps {
		if _addr == addr {
			// 不发自己
			continue
		}
		fn(_addr,connect)
	}
}

func mustWrite(conn net.Conn,word []byte,addr string){
	i, err := conn.Write(word)
	if err != nil {
		fmt.Printf("err-----------%v,i----------->%d", err, i)
		delete(connMaps, addr)
	}
}
