package allowlist

import (
	"encoding/json"
	"fmt"
	"regexp"
)

type AllowList struct {
	AllowedStrings AllowedStrings `json:"allowStrings"`
	AllowedPaths   AllowedPaths   `json:"allowPaths"`
}

type AllowedStrings map[string]AllowedStringEntry
type AllowedPaths []AllowedPathEntry

type AllowedPathEntry struct {
	Regex           *regexp.Regexp
	RegexDefinition string `json:"regex"`
	Reason          string `json:"reason"`
}

type AllowedStringEntry struct {
	RuleName  string                   `json:"rule"`
	Reason    string                   `json:"reason"`
	Locations []AllowListEntryLocation `json:"locations"`
}
type AllowListEntryLocation struct {
	Path   string `json:"path"`
	Commit string `json:"commit"`
}

type regexCompilationError struct {
	attemptedRegex  string
	underlyingError error
}

func (e *regexCompilationError) Error() string {
	return fmt.Sprintf("Failed to compile regex %s for allowed file path. Error: %s", e.attemptedRegex, e.underlyingError.Error())
}

func ParseJson(content []byte) (AllowList, error) {
	var allowListObj AllowList
	err := json.Unmarshal(content, &allowListObj)
	if err != nil {
		return allowListObj, err
	}
	for idx, allowedPath := range allowListObj.AllowedPaths {
		var regexDefinition = allowedPath.RegexDefinition
		if len(regexDefinition) == 0 {
			regexDefinition = "^$"
		}
		if regexDefinition[0] != '^' {
			regexDefinition = "^" + regexDefinition
		}
		if regexDefinition[len(regexDefinition)-1] != '$' {
			regexDefinition = regexDefinition + "$"
		}
		r, err := regexp.Compile(regexDefinition)
		if err != nil {
			return allowListObj, &regexCompilationError{allowedPath.RegexDefinition, err}
		}
		// fmt.Printf("\nCompiled regex: %+v\n", r)
		allowListObj.AllowedPaths[idx].Regex = r

	}
	return allowListObj, err
}
