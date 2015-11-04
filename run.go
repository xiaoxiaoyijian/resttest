package main

import (
	"flag"
	. "github.com/xiaoxiaoyijian/logger"
	"github.com/xiaoxiaoyijian/resttest/utils/file"
	testcase_util "github.com/xiaoxiaoyijian/resttest/utils/testcase"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"time"
)

var (
	testcaseFlag = flag.String("testcase", "", "set testcase file for running")
	repeatFlag   = flag.Int("repeat", 1, "repeat times to run the testcase")
	emailFlag    = flag.String("email", "", "set email settings to send the result report")
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(int64(time.Second))

	flag.Usage = usage
	flag.Parse()

	if *testcaseFlag == "" {
		usage()
		os.Exit(0)
	}

	testcases, err := parseTestcases(*testcaseFlag)
	if err != nil {
		Logger.Error(err.Error())
		os.Exit(1)
	}

	runner := testcase_util.NewRunner(testcases, *repeatFlag)
	if *emailFlag != "" {
		if email := testcase_util.NewEmail(parseEmail(*emailFlag)); email.IsValid() {
			runner.SetEmail(email)
		}
	}
	runner.Run()
}

func usage() {
	println(`run.go is a tool for running testcase using GO language.

Usage:

    go run run.go --testcase=TESTCASE_FILE [--repeat=REPEAT]

Examples:
    go run run.go --testcase=testcase1.json,testcase2.json --repeat=10
`)
}

func parseTestcases(testcaseStr string) ([]*testcase_util.Testcase, error) {
	result := []*testcase_util.Testcase{}

	vals := strings.Split(testcaseStr, ",")
	for _, v := range vals {
		testcase, err := testcase_util.ParseFile(file.FullName(v))
		if err != nil {
			return nil, err
		}

		result = append(result, testcase)
	}

	return result, nil
}

func parseEmail(str string) map[string]string {
	result := make(map[string]string)
	fields := strings.Split(str, ",")
	for _, v := range fields {
		vals := strings.Split(v, ":")
		if len(vals) >= 2 {
			result[vals[0]] = strings.Join(vals[1:], ":")
		}
	}

	return result
}
