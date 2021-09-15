package main

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
