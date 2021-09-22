package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"strings"
)

//AvatarインスタンスがアバターのURLを返すことができない場合に発生するエラー
var ErrNoAvatarURL = errors.New("chat: アバターのURLを取得できません")

type Avatar interface { //①Avatarインターフェースの機能が関数として入ってる
	//指定されたアバターのURLを返す
	//問題発生したらエラー、URLを取得できなかったらErrNoAvatarURLを返す
	GetAvatarURL(c *client) (string, error)
}

//認証サーバーからアバターURLを取得
type AuthAvatar struct{} //②AuthAvatarの構造体を定義

var UserAuthAvatar AuthAvatar //④structからオブジェクトを生成 ⑤へ

func (_ AuthAvatar) GetAvatarURL(c *client) (string, error) { //③インターフェースを満たすようにメソッドを定義
	//クライアントのユーザー情報の中からavatar_urlをとってくる
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

//GravatarからアバターURLを取得
type GravatarAvatar struct{}

var UseGravatar GravatarAvatar

func (_ GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if email, ok := c.userData["email"]; ok { //ユーザー情報にemailが存在する
		if emailStr, ok := email.(string); ok { //string型に直すことができた
			m := md5.New()                                                      //MD5アルゴリズム
			io.WriteString(m, strings.ToLower(emailStr))                        //emailを全て小文字に直す & MD5アルゴリズムに適用
			return fmt.Sprintf("//www.gravatar.com/avatar/%x", m.Sum(nil)), nil //算出したハッシュ値をGravatarURLに埋め込む
		}
	}
	return "", ErrNoAvatarURL
}
