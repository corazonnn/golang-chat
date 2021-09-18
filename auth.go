package main

import "net/http"

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
