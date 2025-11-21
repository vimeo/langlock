package git

import (
	"bytes"
)

type FileDiff struct {
	Paths   []string
	Content []byte
}

func ParseFileDiff(content []byte) *FileDiff {
	/* NOTE: Return nil if there are no changed lines */
	prePathPrefix := []byte("\n---")
	postPathPrefix := []byte("\n+++")
	newLineBytes := []byte("\n")

	// Parse PRE file path
	idxOfLineWithPrePath := bytes.Index(content, prePathPrefix)
	if idxOfLineWithPrePath == -1 {
		return nil
	}
	relativeIdxOfNextLineEnd := bytes.Index(content[idxOfLineWithPrePath+1:], newLineBytes)
	if relativeIdxOfNextLineEnd == -1 {
		return nil
	}
	idxOfEndOfLineWithPrePath := idxOfLineWithPrePath + 1 + relativeIdxOfNextLineEnd
	prePath := string(content[idxOfLineWithPrePath+5 : idxOfEndOfLineWithPrePath])

	// Parse POST file path
	idxOfLineWithPostPath := idxOfEndOfLineWithPrePath
	var idxOfEndOfLineWithPostPath int
	relativeIdxOfNextLineEnd = bytes.Index(content[idxOfLineWithPostPath+1:], newLineBytes)
	if relativeIdxOfNextLineEnd == -1 {
		idxOfEndOfLineWithPostPath = len(content)
	} else {
		idxOfEndOfLineWithPostPath = idxOfLineWithPostPath + 1 + relativeIdxOfNextLineEnd
	}
	if !bytes.Equal(content[idxOfLineWithPostPath:idxOfLineWithPostPath+4], postPathPrefix) {
		return nil
	}
	postPath := string(content[idxOfLineWithPostPath+5 : idxOfEndOfLineWithPostPath])

	var diffContent []byte
	if idxOfEndOfLineWithPostPath+1 < len(content) {
		diffContent = content[idxOfEndOfLineWithPostPath+1:]
	}

	return &FileDiff{
		[]string{prePath, postPath},
		diffContent,
	}
}
