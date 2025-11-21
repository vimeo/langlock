package punctuation

import (
	"regexp"
)

var compiledSeparatorRegex = func() *regexp.Regexp {
	r, err := regexp.Compile(`[[:^alnum:]]+`)
	if err != nil {
		panic(err)
	}
	return r
}()

func CheckHasSeparator(strings []string) bool {
	for _, str := range strings {
		hasNonAlphaNum := compiledSeparatorRegex.MatchString(str)
		if hasNonAlphaNum == true {
			return true
		}
	}
	return false
}
