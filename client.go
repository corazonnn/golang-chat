package main

import (
	"github.com/gorilla/websocket"
)

//ある一つのクライアントへの接続
type client struct {
	socket *websocket.Conn //他の一人のユーザーとのやりとりをする道(websocket)の参照
	send   chan []byte     //受信したメッセージ
	room   *room           //チャットしてるroomへの参照.いる？クライアントとの接続の管理、メッセージのルーティング.clientを管理すものみたいな感じ?
}

func (c *client) read() { //クライアントがwebsocketからデータを読み込む
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil { //readmesageでエラーが出ないなら
			c.room.forward <- msg //msgが読みこれるまで待つ.goroutine使用
		} else {
			break
		}
	}
	c.socket.Close() //websocketを閉じる.通信をやめる
}
func (c *client) write() { //writeってwebsocketからプログラム内に書き出してる
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
