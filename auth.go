package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/stretchr/gomniauth"
	gomniauthcommon "github.com/stretchr/gomniauth/common"
	"github.com/stretchr/objx"
)

type ChatUser interface {
	UniqueID() string //変数UniqueIDを定義するのではなく、UniqueIDを探してくるUniqueID()メソッドを定義している
	AvatarURL() string
}
type chatUser struct {
	gomniauthcommon.User //gomniauth.common.Userという型が埋め込まれ、chatUserは自動的にGomniauthのUserインターフェースを実装したことになってる
	uniqueID             string
}

func (u chatUser) UniqueID() string {
	return u.uniqueID
}

//loginHandlerはサードパーティへのログインの処理を受け持つ.アプリ内ではなくどこか外部へ認証を行うってこと?
//パスの形式：　/auth/{action}/{provider}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/") //segsの中身は [][auth][action][provider]
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":
		//googleやgithubに対応する認証プロバイダのオブジェクトを取得.SharedProviderListから探してくる
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("認証プロバイダの取得に失敗しました：", provider, "-", err)
		}
		//認証プロセスを開始するためのURL.アカウントの選択の画面
		loginUrl, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			log.Fatalln("GetBeginAuthURLの呼び出し中にエラーが発生しました：", provider, "-", err)
		}
		//レスポンスとしてリダイレクトしたい
		w.Header().Set("Location", loginUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)

	case "callback": //認証後のリダイレクト先

		//①認証プロバイダのオブジェクトを取得
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("認証プロバイダーの取得に失敗しました")
		}

		//②
		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			log.Fatalln("認証を完了できませんでした", provider, "-", err)
		}

		//③認証情報を使ってユーザーの情報を取得できる
		user, err := provider.GetUser(creds)
		if err != nil {
			log.Fatalln("ユーザーの取得に失敗しました", provider, "-", err)
		}
		chatUser := &chatUser{User: user}
		m := md5.New()
		io.WriteString(m, strings.ToLower(user.Name()))
		chatUser.uniqueID = fmt.Sprintf("%x", m.Sum(nil))
		avatarURL, err := avatars.GetAvatarURL(chatUser)
		if err != nil {
			log.Fatalln("GetAvatarURLに失敗しました", "-", err)
		}
		//④Nameフィールドの値をエンコード
		authCookieValue := objx.New(map[string]interface{}{
			"userid":     chatUser.uniqueID,
			"name":       user.Name(),
			"avatar_url": avatarURL,
			"email":      user.Email(),
		}).MustBase64()

		//⑤Cookieに保持する
		http.SetCookie(w, &http.Cookie{Name: "auth",
			Value: authCookieValue,
			Path:  "/"})

		//⑥本来のアクセス先である/chatにリダイレクト
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)

	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "アクション%sには非対応です", action)

	}
}

type authHandler struct {
	next http.Handler //wrap用かな
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//Cookieがまだない || ログアウト状態(ログアウトする際にcookie.Value = ""にしてるから)
	if cookie, err := r.Cookie("auth"); err == http.ErrNoCookie || cookie.Value == "" {
		//未認証
		w.Header().Set("Location", "/login")        //http.ResponseWriterに対してHeader呼び出す
		w.WriteHeader(http.StatusTemporaryRedirect) //http.ResponseWriterに対してWriteHeader呼び出す。ログインページにリダイレクト

	} else if err != nil {
		//なんらかの別のエラーが発生
		panic(err.Error())
	} else {
		//成功。ラップされたハンドラを呼び出します
		h.next.ServeHTTP(w, r)
	}
}

//http.Handlerをラップした*authHandlerを生成
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}
