package main

import (
	"flag"
	"github.com/xiaoxiaoyijian/resttest/utils/file"
	testcase_util "github.com/xiaoxiaoyijian/resttest/utils/testcase"
	"log"
	"math/rand"
	"os"
	"runtime"
	"time"
)

var (
	testcaseFlag = flag.String("testcase", "", "set testcase file for running")
	repeatFlag   = flag.Int("repeat", 1, "repeat times to run the testcase")
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

	testcase, err := testcase_util.ParseFile(file.FullName(*testcaseFlag))
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	runner := testcase_util.NewRunner([]*testcase_util.Testcase{testcase}, *repeatFlag)
	runner.Run()
}

func usage() {
	log.Println(`run.go is a tool for running testcase using GO language.

Usage:

    go run run.go --testcase=TESTCASE_FILE [--repeat=REPEAT]

Examples:
	go run run.go --testcase=testcase1.json
	go run run.go --testcase=testcase2.json
`)
}
