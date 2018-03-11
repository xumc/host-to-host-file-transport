package util

import (
	"fmt"
	"os"
)

func CheckErr(err error) {
	if err != nil {
		Log(err)
		os.Exit(1)
	}
}

func Log(v ...interface{}) {
	fmt.Println(v...)
}