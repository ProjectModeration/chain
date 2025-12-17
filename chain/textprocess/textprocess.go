package textprocess

import "strings"

func NormalizeText(s string) string {
	var b strings.Builder
	for _, c := range strings.ToLower(s) {
		if (c >= 'a' && c <= 'z') || c == ' ' {
			b.WriteRune(c)
		}
	}
	return b.String()
}

func ChiScore(s string) float64 {
	return chiSquaredScore(s)
}

func ApplyROT13(s string) string {
	return caesarShift(s, 13)
}

func ApplyCaesar(s string, shift int) string {
	return caesarShift(s, shift)
}
