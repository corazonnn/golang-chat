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
func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	//クライアント作成
	client := new(client)
	//メールアドレスは必ず必要
	// client.userData = map[string]interface{}{"email": "MyEmailAddress@example.com"}
	client.userData = map[string]interface{}{"userid": "0bc83cb571cd1c50ba6f3e8a78ef1346"}
	//作成したクライアントのアバター画像URLを取得する
	url, err := gravatarAvatar.GetAvatarURL(client)
	//正しいURLかどうか確認
	if err != nil {
		t.Error("GravatarAvatar.GetAvatarURLはエラーを返すべきではありません")
	}
	if url != "//www.gravatar.com/avatar/0bc83cb571cd1c50ba6f3e8a78ef1346" {
		t.Errorf("GravatarAvatar.GetAvatarURLが%sという誤った値を返しました", url)
	}
}
