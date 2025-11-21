package patch

import (
	"regexp"
)

type FilePatch struct {
	Paths   []string
	Content string
}

type Patch []FilePatch

func Parse(s string) Patch {
	var results []FilePatch
	//r, _ := regexp.Compile(`(?:^|\n)diff --git [^\n]+\n(?:[^\n]+\n){1,3}--- ([^\n]+)\n\+\+\+ ([^\n]+)\n(.*)(?:\ndiff|$)`)
	r, _ := regexp.Compile(`(?m:^)diff[^\n]*\n(?:(?:old|new|deleted|copy|rename|similarity|dissimilarity|index)[^\n]*\n)*(?:--- ([^\n]*)\n\+\+\+ ([^\n$]*)(?m:$))?`)

	//r, _ := regexp.Compile(`(?:^|\n)diff --git [^\n]+\n(?:[^\n]+\n){1,3}--- ([^\n]+)\n\+\+\+ ([^\n]+)\n`)
	//r, _ := regexp.Regexp(`(?:^|\n)diff --git [^\n]+\n[^\n]+\n--- ([^\n]+)\n\+\+\+ ([^\n]+)\n`)
	locs := r.FindAllStringSubmatchIndex(s, -1)

	for locIdx, loc := range locs {
		if loc[2] == -1 {
			continue
		}
		prePath := s[loc[2]:loc[3]]
		postPath := s[loc[4]:loc[5]]
		paths := make([]string, 0)
		paths = append(paths, prePath)
		if postPath != prePath {
			paths = append(paths, postPath)
		}

		var content string
		var contentIncludiveStartIdx int
		var contentExclusiveEndIdx int
		contentIncludiveStartIdx = loc[5] + 1 // Add one to avoid preceding newline
		if locIdx < len(locs)-1 {
			contentExclusiveEndIdx = locs[locIdx+1][0] - 1 // Remove one to avoid trailing newline
		} else {
			contentExclusiveEndIdx = len(s)
		}
		if contentIncludiveStartIdx >= contentExclusiveEndIdx {
			content = ""
		} else {
			content = s[contentIncludiveStartIdx:contentExclusiveEndIdx]
		}

		newEntry := FilePatch{
			Paths:   paths,
			Content: content,
		}

		results = append(results, newEntry)
	}

	return results

}

// func OldNewPatch(s string) Patch {
// 	var results []FilePatch
// 	r, _ := regexp.Compile(`(?:^|\n)diff --git [^\n]+\n(?:[^\n]+\n){1,3}--- ([^\n]+)\n\+\+\+ ([^\n]+)\n(.*)(?:\ndiff|$)`)
// 	//r, _ := regexp.Compile(`(?:^|\n)diff --git [^\n]+\n(?:[^\n]+\n){1,3}--- ([^\n]+)\n\+\+\+ ([^\n]+)\n`)
// 	//r, _ := regexp.Regexp(`(?:^|\n)diff --git [^\n]+\n[^\n]+\n--- ([^\n]+)\n\+\+\+ ([^\n]+)\n`)
// 	locs := r.FindAllStringSubmatchIndex(s, -1)
//
// 	for locIdx, loc := range locs {
// 		prePath := s[loc[2]:loc[3]]
// 		postPath := s[loc[4]:loc[5]]
// 		paths := make([]string, 0)
// 		paths = append(paths, prePath)
// 		if postPath != prePath {
// 			paths = append(paths, postPath)
// 		}
//
// 		var content string
// 		var contentExclusiveEndIdx int
// 		var contentIncludiveStartIdx int
// 		if locIdx < len(locs)-1 {
// 			contentExclusiveEndIdx = locs[locIdx+1][0] // Next patch includes preceding whitespace
// 		} else {
// 			contentExclusiveEndIdx = len(s)
// 		}
// 		contentIncludiveStartIdx = loc[5] + 1 // Add one to avoid preceding newline
// 		if contentIncludiveStartIdx >= contentExclusiveEndIdx {
// 			content = ""
// 		} else {
// 			content = s[contentIncludiveStartIdx:contentExclusiveEndIdx]
// 		}
//
// 		newEntry := FilePatch{
// 			Paths:   paths,
// 			Content: content,
// 		}
//
// 		results = append(results, newEntry)
// 	}
//
// 	return results
//
// }
