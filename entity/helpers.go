package entity

import (
	"regexp"
	"strconv"
)

func OrderedSlicesEqual[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func IsQuoted(s string) bool {
	regex := regexp.MustCompile(`"[^"]+"`)
	return len(regex.FindString(s)) > 0
}

func IsNumeric(s string) bool {
	if _, err := strconv.ParseFloat(s, 10); err != nil {
		return false
	}
	return true
}

func FloatToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
