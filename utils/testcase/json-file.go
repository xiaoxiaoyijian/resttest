package testcase

import (
	myjson "encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/xiaoxiaoyijian/resttest/utils/file"
	"github.com/xiaoxiaoyijian/resttest/utils/json"
	"strconv"
	"time"
)

func DefaultName() string {
	return fmt.Sprintf("%s_test.json", time.Now().Format("20060102150405"))
}

func CreateJsonFile(testcase *Testcase) error {
	expected, err := testcase.Expected.Encode()
	if err != nil {
		return err
	}

	var template, jsonStr string
	if testcase.Request.Auth != nil {
		valStr, _ := myjson.Marshal(testcase.Request.Auth.Vals)
		template = `{"name":"%s", "request":{"url": "%s", "method":"%s", "auth":{"url":"%s", "vals":%s}}, "response": %s}`
		jsonStr = fmt.Sprintf(template, testcase.Name, testcase.Request.Url, testcase.Request.Method, testcase.Request.Auth.Url, valStr, string(expected))
	} else {
		template = `{"name":"%s", "request":{"url": "%s", "method":"%s"}, "response": %s}`
		jsonStr = fmt.Sprintf(template, testcase.Name, testcase.Request.Url, testcase.Request.Method, string(expected))
	}
	content, err := json.EncodePretty([]byte(jsonStr))
	if err != nil {
		return err
	}

	return file.CreateAndWrite(testcase.Name, content)
}

func ParseFile(filename string) (*Testcase, error) {
	content, err := file.ReadAll(filename)
	if err != nil {
		return nil, err
	}

	return ParseByteArray(content, filename)
}

func ParseByteArray(content []byte, name string) (*Testcase, error) {
	json, err := simplejson.NewJson(content)
	if err != nil {
		return nil, err
	}

	testcase := NewEmptyTestcase()
	testcase.Name = name
	testcase.Expected = json.Get("response")
	testcase.Request.Url, err = json.Get("request").Get("url").String()
	if err != nil {
		return nil, err
	}
	testcase.Request.Method, err = json.Get("request").Get("method").String()
	if err != nil {
		return nil, err
	}

	auth := json.Get("request").Get("auth")
	if auth.Interface() != nil {
		authUrl, err := auth.Get("url").String()
		if err != nil {
			return nil, err
		}
		vals, err := auth.Get("vals").Map()
		if err != nil {
			return nil, err
		}
		authVals := make(map[string]string)
		for k, v := range vals {
			authVals[k] = StringAssert(v)
		}

		testcase.Request.Auth = NewAuth(authUrl, authVals)
	}

	return testcase, nil
}

// StringAssert 转化数据为string类型
func StringAssert(val interface{}) string {
	switch val.(type) {
	case uint:
		return strconv.FormatUint(uint64(val.(uint)), 10)
	case int:
		return strconv.FormatInt(int64(val.(int)), 10)
	case int64:
		return strconv.FormatInt(val.(int64), 10)
	case int32:
		return strconv.FormatInt(int64(val.(int32)), 10)
	case int16:
		return strconv.FormatInt(int64(val.(int16)), 10)
	case int8:
		return strconv.FormatInt(int64(val.(int8)), 10)
	case uint8:
		return strconv.FormatUint(uint64(val.(uint8)), 10)
	case uint16:
		return strconv.FormatUint(uint64(val.(uint16)), 10)
	case uint32:
		return strconv.FormatUint(uint64(val.(uint32)), 10)
	case uint64:
		return strconv.FormatUint(val.(uint64), 10)
	case float64:
		return strconv.FormatFloat(val.(float64), 'f', 'f', 64)
	case float32:
		return strconv.FormatFloat(float64(val.(float32)), 'f', 'f', 32)
	case nil:
		return ""
	case []byte:
		return string(val.([]byte))
	case string:
		return val.(string)
	case bool:
		return strconv.FormatBool(val.(bool))
	}
	return val.(string)
}
