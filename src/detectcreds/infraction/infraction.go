package infraction

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

type Infraction struct {
	Offender     string `json:"offender"`
	RuleName     string `json:"rule"`
	Path         string `json:"path"`
	Commit       string `json:"commit"`
	OffenderHash string `json:"hash"`
}

type InfractionSummary struct {
	Offender     string               `json:"offender"`
	RuleName     string               `json:"rule"`
	OffenderHash string               `json:"hash"`
	Locations    []InfractionLocation `json:"locations"`
}

type InfractionLocation struct {
	Path   string `json:"path"`
	Commit string `json:"commit"`
}

type InfractionSummarySanitized struct {
	RuleName  string               `json:"rule"`
	Locations []InfractionLocation `json:"locations"`
}

// func (il InfractionLocation) String() string {
// 	return fmt.Sprintf("{\"commit\": \"path\"}")
// }

func PrettyPrint(infractions []*Infraction, allowListPath string) string {
	template := `• Potential secret: "%s"
  Secret type: "%s"
  Occurences
%+v
  Remediation
    ○ If the string is indeed sensitive, please do not commit and push it to the
      upstream repo. Instead, remove the secret from your commits before pushing.
    ○ Otherwise, if the string is not actually sensitive, you may add it to the
      branch's allow list by running the following commands in the root of the repo:

          [ -f %s ] || echo '{}' > %s;` +
		` echo "$(jq '.["allowStrings"]["%s"] = %+v' %s)" > %s;` +
		` read -p "$(echo -e '\\n\\nPLEASE provide a brief justification for adding the detected credential to the allow list: ')" LANGLOCK_REASON;` +
		` echo "$(jq --arg r "$LANGLOCK_REASON" '.["allowStrings"]["%s"].reason = $r' %s)" > %s;` +
		` git add %s;` +
		` git commit -m "Update allow list";` +
		` git push

    ○ Alternatively, you can bypass scanning for specific file extensions or
      directory paths by adding a regex to the allow list with the following command:

          [ -f %s ] || echo '{}' > %s;` +
		` read -p "$(echo -e '\\n\\nPLEASE specify the perl-style regular expression to add to the allow list. Scanning will bypass all files whose paths FULLY match the regex: ')" LANGLOCK_PATH_REGEX;` +
		` read -p "$(echo -e '\\n\\nPLEASE provide a brief justification for adding the path regex to the allow list: ')" LANGLOCK_REASON;` +
		` echo "$(jq --arg r1 "$LANGLOCK_PATH_REGEX" --arg r2 "$LANGLOCK_REASON" '.["allowPaths"] += [{"regex": $r1, "reason": $r2}]' %s)" > %s;` +
		` git add %s;` +
		` git commit -m "Update allow list";` +
		` git push
`

	var summary []InfractionSummary
	summary = Summarize(infractions)
	results := make([]string, 0)

	for _, entry := range summary {
		sanitizedEntry := InfractionSummarySanitized{entry.RuleName, entry.Locations}
		sanitizedEntryJson, _ := json.Marshal(&sanitizedEntry)
		sanitizedEntryJsonStr := string(sanitizedEntryJson)
		locationsBulletedList := make([]string, 0)
		for _, loc := range entry.Locations {
			bulletStr := fmt.Sprintf(`    ○ %s - commit %s`, loc.Path, loc.Commit)
			locationsBulletedList = append(locationsBulletedList, bulletStr)
		}
		locationsBulletedListStr := strings.Join(locationsBulletedList, "\n")
		newSection := fmt.Sprintf(
			template,
			entry.Offender,
			entry.RuleName,
			locationsBulletedListStr,
			allowListPath,
			allowListPath,
			entry.OffenderHash,
			sanitizedEntryJsonStr,
			allowListPath,
			allowListPath,
			entry.OffenderHash,
			allowListPath,
			allowListPath,
			allowListPath,
			allowListPath,
			allowListPath,
			allowListPath,
			allowListPath,
			allowListPath,
		)
		results = append(results, newSection)
	}
	result := strings.Join(results, "\n")
	return result

}

func Summarize(infractions []*Infraction) []InfractionSummary {
	var result []InfractionSummary
	infractionSummaryMap := make(map[string]*InfractionSummary)
	for _, infraction := range infractions {
		location := InfractionLocation{
			infraction.Path,
			infraction.Commit,
		}
		_, ok := infractionSummaryMap[string(infraction.Offender)]
		if ok == true {
			infractionSummaryMap[string(infraction.Offender)].Locations = append(infractionSummaryMap[string(infraction.Offender)].Locations, location)
		} else {
			infractionSummaryMap[string(infraction.Offender)] = &InfractionSummary{
				infraction.Offender,
				infraction.RuleName,
				infraction.OffenderHash,
				[]InfractionLocation{location},
			}
		}
	}
	for _, v := range infractionSummaryMap {
		result = append(result, *v)
	}
	return result
}

func ParseJson(content []byte) ([]*Infraction, error) {
	var infractions = make([]*Infraction, 0)
	err := json.Unmarshal(content, &infractions)
	return infractions, err
}

func NewInfraction(offender string, ruleName string, path string) Infraction {

	return Infraction{
		Offender:     offender,
		RuleName:     ruleName,
		Path:         path,
		OffenderHash: HashOffender(offender),
	}

}

func HashOffender(offender string) string {
	sum := sha256.Sum256([]byte(offender))

	dst := make([]byte, hex.EncodedLen(len(sum)))
	hex.Encode(dst, sum[:])

	return string(dst)
}
