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
		return 0
	case reflect.TypeOf(false):
		return false
	case reflect.TypeOf(""):
		return ""
	}

	return nil
}

func AddAny(val1, val2 any) (any, error) {
	if reflect.TypeOf(val1) != reflect.TypeOf(val2) {
		return nil, fmt.Errorf(
			"invalid operation %v(type:%v) + %v(type:%v)",
			val1, reflect.TypeOf(val1),
			val2, reflect.TypeOf(val2),
		)
	}

	switch val1.(type) {
	case int:
		return val1.(int) + val2.(int), nil
	case string:
		return val1.(string) + val2.(string), nil
	default:
		return nil, fmt.Errorf(
			"invalid operation %v(type:%v) + %v(type:%v)",
			val1, reflect.TypeOf(val1),
			val2, reflect.TypeOf(val2),
		)
	}
}

func MulAny(val1, val2 any) (any, error) {
	if reflect.TypeOf(val1) != reflect.TypeOf(val2) {
		return nil, fmt.Errorf(
			"invalid operation %v(type:%v) * %v(type:%v)",
			val1, reflect.TypeOf(val1),
			val2, reflect.TypeOf(val2),
		)
	}

	switch val1.(type) {
	case int:
		return val1.(int) * val2.(int), nil
	default:
		return nil, fmt.Errorf(
			"invalid operation %v(type:%v) * %v(type:%v)",
			val1, reflect.TypeOf(val1),
			val2, reflect.TypeOf(val2),
		)
	}
}

func DivAny(val1, val2 any) (any, error) {
	if reflect.TypeOf(val1) != reflect.TypeOf(val2) {
		return nil, fmt.Errorf(
			"invalid operation %v(type:%v) / %v(type:%v)",
			val1, reflect.TypeOf(val1),
			val2, reflect.TypeOf(val2),
		)
	}

	switch val1.(type) {
	case int:
		return val1.(int) / val2.(int), nil
	default:
		return nil, fmt.Errorf(
			"invalid operation %v(type:%v) / %v(type:%v)",
			val1, reflect.TypeOf(val1),
			val2, reflect.TypeOf(val2),
		)
	}
}

func SubAny(val1, val2 any) (any, error) {
	if reflect.TypeOf(val1) != reflect.TypeOf(val2) {
		return nil, fmt.Errorf(
			"invalid operation %v(type:%v) - %v(type:%v)",
			val1, reflect.TypeOf(val1),
			val2, reflect.TypeOf(val2),
		)
	}

	switch val1.(type) {
	case int:
		return val1.(int) - val2.(int), nil
	default:
		return nil, fmt.Errorf(
			"invalid operation !%v(type:%v) - %v(type:%v)",
			val1, reflect.TypeOf(val1),
			val2, reflect.TypeOf(val2),
		)
	}
}

func NotAny(val1 any) (any, error) {
	boolVal, ok := val1.(bool)
	if ok {
		return !boolVal, nil
	} else {
		return nil, fmt.Errorf(
			"invalid operation !%v(type:%v)",
			val1, reflect.TypeOf(val1),
		)
	}
}

func OrAny(val1, val2 any) (any, error) {
	if reflect.TypeOf(val1) != reflect.TypeOf(val2) {
		return nil, fmt.Errorf(
			"invalid operation %v(type:%v) || %v(type:%v)",
			val1, reflect.TypeOf(val1),
			val2, reflect.TypeOf(val2),
		)
	}

	switch val1.(type) {
	case bool:
		return val1.(bool) || val2.(bool), nil
	default:
		return nil, fmt.Errorf(
			"invalid operation !%v(type:%v) || %v(type:%v)",
			val1, reflect.TypeOf(val1),
			val2, reflect.TypeOf(val2),
		)
	}
}

func AndAny(val1, val2 any) (any, error) {
	if reflect.TypeOf(val1) != reflect.TypeOf(val2) {
		return nil, fmt.Errorf(
			"invalid operation %v(type:%v) || %v(type:%v)",
			val1, reflect.TypeOf(val1),
			val2, reflect.TypeOf(val2),
		)
	}

	switch val1.(type) {
	case bool:
		return val1.(bool) && val2.(bool), nil
	default:
		return nil, fmt.Errorf(
			"invalid operation !%v(type:%v) && %v(type:%v)",
			val1, reflect.TypeOf(val1),
			val2, reflect.TypeOf(val2),
		)
	}
}

func CompareAny(val1, val2 any, compareType string) (any, error) {
	switch compareType {
	case "==":
		return EqualAny(val1, val2)
	case "!=":
		res, err := EqualAny(val1, val2)
		if err != nil {
			return nil, err
		}
		return !res.(bool), err
	case "<":
		return LessAny(val1, val2)
	case "<=":
		less, err := LessAny(val1, val2)
		if err != nil {
			return nil, err
		}
		eq, err := EqualAny(val1, val2)
		if err != nil {
			return nil, err
		}

		return OrAny(less, eq)
	case ">":
		less, err := LessAny(val1, val2)
		if err != nil {
			return nil, err
		}
		eq, err := EqualAny(val1, val2)
		if err != nil {
			return nil, err
		}

		res, _ := OrAny(less, eq)
		return !res.(bool), nil
	case ">=":
		less, err := LessAny(val1, val2)
		if err != nil {
			return nil, err
		}
		
		return !less.(bool), nil
	default:
		panic("unknown compare type")
	}
}

func EqualAny(val1, val2 any) (any, error) {
	if reflect.TypeOf(val1) != reflect.TypeOf(val2) {
		return nil, fmt.Errorf(
			"invalid operation %v(type:%v) compare %v(type:%v)",
			val1, reflect.TypeOf(val1),
			val2, reflect.TypeOf(val2),
		)
	}

	switch val1.(type) {
	case bool:
		return val1.(bool) == val2.(bool), nil
	case int:
		return val1.(int) == val2.(int), nil
	case string:
		return val1.(string) == val2.(string), nil
	default:
		return nil, fmt.Errorf(
			"invalid operation !%v(type:%v) compare %v(type:%v)",
			val1, reflect.TypeOf(val1),
			val2, reflect.TypeOf(val2),
		)
	}
}

func LessAny(val1, val2 any) (any, error) {
	if reflect.TypeOf(val1) != reflect.TypeOf(val2) {
		return nil, fmt.Errorf(
			"invalid operation %v(type:%v) compare %v(type:%v)",
			val1, reflect.TypeOf(val1),
			val2, reflect.TypeOf(val2),
		)
	}

	switch val1.(type) {
	case int:
		return val1.(int) < val2.(int), nil
	default:
		return nil, fmt.Errorf(
			"invalid operation !%v(type:%v) compare %v(type:%v)",
			val1, reflect.TypeOf(val1),
			val2, reflect.TypeOf(val2),
		)
	}
}
