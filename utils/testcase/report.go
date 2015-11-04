package testcase

import (
	"fmt"
	"sync"
)

const (
	PASSED  = 1
	FAILED  = 2
	IGNORED = 3
)

type Report struct {
	Testcases []*Testcase
	Times     int
	Summary   *SummaryReport
	Failed    []*FailReport
	m         sync.Mutex
}

func NewReport(testcases []*Testcase, times int) *Report {
	return &Report{
		Testcases: testcases,
		Times:     times,
		Summary:   &SummaryReport{},
		Failed:    []*FailReport{},
	}
}

func (r *Report) AddCase(testcase *Testcase) {
	r.Testcases = append(r.Testcases, testcase)
}

func (r *Report) SetTimes(times int) {
	r.Times = times
}

func (r *Report) Update(testcase string, result int, reason string) {
	r.m.Lock()
	defer r.m.Unlock()

	r.Summary.Total += 1
	switch result {
	case PASSED:
		r.Summary.Succ += 1
	case FAILED:
		r.Summary.Failed += 1
		r.Failed = append(r.Failed, &FailReport{Testcase: testcase, Reason: reason})
	case IGNORED:
		r.Summary.Ignored += 1
	}
}

func (r *Report) String() string {
	template := `Running testcases %d times, Result of it:
    Summary:
        Total: %d, Succ: %d, Failed:%d, Ignored:%d
    `
	result := fmt.Sprintf(template, r.Times, r.Summary.Total, r.Summary.Succ, r.Summary.Failed, r.Summary.Ignored)

	result += `
    Failed:
    `
	cnt := 0
	for _, v := range r.Failed {
		if cnt == 0 {
			result += fmt.Sprintf("    %s : %s\n", v.Testcase, v.Reason)
		} else {
			result += fmt.Sprintf("        %s : %s\n", v.Testcase, v.Reason)
		}
	}

	result += `
    Testcases List:
    `
	cnt = 0
	for _, v := range r.Testcases {
		if cnt == 0 {
			result += fmt.Sprintf("    %s\n", v.Name)
		} else {
			result += fmt.Sprintf("        %s\n", v.Name)
		}
	}

	return result
}

type SummaryReport struct {
	Total   int
	Succ    int
	Failed  int
	Ignored int
}

type FailReport struct {
	Testcase string
	Reason   string
}
