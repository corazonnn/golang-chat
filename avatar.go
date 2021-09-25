package main

import (
	"errors"
	"io/ioutil"
	"path/filepath"
)

//AvatarインスタンスがアバターのURLを返すことができない場合に発生するエラー
var ErrNoAvatarURL = errors.New("chat: アバターのURLを取得できません")

type Avatar interface { //①Avatarインターフェースの機能が関数として入ってる
	//指定されたアバターのURLを返す
	//問題発生したらエラー、URLを取得できなかったらErrNoAvatarURLを返す
	GetAvatarURL(ChatUser) (string, error)
}
type TryAvatars []Avatar

//Avatarインターフェースを満たしてる. = 個々のGetAvatarURL()を呼びだす代わりにTryAvatarsだけでいい
func (a TryAvatars) GetAvatarURL(u ChatUser) (string, error) {
	for _, avatar := range a {
		if url, err := avatar.GetAvatarURL(u); err == nil {
			return url, nil
		}
	}
	return "", ErrNoAvatarURL
}

//認証サーバーからアバターURLを取得
type AuthAvatar struct{} //②AuthAvatarの構造体を定義

var UserAuthAvatar AuthAvatar //④structからオブジェクトを生成 ⑤へ

func (_ AuthAvatar) GetAvatarURL(u ChatUser) (string, error) { //③インターフェースを満たすようにメソッドを定義
	//クライアントのユーザー情報の中からavatar_urlをとってくる
	url := u.AvatarURL()
	if url != "" {
		return url, nil
	}
	return "", ErrNoAvatarURL
}

//GravatarからアバターURLを取得
type GravatarAvatar struct{}

var UseGravatar GravatarAvatar

func (_ GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + u.UniqueID(), nil
}

type FileSystemAvatar struct{}

var UseFileSystemAvatar FileSystemAvatar

func (_ FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {
	if files, err := ioutil.ReadDir("avatars"); err == nil { //ディレクトリの内容を一括して読み込む
		for _, file := range files {
			if file.IsDir() { //  /avatars/ファイル　だけとは限らない. avatars/ディレクトリ　の可能性もあるから排除
				continue
			}
			useridStr := u.UniqueID()
			if match, _ := filepath.Match(useridStr+"*", file.Name()); match { //useridStrの部分がマッチしてるかどうか.jpeg とか.pngとかは気にせず
				return "/avatars/" + file.Name(), nil
			}
		}
	}

	return "", ErrNoAvatarURL
}
