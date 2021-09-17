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
}

//以上終了を発生させないために,Traceメソッドを呼び出す前にOffメソッドによる「サイレント」なTraceを取得=どこにも出力しないtraceを得る
func TestOff(t *testing.T) {
	var silentTracer Tracer = Off()
	silentTracer.Trace("データ")
}
