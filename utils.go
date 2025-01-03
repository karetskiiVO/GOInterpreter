package main

import (
	"fmt"
	"reflect"
	"strings"
)

func CloneAny(val any) any {
	switch val.(type) {
	case int:
		return val.(int)
	case string:
		return strings.Clone(val.(string))
	case bool:
		return val.(bool)
	default:
		panic(fmt.Sprintf("unknown type: %v", reflect.TypeOf(val).String()))
	}
}

func ReflectType(typename string) (reflect.Type, error) {
	switch typename {
	case "int":
		return reflect.TypeOf(1), nil
	case "bool":
		return reflect.TypeOf(true), nil
	case "string":
		return reflect.TypeOf("string"), nil
	default:
		return nil, fmt.Errorf("unknown type %v", typename)
	}
}

func NewVariable(Type reflect.Type) any {
	switch Type {
	case reflect.TypeOf(0):
		return 0;
	case reflect.TypeOf(false):
		return false;
	case reflect.TypeOf(""):
		return ""
	}

	return nil
}