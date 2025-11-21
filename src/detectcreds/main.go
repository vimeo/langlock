package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime/debug"

	"detectcreds/scan"

	"detectcreds/allowlist"
	"detectcreds/git"
	"detectcreds/infraction"

	"gopkg.in/alecthomas/kingpin.v2"
)

const EXIT_SUCCESS = 0
const EXIT_FAIL_CANT_LOAD_INPUT_FILE = 1
const EXIT_FAIL_CANT_PARSE_INPUT_FILE_AS_JSON_INFRACTIONS = 2
const EXIT_FAIL_CANT_LOAD_ALLOW_LIST_FILE = 3
const EXIT_FAIL_CANT_PARSE_ALLOW_LIST_FILE_AS_JSON = 4
const EXIT_FAIL_UNEXPECTED_RUNTIME_ERROR = 5

const INPUT_TYPE_LOG = "log"
const INPUT_TYPE_PLAIN = "plain"
const INPUT_TYPE_INFRACTIONS = "infractions"

var INPUT_TYPES = []string{
	INPUT_TYPE_LOG,
	INPUT_TYPE_PLAIN,
	INPUT_TYPE_INFRACTIONS,
}

const OUTPUT_TYPE_LIST = "list"
const OUTPUT_TYPE_SUMMARY = "summary"
const OUTPUT_TYPE_PRETTY = "pretty"

var OUTPUT_TYPES = []string{
	OUTPUT_TYPE_LIST,
	OUTPUT_TYPE_SUMMARY,
	OUTPUT_TYPE_PRETTY,
}

func parseArgs() (string, string, string, string, string, int) {
	var (
		inputFile = kingpin.Arg(
			"input-file",
			"Path of the input file.",
		).Required().String()

		inputType = kingpin.Flag(
			"in-type",
			"Type of the input file (ie. whether to treat the input file"+
				"as unstructured plaintext, as a git log, as a JSON list of infractions, etc.)",
		).Short('i').Default(INPUT_TYPE_PLAIN).Enum(INPUT_TYPES...)

		outputType = kingpin.Flag(
			"out-type",
			"Type of the output file.",
		).Short('o').Default(OUTPUT_TYPE_LIST).Enum(OUTPUT_TYPES...)

		allowListName = kingpin.Flag(
			"allow-list-name",
			"The path and name of the allow list file. "+
				"This value is for labeling output metadata only.",
		).String()

		allowListFile = kingpin.Flag(
			"allow-list",
			"Path of a JSON allow-list file.",
		).Short('a').String()

		numThreads = kingpin.Flag(
			"threads",
			"Number of parallel threads to use during secret scanning.",
		).Short('n').Default("1").Int()
	)

	kingpin.CommandLine.Help = "Scan a git diff for secrets"

	kingpin.Parse()

	return *inputFile, *inputType, *outputType, *allowListName, *allowListFile, *numThreads
}

func main() {

	/***********************************************************
	Setup error handling
	***********************************************************/
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(
				os.Stderr,
				"Panicking: unexpected runtime error in golang program:\n%+v\nStack:\n%s",
				r,
				string(debug.Stack()),
			)
			os.Exit(EXIT_FAIL_UNEXPECTED_RUNTIME_ERROR)
		}
	}()

	/***********************************************************
	Parse command line options
	***********************************************************/
	argInputFile,
		argInputType,
		argOutputType,
		argAllowListName,
		argAllowListFile,
		argNumThreads := parseArgs()

	/***********************************************************
	Invoke core logic/processing
	***********************************************************/
	exitStatus, message := ScanForSecrets(
		argInputFile,
		argInputType,
		argOutputType,
		argAllowListName,
		argAllowListFile,
		argNumThreads,
	)

	/***********************************************************
	Output/print results
	***********************************************************/
	var stream io.Writer
	if exitStatus == EXIT_SUCCESS {
		stream = os.Stdout
	} else {
		stream = os.Stderr
	}

	fmt.Fprintf(stream, message)
	os.Exit(exitStatus)
}

func ScanForSecrets(
	argInputFile string,
	argInputType string,
	argOutputType string,
	argAllowListName string,
	argAllowListFile string,
	argNumThreads int,
) (exitStatus int, message string) {

	var err error

	/***********************************************************
	Load input file
	***********************************************************/
	var content []byte

	content, err = ioutil.ReadFile(argInputFile)
	if err != nil {
		return EXIT_FAIL_CANT_LOAD_INPUT_FILE,
			fmt.Sprintf(
				"Could not load input file %s: %v",
				argInputFile,
				err,
			)
	}

	/***********************************************************
	Find infractions in input file
	***********************************************************/
	var infractions []*infraction.Infraction
	if argInputType == INPUT_TYPE_LOG {
		log := git.ParseLog(content)
		infractions = scan.ScanGitLogForSecrets(log, argNumThreads)
	} else if argInputType == INPUT_TYPE_PLAIN {
		infractions = scan.ScanPlaintextForSecrets(content, argNumThreads)
		for _, infraction := range infractions {
			infraction.Path = argInputFile
		}
	} else { // argInputType == INPUT_TYPE_INFRACTIONS
		infractions, err = infraction.ParseJson(content)
		if err != nil {
			return EXIT_FAIL_CANT_PARSE_INPUT_FILE_AS_JSON_INFRACTIONS,
				fmt.Sprintf(
					"Could not parse input file %s as JSON list of infractions: %v",
					argInputFile,
					err,
				)
		}
	}

	/***********************************************************
	Remove infractions that are pemitted via the allow list
	***********************************************************/
	if argAllowListFile != "" {
		allowListContent, err := ioutil.ReadFile(argAllowListFile)
		if err != nil {
			return EXIT_FAIL_CANT_LOAD_ALLOW_LIST_FILE,
				fmt.Sprintf(
					"Could not load allowlist file %s: %v",
					argAllowListFile,
					err,
				)
		}
		if len(allowListContent) > 0 {
			allowList, err := allowlist.ParseJson(allowListContent)
			if err != nil {
				return EXIT_FAIL_CANT_PARSE_ALLOW_LIST_FILE_AS_JSON,
					fmt.Sprintf(
						"Could not parse input file %s as JSON allow list: %v",
						argAllowListFile,
						err,
					)
			}
			infractions = scan.FilterByAllowedStrings(infractions, allowList.AllowedStrings)
			infractions = scan.FilterByAllowedPaths(infractions, allowList.AllowedPaths)
		}
	}

	/***********************************************************
	Output either a JSON infractions list or suggested JSON allow list
	***********************************************************/
	var report = ""
	if len(infractions) > 0 {
		if argOutputType == OUTPUT_TYPE_LIST {
			infractionsJsonList, _ := json.Marshal(&infractions)
			report = string(infractionsJsonList)
		} else if argOutputType == OUTPUT_TYPE_PRETTY {
			infractionsPretty := infraction.PrettyPrint(infractions, argAllowListName)
			report = string(infractionsPretty)
		} else { // argOutputType == OUTPUT_TYPE_SUMMARY
			infractionsSummary := infraction.Summarize(infractions)
			infractionsSummaryJson, _ := json.Marshal(&infractionsSummary)
			report = string(infractionsSummaryJson)
		}
	}

	return EXIT_SUCCESS, report

}
