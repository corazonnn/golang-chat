package main

import (
	"time"

	"github.com/gorilla/websocket"
)

//ある一つのクライアントへの接続
type client struct {
	socket   *websocket.Conn //他の一人のユーザーとのやりとりをする道(websocket)の参照
	send     chan *message   //受信したメッセージ
	room     *room           //チャットしてるroomへの参照.いる？クライアントとの接続の管理、メッセージのルーティング.clientを管理すものみたいな感じ?
	userData map[string]interface{}
}

func (c *client) read() { //クライアントがwebsocketからデータを読み込む
	for {
		var msg *message
		//クライアント側から送られたJSON文字列をmessageオブジェクトに変換.クライアント側からはmessageのみ送られる.時間と名前は自分でとってくる
		if err := c.socket.ReadJSON(&msg); err == nil { //chat.htmlからsocket.sendによってwebsocket上に乗せられたデータの参照
			msg.When = time.Now()
			msg.Name = c.userData["name"].(string)
			msg.AvatarURL, _ = c.room.avatar.GetAvatarURL(c)   //Gravatar
			if avatarURL, ok := c.userData["avatar_url"]; ok { //Googleアカウント
				msg.AvatarURL = avatarURL.(string)
			}
			c.room.forward <- msg //msgが読みこれるまで待つ.goroutine使用
		} else {
			break
		}
	}
	c.socket.Close() //websocketを閉じる.通信をやめる
}
func (c *client) write() { //writeってwebsocketからプログラム内に書き出してる
	for msg := range c.send {
		//messageオブジェクトをサーバー側から、JSON形式に変換してクライアント側に送る
		if err := c.socket.WriteJSON(msg); err != nil { //これって具体的にはどこに飛ばしてんの??
			break
		}
	}
	c.socket.Close()
}
