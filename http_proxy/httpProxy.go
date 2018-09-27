package http_proxy

import (
	"log"
	"net"
	"fmt"
	"io"
	"net/textproto"
	"bufio"
	"strings"
)

//开启一个代理服务
func Start()  {
	log.SetFlags(log.LstdFlags|log.Lshortfile)
	l, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("proxy listent in 8081 success, waiting for connection ...")
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("accept error Retry")
			continue
		}
		go handleClientRequest(conn)
	}
}
// 处理客户端与代理服务器连接
func handleClientRequest(client net.Conn)  {
	if client == nil {
		return
	}
	defer client.Close()
	who := client.RemoteAddr()
	fmt.Println("handle RemoteAddr：", who)


	buf := make([]byte, 0, 1024)
	n, err := client.Read(buf[:])
	if err != nil {
		log.Println(err)
		return
	}
	//获得了请求的host和port，就开始拨号吧
	server, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Println(err)
		return
	}
	server.Write(buf[:n])
	//进行转发
	go io.Copy(server, client)
	io.Copy(client, server)

}

// 获得客户端 request header 信息
func getClientInfo(client net.Conn) {
	var brc = bufio.NewReader(client)
	tp := textproto.NewReader(brc)
	// First line: GET /index.html HTTP/1.0
	var requestLine, err = tp.ReadLine()
	if err != nil {

	}
	fmt.Println(requestLine)
	// Subsequent lines: Key: value.
	mimeHeader, err := tp.ReadMIMEHeader()
	if err != nil {

	}
	fmt.Println(mimeHeader)
}


func parseRequestLine(line string) (method, requestURI, proto string, ok bool) {
	s1 := strings.Index(line, " ")
	s2 := strings.Index(line[s1+1:], " ")
	if s1 < 0 || s2 < 0 {
		return
	}
	s2 += s1 + 1
	return line[:s1], line[s1+1 : s2], line[s2+1:], true
}


