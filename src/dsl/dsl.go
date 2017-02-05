package dsl

import (
	"os"

	"github.com/yuin/gopher-lua"
)

type Config struct {
	LogFd      *os.File
	StorageDir string
	dbpath     string
}

func Register(L *lua.LState, config *Config) {

	newHelpers := func(L *lua.LState) int {
		ud := L.NewUserData()
		ud.Value = config
		L.SetMetatable(ud, L.GetTypeMetatable("helpers"))
		L.Push(ud)
		return 1
	}

	newStorage := func(L *lua.LState) int {
		ud := L.NewUserData()
		ud.Value = config
		L.SetMetatable(ud, L.GetTypeMetatable("storage"))
		L.Push(ud)
		return 1
	}

	helpers := L.NewTypeMetatable("helpers")
	L.SetGlobal("helpers", helpers)
	L.SetField(helpers, "new", L.NewFunction(newHelpers))
	L.SetField(helpers, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"find":   config.dslFind,
		"copy":   config.dslCopy,
		"remove": config.dslRemove,
		"alive":  config.dslAlive,
		"log":    config.dslLog,
		"sleep":  config.dslSleep,
		"http":   config.dslHttp,
	}))

	storage := L.NewTypeMetatable("storage")
	L.SetGlobal("storage", storage)
	L.SetField(storage, "new", L.NewFunction(newStorage))
	L.SetField(storage, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"get":    config.dslStorageGet,
		"set":    config.dslStorageSet,
		"expire": config.dslStorageExpire,
	}))
}
