package textprocess

func rot13(s string) string {
	return caesarShift(s, 13)
}

func DetectROT13(s string) (bool, string, int) {
	if len(s) < 6 {
		return false, "", 0
	}

	decoded := rot13(s)

	scoreOriginal := chiSquaredScore(s)
	scoreDecoded := chiSquaredScore(decoded)

	// how much better decoded is vs original
	gap := scoreOriginal - scoreDecoded

	if gap <= 20 {
		return false, "", 0
	}

	conf := int((gap / (gap + scoreDecoded)) * 100)
	if conf > 99 {
		conf = 99
	}
	if conf < 0 {
		conf = 0
	}

	return true, decoded, conf
}
