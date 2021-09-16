package trace

import (
	"bytes"
	"testing"
)

//ユニットテスト
func TestNew(t *testing.T) {
	var buf bytes.Buffer //出力されるデータがbytes.Bufferに保持されている
	tracer := New(&buf)
	if tracer == nil {
		t.Error("newからの戻り値がnilです")
	} else {
		tracer.Trace("こんにちは、traceパッケージ")
		if buf.String() != "こんにちは、traceパッケージ\n" {
			t.Errorf("'%s'という誤った文字列が出力されました", buf.String())
		}
	}
	t.Error("まだテストを作成していません")
}
