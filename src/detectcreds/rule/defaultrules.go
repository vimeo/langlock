package rule

import (
	"regexp"

	"detectcreds/rule/digit"
	"detectcreds/rule/english"
	"detectcreds/rule/hex"
	"detectcreds/rule/punctuation"
)

func compileRegexOrPanic(s string) *regexp.Regexp {
	r, err := regexp.Compile(s)
	if err != nil {
		panic(err)
	}
	return r
}

var genericCredentialRegexPrefix = `((?si)api_key|apikey|secret|password|pass|pw|key|token|oauth|bearer|authorization|auth|login|cred|credential|private|signature)((?s).{0,20}?)?`

var genericCredBase62OrHexRule = Rule{
	Name:         "Generic credential - base 64/62/36 or hex",
	ContentRegex: compileRegexOrPanic(genericCredentialRegexPrefix + `([0-9a-zA-Z~_\/.+=-]{16,64})`),
	ContentChecks: []ContentCheck{
		// ContentCheck{english.CheckIsNotEnglish, []int{3}},
		// ContentCheck{hex.CheckIsNotHex, []int{3}},
		ContentCheck{
			func(strings []string) bool {
				return english.CheckIsNotEnglish(strings) || !hex.CheckIsNotHex(strings)
			},
			[]int{3}},
		ContentCheck{digit.CheckHasDigit, []int{3}},
		ContentCheck{punctuation.CheckHasSeparator, []int{2}},

		//distribution.NewDistributionCheck(distribution.NewEvenRuneDistributionRule(
		//	[]int{1},
		//	[]string{"abcdefghijklmnopqrstuvwxyz", "ABCDEFGHIJKLMNOPQRSTUVWXYZ", "0123456789"},
		//	[]float64{90.0, 90.0, 90.0},
		//)),
	},
	MatchFormatter: func(match string, groups []string) string { return groups[3] },
}

// var genericCredLowercaseAlphaNumBase36Rule = Rule{
// Name:         "Generic credential - lowercase alphanumeric / base 36",
// ContentRegex: compileRegexOrPanic(genericCredentialRegexPrefix + `([0-9a-z~_\/.+=-]{16,64})`),
// ContentChecks: []func([]string) bool{
// distribution.NewDistributionCheck(distribution.NewEvenRuneDistributionRule(
// []int{1},
// []string{"abcdefghijklmnopqrstuvwxyz", "0123456789"},
// []float64{50.0, 50.0},
// )),
// english.CheckIsNotEnglish,
// },
// }
//
// var genericCredUppercaseAlphaNumBase36Rule = Rule{
// Name:         "Generic credential - uppercase alphanumeric / base 36",
// ContentRegex: compileRegexOrPanic(genericCredentialRegexPrefix + `([0-9A-Z~_\/.+=-]{16,64})`),
// ContentChecks: []func([]string) bool{
// distribution.NewDistributionCheck(distribution.NewEvenRuneDistributionRule(
// []int{1},
// []string{"ABCDEFGHIJKLMNOPQRSTUVWXYZ", "0123456789"},
// []float64{50.0, 50.0},
// )),
// english.CheckIsNotEnglish,
// },
// }

var genericCredLowercaseHexRule = Rule{
	Name:         "Generic credential - lowercase hex",
	ContentRegex: compileRegexOrPanic(genericCredentialRegexPrefix + `([0-9a-fx~_\/.+=-]{16,64})`),
	// ContentChecks: []func([]string) bool{
	// 	distribution.NewDistributionCheck(distribution.NewEvenRuneDistributionRule(
	// 		[]int{1},
	// 		[]string{"abcdef", "0123456789"},
	// 		[]float64{50.0, 50.0},
	// 	)),
	// },
}
var genericCredHexRule = Rule{
	Name:         "Generic credential - hex",
	ContentRegex: compileRegexOrPanic(genericCredentialRegexPrefix + `([0-9A-Fx~_\/.+=-]{16,64}|[0-9a-fx~_\/.+=-]{16,64})`),
}
var genericCredUppercaseHexRule = Rule{
	Name:         "Generic credential - uppercase hex",
	ContentRegex: compileRegexOrPanic(genericCredentialRegexPrefix + `([0-9A-Fx~_\/.+=-]{16,64})`),
	// ContentChecks: []func([]string) bool{
	// 	distribution.NewDistributionCheck(distribution.NewEvenRuneDistributionRule(
	// 		[]int{1},
	// 		[]string{"ABCDEF", "0123456789"},
	// 		[]float64{50.0, 50.0},
	// 	)),
	// },
}

var DefaultRules = []Rule{
	genericCredBase62OrHexRule,
	// genericCredLowercaseAlphaNumBase36Rule,
	// genericCredUppercaseAlphaNumBase36Rule,
	// genericCredLowercaseHexRule,
	// genericCredUppercaseHexRule,
	//genericCredHexRule,

	//{
	//	Name: "Combined megarule",
	//	ContentRegex: compileRegexOrPanic(strings.Join([]string{
	//		`((?i:postgresql|postgres|psql|mongodb|mongo|mdb|mysql|sql|https|http|sftp|ftp)://[a-zA-Z0-9]+:[a-zA-Z0-9]+@[/a-zA-Z0-9?=&_.,-]+)`,
	//		`((A3T[A-Z0-9]|AKIA|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{16})`,
	//		`((?i)aws(.{0,20})?(?-i)[0-9a-zA-Z\/+]{40})`,
	//		`((?i)aws(.{0,20})?(?-i)['\"][0-9a-zA-Z\/+]{40}['\"])`,
	//		`(amzn\.mws\.[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})`,
	//		`(da2-[a-z0-9]{26})`,
	//		`(EAACEdEose0cBA[0-9A-Za-z]+)`,
	//		`((xox[pborsa]-[0-9]{12}-[0-9]{12}-[0-9]{12}-[a-z0-9]{32}))`,
	//		`(curl .{0,30}(-u|--user) \S+:\S+)`,
	//		`((A3T[A-Z0-9]|AKIA|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{16})`,
	//		`((?i)aws(.{0,20})?(?-i)[0-9a-zA-Z\/+]{40})`,
	//		`((?i)aws(.{0,20})?(?-i)['\"][0-9a-zA-Z\/+]{40}['\"])`,
	//		`(amzn\.mws\.[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})`,
	//		`(da2-[a-z0-9]{26})`,
	//		`((?i)(facebook|fb)(.{0,20})?(?-i)['\"][0-9a-f]{32}['\"])`,
	//		`((?i)(facebook|fb)(.{0,20})?['\"][0-9]{13,17}['\"])`,
	//		`(EAACEdEose0cBA[0-9A-Za-z]+)`,
	//		`((?i)(?:fastly).{0,40}\b([A-Za-z0-9_-]{32})\b)`,
	//		`((?i)github(.{0,20})?(?-i)[0-9a-zA-Z]{35,40})`,
	//		`((?i)(?:jira).{0,40}\b([a-zA-Z-0-9]{24})\b)`,
	//		`((?i)(?:jira).{0,40}\b([a-zA-Z-0-9]{5,24}\\@[a-zA-Z-0-9]{3,16}\\.com)\b)`,
	//		`((?i)linkedin(.{0,20})?(?-i)[0-9a-z]{12})`,
	//		`((?i)linkedin(.{0,20})?[0-9a-z]{16})`,
	//		`(xox[baprs]-([0-9a-zA-Z]{10,48})?)`,
	//		`((xox[pborsa]-[0-9]{12}-[0-9]{12}-[0-9]{12}-[a-z0-9]{32}))`,
	//		`(-----BEGIN ((EC|PGP|DSA|RSA|OPENSSH|ENCRYPTED) )?PRIVATE KEY( BLOCK)?-----)`,
	//		`(AIza[0-9A-Za-z\\-_]{35})`,
	//		`("type": "service_account")`,
	//		`([0-9]+-[0-9A-Za-z_]{32}\\.apps\\.googleusercontent\\.com)`,
	//		`(ya29\\.[0-9A-Za-z\\-_]+)`,
	//		`((?i)heroku(.{0,20})?[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})`,
	//		`([hH][eE][rR][oO][kK][uU].*[0-9A-F]{8}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{12})`,
	//		`((?i)(mailchimp|mc)(.{0,20})?[0-9a-f]{32}-us[0-9]{1,2})`,
	//		`(((?i)(mailgun|mg)(.{0,20})?)?key-[0-9a-z]{32})`,
	//		`([a-z0-9-]{1,40}\.okta(?:preview|-emea){0,1}\.com)`,
	//		`(00[a-zA-Z0-9_-]{40})`,
	//		`(access_token\$production\$[0-9a-z]{16}\$[0-9a-f]{32})`,
	//		`(sk_live_[0-9a-z]{32})`,
	//		`(SG\.[\w_]{16,32}\.[\w_]{16,64})`,
	//		`(https://hooks.slack.com/services/T[a-zA-Z0-9_]{8}/B[a-zA-Z0-9_]{8,12}/[a-zA-Z0-9_]{24})`,
	//		`((?i)stripe(.{0,20})?[sr]k_live_[0-9a-zA-Z]{24})`,
	//		`(sq0atp-[0-9A-Za-z\-_]{22})`,
	//		`(sq0csp-[0-9A-Za-z\\-_]{43})`,
	//		`([0-9]+:AA[0-9A-Za-z\\-_]{33})`,
	//		`((?i)twilio(.{0,20})?SK[0-9a-f]{32})`,
	//		`((?i)twitter(.{0,20})?[0-9a-z]{18,25})`,
	//		`((?i)twitter(.{0,20})?[0-9a-z]{35,44})`,
	//		`([tT][wW][iI][tT][tT][eE][rR].*[1-9][0-9]+-[0-9a-zA-Z]{40})`,
	//	}, "|")),
	//},
	{
		Name:         "DB/FTP/SFTP/HTTP/HTTPs Connection String",
		ContentRegex: compileRegexOrPanic(`(?i:postgresql|postgres|psql|mongodb|mongo|mdb|mysql|sql|https|http|sftp|ftp)://[a-zA-Z0-9]+:[a-zA-Z0-9]+@[/a-zA-Z0-9?=&_.,-]+`),
	},
	//	{
	//		Name:         "HTTP/HTTPS Connection String",
	//		ContentRegex: compileRegexOrPanic(`(?i:https|http)://[a-zA-Z0-9]+:[a-zA-Z0-9]+@[/a-zA-Z0-9?=&_.,-]+`),
	//	},
	//	{
	//		Name:         "MongoDB Connection String",
	//		ContentRegex: compileRegexOrPanic(`(?i:mongodb|mongo|mdb)://[a-zA-Z0-0]+:[a-zA-Z0-9]+@[/a-zA-Z0-9?=&_.,-]+`),
	//	},
	//	{
	//		Name:         "FTP/SFTP Connection String",
	//		ContentRegex: compileRegexOrPanic(`(?i:sftp|ftp)://[a-zA-Z0-9]+:[a-zA-Z0-9]+@[/a-zA-Z0-9?=&_.,-]+`),
	//	},
	{
		Name:           "Curl with Password",
		ContentRegex:   compileRegexOrPanic(`(?s)curl .{0,30}(-u|--user) (\S+:\S+)`),
		MatchFormatter: func(match string, groups []string) string { return groups[2] },
	},
	{
		Name:      "File netrc",
		PathRegex: compileRegexOrPanic(`.*netrc$`),
	},
	{
		Name:      "File etc/passwd",
		PathRegex: compileRegexOrPanic(`.*etc/passwd.*`),
	},
	{
		Name:      "File etc/shadow",
		PathRegex: compileRegexOrPanic(`.*etc/shadow.*`),
	},
	{
		Name:      "File private ssh key id_rsa",
		PathRegex: compileRegexOrPanic(`.*id_rsa$`),
	},
	{
		Name:         "AWS Manager ID",
		ContentRegex: compileRegexOrPanic(`(A3T[A-Z0-9]|AKIA|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{16}`),
	},
	{
		Name:         "AWS Secret Key v1",
		ContentRegex: compileRegexOrPanic(`(?is)aws(.{0,20})?(?-i)[0-9a-zA-Z\/+]{40}`),
	},
	{
		Name:         "AWS Secret Key v2",
		ContentRegex: compileRegexOrPanic(`(?is)aws(.{0,20})?(?-i)['\"][0-9a-zA-Z\/+]{40}['\"]`),
	},
	{
		Name:         "AWS MWS key",
		ContentRegex: compileRegexOrPanic(`amzn\.mws\.[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`),
	},
	{
		Name:         "AWS AppSync GraphQL Key",
		ContentRegex: compileRegexOrPanic(`da2-[a-z0-9]{26}`),
	},
	{
		Name:           "Facebook Secret Key",
		ContentRegex:   compileRegexOrPanic(`(?is)(facebook|fb)(.{0,20})?(?-i)['\"]([0-9a-f]{32})['\"]`),
		MatchFormatter: func(match string, groups []string) string { return groups[3] },
	},
	{
		Name:           "Facebook Client ID",
		ContentRegex:   compileRegexOrPanic(`(?is)(facebook|fb)(.{0,20})?['\"]([0-9]{13,17})['\"]`),
		MatchFormatter: func(match string, groups []string) string { return groups[3] },
	},
	{
		Name:         "Facebook Access Token",
		ContentRegex: compileRegexOrPanic(`EAACEdEose0cBA[0-9A-Za-z]+`),
	},
	{
		Name:           "Fastly Personal Token",
		ContentRegex:   compileRegexOrPanic(`(?is)(?:fastly).{0,40}\b([A-Za-z0-9_-]{32})\b`),
		MatchFormatter: func(match string, groups []string) string { return groups[1] },
	},
	{
		Name:           "Github",
		ContentRegex:   compileRegexOrPanic(`(?i)github((?s).{0,20})?(?-i)([0-9a-zA-Z]{35,40})`),
		MatchFormatter: func(match string, groups []string) string { return groups[2] },
	},
	{
		Name:           "Jira - 1",
		ContentRegex:   compileRegexOrPanic(`(?is)(?:jira).{0,40}\b([a-zA-Z-0-9]{24})\b`),
		MatchFormatter: func(match string, groups []string) string { return groups[1] },
	},
	{
		Name:           "Jira - 2",
		ContentRegex:   compileRegexOrPanic(`(?is)(?:jira).{0,40}\b([a-zA-Z-0-9]{5,24}\\@[a-zA-Z-0-9]{3,16}\\.com)\b`),
		MatchFormatter: func(match string, groups []string) string { return groups[1] },
	},
	{
		Name:           "LinkedIn Client ID",
		ContentRegex:   compileRegexOrPanic(`(?i)linkedin((?s).{0,20})?(?-i)([0-9a-z]{12})`),
		MatchFormatter: func(match string, groups []string) string { return groups[2] },
		ContentChecks: []ContentCheck{
			ContentCheck{english.CheckIsNotEnglish, []int{2}},
		},
	},
	{
		Name:           "LinkedIn Secret Key",
		ContentRegex:   compileRegexOrPanic(`(?i)linkedin((?s).{0,20})?([0-9a-z]{16})`),
		MatchFormatter: func(match string, groups []string) string { return groups[2] },
		ContentChecks: []ContentCheck{
			ContentCheck{english.CheckIsNotEnglish, []int{2}},
		},
	},
	{
		Name:         "Slack",
		ContentRegex: compileRegexOrPanic(`xox[baprs]-([0-9a-zA-Z]{10,48})?`),
	},
	{
		Name:         "Slack token",
		ContentRegex: compileRegexOrPanic(`(xox[pborsa]-[0-9]{12}-[0-9]{12}-[0-9]{12}-[a-z0-9]{32})`),
	},
	{
		// Largest commonish key file size is ~12,749 chars (produced from `ssh-keygen -b 16384`)
		// Smallest commish key file is ~1,064 chars (produced from `ssh-keygen -b 1024`)
		// But fun fact, GoLang does not support regex {n,m} quantifier values > 1000. See here: https://github.com/golang/go/issues/7252
		Name:         "Asymmetric Private Key",
		ContentRegex: compileRegexOrPanic(`(?s)-----BEGIN ((EC|PGP|DSA|RSA|OPENSSH|ENCRYPTED) )?PRIVATE KEY( BLOCK)?-----[a-zA-Z0-9\n+/\\=]{900,1000}`),
		// ContentRegex: compileRegexOrPanic(`(?s)-----BEGIN ((EC|PGP|DSA|RSA|OPENSSH|ENCRYPTED) )?PRIVATE KEY( BLOCK)?-----[a-zA-Z0-9\n+/=]{900,14000}-----END ((EC|PGP|DSA|RSA|OPENSSH|ENCRYPTED) )?PRIVATE KEY( BLOCK)?-----`),
	},
	{
		Name:         "Google or YouTube or Gmail or GDrive API key",
		ContentRegex: compileRegexOrPanic(`AIza[0-9A-Za-z\\-_]{35}`),
	},
	// {
	// NOTE - commenting this out. The fields (e.g. private key) within the service account json should trigger other detection rules
	// ALSO NOTE - if you do indeed uncomment this rule, you likely to adjust it to match on the unique part, otherwise any allow list may block unrelated entries
	// Name:         "Google (GCP) Service Account",
	// ContentRegex: compileRegexOrPanic(`"type": "service_account"`),
	// },
	{
		Name:         "Google Cloud Platform or Google Drive or GMail or YouTube OAuth",
		ContentRegex: compileRegexOrPanic(`[0-9]+-[0-9A-Za-z_]{32}\\.apps\\.googleusercontent\\.com`),
	},
	{ // NOTE: highly skeptical that this regex is correct. I have not been able to verify that these tokens start with "ya29"
		Name:         "Google OAuth Access Token",
		ContentRegex: compileRegexOrPanic(`ya29\\.[0-9A-Za-z\\-_]+`),
	},
	{
		Name:           "Heroku API key",
		ContentRegex:   compileRegexOrPanic(`(?i)heroku((?s).{0,20})?([0-9A-Fa-f]{8}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{12})`),
		MatchFormatter: func(match string, groups []string) string { return groups[2] },
	},
	{
		Name:           "Heroku API key alt.",
		ContentRegex:   compileRegexOrPanic(`(?i)heroku((?s).{0,20})?([0-9A-Fa-f]{8}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{12})`),
		MatchFormatter: func(match string, groups []string) string { return groups[2] },
	},
	{
		Name:           "MailChimp API key",
		ContentRegex:   compileRegexOrPanic(`(?i)(mailchimp|mc)((?s).{0,20})?([0-9a-f]{32}-us[0-9]{1,2})`),
		MatchFormatter: func(match string, groups []string) string { return groups[3] },
	},
	{
		Name:           "Mailgun API key",
		ContentRegex:   compileRegexOrPanic(`(?i)(mailgun|mg)((?s).{0,20})?(key-[0-9a-z]{32})`),
		MatchFormatter: func(match string, groups []string) string { return groups[3] },
	},
	{
		Name:         "Okta - 1",
		ContentRegex: compileRegexOrPanic(`[a-z0-9-]{1,40}\.okta(?:preview|-emea){0,1}\.com`),
	},
	{
		Name:           "Okta - 2",
		ContentRegex:   compileRegexOrPanic(`(?si)okta(.{0,20})?(00[a-zA-Z0-9_-]{40})\b`),
		MatchFormatter: func(match string, groups []string) string { return groups[2] },
	},
	{
		Name:         "PayPal Braintree access token",
		ContentRegex: compileRegexOrPanic(`access_token\$production\$[0-9a-z]{16}\$[0-9a-f]{32}`),
	},
	{
		Name:         "Picatic API key",
		ContentRegex: compileRegexOrPanic(`sk_live_[0-9a-z]{32}`),
	},
	{
		Name:         "SendGrid API Key",
		ContentRegex: compileRegexOrPanic(`SG\.[\w_]{16,32}\.[\w_]{16,64}`),
	},
	{
		Name:         "Slack Webhook",
		ContentRegex: compileRegexOrPanic(`https://hooks.slack.com/services/T[a-zA-Z0-9_]{8}/B[a-zA-Z0-9_]{8,12}/[a-zA-Z0-9_]{24}`),
	},
	{
		Name:           "Stripe API key",
		ContentRegex:   compileRegexOrPanic(`(?i)stripe((?s).{0,20})?([sr]k_live_[0-9a-zA-Z]{24})`),
		MatchFormatter: func(match string, groups []string) string { return groups[2] },
	},
	{
		Name:         "Square access token",
		ContentRegex: compileRegexOrPanic(`sq0atp-[0-9A-Za-z\-_]{22}`),
	},
	{
		Name:         "Square OAuth secret",
		ContentRegex: compileRegexOrPanic(`sq0csp-[0-9A-Za-z\\-_]{43}`),
	},
	{
		Name:         "Telegram Bot API Key",
		ContentRegex: compileRegexOrPanic(`[0-9]+:AA[0-9A-Za-z\\-_]{33}`),
	},
	{
		Name:           "Twilio API key",
		ContentRegex:   compileRegexOrPanic(`(?i)twilio((?s).{0,20})?(SK[0-9a-f]{32})`),
		MatchFormatter: func(match string, groups []string) string { return groups[2] },
	},
	{
		Name:           "Twitter Client ID",
		ContentRegex:   compileRegexOrPanic(`(?i)twitter((?s).{0,20})?([0-9a-z]{18,25})`),
		MatchFormatter: func(match string, groups []string) string { return groups[2] },
		ContentChecks: []ContentCheck{
			ContentCheck{english.CheckIsNotEnglish, []int{2}},
		},
	},
	{
		Name:           "Twitter Secret Key",
		ContentRegex:   compileRegexOrPanic(`(?i)twitter((?s).{0,20})?([0-9a-z]{35,44})`),
		MatchFormatter: func(match string, groups []string) string { return groups[2] },
		ContentChecks: []ContentCheck{
			ContentCheck{english.CheckIsNotEnglish, []int{2}},
		},
	},
	{
		Name:           "Twitter Access Token",
		ContentRegex:   compileRegexOrPanic(`(?i)twitter.{0,20}([1-9][0-9]+-[0-9a-zA-Z]{40})`),
		MatchFormatter: func(match string, groups []string) string { return groups[1] },
		ContentChecks: []ContentCheck{
			ContentCheck{english.CheckIsNotEnglish, []int{1}},
		},
	},
}
