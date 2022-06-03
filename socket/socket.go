package socket

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	WriteBufferPool: &sync.Pool{},
	ReadBufferSize:  1024,
}

// GetWSConnectClient Проверка подключения к сокету от Клиента
func GetWSConnectClient(url string) (*websocket.Conn, error) {
	var conn, _, err = websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		panic("Подключение на работает")
		return nil, err
	}
	return conn, nil
}

// WriteSocket запись полученного ответа от клиента
func WriteSocket(conn *websocket.Conn, payload []byte) {
	if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
		log.Printf("Ошибка в файле socket.go в WriteSocket # Ошибка %v\n", err)
	}
}

// ReadSocketClient Чтение данных из сокета клиента и отправка их на postUser
func ReadSocketClient(w http.ResponseWriter, r *http.Request) []byte {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Ошибка в ReadSocketServer в Upgrade #%v\n", err)
		return nil
	}
	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Printf("Ошибка в ReadSocketClient #%v\n\n", err)
		return nil
	}
	log.Println("Прочтено из WebSocketClient " + string(msg))
	return msg
}
