package main

import (
	"GameServer/controller"
	controllerdb "GameServer/controller-db"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/logrusorgru/aurora"
)

func main() {

	controllerdb.ConnectDB()
	controllerdb.CreatorInit()

	// UDP 서버 소켓 생성
	addr, err := net.ResolveUDPAddr("udp", ":8050")
	if err != nil {
		fmt.Println("Error : resolving UDP address:", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error : listening:", err)
		return
	}
	defer conn.Close()

	fmt.Println(aurora.Green("============= Game Server Started ============="))

	for {
		controller.GetRequest(conn)
	}
}

func GetOutboundIP() {
	conn, err := net.Dial("udp", "0.0.0.0:8050")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)

	// fmt.Printf("Server Address : %s\n", localAddr.String())
	fmt.Println(
		aurora.Sprintf(
			aurora.Gray(12, "Server Address : %s"), localAddr.String()))
}

func ExitTask(sigChan chan os.Signal) {
	sig := <-sigChan
	fmt.Printf("\n\nReceived Signal : %s\n", sig)
	fmt.Println(aurora.Green("============= Game Server Closed ============="))

	os.Exit(0)
}
