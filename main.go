package main

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

//【目的】ファイルからテンプレートを作成し、データを出力する
type templateHandler struct { //テンプレートの構造
	once     sync.Once //一度だけ実行
	filename string
	templ    *template.Template //パースしたファイルの内容を入れる
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, u *http.Request) {
	t.once.Do(func() {
		t.templ =
			template.Must(template.ParseFiles(filepath.Join("templates",
				t.filename)))
	})
	t.templ.Execute(w, nil)
}

func main() {
	r := newRoom()
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r) //この書き方はありなのか？？
	//チャットルームを開始
	go r.run()
	//webサーバーの起動
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
