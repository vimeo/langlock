package git

type Commit struct {
	Hash string
	Diff *Diff
}

func ParseCommit(content []byte) *Commit {
	// Get commit ido
	commitHash := string(content[7:47])
	diff := ParseDiff(content)
	return &Commit{
		commitHash,
		diff,
	}
}
