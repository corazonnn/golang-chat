package main

import "time"

//クライアントとサーバーで送受信するデータ
type message struct {
	Name      string
	Message   string
	When      time.Time
	AvatarURL string
}
