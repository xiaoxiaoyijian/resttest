package json

import (
	"github.com/bitly/go-simplejson"
	. "github.com/xiaoxiaoyijian/resttest/utils/logger"
)

func EncodePretty(b []byte) ([]byte, error) {
	json, err := simplejson.NewJson(b)
	if err != nil {
		return nil, err
	}

	return json.EncodePretty()
}

func PrintJson(json *simplejson.Json) {
	content, err := json.EncodePretty()
	if err != nil {
		Logger.Error(err.Error())
		return
	}

	Logger.Info(string(content))
	return
}
