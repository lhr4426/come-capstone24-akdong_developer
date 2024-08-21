package main

import (
	"GameServer/controller"
	controllerdb "GameServer/controller-db"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/logrusorgru/aurora"
)

func main() {

}

// 240821 update (udp server and http server)

func startUdp() {
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

	go GetOutboundIP()
	go controller.ConsoleController(conn)
	go controllerdb.ConnectDB()

	// 클라이언트 요청 처리
	for {
		// 클라이언트 요청
		controller.GetRequest(conn)
	}
}

func startHttp() {
	fpLog, err := os.OpenFile("logfile.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer fpLog.Close()

	// 표준로거를 파일로그로 변경
	log.SetOutput(fpLog)

	router := gin.Default()

	configs.ConnectDB()

	routes.MapRoute(router)

	log.Println("Server is starting...")
	err = router.Run("0.0.0.0:8070")
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
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
