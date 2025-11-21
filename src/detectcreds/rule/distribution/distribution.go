package distribution

type runeGroup struct {
	runes              map[rune]struct{}
	expectedPercentage float64
	tolerancePercent   float64
}

type RuneDistributionRule struct {
	regexGroupIdxs     []int
	expectedRuneGroups []runeGroup
}

func stringToRuneSet(content string) map[rune]struct{} {
	result := make(map[rune]struct{})
	for _, char := range content {
		result[char] = struct{}{}
	}
	return result
}

func NewEvenRuneDistributionRule(regexGroupIdxs []int, runeGroupStrings []string, tolerances []float64) RuneDistributionRule {
	var runeGroups = make([]runeGroup, len(runeGroupStrings))
	totalNumDifferentRunes := 0
	for _, runeGroupString := range runeGroupStrings {
		totalNumDifferentRunes += len(runeGroupString)
	}
	for idx, runeGroupString := range runeGroupStrings {
		runeGroups[idx] = runeGroup{
			runes:              stringToRuneSet(runeGroupString),
			expectedPercentage: 100.0 * float64(len(runeGroupString)) / float64(totalNumDifferentRunes),
			tolerancePercent:   tolerances[idx],
		}
	}
	return RuneDistributionRule{
		regexGroupIdxs:     regexGroupIdxs,
		expectedRuneGroups: runeGroups,
	}
}

func NewDistributionCheck(rule RuneDistributionRule) func([]string) bool {
	return func(samples []string) bool {
		return checkDistribution(samples[1:], rule)
	}

}

func checkDistribution(samples []string, rule RuneDistributionRule) bool {
	countPerRuneGroup := make([]int, len(rule.expectedRuneGroups))
	totalCount := 0
	for regexGroupIdx := range rule.regexGroupIdxs {
		sample := samples[regexGroupIdx]
		for _, char := range sample {
			for idx, expectedRuneGroup := range rule.expectedRuneGroups {
				_, ok := expectedRuneGroup.runes[char]
				if ok == true {
					countPerRuneGroup[idx] += 1
					totalCount += 1
				}
			}
		}
	}
	for runeGroupIdx, runeGroup := range rule.expectedRuneGroups {
		actualCount := countPerRuneGroup[runeGroupIdx]
		actualPercentage := 100.0 * float64(actualCount) / float64(totalCount)
		if actualPercentage < runeGroup.expectedPercentage*(100.0-runeGroup.tolerancePercent)/100.0 ||
			actualPercentage > runeGroup.expectedPercentage*(100.0+runeGroup.tolerancePercent)/100.0 {
			return false
		}
	}
	return true
}
