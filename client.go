// client
package main

import (
	protocol "./protocol"
	utils "./utils"
	"fmt"
	"net"
	"os"
	"runtime"
	"time"
)

var cmdClient = &Command{
	Run:               runClient,
	UsageLine:         "client",
	Short:             "start up a client",
	Long:              "start up a client",
	DefaultParameters: ":1200",
}

// func main() {
func runClient(cmd *Command, args []string) bool {

	fmt.Fprintf(os.Stderr, "number of cpus:%d:\n", utils.NCPU)
	fmt.Fprintf(os.Stderr, "number of cpus as reported by go:%d:\n", runtime.NumCPU())
	runtime.GOMAXPROCS(utils.NCPU)

	service := cmd.DefaultParameters

	if len(args) > 0 {
		service = args[0]
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	utils.CheckError(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	utils.CheckError(err)

	defer conn.Close()
	utils.Log("connect success")

	go sender(conn)
	for {
		time.Sleep(1 * 1e9)
	}

	return true
}

func sender(conn net.Conn) {
	for i := 0; i < 1000; i++ {
		words := "{\"Id\":1,\"Name\":\"golang\",\"Message\":\"message\"}"
		conn.Write(protocol.Packet([]byte(words)))
	}

	utils.Log("send over")
}
