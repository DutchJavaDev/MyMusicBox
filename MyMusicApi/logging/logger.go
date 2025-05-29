package logging

import "fmt"

var OutputLog bool

func Info(a any) {
	log(fmt.Sprintf("\033[32m[API Info]\033[0m %s", a))
}
func Error(a any) {
	log(fmt.Sprintf("\033[31m[API Error]\033[0m %s", a))
}
func Warning(a any) {
	log(fmt.Sprintf("\033[33m[API Warning]\033[0m %s", a))
}

func log(text string) {
	if OutputLog {
		fmt.Println(text)
	}
}
