package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
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
			return fmt.Sprintf("//www.gravatar.com/avatar/" + useridStr), nil //算出したハッシュ値をGravatarURLに埋め込む //fmt.printfは標準出力に表示。ただこの文字列をstring型として変数に受け取りたい。その場合,fmt.Sprintfを使う
		}
	}
	return "", ErrNoAvatarURL
}

type FileSystemAvatar struct{}

var UseFileSystemAvatar FileSystemAvatar

func (_ FileSystemAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			if files, err := ioutil.ReadDir("avatars"); err == nil { //ディレクトリの内容を一括して読み込む
				for _, file := range files {
					if file.IsDir() { //  /avatars/ファイル　だけとは限らない. avatars/ディレクトリ　の可能性もあるから排除
						continue
					}
					fmt.Println("file is", file)
					if match, _ := filepath.Match(useridStr+"*", file.Name()); match { //useridStrの部分がマッチしてるかどうか.jpeg とか.pngとかは気にせず
						return "/avatars/" + file.Name(), nil
					}
				}
			}
		}
	}
	return "", ErrNoAvatarURL
}
