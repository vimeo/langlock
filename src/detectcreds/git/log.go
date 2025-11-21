package git

type Log []*Commit

func ParseLog(content []byte) *Log {
	var result = make([]*Commit, 0)
	unparsedCommits := IndexSectionsByLinePrefix(content, []byte("commit"))
	for _, unparsedCommit := range unparsedCommits {
		result = append(result, ParseCommit(unparsedCommit))
	}
	log := Log(result)
	return &log
}
