package main

import (
	"testing"
)

type testCaseScanForSecrets struct {
	argInputFile        string
	argInputType        string
	argOutputType       string
	argAllowListName    string
	argAllowListFile    string
	argNumThreads       int
	expectedExitStatus  int
	expectedExitMessage string
}

var expectedResult1 = `[{"offender":"aws_access_key_id=AKIAW5JVYFDT5ZB8FAKE","rule":"Generic credential - base 64/62/36 or hex","path":"src/Psalm/Type/Atomic/TNonEmptyLowercaseString.php","commit":"1b81ce9921dd5d2b6a58870ccc1e9ca7a70a035c","hash":"994d3a39fe234af3e931eb26d39568c926cba78c7670dbaa9a741cc0d5aaf8e3"},{"offender":"0xCE572DC620C806C230EB68FBD371D1","rule":"Generic credential - base 64/62/36 or hex","path":"src/Psalm/Internal/Analyzer/ClassLikeAnalyzer.php","commit":"3687d34a5a499669d1300d17329dd2a42258be18","hash":"83e833c1b979518ee88e6d49a4ac60fd24df3b81ae0812ad14043f7b21d79ffc"},{"offender":"0xce572dc620c806c230eb68fbd371d1","rule":"Generic credential - base 64/62/36 or hex","path":"src/Psalm/Internal/Analyzer/ClassLikeAnalyzer.php","commit":"3687d34a5a499669d1300d17329dd2a42258be18","hash":"2581719ce2aebc4aa2cec487feaa3ce5edc6e0852c34ac9bd00c7dfbbb057484"},{"offender":"AmiMmpOpxeaOkG8t","rule":"Generic credential - base 64/62/36 or hex","path":"src/Psalm/Internal/Analyzer/Statements/Expression/Call/ArgumentsAnalyzer.php","commit":"7ef3d4711fc02cf8d48eb5a3894d8b79347fd86a","hash":"f7a5c64ed96e42bb0d264895ae1bbd41bbe81161c1744977d078b9aee8caa89e"},{"offender":"postgresql://other:fakeDbPassword@localhost/otherdb?connect_timeout=10\u0026application_name=myapp","rule":"DB/FTP/SFTP/HTTP/HTTPs Connection String","path":"docs/manipulating_code/fixing.md","commit":"38f74815d6c5b7ecfe7055ef36151eb808d6c712","hash":"c6f2da76bc8ba67c34cf2e91418ab1971d98ec9e9e88228a49121be436676d32"},{"offender":"src/Psalm/TEST/CASE/FAKE/CRED/id_rsa","rule":"File private ssh key id_rsa","path":"src/Psalm/TEST/CASE/FAKE/CRED/id_rsa","commit":"4a5f74c091caf38c36d76e1456b4038e67aa10d6","hash":"d9c3c5b28cbe01f86a3e93b4e743a6e7074f883bde070f37fed1170311f8c3aa"},{"offender":"src/TESTCASEFAKECRED/id_rsa","rule":"File private ssh key id_rsa","path":"src/TESTCASEFAKECRED/id_rsa","commit":"4a5f74c091caf38c36d76e1456b4038e67aa10d6","hash":"8b0a73a73b68b46cd1d23284183db20ab9970f7e0825185d5c8dfcc682e27c76"},{"offender":"AKIAW5JVYFDT5ZB8FAKE","rule":"AWS Manager ID","path":"src/Psalm/Type/Atomic/TNonEmptyLowercaseString.php","commit":"1b81ce9921dd5d2b6a58870ccc1e9ca7a70a035c","hash":"e104e5d343812376bbe4f1a9847c4781bea6f781f74d9965b5bb6c6bf3f044e8"},{"offender":"AIzaSyB9V5n21zmPBaIuO1SHoxZHnipZb76Azz0","rule":"Google or YouTube or Gmail or GDrive API key","path":"src/Psalm/Internal/Type/TypeCombination.php","commit":"74eea18563e26652bbbfe0deb7b12a5fc1abc9e3","hash":"83c6afcf97d7603a821ad92a72f717b8743dbb9a4bb1fb63b4961c237fce8e77"}]`

var expectedResult2 = `[{"offender":"0xCE572DC620C806C230EB68FBD371D1","rule":"Generic credential - base 64/62/36 or hex","path":"src/Psalm/Internal/Analyzer/ClassLikeAnalyzer.php","commit":"3687d34a5a499669d1300d17329dd2a42258be18","hash":"83e833c1b979518ee88e6d49a4ac60fd24df3b81ae0812ad14043f7b21d79ffc"},{"offender":"0xce572dc620c806c230eb68fbd371d1","rule":"Generic credential - base 64/62/36 or hex","path":"src/Psalm/Internal/Analyzer/ClassLikeAnalyzer.php","commit":"3687d34a5a499669d1300d17329dd2a42258be18","hash":"2581719ce2aebc4aa2cec487feaa3ce5edc6e0852c34ac9bd00c7dfbbb057484"},{"offender":"AmiMmpOpxeaOkG8t","rule":"Generic credential - base 64/62/36 or hex","path":"src/Psalm/Internal/Analyzer/Statements/Expression/Call/ArgumentsAnalyzer.php","commit":"7ef3d4711fc02cf8d48eb5a3894d8b79347fd86a","hash":"f7a5c64ed96e42bb0d264895ae1bbd41bbe81161c1744977d078b9aee8caa89e"},{"offender":"postgresql://other:fakeDbPassword@localhost/otherdb?connect_timeout=10\u0026application_name=myapp","rule":"DB/FTP/SFTP/HTTP/HTTPs Connection String","path":"docs/manipulating_code/fixing.md","commit":"38f74815d6c5b7ecfe7055ef36151eb808d6c712","hash":"c6f2da76bc8ba67c34cf2e91418ab1971d98ec9e9e88228a49121be436676d32"},{"offender":"src/Psalm/TEST/CASE/FAKE/CRED/id_rsa","rule":"File private ssh key id_rsa","path":"src/Psalm/TEST/CASE/FAKE/CRED/id_rsa","commit":"4a5f74c091caf38c36d76e1456b4038e67aa10d6","hash":"d9c3c5b28cbe01f86a3e93b4e743a6e7074f883bde070f37fed1170311f8c3aa"},{"offender":"src/TESTCASEFAKECRED/id_rsa","rule":"File private ssh key id_rsa","path":"src/TESTCASEFAKECRED/id_rsa","commit":"4a5f74c091caf38c36d76e1456b4038e67aa10d6","hash":"8b0a73a73b68b46cd1d23284183db20ab9970f7e0825185d5c8dfcc682e27c76"},{"offender":"AIzaSyB9V5n21zmPBaIuO1SHoxZHnipZb76Azz0","rule":"Google or YouTube or Gmail or GDrive API key","path":"src/Psalm/Internal/Type/TypeCombination.php","commit":"74eea18563e26652bbbfe0deb7b12a5fc1abc9e3","hash":"83c6afcf97d7603a821ad92a72f717b8743dbb9a4bb1fb63b4961c237fce8e77"}]`

var testCases = []testCaseScanForSecrets{
	{"sampledata/gitlog_0creds.txt", INPUT_TYPE_LOG, OUTPUT_TYPE_LIST, "allow-list.json", "", 1, EXIT_SUCCESS, ""},
	{"sampledata/gitlog_8creds.txt", INPUT_TYPE_LOG, OUTPUT_TYPE_LIST, "allow-list.json", "", 1, EXIT_SUCCESS, expectedResult1},
	{"sampledata/gitlog_8creds.txt", INPUT_TYPE_LOG, OUTPUT_TYPE_LIST, "allow-list.json", "sampledata/allowlist.json", 1, EXIT_SUCCESS, expectedResult2},
}

func TestScanForSecrets(t *testing.T) {
	for _, testCase := range testCases {
		actualExitStatus, actualExitMessage := ScanForSecrets(
			testCase.argInputFile,
			testCase.argInputType,
			testCase.argOutputType,
			testCase.argAllowListName,
			testCase.argAllowListFile,
			testCase.argNumThreads,
		)
		if actualExitMessage != testCase.expectedExitMessage {
			t.Errorf(
				"Unexpected exit message.\n\nExpected:\n%s\n\nActual:\n%s\n",
				testCase.expectedExitMessage,
				actualExitMessage,
			)
		}
		if actualExitStatus != testCase.expectedExitStatus {
			t.Errorf(
				"Unexpected exit status. Expected: %d Actual: %d\n",
				testCase.expectedExitStatus,
				actualExitStatus,
			)
		}
	}
}
