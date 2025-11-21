package git

import "bytes"

func IndexSectionsByLinePrefix(content []byte, linePrefix []byte) [][]byte {
	sectionsFound := make([][]byte, 0)
	startIndices := make([]int, 0)
	remainingContent := content
	prevAbsoluteIdx := -1
	for {
		relativeIdx := bytes.Index(remainingContent, linePrefix)
		if relativeIdx == -1 {
			break
		}
		absoluteIdx := prevAbsoluteIdx + 1 + relativeIdx
		if absoluteIdx == 0 || content[absoluteIdx-1] == '\n' {
			startIndices = append(startIndices, absoluteIdx)
		}
		prevAbsoluteIdx = absoluteIdx
		remainingContent = content[absoluteIdx+1:]
	}
	var sectionInclusiveStart int
	var sectionExclusiveEnd int
	for i := range startIndices {
		sectionInclusiveStart = startIndices[i]
		if i < len(startIndices)-1 {
			sectionExclusiveEnd = startIndices[i+1] - 1 // Subtract one to remove trailing newline
		} else {
			sectionExclusiveEnd = len(content)
		}
		sectionsFound = append(sectionsFound, content[sectionInclusiveStart:sectionExclusiveEnd])
	}

	return sectionsFound
}
