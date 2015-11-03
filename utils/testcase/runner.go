package testcase

import (
	"github.com/bitly/go-simplejson"
	"github.com/xiaoxiaoyijian/resttest/utils/http"
	"github.com/xiaoxiaoyijian/resttest/utils/json"
	. "github.com/xiaoxiaoyijian/resttest/utils/logger"
	"runtime"
	"sync"
)

type Runner struct {
	testcases  []*Testcase
	times      int
	httpClient *http.AuthClient
}

func NewRunner(testcases []*Testcase, times int) *Runner {
	if times <= 0 {
		times = 1
	}

	return &Runner{
		testcases:  testcases,
		times:      times,
		httpClient: http.NewAuthClient(),
	}
}

func (this *Runner) AddCase(testcase *Testcase) {
	this.testcases = append(this.testcases, testcase)
}

func (this *Runner) SetTimes(times int) {
	this.times = times
}

func (this *Runner) Run() {
	if len(this.testcases) == 0 || this.times <= 0 {
		Logger.Info("No testcases to run.")
		return
	}

	Logger.Infof("Start to run testcases: %d times.", this.times)
	defer Logger.Infof("Running testcase %d times done.", this.times)

	var wg sync.WaitGroup
	for i := 0; i < this.times; i++ {
		wg.Add(1)

		go func(n int) {
			defer func() {
				wg.Done()
				handler_recover()
			}()

			var myErr error
			for _, v := range this.testcases {
				Logger.Infof("Running testcase[%d]: %s, REST url: %s", n, v.Name, v.Request.Url)
				myErr = this.runTestcase(v, n)
				if myErr != nil {
					Logger.Infof("Testcase %s[%d] FAILED : %s", v.Name, n, myErr.Error())
					return
				} else {
					Logger.Infof("Testcase %s[%d] PASSED!", v.Name, n)
				}
			}

		}(i)
	}
	wg.Wait()
}

func (this *Runner) runTestcase(testcase *Testcase, n int) error {
	if testcase.Request.Auth != nil {
		if err := this.httpClient.Auth(testcase.Request.Auth.Url, testcase.Request.Auth.Vals); err != nil {
			Logger.Errorf("Auth failed: %s", err.Error())
			return err
		}
	}

	content, err := this.httpClient.Get(testcase.Request.Url)
	if err != nil {
		Logger.Errorf("HTTP GET failed: %s", err.Error())
		return err
	}

	response, err := simplejson.NewJson(content)
	if err != nil {
		Logger.Errorf("Json parse error: %s", err.Error())
		return err
	}

	return json.Compare(response, testcase.Expected)
}

func handler_recover() {
	if err := recover(); err != nil {
		var st = func(all bool) string {
			// Reserve 1K buffer at first
			buf := make([]byte, 512)

			for {
				size := runtime.Stack(buf, all)
				// The size of the buffer may be not enough to hold the stacktrace,
				// so double the buffer size
				if size == len(buf) {
					buf = make([]byte, len(buf)<<1)
					continue
				}
				break
			}

			return string(buf)
		}
		Logger.Errorf("panic:%v\nstack:%v", err, st(false))
	}
}
