package patch

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParse(t *testing.T) {

	for _, testCase := range testCases {
		actualOutput := Parse(testCase.Input)
		if !cmp.Equal(actualOutput, testCase.ExpectedOutput) {
			t.Errorf("Incorrect patch created.\n_HAVE_:\n--%+v--\n_WANT_:\n--%+v--", actualOutput, testCase.ExpectedOutput)
		}
	}
}

type TestCaseNewPatch struct {
	Input          string
	ExpectedOutput Patch
}

var testCases = []TestCaseNewPatch{
	{
		Input: `diff --git a/sample_binary_file.bin b/sample_binary_file.bin
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
+a
diff --git a/small_secrets.txt b/small_secrets.txt
index 554f1ee..680fc01 100644
--- small_secrets.txt
+++ small_secrets.txt
@@ -41,3 +41,4 @@ aws_access_key_id=AKIAW5PVYFYB3ZBF9MPP
 aws_access_key_id=AKIAW5PVYFYB3ZBF9MPP
 aws_access_key_id=AKIAW5PVYFYB3ZBF9MPP
 aws_access_key_id=AKIAW5PVYFYB3ZBF9MPP
+aws_access_key_id=AKIAW5PVYFYB3ZBF9MPP`,
		ExpectedOutput: Patch{
			FilePatch{
				Paths: []string{"small_safe.txt"},
				Content: `@@ -1,2 +1,3 @@
 a
 a
+a`,
			},
			FilePatch{
				Paths: []string{"small_secrets.txt"},
				Content: `@@ -41,3 +41,4 @@ aws_access_key_id=AKIAW5PVYFYB3ZBF9MPP
 aws_access_key_id=AKIAW5PVYFYB3ZBF9MPP
 aws_access_key_id=AKIAW5PVYFYB3ZBF9MPP
 aws_access_key_id=AKIAW5PVYFYB3ZBF9MPP
+aws_access_key_id=AKIAW5PVYFYB3ZBF9MPP`,
			},
		},
	},
	{
		Input: `diff --git a/a.txt b/a.txt
index 7898192..7e8a165 100644
--- a/a.txt
+++ b/a.txt
@@ -1 +1,2 @@
 a
+a
diff --git a/b.txt b/b.txt
old mode 100644
new mode 100755
diff --git a/c.txt b/c.txt
index 6178079..73603e1 100644
--- a/c.txt
+++ b/c.txt
@@ -1 +1,2 @@
 b
+b`,
		ExpectedOutput: Patch{
			FilePatch{
				Paths: []string{"a/a.txt", "b/a.txt"},
				Content: `@@ -1 +1,2 @@
 a
+a`,
			},
			FilePatch{
				Paths: []string{"a/c.txt", "b/c.txt"},
				Content: `@@ -1 +1,2 @@
 b
+b`,
			},
		},
	},
}
