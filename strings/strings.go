package strings

func Repeat(input string, times int) string {
	if times <= 0 {
		return ""
	}

	result := ""
	for i := 0; i < times; i++ {
		result += input
	}

	return result
}
