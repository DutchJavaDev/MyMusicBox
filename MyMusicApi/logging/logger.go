package logging

import "fmt"

var OutputLog bool

func Log(a any) {
	if OutputLog {
		fmt.Println(a)
	}
}
