package main

import (
	"fmt"
	"net/http"
	"runtime"
)

type syntaxErr struct {
	s string
}

func (e *syntaxErr) Error() string {
	return e.s
}

func prependInfo(err error, file string, line int, ok bool) error {
	if ok {
		return fmt.Errorf("%s:%d => %s", file, line, err)
	} else {
		return fmt.Errorf("%s (failed getting runtime info)", err)
	}
}

func funcName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}

// func dumpStack() {
//   buf := make([]byte, 1<<16)
//   runtime.Stack(buf, true)
//   fmt.Printf("%s", buf)
// }

// func newErrWithInfo(s string) error {
//   _, file, line, ok := runtime.Caller(1)
//   return prependInfo(errors.New(s), file, line, ok)
// }

// func errWithInfo(err error) error {
//   _, file, line, ok := runtime.Caller(1)
//   return prependInfo(err, file, line, ok)
// }

func httpGenericErr(w http.ResponseWriter) {
	http.Error(w, "Something went wrong, check the logs", http.StatusInternalServerError)
}
