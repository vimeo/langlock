package scan

import (
	"detectcreds/allowlist"
	"detectcreds/git"
	"detectcreds/infraction"
	"detectcreds/rule"
)

func ScanGitLogForSecrets(log *git.Log, numThreads int) []*infraction.Infraction {
	scan := func(rules []rule.Rule) []*infraction.Infraction {
		var infractions = make([]*infraction.Infraction, 0)
		for _, r := range rules {
			for _, commit := range *log {
				for _, fileDiff := range commit.Diff.FileDiffs {
					newInfractions := r.Check(fileDiff.Content, fileDiff.Paths)
					for _, infra := range newInfractions {
						infra.Commit = commit.Hash
					}
					infractions = append(infractions, newInfractions...)
				}
			}
		}
		return infractions
	}
	return parallelize(numThreads, rule.DefaultRules, scan)
}

func ScanPlaintextForSecrets(content []byte, numThreads int) []*infraction.Infraction {

	scan := func(rules []rule.Rule) []*infraction.Infraction {
		var infras []*infraction.Infraction
		for _, r := range rules {
			infras = append(infras, r.Check(content, []string{})...)
		}
		return infras

	}

	return parallelize(numThreads, rule.DefaultRules, scan)
}

func parallelize(wantedNumThreads int, totalRuleList []rule.Rule, funct func([]rule.Rule) []*infraction.Infraction) []*infraction.Infraction {

	if wantedNumThreads == 0 {
		return funct(totalRuleList)
	}

	c := make(chan []*infraction.Infraction)
	goRountineFunc := func(rules []rule.Rule) {
		c <- funct(rules)
	}

	var results []*infraction.Infraction

	numJobsPerThread := len(totalRuleList) / wantedNumThreads
	numAdditionalJobs := len(totalRuleList) % wantedNumThreads

	jobIdx := 0
	numActiveThreads := 0
	for idx := 0; idx < wantedNumThreads; idx++ {
		numJobsForCurrentThread := numJobsPerThread
		if numAdditionalJobs > 0 {
			numJobsForCurrentThread++
			numAdditionalJobs--
		}
		if numJobsForCurrentThread > 0 {
			go goRountineFunc(totalRuleList[jobIdx : jobIdx+numJobsForCurrentThread])
			numActiveThreads++
		}
		jobIdx = jobIdx + numJobsForCurrentThread
	}

	for numActiveThreads > 0 {
		resultsForThread := <-c
		results = append(results, resultsForThread...)
		numActiveThreads--
	}

	return results
}

func FilterByAllowedStrings(infractions []*infraction.Infraction, allowedStrings allowlist.AllowedStrings) []*infraction.Infraction {
	var unallowedInfractions []*infraction.Infraction
	for _, infraction := range infractions {
		_, ok := allowedStrings[string(infraction.OffenderHash)]
		if ok == false {
			unallowedInfractions = append(unallowedInfractions, infraction)
		}

	}
	return unallowedInfractions
}

func FilterByAllowedPaths(infractions []*infraction.Infraction, allowedPaths allowlist.AllowedPaths) []*infraction.Infraction {
	var unallowedInfractions []*infraction.Infraction
	for _, infraction := range infractions {
		isPathAllowed := false
		for _, allowedPath := range allowedPaths {
			// fmt.Printf("\ninfraction %+v\n", infraction)
			// fmt.Printf("infraction.Path %+v\n", infraction.Path)
			// fmt.Printf("allowedPath %+v\n\n", allowedPath)
			isPathAllowed = allowedPath.Regex.MatchString(infraction.Path)
			if isPathAllowed == true {
				break
			}
		}
		if isPathAllowed == false {
			unallowedInfractions = append(unallowedInfractions, infraction)
		}
	}
	return unallowedInfractions
}
