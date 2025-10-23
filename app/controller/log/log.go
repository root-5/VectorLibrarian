package log

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"
)

func Info(msg string) {
	fmt.Printf("\x1b[48;5;520m%s\x1b[0m	%s\n", time.Now().Format("2006-01-02 15:04:05"), msg)
}

func Error(err error) {
	pc, file, line, ok := runtime.Caller(1)
	if ok {
		fmt.Printf("\x1b[48;5;520m%s\x1b[0m	\x1b[48;5;160m ERROR \x1b[0m\n			%s:%d  %s\n			%v\n",
			time.Now().Format("2006-01-02 15:04:05"),
			filepath.Base(file),
			line,
			runtime.FuncForPC(pc).Name(),
			err,
		)
	}
}