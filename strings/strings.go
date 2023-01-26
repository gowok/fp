package strings

func Repeat(times int, input string) string {
	result := ""
	for i := 0; i < times; i++ {
		result += input
	}

	return result
}
