package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/stretchr/gomniauth"
)

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
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "アクション%sには非対応です", action)

	}
}

type authHandler struct {
	next http.Handler //wrap用かな
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
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
