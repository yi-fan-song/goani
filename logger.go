package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

type Logger struct {
	IncludeTime bool
	UseColor    bool

	Out io.Writer
}

const (
	RED    = "31m"
	GREEN  = "32m"
	YELLOW = "33m"
	BLUE   = "34m"
	PURPLE = "35m"
	CYAN   = "36m"
)

func DefaultLogger() *Logger {
	return &Logger{
		IncludeTime: true,
		UseColor:    true,
		Out:         os.Stdout,
	}
}

func (l Logger) write(s string, color string) {
	if l.IncludeTime {
		s = "[" + time.Now().Format(time.Stamp) + "] " + s
	}
	if l.UseColor {
		s = "\033[0;" + color + s + "\033[0m"
	}
	fmt.Fprintln(l.Out, s)
}

func (l Logger) Log(s string) {
	l.write(s, CYAN)
}

func (l Logger) Info(s string) {
	l.write(s, GREEN)
}

func (l Logger) Warn(s string) {
	l.write(s, YELLOW)
}

func (l Logger) Error(s string) {
	l.write(s, RED)
}

func (l Logger) Fatal(s string) {
	l.write(s, RED)
}

type CombinedWriter struct {
	Writer1, Writer2 io.Writer
}

func (cw CombinedWriter) Write(p []byte) (n int, err error) {
	n1, err1 := cw.Writer1.Write(p)
	n2, err2 := cw.Writer2.Write(p)

	if err1 != nil && err2 != nil {
		return 0, fmt.Errorf("%v \n%v", err1, err2)
	}
	if err1 != nil {
		return n1, err1
	}
	if err2 != nil {
		return n2, err2
	}
	if n1 != n2 {
		return 0, fmt.Errorf("amount written to both writers don't match")
	}

	return n1, nil
}
