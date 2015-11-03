package json

import (
	"github.com/bitly/go-simplejson"
	"log"
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
		log.Println(err.Error())
		return
	}

	log.Println(string(content))
	return
}
