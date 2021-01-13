package main

// 输出文件说明：
// tcp_dump -i tcp -s 127.0.0.1:5017 -o ./data
// 文件保存的默认名称为tcp_127.0.0.1:5017_20210113_102021678.raw
// go run socket_sink.go -i tcp -s 192.168.1.105:5017 -o ./data

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"
)

func checkError(err error) {
	if err != nil {
		log.Printf("Error: [%s]", err.Error())
		os.Exit(1)
	}
}

var (
	Help        bool   // help info
	Interface   string // 支持的数据类型[tcp udp serial socketcan]
	HostUrl     string // 服务地址
	OutFileName string // 输出文件名称
	OutFileDir  string // 输出文件目录
	OutFilePath string // 输出文件路径
)

func create_file(filepath string) *os.File {
	fileObj, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Failed to open the file", err.Error())
		os.Exit(2)
	}
	// defer fileObj.Close()

	return fileObj
}

// 打印帮助信息
func usage() {
	fmt.Println(`Usage: client [-hi:s:o:]`)
	flag.PrintDefaults()
	os.Exit(0)
}

func main() {
	flag.BoolVar(&Help, "h", false, "useage help")
	flag.StringVar(&Interface, "i", "tcp", "select interface[tcp|udp]")
	flag.StringVar(&HostUrl, "s", "127.0.0.1:9090", "connect service url")
	flag.StringVar(&OutFileDir, "o", "./data", "output directory")
	flag.Parse()

	if Help {
		usage()
	}

	log.SetFlags(0)

	conn, err := net.Dial(Interface, HostUrl)
	checkError(err)
	defer conn.Close()

	time_name := time.Now().Format("20060102_150405")
	OutFileName = fmt.Sprintf("%s_%s_%s.raw", Interface, HostUrl, time_name)

	if _, err := os.Stat(OutFileDir); os.IsNotExist(err) {
		if err := os.MkdirAll(OutFileDir, os.ModePerm); err != nil {
			log.Printf("MkdirAll Failed:[%v]", err)
			os.Exit(1)
		}
	}

	OutFilePath, _ = filepath.Abs(filepath.Join(OutFileDir, OutFileName))
	log.Println(OutFilePath)

	fileObj, err := os.OpenFile(OutFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Printf("OpenFile Failed:[%v]", err)
		os.Exit(2)
	}
	defer fileObj.Close()

	for {
		raw_data := make([]byte, 512)
		if len, err := conn.Read(raw_data); err != nil {
			log.Printf("Read Error:[%v] [%d]", err, len)
		} else {
			valid_data := raw_data[0:len]
			// log.Printf("%d:%x", len, valid_data)

			hex_data := hex.EncodeToString(valid_data)

			// decoded, err := hex.DecodeString(hex_data)
			// log.Printf("%s\n", decoded)

			time_stamp := time.Now().UnixNano()
			line := fmt.Sprintf("%v,%v,%s\n", time_stamp, len, hex_data)
			log.Printf(line)

			if _, err := fileObj.Write([]byte(line)); err != nil {
				log.Printf("Write Error[%v]", err)
			}
		}
	}
}
