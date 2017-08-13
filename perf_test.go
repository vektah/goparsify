package goparsify

import "testing"

func BenchmarkAny(b *testing.B) {
	p := Any("hello", "goodbye", "help")

	for i := 0; i < b.N; i++ {
		_, _ = Run(p, "hello")
		_, _ = Run(p, "hello world")
		_, _ = Run(p, "good boy")
		_, _ = Run(p, "help me")
	}
}
