package textprocess

import "strings"

func caesarShift(s string, shift int) string {
	res := []rune{}

	for _, c := range s {
		if c >= 'a' && c <= 'z' {
			res = append(res, 'a'+(c-'a'-rune(shift)+26)%26)
		} else if c >= 'A' && c <= 'Z' {
			res = append(res, 'A'+(c-'A'-rune(shift)+26)%26)
		} else {
			res = append(res, c)
		}
	}
	return string(res)
}

func looksEnglish(s string) int {
	score := 0
	lower := strings.ToLower(s)

	common := []string{" the ", " and ", " to ", " of "}
	for _, w := range common {
		if strings.Contains(lower, w) {
			score += 5
		}
	}

	vowels := 0
	letters := 0
	for _, c := range lower {
		if c >= 'a' && c <= 'z' {
			letters++
			if strings.ContainsRune("aeiou", c) {
				vowels++
			}
		} else if c != ' ' {
			score -= 2
		}
	}

	if letters > 0 {
		ratio := float64(vowels) / float64(letters)
		if ratio > 0.3 && ratio < 0.6 {
			score += 5
		}
	}

	return score
}

var englishFreq = map[rune]float64{
	'a': 8.17, 'b': 1.49, 'c': 2.78, 'd': 4.25,
	'e': 12.70, 'f': 2.23, 'g': 2.02, 'h': 6.09,
	'i': 6.97, 'j': 0.15, 'k': 0.77, 'l': 4.03,
	'm': 2.41, 'n': 6.75, 'o': 7.51, 'p': 1.93,
	'q': 0.10, 'r': 5.99, 's': 6.33, 't': 9.06,
	'u': 2.76, 'v': 0.98, 'w': 2.36, 'x': 0.15,
	'y': 1.97, 'z': 0.07,
}

func DetectCaesar(s string) (bool, int, string, int) {
	if len(s) < 6 {
		return false, 0, "", 0
	}

	bestScore := 1e9
	secondBest := 1e9
	bestShift := 0
	bestText := ""

	for i := 1; i < 26; i++ {
		decoded := caesarShift(s, i)
		score := chiSquaredScore(decoded)

		if score < bestScore {
			secondBest = bestScore
			bestScore = score
			bestShift = i
			bestText = decoded
		} else if score < secondBest {
			secondBest = score
		}
	}

	gap := secondBest - bestScore

	if gap <= 20 {
		return false, 0, "", 0
	}

	// map gap → confidence %
	conf := int((gap / (gap + bestScore)) * 100)
	if conf > 99 {
		conf = 99
	}
	if conf < 0 {
		conf = 0
	}

	return true, bestShift, bestText, conf
}

func chiSquaredScore(s string) float64 {
	counts := map[rune]int{}
	total := 0

	for _, c := range strings.ToLower(s) {
		if c >= 'a' && c <= 'z' {
			counts[c]++
			total++
		}
	}

	if total == 0 {
		return 1e9
	}

	score := 0.0
	for c, expected := range englishFreq {
		observed := float64(counts[c])
		exp := expected * float64(total) / 100.0
		diff := observed - exp
		score += (diff * diff) / exp
	}

	return score
}
