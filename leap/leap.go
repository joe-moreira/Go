package leap

func IsLeapYear(y int) bool {
	if (y % 4) == 0 {
		if (y % 100) == 0 {
			if (y % 400) == 0 {
				return true
			}
			return false
		}
		return true
	}
	return false
}
