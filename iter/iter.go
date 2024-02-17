package iter

import coro "github.com/nobishino/gocoro/coro"

// Seqは個々の値のシーケンスに対するイテレータです。seq(yield)として呼び出されると、seqはシーケンス内の各値vに対してyield(v)を呼び出し、yieldがfalseを返した場合は早期に停止します。
type Seq[V any] func(yield func(V) bool)

// Seq2は値のペアのシーケンスに対するイテレータです。一般的にはキーと値のペアです。
// seq(yield)として呼び出されると、seqはシーケンス内の各ペア(k, v)に対してyield(k, v)を呼び出し、yieldがfalseを返した場合は早期に停止します。
type Seq2[K, V any] func(yield func(K, V) bool)

// Pullは「プッシュスタイル」のイテレータシーケンスseqを、「プルスタイル」のイテレータとしてアクセスするための2つの関数nextとstopに変換します。
// Nextはシーケンス内の次の値と、その値が有効かどうかを示すブール値を返します。シーケンスが終了した場合、nextはゼロ値のVとfalseを返します。
// シーケンスの終わりに到達した後や、stopを呼び出した後にnextを呼び出すことは有効です。これらの呼び出しは、ゼロ値のVとfalseを続けて返します。
// Stopはイテレーションを終了します。次の値に興味がなくなり、nextがまだシーケンスの終了を示していない場合に呼び出す必要があります。stopを複数回呼び出すことや、nextが既にfalseを返した場合に呼び出すことは有効です。
// 複数のゴルーチンから同時にnextまたはstopを呼び出すことはエラーです。
func Pull[V any](seq Seq[V]) (next func() (V, bool), stop func()) {
	copush := func(more bool, yield func(V) bool) V {
		if more {
			seq(yield)
		}
		var zero V
		return zero
	}
	resume := coro.New(copush)
	next = func() (V, bool) {
		return resume(true)
	}
	stop = func() {
		resume(false)
	}
	return next, stop
}

// Pull2は「プッシュスタイル」のイテレータシーケンスseqを、「プルスタイル」のイテレータとしてアクセスするための2つの関数nextとstopに変換します。
// Nextはシーケンス内の次のペアと、そのペアが有効かどうかを示すブール値を返します。
// シーケンスが終了した場合、nextはゼロ値のペアとfalseを返します。シーケンスの終わりに到達した後や、stopを呼び出した後にnextを呼び出すことは有効です。これらの呼び出しは、ゼロ値のペアとfalseを続けて返します。
// Stopはイテレーションを終了します。次の値に興味がなくなり、nextがまだシーケンスの終了を示していない場合に呼び出す必要があります。stopを複数回呼び出すことや、nextが既にfalseを返した場合に呼び出すことは有効です。
// 複数のゴルーチンから同時にnextまたはstopを呼び出すことはエラーです。
func Pull2[K, V any](seq Seq2[K, V]) (next func() (K, V, bool), stop func()) {
	type KV struct {
		K K
		V V
	}
	copush := func(more bool, yield func(KV) bool) KV {
		if more {
			seq(func(k K, v V) bool {
				return yield(KV{k, v})
			})
		}
		var zero KV
		return zero
	}
	resume := coro.New(copush)
	next = func() (K, V, bool) {
		kv, ok := resume(true)
		return kv.K, kv.V, ok
	}
	stop = func() {
		resume(false)
	}
	return next, stop
}
