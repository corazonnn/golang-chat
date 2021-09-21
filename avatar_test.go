package main

import "testing"

//GetAvatarURLのテスト
func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	//適当に今作成したクライアントをGetAvatarURL()に渡して,ErrNoAvatarURLが変えることの確認
	client := new(client)
	url, err := authAvatar.GetAvatarURL(client)
	if err != ErrNoAvatarURL {
		t.Error("値が存在しない場合、AuthAvatar.GetAvatarURLはErrNoAvatarURLを返すべきです")
	}
	//値をセットした上でGetAvatarURL()に渡して,値が返ってくることを確認する
	testUrl := "http://url-to-avatar/"
	client.userData = map[string]interface{}{"avatar_url": testUrl}
	url, err = authAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("値が存在する場合、AuthAvatar.GetAvatarURLはエラーを返すべきではありません")
	} else {
		if url != testUrl {
			t.Error("AuthAvatar.GetAvatarURLは正しいURLを返すべきです")
		}
	}
}
