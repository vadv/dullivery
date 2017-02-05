package dsl

import (
	"fmt"
	"time"

	"github.com/yuin/gopher-lua"
)

func (d *Config) dslLog(L *lua.LState) int {
	msg := L.CheckString(1)
	logTime := time.Now().Format("02/01/2006 15:04:05")
	fmt.Fprintf(d.LogFd, "[%s] %s\n", logTime, msg)
	return 0
}
