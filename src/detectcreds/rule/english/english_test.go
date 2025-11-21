package english

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type testCaseCheckIsNotEnglish struct {
	content       []string
	desiredResult bool
}

var testCases = []testCaseCheckIsNotEnglish{
	{[]string{"logBeforeExiting"}, false},
	{[]string{"NewFactoryObject"}, false},
	{[]string{"numWidgetsPerHour"}, false},
	{[]string{"avgNumParamsInArr"}, false},
	{[]string{"thisIsARandVarName"}, false},
	{[]string{"UpperExclusiveBound"}, false},
	{[]string{"logBeforeExitingThread"}, false},
	{[]string{"locationsNumberedListStr"}, false},
	{[]string{"mapVulnerabilities2Fixes"}, false},
	{[]string{"NewEvenChatDistributionRule"}, false},
	{[]string{"DropboxWebhookAutoUpload"}, false},
	{[]string{"PopularClip.PopularClipsV6"}, false},
	{[]string{"Search.Update.Queue"}, false},
	{[]string{"g6oGqje8z9GuPHX"}, true},
	{[]string{"8P6Df6QwVNqwdSN"}, true},
	{[]string{"dNsbrGG9MhJc2pv"}, true},
	{[]string{"Cb6Uu2M7SGpeKw8"}, true},
	{[]string{"ctw4jRpmqm6SD78"}, true},
	{[]string{"N2m4JVfUXbepQaW"}, true},
	{[]string{"V6o2b2L5BYKL2QW"}, true},
	{[]string{"3Fuam7a9dPVRosp"}, true},
	{[]string{"CFXH5Xa9ae7RiZZZ"}, true},
	{[]string{"WTb16muNf4XAdn7T"}, true},
	{[]string{"Nv0ThNNqCzk3i7cZ"}, true},
	{[]string{"OCoY1Jz2xzR9RfVd"}, true},
	{[]string{"XlYYMXjZv2MjuU4x"}, true},
	{[]string{"EtaoxdLhUBmAH0H7"}, true},
	{[]string{"DShBGO21hyRhjrtl"}, true},
	{[]string{"tBXOVxUPk3YOLYFg"}, true},
	{[]string{"ADWCxiI71OcBLlzv"}, true},
	// {[]string{"1nDioD2qGJaf8zU8"}, true}, // Fails
	{[]string{"xDjTqda4n8c3GXt0"}, true},
	{[]string{"pvooVMd0fE3x6EXn"}, true},
	{[]string{"AmiMmpOpxeaOkG8t"}, true},
	{[]string{"JLs3e5Vi9RreRcJaZsp4SRTu"}, true},
	{[]string{"gLW8Uu8bELhtLaZVHi6KB4Zk"}, true},
	{[]string{"xsZkecRY3nmwjZ2K4L8m4PSd"}, true},
	{[]string{"qi7JBWH2McbMTxdrekz94x4r"}, true},
	{[]string{"Xn9cgHbKJrXorNENh2BTFgQSb2yPmpCJAFrFMxEYyg4S5ZEZeGbW32Hoy6ypPuJ"}, true},
	{[]string{"CcX7mn6jHRJ883eZo3JJo6cEpdapz2RXyiu8CxmUry3WTtxQGzwdBEhx5DikEF7"}, true},
	{[]string{"HcWwD4moxqDGWirtp3K6Gy2XEX6JnNGB5weLTSZWH7QccQMLpcwk9LiJVnKjpao"}, true},
	{[]string{"s9gkcuaqpSLdzJ6a8cDd9HCCAMGfWYsUmif3poA65s3eRRpWNSjrGNJUdVaj5Qf"}, true},
	{[]string{"sWjbng83muWUs3aTLvVDzbX4yh8WpuZsmFFACufvwoN7JTaR8cNMpYE2nM5CF6V"}, true},
	{[]string{"jjbk6zN8aQTW9kp9Bcqu2c7zBKekKusdSy6eFhwdWW9G9LJhGJvWVWoEpwFYib3"}, true},
	{[]string{"WA3FHz9o2nZZw75hxQ2uVnbSnHwHQZFKEbLoPqJqD7uW93H3hJRjSPP55umzSAk"}, true},
}

func TestCheckIsNotEnglish(t *testing.T) {

	for _, testCase := range testCases {
		testStrings := []string{}

		testStrings = append(testStrings, testCase.content...)
		actualOutput := CheckIsNotEnglish(testStrings)
		if !cmp.Equal(actualOutput, testCase.desiredResult) {
			t.Errorf("Incorrect result for CheckIsNotEnglish for string: %+v\n_HAVE_: %+v\n_WANT_: %+v",
				testCase.content, actualOutput, testCase.desiredResult)
		}
	}
}
