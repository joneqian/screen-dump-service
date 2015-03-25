//server
package main

import (
	protocol "./protocol"
	utils "./utils"
	"fmt"
	"net"
	"os"
	"runtime"
)

var cmdServer = &Command{
	Run:               runServer,
	UsageLine:         "server",
	Short:             "start up a server",
	Long:              "start up a server",
	DefaultParameters: ":1200",
}

// func main() {
func runServer(cmd *Command, args []string) bool {

	fmt.Fprintf(os.Stderr, "number of cpus:%d:\n", utils.NCPU)
	fmt.Fprintf(os.Stderr, "number of cpus as reported by go:%d:\n", runtime.NumCPU())
	runtime.GOMAXPROCS(utils.NCPU)

	service := cmd.DefaultParameters

	if len(args) > 0 {
		service = args[0]
	}

	netListen, err := net.Listen("tcp", service)
	utils.CheckError(err)

	defer netListen.Close()

	utils.Log("Waiting for clients")
	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}

		utils.Log(conn.RemoteAddr().String(), " tcp connect success")
		go handleConnection(conn)
	}

	return true

}

func handleConnection(conn net.Conn) {
	//声明一个临时缓冲区，用来存储被截断的数据
	tmpBuffer := make([]byte, 0)

	//声明一个管道用于接收解包的数据
	readerChannel := make(chan []byte, 16)
	go reader(readerChannel)

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			utils.Log(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}

		tmpBuffer = protocol.Unpack(append(tmpBuffer, buffer[:n]...), readerChannel)
	}
}

func reader(readerChannel chan []byte) {
	for {
		select {
		case data := <-readerChannel:
			utils.Log(string(data))
		}
	}
}
