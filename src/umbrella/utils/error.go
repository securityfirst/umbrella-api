package utils

import (
	"fmt"
	"log"
	"runtime"

	"github.com/fatih/color"
)

func CheckErr(err error) {
	if err != nil {
		info := color.New(color.FgGreen).SprintFunc()
		pc, fn, line, _ := runtime.Caller(1)
		log.Printf(info(fmt.Sprintf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), fn, line, err)))
	}
}

func redLog(toLog interface{}) {
	info := color.New(color.FgRed).SprintFunc()
	pc, fn, line, _ := runtime.Caller(1)
	log.Printf(info(fmt.Sprintf("%s[%s:%d] %v", runtime.FuncForPC(pc).Name(), fn, line, fmt.Sprint(toLog))))
}
