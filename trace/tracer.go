package trace

import (
	"fmt"
	"io"
)

//Tracerはログを記録できるオブジェクトを表すインターフェース
type Tracer interface {
	Trace(...interface{}) //任意の引数の型を何個でも取れる
}
type tracer struct {
	out io.Writer //ここに情報が出力されていく
}

func New(w io.Writer) Tracer { //出力する場所(標準出力、ファイル、csv等)を指定できる
	return &tracer{out: w}
}

//Tracerインターフェースで要求しているメソッドを定義
func (t *tracer) Trace(a ...interface{}) {
	t.out.Write([]byte(fmt.Sprint(a...)))
	t.out.Write([]byte("\n"))
}

type nilTracer struct{}

//なんのためのTraceメソッド??roomのstructで、trace の型はtrace.Tracerになっている。Tracerインターフェースは、Traceメソッドを実装している必要があるので空っぽだけどTraceメソッドをつけた
func (t *nilTracer) Trace(a ...interface{}) {}

func Off() Tracer {
	return &nilTracer{}
}
