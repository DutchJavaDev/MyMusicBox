package logging

import (
	"fmt"
)

func Info(a any) {
	log(fmt.Sprintf("\033[32m[API Info]\033[0m %s", a))
}

func Warning(a any) {
	log(fmt.Sprintf("\033[33m[API Warning]\033[0m %s", a))
}

func Error(a any) {
	log(fmt.Sprintf("\033[31m[API Error]\033[0m %s", a))
}

func ErrorStackTrace(e error) {
	fmt.Printf("e: %v\n", e)
}

func log(text string) {
	fmt.Println(text)
}
