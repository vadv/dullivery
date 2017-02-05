package dsl

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/yuin/gopher-lua"
)

func (d *Config) dslHttp(L *lua.LState) int {
	verb := strings.ToLower(L.CheckString(1))
	url := L.CheckString(2)
	switch verb {
	case "get":
		response, err := http.Get(url)
		if err != nil {
			L.RaiseError("http error: %s\n", err.Error())
			return 0
		}
		defer response.Body.Close()
		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			L.RaiseError("http read response error: %s\n", err.Error())
			return 0
		}
		// write response
		result := L.NewTable()
		L.SetField(result, "code", lua.LNumber(response.StatusCode))
		L.SetField(result, "body", lua.LString(string(data)))
		L.Push(result)
		return 1
	default:
		L.RaiseError("unsupported http verb: %s", verb)
		return 0
	}
}
