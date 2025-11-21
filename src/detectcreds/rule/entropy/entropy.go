package entropy

import "math"

type EntropyRange struct {
	LowerInclusiveBound float64
	UpperInclusiveBound float64
}

func calcInformationEntropy(str string) float64 {

	var charCountMap map[rune]int

	for _, char := range str {
		count, ok := charCountMap[char]
		if ok == false {
			charCountMap[char] = 1
		} else {
			charCountMap[char] = count + 1
		}
	}

	var entropy float64
	strLen := float64(len(str))
	entropy = 0.0

	for _, count := range charCountMap {
		probability := float64(count) / strLen
		entropy += probability * math.Log2(probability)
	}
	return entropy
}

func NewEntropyCheck(entropies []EntropyRange) func([]string) bool {
	return func(strings []string) bool {
		return checkEntropies(strings[1:], entropies)
	}
}

func checkEntropies(strings []string, entropies []EntropyRange) bool {
	if len(entropies) > len(strings) {
		return false
	}
	for idx, entropy := range entropies {
		strEntropy := calcInformationEntropy(strings[idx])
		lowerBoundSatisfied := entropy.LowerInclusiveBound < 0 || (entropy.LowerInclusiveBound >= 0 && strEntropy >= entropy.LowerInclusiveBound)
		upperBoundSatisfied := entropy.UpperInclusiveBound < 0 || (entropy.UpperInclusiveBound >= 0 && strEntropy <= entropy.UpperInclusiveBound)

		if !lowerBoundSatisfied || !upperBoundSatisfied {
			return false
		}
	}
	return true
}
