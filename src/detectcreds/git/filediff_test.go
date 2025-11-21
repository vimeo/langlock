package git

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseFileDiff(t *testing.T) {

	for _, testCase := range testCases {
		actualOutput := ParseFileDiff(testCase.Input)
		if !cmp.Equal(actualOutput, testCase.ExpectedOutput) {
			t.Errorf("Incorrect parsing of FileDiff byte slice.\n_HAVE_:\n--%+v--\n_WANT_:\n--%+v--", actualOutput, testCase.ExpectedOutput)
		}
	}
}

type TestCaseNewPatch struct {
	Input          []byte
	ExpectedOutput *FileDiff
}

var testCases = []TestCaseNewPatch{
	{
		Input: []byte(`diff --git a/sample_binary_file.bin b/sample_binary_file.bin
new file mode 100755
index 0000000..a57718c
Binary files /dev/null and b/sample_binary_file.bin differ
diff --git a/small_safe.txt b/small_safe.txt
index 7e8a165..16f18f3 100644
--- small_safe.txt
+++ small_safe.txt
@@ -1,2 +1,3 @@
 a
 a
+a`),
		ExpectedOutput: &FileDiff{
			Paths: []string{"small_safe.txt", "small_safe.txt"},
			Content: []byte(`@@ -1,2 +1,3 @@
 a
 a
+a`),
		},
	},
	{
		Input: []byte(`diff --git a/b.txt b/b.txt
old mode 100644
new mode 100755`),
		ExpectedOutput: nil,
	},
	{
		Input: []byte(`diff --git a/a.txt b/a.txt
index 7898192..7e8a165 100644
--- a/a.txt
+++ b/b.txt`),
		ExpectedOutput: &FileDiff{
			Paths:   []string{"a/a.txt", "b/b.txt"},
			Content: nil,
		},
	},
}
