package main

import (
	"book-chat/trace"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type room struct {
	forward chan []byte      //forwardは他のクライアントに転送するためのメッセージを保持するチャネル
	join    chan *client     //チャットルームに参加しようとしてるクライアントのためのチャネル
	leave   chan *client     //チャットルームから退席しようとしているクライアントのためのチャネル
	clients map[*client]bool //在室している全てのクライアントが保持される
	tracer  trace.Tracer     //チャット上で行われた操作のログを受け取る
}

//newRoomではすぐに利用できるチャットルームを生成して返す
func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() { //チャットルームを開始
	for {
		select { //goroutineを使用する前にselectを置いている??複数のチャネルを同時に待ち状態にしたい時
		case client := <-r.join: //並行処理.r.joinに値が入ってきた時
			//参加
			r.clients[client] = true
			r.tracer.Trace("新しいクライアントが参加しました")
		case client := <-r.leave:
			//退室
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("クライアントが退室しました")
		case msg := <-r.forward:
			r.tracer.Trace("メッセージを受信しました：", string(msg)) //これから全員の.sendに入れていく＝送信していく
			//全てのクライアントにメッセージを転送
			for client := range r.clients {
				select {
				case client.send <- msg:
					//メッセージを受信
					r.tracer.Trace("--クライアントに送信されました")
				default:
					//送信に失敗
					delete(r.clients, client)
					close(client.send)
					r.tracer.Trace("--送信に失敗しました。クライアントをクリーンアップします")
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
	socket, err := upgrader.Upgrade(w, req, nil) //何してる？？HTTP通信からwebsocket通信に更新(具体的にはハンドシェイクしてる)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	client := &client{ //クライアントの構造体のclientを生成.クライアント作ってるけどこれは誰??
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client                     //新しいクライアントが生成されたら、そいつをroomにjoinさせる
	defer func() { r.leave <- client }() //最後は退室させる.いつ呼ばれるの？？
	go client.write()                    //メッセージを受信状態.goroutineとして実行
	client.read()
}
