package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type room struct {
	forward chan []byte      //forwardは他のクライアントに転送するためのメッセージを保持するチャネル
	join    chan *client     //チャットルームに参加しようとしてるクライアントのためのチャネル
	leave   chan *client     //チャットルームから退席しようとしているクライアントのためのチャネル
	clients map[*client]bool //在室している全てのクライアントが保持される
}

func (r *room) run() {
	for {
		select { //goroutineを使用する前にselectを置いている??
		case client := <-r.join: //並行処理.r.joinに値が入ってきた時
			//参加
			r.clients[client] = true
		case client := <-r.leave:
			//退室
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			//全てのクライアントにメッセージを転送
			for client := range r.clients {
				select {
				case client.send <- msg:
					//メッセージを受信
				default:
					//送信に失敗
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024 //ソケットのバッファサイズ
	messageBufferSize = 256  //メッセージのバッファサイズ
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

//*room型をhttp.Handler型に適合(ServeHTTPメソッドを定義)することで,*roomはHTTPハンドラとして扱えるようになる
func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//websocketを利用するためにUpgradeする必要がある
	socket, err := upgrader.Upgrade(w, req, nil) //何してる？？websocketコネクションの取得
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	client := &client{ //クライアントの構造体のclientを生成
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client                     //新しいクライアントが生成されたら、そいつをroomにjoinさせる
	defer func() { r.leave <- client }() //最後は退室させる.いつ呼ばれるの？？
	go client.write()                    //メッセージを受信状態
	client.read()
}
