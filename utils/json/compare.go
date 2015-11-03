package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"reflect"
)

func Compare(A, B *simplejson.Json) error {
	return compareInterface(A.Interface(), B.Interface())
}

func NotMatchError(typeOf string, A, B interface{}) error {
	return errors.New(fmt.Sprintf("%s not matched, result: %v, expected: %v", typeOf, A, B))
}

func compareInterface(A, B interface{}) error {
	if A == nil && B == nil {
		return nil
	}
	if A == nil || B == nil {
		return NotMatchError("Nil", A, B)
	}

	typeA := getType(A)
	typeB := getType(B)

	if typeA != typeB {
		return NotMatchError("Type", typeA, typeB)
	}

	switch typeA {
	case "string":
		if A.(string) != B.(string) {
			return NotMatchError(typeA, A, B)
		}
		return nil

	case "number":
		if A.(json.Number).String() != B.(json.Number).String() {
			return NotMatchError(typeA, A, B)
		}
		return nil

	case "map[string]interface{}":
		vA := A.(map[string]interface{})
		vB := B.(map[string]interface{})
		return compareMap(vA, vB)

	case "[]interface{}":
		vA := A.([]interface{})
		vB := B.([]interface{})
		return compareArray(vA, vB)
	}

	return NotMatchError("Unknown Type", typeA, typeB)
}

func getType(data interface{}) string {
	switch data.(type) {
	case map[string]interface{}:
		return "map[string]interface{}"
	case []interface{}:
		return "[]interface{}"
	case json.Number:
		return "number"
	}

	return reflect.TypeOf(data).Name()
}

func compareMap(A, B map[string]interface{}) error {
	if len(A) != len(B) {
		return NotMatchError("Map Length", len(A), len(B))
	}

	var err error
	for k, v := range A {
		if err = compareInterface(v, B[k]); err != nil {
			return err
		}
	}
	return nil
}

func compareArray(A, B []interface{}) error {
	if len(A) != len(B) {
		return NotMatchError("Array Length", len(A), len(B))
	}

	var err error
	for k, v := range A {
		if err = compareInterface(v, B[k]); err != nil {
			return err
		}
	}
	return nil
}
