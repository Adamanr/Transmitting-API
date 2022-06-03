package server

import (
	"github.com/gin-gonic/gin"
	"log"
	"test/cmd"
	"test/socket"
)

var gs = *cmd.GetGlobalStorage()

// Handler Маршрутизатор
func Handler() {
	gins := gin.New()

	sockets := gins.Group("/socket")
	sockets.GET("/cs", SocketClient) // Сокет Клиент -> Сервер
	sockets.GET("/sc", SocketServer) // Сокет Сервер -> Клиент

	users := gins.Group("/users")
	//users.Use(config.GetMessageTokenCheck)
	users.POST("/getUser/:msgr", gs.GetUsers) //Получает данные от клиента и отправляет их на сокет
	users.POST("/getData", gs.GetData)        // Принимает данные от сервера и отправляет их на сокет

	gins.POST("/server", gs.Server)
	gins.POST("/client", gs.Client)

	if err := gins.Run("0.0.0.0:8080"); err != nil {
		log.Printf("Сервер не запущен из-за ошибки: %v\n", err)
		return
	}
}

// SocketClient - открывает сокет и читает ответ от Клиента
func SocketClient(ctx *gin.Context) {
	go gs.PostData(socket.ReadSocketClient(ctx.Writer, ctx.Request), "http://127.0.0.1:8080/client")
}

// SocketServer - открывает сокет и читает ответ от Сервера
func SocketServer(ctx *gin.Context) {
	go gs.PostData(socket.ReadSocketClient(ctx.Writer, ctx.Request), "http://127.0.0.1:8080/server")
}
