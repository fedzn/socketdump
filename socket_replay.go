package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func process(conn net.Conn) {
	defer conn.Close()

	lines := load_data(LogFilePath)
	t := time.Now().UnixNano()
	log.Printf("new client [%v,%v]", t, conn)

	for index, line := range lines {
		vs := strings.Split(line, ",")
		log.Printf("%d:%s %s %s", index, vs[0], vs[1], vs[2])

		hex_data := vs[2]
		decoded, err := hex.DecodeString(hex_data)
		if err != nil {
			log.Printf("DecodeString [%v]", err)
			continue
		}

		if len, err := conn.Write(decoded); err != nil {
			log.Printf("Write [%d] Error[%v]", len, err)
			break
		}

		time.Sleep(time.Duration(100) * time.Millisecond)
	}
}

func load_data(file_path string) (lines []string) {
	f, err := os.Open(file_path)
	if err != nil {
		log.Printf("Open Failed:[%v]", err)
		return
	}
	defer f.Close()

	br := bufio.NewReader(f)
	for {
		s, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}

		line := string(s)
		// fmt.Println(line)
		lines = append(lines, line)
	}

	return lines
}

var (
	Help        bool   // help info
	Interface   string // 支持的数据类型[tcp udp serial socketcan]
	Port        uint   // 服务端口
	LogFilePath string // 输入文件路径
)

// 打印帮助信息
func usage() {
	fmt.Println(`Usage: client [-hi:p:l:]`)
	flag.PrintDefaults()
	os.Exit(0)
}

func main() {
	flag.BoolVar(&Help, "h", false, "useage help")
	flag.StringVar(&Interface, "i", "tcp", "select interface[tcp|udp]")
	flag.UintVar(&Port, "p", 5017, "listen port")
	flag.StringVar(&LogFilePath, "l", "./data", "input logfile path")
	flag.Parse()

	if Help {
		usage()
	}

	log.SetFlags(0)

	// 监听TCP 服务端口
	listener, err := net.Listen(Interface, fmt.Sprintf(":%d", Port))
	if err != nil {
		log.Printf("Listen Failed:[%v]", err)
		os.Exit(1)
	}

	for {
		// 建立socket连接
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Listen.Accept Failed:[%v]", err)
			continue
		}

		// 业务处理逻辑
		go process(conn)
	}
}
