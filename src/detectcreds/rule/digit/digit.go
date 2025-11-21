package digit

import (
	"regexp"
)

var compiledDigitRegex = func() *regexp.Regexp {
	r, err := regexp.Compile(`[0-9]+`)
	if err != nil {
		panic(err)
	}
	return r
}()

func CheckHasDigit(strings []string) bool {
	for _, str := range strings {
		hasDigit := compiledDigitRegex.MatchString(str)
		if hasDigit == true {
			return true
		}
	}
	return false
}
