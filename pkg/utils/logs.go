package utils

import "runtime"

func CurrentFunction() string {
	counter, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(counter).Name()
}
