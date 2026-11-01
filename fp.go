package fp

func Flip[T any, U any, V any](f func(T, U) V) func(U, T) V {
	return func(u U, t T) V { return f(t, u) }
}
