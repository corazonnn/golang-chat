package trace

import "io"

//Tracerはログを記録できるオブジェクトを表すインターフェース
type Tracer interface {
	Trace(...interface{}) //任意の引数の型を何個でも取れる
}

func New(w io.Writer) Tracer { return nil }
