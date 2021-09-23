package main

import (
	"errors"
	"fmt"
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
	if userid, ok := c.userData["userid"]; ok { //ユーザー情報にemailが存在する
		if useridStr, ok := userid.(string); ok { //string型に直すことができた
			return fmt.Sprintf("//www.gravatar.com/avatar/" + useridStr), nil //算出したハッシュ値をGravatarURLに埋め込む
		}
	}
	return "", ErrNoAvatarURL
}
