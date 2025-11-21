package rule

import (
	"regexp"

	"detectcreds/infraction"
	"detectcreds/rule/entropy"
)

type Rule struct {
	Name         string
	ContentRegex *regexp.Regexp
	PathRegex    *regexp.Regexp
	Entropies    []entropy.EntropyRange
	//Distributions []distribution.DistributionRule
	ContentChecks  []ContentCheck
	MatchFormatter func(string, []string) string
}

type ContentCheck struct {
	Checker      func([]string) bool
	GroupIndices []int
}

//func (r Rule) Check(content string, paths []string) []infraction.Infraction {
//	var findings []infraction.Infraction
//
//	locs := r.ContentRegex.FindAllStringIndex(content, -1)
//	for _, loc := range locs {
//		match := content[loc[0]:loc[1]]
//		groups := r.ContentRegex.FindStringSubmatch(match)
//		if len(r.Entropies) != 0 && !entropy.TriggerEntropies(groups, r.Entropies) {
//			continue
//		}
//
//		var path string
//		if len(paths) > 0 {
//			path = paths[len(paths)-1]
//		}
//		newInfraction := infraction.NewInfraction(match, r.Name, path)
//		findings = append(findings, newInfraction)
//	}
//	return findings
//}

func (r Rule) checkContentViolation(content []byte) []*infraction.Infraction {
	var findings []*infraction.Infraction
	var locs [][]int
	if r.ContentRegex != nil {
		locs = r.ContentRegex.FindAllIndex(content, -1)
		for _, loc := range locs {
			match := string(content[loc[0]:loc[1]])
			groups := r.ContentRegex.FindStringSubmatch(match)
			// if len(r.Entropies) != 0 && !entropy.TriggerEntropies(groups, r.Entropies) {
			// 	continue
			// }
			anyFailedContentChecks := false
			for _, contentCheck := range r.ContentChecks {
				var stringsForContentCheck []string
				for _, groupIdx := range contentCheck.GroupIndices {
					stringsForContentCheck = append(stringsForContentCheck, groups[groupIdx])
				}
				if !contentCheck.Checker(stringsForContentCheck) {
					anyFailedContentChecks = true
					break
				}
			}
			if anyFailedContentChecks == true {
				continue
			}

			if r.MatchFormatter != nil {
				match = r.MatchFormatter(match, groups)
			}
			newInfraction := infraction.NewInfraction(match, r.Name, "")
			findings = append(findings, &newInfraction)
		}
	}
	return findings

}

func (r Rule) checkPathViolation(path string) bool {
	if r.PathRegex != nil && r.PathRegex.MatchString(path) {
		return true
	}
	return false
}

func (r Rule) Check(content []byte, paths []string) []*infraction.Infraction {
	var findings []*infraction.Infraction
	var path string
	if len(paths) > 0 {
		path = paths[len(paths)-1]
	}

	if r.ContentRegex != nil && r.PathRegex != nil {
		if r.checkPathViolation(path) {
			findings = r.checkContentViolation(content)
		}
	} else if r.ContentRegex != nil {
		findings = r.checkContentViolation(content)
	} else if r.PathRegex != nil {
		if r.checkPathViolation(path) {
			newInfraction := infraction.NewInfraction(path, r.Name, path)
			findings = []*infraction.Infraction{&newInfraction}
		}
	}

	for idx, _ := range findings {
		findings[idx].Path = path
	}

	return findings
}
