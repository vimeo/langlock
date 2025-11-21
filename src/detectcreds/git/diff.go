package git

type Diff struct {
	FileDiffs []*FileDiff
}

func ParseDiff(content []byte) *Diff {
	diffPrefix := []byte("diff")
	var fileDiffs = make([]*FileDiff, 0)
	unparsedFileDiffs := IndexSectionsByLinePrefix(content, diffPrefix)
	for _, unparsedFileDiff := range unparsedFileDiffs {
		parsedFileDiff := ParseFileDiff(unparsedFileDiff)
		if parsedFileDiff != nil {
			fileDiffs = append(fileDiffs, parsedFileDiff)
		}
	}
	return &Diff{
		fileDiffs,
	}
}
