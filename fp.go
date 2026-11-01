package fp

func Flip[T any, U any, V any](f func(T, U) V) func(U, T) V {
	return func(u U, t T) V { return f(t, u) }
}

func Pipe[T any](x T, fns ...func(T) T) T {
	for _, fn := range fns {
		x = fn(x)
	}
	return x
}

func Compose[T any](fns ...func(T) T) func(T) T {
	return func(x T) T {
		for i := len(fns) - 1; i >= 0; i-- {
			x = fns[i](x)
		}
		return x
	}
}
