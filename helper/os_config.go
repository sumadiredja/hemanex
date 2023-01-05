package helper

import (
	"fmt"
	"runtime"
)

func IsWindows() bool {
	fmt.Println(runtime.GOOS)
	return runtime.GOOS == "windows"
}
