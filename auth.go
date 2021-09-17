package main

import "net/http"

type authHandler struct {
	next http.Handler //wrap用かな
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		//未認証
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)

	} else if err != nil {
		//なんらかの別のエラーが発生
		panic(err.Error())
	} else {
		//成功。ラップされたハンドラを呼び出します
		h.next.ServeHTTP(w, r)
	}
}

//これわからん
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}
