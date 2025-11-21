package scan

import (
	"io/ioutil"
	"testing"

	"detectcreds/git"
)

type testCaseScanGitLogForSecrets struct {
	pathToDataFile         string
	expectedNumInfractions int
}

var testCases = []testCaseScanGitLogForSecrets{
	{"../sampledata/gitlog_0creds.txt", 0},
	{"../sampledata/gitlog_8creds.txt", 9}, // Actually has 8 creds, but the AWS key gets counted twice. Each cred labeled with "TEST" comment
}

func TestScanGitLogForSecrets(t *testing.T) {
	for _, testCase := range testCases {
		sampleData, err := ioutil.ReadFile(testCase.pathToDataFile)
		if err != nil {
			t.Errorf("Could not load data file for test case: %s", testCase.pathToDataFile)
		}
		var log *git.Log
		log = git.ParseLog(sampleData)
		// Test different number of threads
		threadCounts := []int{1, 2, 3, 4}
		for _, threadCount := range threadCounts {
			actualNumInfractionsFound := len(ScanGitLogForSecrets(log, threadCount))
			if actualNumInfractionsFound != testCase.expectedNumInfractions {
				t.Errorf("Incorrect number of infractions found for ScanGitLogForSecrets for data file %s with thread count %d: \n_HAVE_: %d\n_WANT_: %d",
					testCase.pathToDataFile, threadCount, actualNumInfractionsFound, testCase.expectedNumInfractions)
			}
		}
	}
}
