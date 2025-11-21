package hex

import (
	"regexp"
)

var compiledHexRegex = func() *regexp.Regexp {
	r, err := regexp.Compile(`^([0-9a-fx~_\/.+=-]{16,64}|[0-9A-Fx~_\/.+=-]{16,64})$`)
	if err != nil {
		panic(err)
	}
	return r
}()

func CheckIsNotHex(strings []string) bool {
	foundNonHexLetter := false
	for _, str := range strings {
		isHex := compiledHexRegex.MatchString(str)
		if isHex == false {
			foundNonHexLetter = true
		}
	}
	return foundNonHexLetter
}
