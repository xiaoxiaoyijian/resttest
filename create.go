package main

import (
	"flag"
	. "github.com/xiaoxiaoyijian/logger"
	"github.com/xiaoxiaoyijian/resttest/utils/http"
	testcase_util "github.com/xiaoxiaoyijian/resttest/utils/testcase"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"time"
)

var (
	urlFlag    = flag.String("url", "", "set request url")
	outputFlag = flag.String("output", "", "set output json filename")

	authUrlFlag = flag.String("authurl", "", "set url for authorize")
	authValFlag = flag.String("authval", "", "set values for authorize")
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(int64(time.Second))

	flag.Usage = usage
	flag.Parse()

	if *urlFlag == "" || (*authUrlFlag != "" && *authValFlag == "") {
		usage()
		os.Exit(0)
	}

	Logger.Infof("Start to create testcase: %s REST url: %s.", *outputFlag, *urlFlag)
	defer Logger.Infof("Create testcase: %s done.", *outputFlag)

	var auth *testcase_util.Auth
	client := http.NewAuthClient()
	if *authUrlFlag != "" {
		authValues := parseAuthVals(*authValFlag)
		auth = testcase_util.NewAuth(*authUrlFlag, authValues)

		err := client.Auth(auth.Url, auth.Vals)
		checkErr(err)
	}
	result, err := client.Get(*urlFlag)
	checkErr(err)

	testcase, err := testcase_util.NewTestcase(*outputFlag, *urlFlag, testcase_util.GET, result, auth)
	checkErr(err)

	err = testcase.ToJsonFile()
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		Logger.Error(err.Error())
		os.Exit(1)
	}
}

func usage() {
	println(`create.go is a tool for creating testcase using GO language, optional you can with authorization firstly.

Usage:

    go run create.go --url=URL [--output=OUTPUT] [--authurl=AUTHURL --authvals=]
}
`)
}

func parseAuthVals(valStr string) map[string]string {
	authValues := make(map[string]string)
	fields := strings.Split(valStr, ";")
	for _, v := range fields {
		vals := strings.Split(v, ":")
		if len(vals) == 2 {
			authValues[vals[0]] = vals[1]
		}
	}

	return authValues
}
