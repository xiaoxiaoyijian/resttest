package testcase

import (
	"github.com/bitly/go-simplejson"
)

const (
	GET  = "get"
	POST = "post"
)

type Testcase struct {
	Name     string
	Request  *Request
	Expected *simplejson.Json
}

func NewEmptyTestcase() *Testcase {
	return &Testcase{
		Request: NewRequest("", GET, nil),
	}
}

func NewTestcase(name, url, method string, expected []byte, auth *Auth) (*Testcase, error) {
	expectedJson, err := simplejson.NewJson(expected)
	if err != nil {
		return nil, err
	}

	testcase := &Testcase{
		Name:     name,
		Request:  NewRequest(url, method, auth),
		Expected: expectedJson,
	}

	if testcase.Name == "" {
		testcase.Name = DefaultName()
	}

	if testcase.Request.Method == "" {
		testcase.Request.Method = GET
	}

	return testcase, nil
}

func (this *Testcase) ToJsonFile() error {
	return CreateJsonFile(this)
}

type Auth struct {
	Url  string
	Vals map[string]string
}

func NewAuth(url string, vals map[string]string) *Auth {
	return &Auth{
		Url:  url,
		Vals: vals,
	}
}

type Request struct {
	Url    string
	Method string
	Auth   *Auth
}

func NewRequest(url, method string, auth *Auth) *Request {
	return &Request{
		Url:    url,
		Method: method,
		Auth:   auth,
	}
}
