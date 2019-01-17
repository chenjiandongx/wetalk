package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

var (
	// 连接客户端
	clients = make(map[*websocket.Conn]bool)
	// 广播通道
	broadcast = make(chan string)
	// 服务端端口
	serverPort int
	// websockets Upgrader
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	// 服务端命令
	cmdServer = &cobra.Command{
		Use:   "server ",
		Short: "start websockets server",
		Run: func(cmd *cobra.Command, args []string) {
			startServer(serverPort)
		},
	}
)

func init() {
	cmdServer.Flags().IntVarP(&serverPort, "port", "p", 8087, "server port")
}

func startServer(port int) {

	http.HandleFunc("/", handleConnections)

	go handleMessages()

	log.Println(fmt.Sprintf("Websocket server started on :%d", port))
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade 初始化
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	// 客户端注册
	clients[ws] = true

	for {
		// messageType = 1
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		// 广播信息
		broadcast <- string(msg)
	}
}

func handleMessages() {
	for {
		// 接收下一条广播信息
		msg := <-broadcast
		// 广播信息到所有客户端
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
