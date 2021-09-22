package main

import "errors"

//AvatarインスタンスがアバターのURLを返すことができない場合に発生するエラー
var ErrNoAvatarURL = errors.New("chat: アバターのURLを取得できません")

type Avatar interface { //①Avatarインターフェースの機能が関数として入ってる
	//指定されたアバターのURLを返す
	//問題発生したらエラー、URLを取得できなかったらErrNoAvatarURLを返す
	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct{} //②AuthAvatarの構造体を定義

var UserAuthAvatar AuthAvatar //④structからオブジェクトを生成

func (_ AuthAvatar) GetAvatarURL(c *client) (string, error) { //③インターフェースを満たすようにメソッドを定義
	//クライアントのユーザー情報の中からavatar_urlをとってくる
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

//④以降は、このメソッドを使うための処理を記述していく
//④structからオブジェクトを生成
//⑤Avatarインターフェースを持つオブジェクトに定義
//⑥オブジェクト.メソッドで使う
