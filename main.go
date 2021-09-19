package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

//【目的】ファイルからテンプレートを作成し、データを出力する
type templateHandler struct { //テンプレートの構造
	once     sync.Once //一度だけ実行
	filename string
	templ    *template.Template //パースしたファイルの内容を入れる
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ =
			template.Must(template.ParseFiles(filepath.Join("templates",
				t.filename)))
	})
	//なぜ
	data := map[string]interface{}{
		"Host": r.Host,
	}
	//data=>map[Host:localhost:8080]

	//もし認証されているなら
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	//data => map[Host:localhost:8080 UserData:map[name:相川佑也]]
	t.templ.Execute(w, data) //第二引数は、t.templの{{}}の中に何を入れているのか
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス") //フラグ addr を宣言し、そのデフォルト値を ":8080" とし、フラグの短い説明を与えている
	flag.Parse()                                             //コマンドラインの引数のフラグが解析され、フラグが変数にバインドされる
	//Gomniauthのセットアップ
	gomniauth.SetSecurityKey("corazon73647364")
	gomniauth.WithProviders(
		facebook.New("350210712146-9kfrktdjogrie9ru09mnclesuatlt4ib.apps.googleusercontent.com", "i0hMhZpjve1RfMzwWd4Bv8nM", "http://localhost:8080/auth/callback/facebook"),
		github.New("350210712146-9kfrktdjogrie9ru09mnclesuatlt4ib.apps.googleusercontent.com", "i0hMhZpjve1RfMzwWd4Bv8nM", "http://localhost:8080/auth/callback/github"),
		google.New("350210712146-9kfrktdjogrie9ru09mnclesuatlt4ib.apps.googleusercontent.com", "i0hMhZpjve1RfMzwWd4Bv8nM", "http://localhost:8080/auth/callback/google"),
	)
	r := newRoom()
	//①"/chat"にアクセス②MustAuth内でtemplateHandlerをラップしたauthHandlerが生成③authHandlerが生成されたことでauthHandlerのServeHTTPが呼ばれる(authというcookieの有無をチェック)
	//④認証成功したら、templateHandlerのServeHTTPが呼ばれる
	//④認証失敗したら、http.ResponseWriterに対してHeader,WriteHandlerを呼び出し、ログインページにリダイレクト
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"}) //MustAuthはいらない
	http.HandleFunc("/auth/", loginHandler)

	http.Handle("/room", r)
	//チャットルームを開始
	go r.run()
	//webサーバーの起動
	log.Println("webサーバーを開始します。ポート：", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
