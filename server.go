//server
package main

import (
	"./packet"
	"./tcpsession"
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

func sendFinish(s interface{}, wpk *packet.Wpacket) {
	session := s.(*tcpsession.Tcpsession)
	session.Close()
}

func processClient(session *tcpsession.Tcpsession, rpk *packet.Rpacket) {
	session.Send(packet.NewWpacket(rpk.Buffer(), rpk.IsRaw()), sendFinish)
}

func sessionClose(session *tcpsession.Tcpsession) {
	fmt.Printf("client disconnect\n")
}

// func main() {
func runServer(cmd *Command, args []string) bool {

	fmt.Fprintf(os.Stderr, "number of cpus:%d:\n", NCPU)
	fmt.Fprintf(os.Stderr, "number of cpus as reported by go:%d:\n", runtime.NumCPU())
	runtime.GOMAXPROCS(NCPU)

	service := cmd.DefaultParameters

	if len(args) > 0 {
		service = args[0]
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to accept:retrying:%t:\n", err)
			continue
		}
		session := tcpsession.NewTcpSession(conn, true)
		go tcpsession.ProcessSession(session, processClient, sessionClose)
	}

	return true

}
