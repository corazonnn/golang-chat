package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

func uploaderHandler(w http.ResponseWriter, req *http.Request) {
	userId := req.FormValue("userid")               //htmlのformから取り出す(認証した時点でハッシュ値を算出していて、それをuseridに入れてる)
	file, header, err := req.FormFile("avatarFile") //reqestの中から取り出すからレシーバに指定してる
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file) //ファイルを読み込んで[]byteにしてる
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	filename := filepath.Join("avatars", userId+filepath.Ext(header.Filename)) //filepath.Ext:ファイル名の拡張子を返す
	err = ioutil.WriteFile(filename, data, 0777)                               //filenameに権限が0777、中身はdataで書き込む
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	io.WriteString(w, "成功")

}
