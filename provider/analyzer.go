package provider

import "reflect"

type ErrNotAFunction struct{}

func (e ErrNotAFunction) Error() string {
	return "not a function"
}

func analyzeFunction(constructFunction any) ([]reflect.Type, []reflect.Type, error) {
	if reflect.TypeOf(constructFunction).Kind() != reflect.Func {
		return nil, nil, ErrNotAFunction{}
	}

	constructor := reflect.ValueOf(constructFunction)
	var args []reflect.Type

	for i := 0; i < constructor.Type().NumIn(); i++ {
		args = append(args, constructor.Type().In(i))
	}

	var returns []reflect.Type

	for i := 0; i < constructor.Type().NumOut(); i++ {
		returns = append(returns, constructor.Type().Out(i))
	}

	return args, returns, nil
}

func analyzeMethod(constructMethod any) ([]reflect.Type, []reflect.Type, error) {
	return analyzeFunction(constructMethod)
}
