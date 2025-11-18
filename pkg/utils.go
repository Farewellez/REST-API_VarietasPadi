package pkg

// Result adalah tipe data yang menggabungkan nilai hasil dan error,
// mirip konsep Either/Result type di FP untuk error handling murni.
type Result[T any] struct {
	Value T
	Err   error
}

// Map adalah Higher-Order Function.
// Ia menerima slice 'in' dan fungsi 'f', lalu mengembalikan slice baru
// yang merupakan hasil penerapan 'f' ke setiap elemen 'in'.
func Map[T any, R any](in []T, f func(T) R) []R {
	out := make([]R, len(in))
	for i, v := range in {
		out[i] = f(v)
	}
	return out
}

// Filter adalah Higher-Order Function.
// Ia menerima slice 'in' dan fungsi 'pred' (predicate, fungsi pengecek),
// dan mengembalikan slice baru yang hanya berisi elemen yang memenuhi 'pred'.
func Filter[T any](in []T, pred func(T) bool) []T {
	out := []T{}
	for _, v := range in {
		if pred(v) {
			out = append(out, v)
		}
	}
	return out
}
