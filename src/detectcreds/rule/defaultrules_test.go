package rule

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type testCaseForDefaultRules struct {
	rule                  Rule
	input                 []byte
	desiredNumInfractions int
}

var testCasesForGenericCredBase62Rule = []testCaseForDefaultRules{
	{genericCredBase62OrHexRule, []byte("blah.blah:ssh-rsa AAAFB3NzaC1yc7EAAAADAQABAAABAQDCwGSIN7uopJ63XbSaGwSU/b1MrtX7iNf2sZsTj8agW09yptoQ9Tl1Kya8XM305hbwY7otgDb0Yh4UZG6EUchE9wrqqztQhkgR00G/kUxtpozRR0Nv+19yx+9Yo3CCetuynoLnfeqUnntKYzI8yX54svzYgCbfBviwCPK0Xik2EJ/brRN79tcfVWFATNRL0aYAFsTKpbc8W7KDk8IJwFvSItj4MhZaScGFgR3cpPCPmVrlOsdASCcerC5ryUkZ1Ly09C9LbJL0xjvsWfa5OIL8b034Lwlyr4r9QtPWkV0PBb5cmfwit3VpnxPjyveuVbrhVYNMMp5x7UcbeYl3FnV3 blah.blahblah@BLAH-00844.local"), 0}, // Note that, even though it contains the substring "pw," we do not count it because there is no separator
	{genericCredBase62OrHexRule, []byte("FooBaseApiKey = logBeforeExitingThread()"), 0},
	{genericCredBase62OrHexRule, []byte("secretToken is 'ogBeforeExiting'"), 0},
	{genericCredBase62OrHexRule, []byte("credential: {\"NewFactoryObject\""), 0},
	{genericCredBase62OrHexRule, []byte("Password: \n numWidgetsPerHour"), 0},
	{genericCredBase62OrHexRule, []byte("pw == avg-Num-Params-nArr"), 0},
	{genericCredBase62OrHexRule, []byte("sEcReT this_Is_A_Rand_Var_Name"), 0},
	{genericCredBase62OrHexRule, []byte("FacebookCred = UpperExclusiveBound"), 0},
	{genericCredBase62OrHexRule, []byte("FacebookCred = upperExclusiveBound"), 0},
	{genericCredBase62OrHexRule, []byte("The password is not log_before_exiting_thread"), 0},
	{genericCredBase62OrHexRule, []byte("Nor is the secret locations-numbered-list-str"), 0},
	{genericCredBase62OrHexRule, []byte("token: mapVulnerabilities2Fixes"), 0},
	{genericCredBase62OrHexRule, []byte("pw: NewEvenChatDistributionRule"), 0},
	{genericCredBase62OrHexRule, []byte("tokens: 1\n+      Dropbox.WebhookAutoUpload"), 0},
	{genericCredBase62OrHexRule, []byte("tokens: 60\n+      PopularClip.PopularClipsV6"), 0},
	{genericCredBase62OrHexRule, []byte("tokens: 40\n+      Search.Update.Queue"), 0},
	{genericCredBase62OrHexRule, []byte("FooBaseApiKey = 0xABCDEF0123456788982"), 1},
	{genericCredBase62OrHexRule, []byte("FooBaseApiKey = g6oGqje8z9GuPHmX"), 1},
	{genericCredBase62OrHexRule, []byte("secretToken is 8P6Df6QwVNqwdSNn"), 1},
	{genericCredBase62OrHexRule, []byte("credential: {\"dNsbrGG9MhJc2pvzV\""), 1},
	{genericCredBase62OrHexRule, []byte("Password: Cb6Uu2M7SGpeKwd8"), 1},
	{genericCredBase62OrHexRule, []byte("pw == ctw-4jRp-mqm-6SD-78E"), 1},
	{genericCredBase62OrHexRule, []byte("password. N2_m4_JVfGUX_bepQaW"), 1},
	{genericCredBase62OrHexRule, []byte("key <V6o2b2L5BmYKL2QW>"), 1},
	{genericCredBase62OrHexRule, []byte("apiKey --> 3Fuam7a9d2PVRosp"), 1},
	{genericCredBase62OrHexRule, []byte("OAUTH2.0_token: JLs3e5Vi9RreRcJaZsp4SRTu"), 1},
	{genericCredBase62OrHexRule, []byte("encryption_key='gLW8Uu8bELhtLaZVHi6KB4Zk'"), 1},
	{genericCredBase62OrHexRule, []byte("FooBaseApiKey = 'xjSNOT0ujQF4a72rEujc9LYNO'"), 1},
	{genericCredBase62OrHexRule, []byte("apiKey == \"xxsZkecRY3nmwjZ2K4L8m4PSd\""), 1},
	{genericCredBase62OrHexRule, []byte("apiKey == 'qi7JBWH2McbMTxdrekz94x4r=="), 1},
	{genericCredBase62OrHexRule, []byte("apiKey == 'Xn9cgHbKJrXorNENh2BTFgQSb2yPmpCJAFrFMxEYyg4S5ZEZeGbW32Hoy6ypPuJ"), 1},
	{genericCredBase62OrHexRule, []byte("apiKey == 'HcWwD4moxqDGWirtp3K6Gy2XEX6JnNGB5weLTSZWH7QccQMLpcwk9LiJVnKjpao"), 1},
	{genericCredBase62OrHexRule, []byte("apiKey == 's9gkcuaqpSLdzJ6a8cDd9HCCAMGfWYsUmif3poA65s3eRRpWNSjrGNJUdVaj5Qf"), 1},
	{genericCredBase62OrHexRule, []byte("apiKey == 'sWjbng83muWUs3aTLvVDzbX4yh8WpuZsmFFACufvwoN7JTaR8cNMpYE2nM5CF6V"), 1},
	{genericCredBase62OrHexRule, []byte("apiKey == 'jjbk6zN8aQTW9kp9Bcqu2c7zBKekKusdSy6eFhwdWW9G9LJhGJvWVWoEpwFYib3"), 1},
	{genericCredBase62OrHexRule, []byte("apiKey == 'WA3FHz9o2nZZw75hxQ2uVnbSnHwHQZFKEbLoPqJqD7uW93H3hJRjSPP55umzSAk"), 1},
	{genericCredBase62OrHexRule, []byte("apiKey \n 'CcX7mn6jHRJ883eZo3JJo6cEpdapz2RXyiu8CxmUry3WTtxQGzwdBEhx5DikEF7"), 1},
	{genericCredBase62OrHexRule, []byte("apiKey: \n { \n 'CcX7mn6jHRJ883eZo3JJo6cEpdapz2RXyiu8CxmUry3WTtxQGzwdBEhx5DikEF7"), 1},
	{genericCredLowercaseHexRule, []byte("FooBaseApiKey = logBeforeExitingThread()"), 0},
	{genericCredLowercaseHexRule, []byte("secretToken is 'ogBeforeExiting'"), 0},
	{genericCredLowercaseHexRule, []byte("credential: {\"NewFactoryObject\""), 0},
	{genericCredLowercaseHexRule, []byte("Password: \n numWidgetsPerHour"), 0},
	{genericCredLowercaseHexRule, []byte("pw == avg-Num-Params-nArr"), 0},
	{genericCredLowercaseHexRule, []byte("foobaseapikey = abb0119765f2e3da09ec"), 1},
	{genericCredLowercaseHexRule, []byte("foobaseapikey = 0xabb0119765f2e3da09ec()"), 1},
	{genericCredLowercaseHexRule, []byte("secretToken is '673eff803468c21dbcba'"), 1},
	{genericCredLowercaseHexRule, []byte("SecretToken is '673eFF803468c21dbcba'"), 0}, // Not finding because case is mixed
	{genericCredLowercaseHexRule, []byte("secretToken is '673e-ff80-3468-c21-dbcba'"), 1},
	{genericCredLowercaseHexRule, []byte("secretToken is '673e_ff80_3468_c21_dbcba'"), 1},
	{genericCredUppercaseHexRule, []byte("FooBaseApiKey = logBeforeExitingThread()"), 0},
	{genericCredUppercaseHexRule, []byte("secretToken is 'ogBeforeExiting'"), 0},
	{genericCredUppercaseHexRule, []byte("credential: {\"NewFactoryObject\""), 0},
	{genericCredUppercaseHexRule, []byte("Password: \n numWidgetsPerHour"), 0},
	{genericCredUppercaseHexRule, []byte("pw == avg-Num-Params-nArr"), 0},
	{genericCredUppercaseHexRule, []byte("foobaseapikey = abb0119765f2e3da09ec"), 0},
	{genericCredUppercaseHexRule, []byte("foobaseapikey = ABB0119765F2E3DA09EC"), 1},
	{genericCredUppercaseHexRule, []byte("foobaseapikey = 0xABB0119765F2E3DA09EC()"), 1},
	{genericCredUppercaseHexRule, []byte("secretToken is '673EFF803468C21DBCBA'"), 1},
	{genericCredUppercaseHexRule, []byte("SecretToken is '673eFF803468c21dbcba'"), 0}, // Not finding because case is mixed
	{genericCredUppercaseHexRule, []byte("secretToken is '673E-FF80-3468-C21-DBCBA'"), 1},
	{genericCredUppercaseHexRule, []byte("secretToken is '673E_FF80_3468_C21_DBCBA'"), 1},
	{genericCredUppercaseHexRule, []byte("FooBaseApiKey = 0xABCDEF0123456788982"), 1},
}

func TestGenericCredBase62RRule(t *testing.T) {
	for _, testCase := range testCasesForGenericCredBase62Rule {
		actualOutput := testCase.rule.Check(testCase.input, nil)
		if !cmp.Equal(len(actualOutput), testCase.desiredNumInfractions) {
			t.Errorf("Incorrect number of infractions found for rule <%s> when applied to string: %+v\n_HAVE_: %+v\n_WANT_: %+v",
				testCase.rule.Name, string(testCase.input), len(actualOutput), testCase.desiredNumInfractions)
		}
	}
}
