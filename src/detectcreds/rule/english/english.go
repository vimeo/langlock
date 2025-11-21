package english

import (
	"math"
	"unicode"
)

var hundredMostCommonEnglishLetterPairs = map[[2]rune]float64{
	[2]rune{'t', 'h'}: 0.0330,
	[2]rune{'h', 'e'}: 0.0302,
	[2]rune{'a', 'n'}: 0.0181,
	[2]rune{'i', 'n'}: 0.0179,
	[2]rune{'e', 'r'}: 0.0169,
	[2]rune{'n', 'd'}: 0.0146,
	[2]rune{'r', 'e'}: 0.0133,
	[2]rune{'e', 'd'}: 0.0126,
	[2]rune{'e', 's'}: 0.0115,
	[2]rune{'o', 'u'}: 0.0115,
	[2]rune{'t', 'o'}: 0.0115,
	[2]rune{'h', 'a'}: 0.0114,
	[2]rune{'e', 'n'}: 0.0111,
	[2]rune{'e', 'a'}: 0.0110,
	[2]rune{'s', 't'}: 0.0109,
	[2]rune{'n', 't'}: 0.0106,
	[2]rune{'o', 'n'}: 0.0106,
	[2]rune{'a', 't'}: 0.0104,
	[2]rune{'h', 'i'}: 0.0097,
	[2]rune{'a', 's'}: 0.0095,
	[2]rune{'i', 't'}: 0.0093,
	[2]rune{'n', 'g'}: 0.0092,
	[2]rune{'i', 's'}: 0.0086,
	[2]rune{'o', 'r'}: 0.0084,
	[2]rune{'e', 't'}: 0.0083,
	[2]rune{'o', 'f'}: 0.0080,
	[2]rune{'t', 'i'}: 0.0076,
	[2]rune{'a', 'r'}: 0.0075,
	[2]rune{'t', 'e'}: 0.0075,
	[2]rune{'s', 'e'}: 0.0074,
	[2]rune{'m', 'e'}: 0.0068,
	[2]rune{'s', 'a'}: 0.0067,
	[2]rune{'n', 'e'}: 0.0066,
	[2]rune{'w', 'a'}: 0.0066,
	[2]rune{'v', 'e'}: 0.0065,
	[2]rune{'l', 'e'}: 0.0064,
	[2]rune{'n', 'o'}: 0.0060,
	[2]rune{'t', 'a'}: 0.0059,
	[2]rune{'a', 'l'}: 0.0057,
	[2]rune{'d', 'e'}: 0.0057,
	[2]rune{'o', 't'}: 0.0057,
	[2]rune{'s', 'o'}: 0.0057,
	[2]rune{'d', 't'}: 0.0056,
	[2]rune{'l', 'l'}: 0.0056,
	[2]rune{'t', 't'}: 0.0056,
	[2]rune{'e', 'l'}: 0.0055,
	[2]rune{'r', 'o'}: 0.0055,
	[2]rune{'a', 'd'}: 0.0052,
	[2]rune{'d', 'i'}: 0.0050,
	[2]rune{'e', 'w'}: 0.0050,
	[2]rune{'r', 'a'}: 0.0050,
	[2]rune{'r', 'i'}: 0.0050,
	[2]rune{'s', 'h'}: 0.0050,
	[2]rune{'h', 'o'}: 0.0049,
	[2]rune{'s', 'i'}: 0.0049,
	[2]rune{'d', 'a'}: 0.0048,
	[2]rune{'e', 'e'}: 0.0048,
	[2]rune{'o', 'm'}: 0.0048,
	[2]rune{'u', 't'}: 0.0048,
	[2]rune{'b', 'e'}: 0.0047,
	[2]rune{'e', 'm'}: 0.0047,
	[2]rune{'l', 'i'}: 0.0047,
	[2]rune{'o', 'w'}: 0.0046,
	[2]rune{'c', 'o'}: 0.0045,
	[2]rune{'e', 'c'}: 0.0045,
	[2]rune{'m', 'a'}: 0.0044,
	[2]rune{'u', 'r'}: 0.0044,
	[2]rune{'w', 'h'}: 0.0044,
	[2]rune{'s', 's'}: 0.0043,
	[2]rune{'r', 't'}: 0.0042,
	[2]rune{'d', 'o'}: 0.0041,
	[2]rune{'e', 'i'}: 0.0041,
	[2]rune{'l', 'o'}: 0.0041,
	[2]rune{'f', 'o'}: 0.0040,
	[2]rune{'l', 'a'}: 0.0040,
	[2]rune{'n', 'a'}: 0.0040,
	[2]rune{'a', 'i'}: 0.0039,
	[2]rune{'w', 'e'}: 0.0039,
	[2]rune{'w', 'i'}: 0.0039,
	[2]rune{'c', 'e'}: 0.0038,
	[2]rune{'c', 'h'}: 0.0038,
	[2]rune{'f', 't'}: 0.0037,
	[2]rune{'i', 'l'}: 0.0037,
	[2]rune{'i', 'm'}: 0.0037,
	[2]rune{'r', 's'}: 0.0037,
	[2]rune{'l', 'd'}: 0.0036,
	[2]rune{'o', 'o'}: 0.0036,
	[2]rune{'u', 'n'}: 0.0036,
	[2]rune{'y', 'o'}: 0.0036,
	[2]rune{'d', 's'}: 0.0035,
	[2]rune{'g', 'h'}: 0.0035,
	[2]rune{'u', 's'}: 0.0035,
	[2]rune{'t', 's'}: 0.0034,
	[2]rune{'u', 'l'}: 0.0034,
	[2]rune{'a', 'c'}: 0.0033,
	[2]rune{'e', 'h'}: 0.0033,
	[2]rune{'e', 'o'}: 0.0033,
	[2]rune{'i', 'd'}: 0.0033,
	[2]rune{'n', 'i'}: 0.0033,
	[2]rune{'n', 's'}: 0.0033,
	[2]rune{'h', 't'}: 0.0032,
	[2]rune{'i', 'c'}: 0.0032,
	[2]rune{'c', 'a'}: 0.0031,
	[2]rune{'l', 'y'}: 0.0031,
	[2]rune{'t', 'w'}: 0.0031,
	[2]rune{'e', 'f'}: 0.0030,
	[2]rune{'p', 'e'}: 0.0030,
	[2]rune{'k', 'e'}: 0.0029,
	[2]rune{'m', 'o'}: 0.0029,
	[2]rune{'w', 'o'}: 0.0029,
	[2]rune{'e', 'p'}: 0.0028,
	[2]rune{'g', 'e'}: 0.0028,
	[2]rune{'o', 's'}: 0.0028,
	[2]rune{'t', 'r'}: 0.0028,
	[2]rune{'i', 'r'}: 0.0027,
	[2]rune{'a', 'm'}: 0.0026,
	[2]rune{'a', 'y'}: 0.0026,
	[2]rune{'e', 'y'}: 0.0026,
	[2]rune{'o', 'l'}: 0.0026,
	[2]rune{'d', 'h'}: 0.0025,
	[2]rune{'a', 'f'}: 0.0025,
	[2]rune{'i', 'g'}: 0.0025,
	[2]rune{'m', 'i'}: 0.0025,
	[2]rune{'n', 'c'}: 0.0025,
	[2]rune{'e', 'v'}: 0.0024,
	[2]rune{'g', 'a'}: 0.0024,
	[2]rune{'i', 'o'}: 0.0024,
	[2]rune{'s', 'w'}: 0.0024,
	[2]rune{'e', 'b'}: 0.0023,
	[2]rune{'f', 'i'}: 0.0023,
	[2]rune{'g', 'o'}: 0.0023,
	[2]rune{'i', 'e'}: 0.0023,
	[2]rune{'p', 'a'}: 0.0023,
	[2]rune{'b', 'u'}: 0.0022,
	[2]rune{'b', 'u'}: 0.0021,
	[2]rune{'p', 'o'}: 0.0021,
	[2]rune{'r', 'y'}: 0.0021,
	[2]rune{'a', 'b'}: 0.0020,
	[2]rune{'a', 'p'}: 0.0020,
	[2]rune{'a', 'v'}: 0.0020,
	[2]rune{'d', 'b'}: 0.0020,
	[2]rune{'f', 'e'}: 0.0020,
	[2]rune{'r', 'd'}: 0.0020,
	[2]rune{'s', 'p'}: 0.0020,
	[2]rune{'s', 'u'}: 0.0020,
	[2]rune{'y', 't'}: 0.0020,
	[2]rune{'b', 'o'}: 0.0019,
	[2]rune{'d', 'w'}: 0.0019,
	[2]rune{'y', 's'}: 0.0019,
	[2]rune{'a', 'g'}: 0.0018,
	[2]rune{'c', 'k'}: 0.0018,
	[2]rune{'g', 'i'}: 0.0018,
	[2]rune{'m', 'y'}: 0.0018,
	[2]rune{'o', 'd'}: 0.0018,
	[2]rune{'p', 'r'}: 0.0018,
	[2]rune{'y', 'a'}: 0.0018,
	[2]rune{'b', 'l'}: 0.0017,
	[2]rune{'i', 'f'}: 0.0017,
	[2]rune{'s', 'c'}: 0.0017,
	[2]rune{'t', 'l'}: 0.0017,
	[2]rune{'t', 'u'}: 0.0017,
	[2]rune{'d', 'n'}: 0.0016,
	[2]rune{'f', 'r'}: 0.0016,
	[2]rune{'g', 't'}: 0.0016,
	[2]rune{'n', 'h'}: 0.0016,
	[2]rune{'o', 'a'}: 0.0016,
	[2]rune{'r', 'n'}: 0.0016,
	[2]rune{'t', 'y'}: 0.0016,
	[2]rune{'u', 'p'}: 0.0016,
	[2]rune{'c', 't'}: 0.0015,
	[2]rune{'e', 'g'}: 0.0015,
	[2]rune{'l', 't'}: 0.0015,
	[2]rune{'o', 'p'}: 0.0015,
	[2]rune{'p', 'l'}: 0.0015,
}

var cumulativeFreqOfMostCommonEnglishLetterPairsInEnglish = func(data map[[2]rune]float64) float64 {
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum
}(hundredMostCommonEnglishLetterPairs)

var cumulativeFreqOfMostCommonEnglishLetterPairsInRandomAlphaString = float64(len(hundredMostCommonEnglishLetterPairs)) * 1.0 / (26.0 * 26.0)

func CheckIsNotEnglish(strings []string) bool {
	sampleFreq := calcFreqOfCommonEnglishLetterPairsInSampleText(strings)

	distFromSampleToEnglish := math.Abs(sampleFreq - cumulativeFreqOfMostCommonEnglishLetterPairsInEnglish)
	distFromSampleToRandom := math.Abs(sampleFreq - cumulativeFreqOfMostCommonEnglishLetterPairsInRandomAlphaString)

	result := distFromSampleToEnglish > distFromSampleToRandom
	return result
}

func calcFreqOfCommonEnglishLetterPairsInSampleText(strings []string) float64 {
	numLetterPairsInSample := 0
	numCommonEnglishLetterPairsInSample := 0
	var pair [2]rune

	for _, content := range strings {
		for idx := 0; idx < len(content)-1; idx++ {
			char := rune(content[idx])
			nextChar := rune(content[idx+1])

			if !unicode.IsLetter(char) ||
				!unicode.IsLetter(nextChar) ||
				unicode.IsLower(char) && unicode.IsUpper(nextChar) {
				continue
			}

			pair[0] = unicode.ToLower(char)
			pair[1] = unicode.ToLower(nextChar)
			numLetterPairsInSample += 1

			_, ok := hundredMostCommonEnglishLetterPairs[pair]
			if ok == true {
				numCommonEnglishLetterPairsInSample += 1
			}
		}
	}

	avgFreqOfCommonEnglishPairsInSample := float64(numCommonEnglishLetterPairsInSample) / float64(numLetterPairsInSample)
	return avgFreqOfCommonEnglishPairsInSample
}
