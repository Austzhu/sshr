package public

import (
	"fmt"
	"log"
)

func Recover(err *error) func() {
	return func() {
		if rec := recover(); rec != nil {
			*err = fmt.Errorf("%v", rec)
		}
	}
}

func Die(msg string, err error) {
	if err != nil {
		panic(fmt.Sprintf("%s fail: %v\n", msg, err))
	}
}

func Warn(msg string, err error) {
	if err != nil {
		log.Printf("%s fail: %v\n", msg, err)
	}
}
