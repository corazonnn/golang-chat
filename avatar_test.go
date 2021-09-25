package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	gomniauthtest "github.com/stretchr/gomniauth/test"
)

//GetAvatarURLのテスト
func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	//適当に今作成したクライアントをGetAvatarURL()に渡して,ErrNoAvatarURLが変えることの確認
	// client := new(client)
	// url, err := authAvatar.GetAvatarURL(client)
	//単にクライアントを生成してgetavatarURL()につっこむのではなく、段階を経て突っ込む
	testUser := &gomniauthtest.TestUser{}
	testUser.On("avatarURL").Return("", ErrNoAvatarURL)
	testChatUser := &chatUser{User: testUser}
	url, err := authAvatar.GetAvatarURL(testChatUser)
	if err != ErrNoAvatarURL {
		t.Error("値が存在しない場合、AuthAvatar.GetAvatarURLはErrNoAvatarURLを返すべきです")
	}
	testUrl := "http://url-to-avatar/"
	testUser = &gomniauthtest.TestUser{}
	testChatUser.User = testUser
	testUser.On("avatarURL").Return(testUrl, nil)
	url, err = authAvatar.GetAvatarURL(testChatUser)
	if err != nil {
		t.Error("値が存在する場合、AuthAvatar.GetAvatarURLはエラーを返すべきではありません")
	} else {
		if url != testUrl {
			t.Error("AuthAvatar.GetAvatarURLは正しいURLを返すべきです")
		}
	}
}
func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	// client := new(client) //クライアント作成
	// client.userData = map[string]interface{}{"userid": "0bc83cb571cd1c50ba6f3e8a78ef1346"} //メールアドレスの代わりにuserid使用
	// url, err := gravatarAvatar.GetAvatarURL(client) //作成したクライアントのアバター画像URLを取得する
	user := &chatUser{uniqueID: "abc"}
	url, err := gravatarAvatar.GetAvatarURL(user)
	if err != nil { //正しいURLかどうか確認
		t.Error("GravatarAvatar.GetAvatarURLはエラーを返すべきではありません")
	}
	if url != "//www.gravatar.com/avatar/0bc83cb571cd1c50ba6f3e8a78ef1346" {
		t.Errorf("GravatarAvatar.GetAvatarURLが%sという誤った値を返しました", url)
	}
}
func TestFileSystemAvatar(t *testing.T) {
	//テスト用のアバターのファイルを生成
	filename := filepath.Join("avatars", "abc.jpg")
	ioutil.WriteFile(filename, []byte{}, 0777)
	//作成したファイルはテスト終了後に削除する
	defer func() { os.Remove(filename) }()

	//
	var fileSystemAvatar FileSystemAvatar
	// client := new(client)
	// client.userData = map[string]interface{}{"userid": "abc"}
	user := &chatUser{uniqueID: "abc"}
	url, err := fileSystemAvatar.GetAvatarURL(user)
	if err != nil {
		t.Error("FileSystemAvatar.GetAvatarURLはエラーを返すべきではありません")
	}
	if url != "/avatars/abc.jpg" {
		t.Errorf("FileSystemAvatar.GetAvatarURL%sという誤った値を返しました", url)
	}

}
