package dsl

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/yuin/gopher-lua"
)

// в дальнейщем можно сделать лок именнованым: map[string]Mutex
var globalStorageLock = &sync.Mutex{}

type dslStorageValue struct {
	Value string `json:"value"`
	SetAt int64  `json:"set_at"`
	Ttl   int64  `json:"ttl"`
}

type dslStorage struct {
	Data     map[string]*dslStorageValue `json:"data"`
	Filename string                      `json:"-"`
}

func (c *Config) Storage(name string) (*dslStorage, error) {
	return loadDslStorage(filepath.Join(c.StorageDir, name+".json"))
}

func loadDslStorage(filename string) (*dslStorage, error) {
	d := &dslStorage{Data: make(map[string]*dslStorageValue, 0), Filename: filename}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			data = []byte("{}")
		} else {
			return nil, err
		}
	}
	if err := json.Unmarshal(data, d); err != nil {
		return nil, err
	}
	return d, d.deleteExpired()
}

func (d *dslStorage) save() error {
	data, err := json.Marshal(d)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(d.Filename, data, 0640)
}

func (val *dslStorageValue) isExpired() bool {
	if val.Ttl == 0 {
		return false
	}
	return val.SetAt > time.Now().Unix()+val.Ttl
}

func (d *dslStorage) deleteExpired() error {
	for key, val := range d.Data {
		if val.isExpired() {
			delete(d.Data, key)
		}
	}
	return d.save()
}

func (d *dslStorage) get(key string) (*dslStorageValue, bool) {
	globalStorageLock.Lock()
	defer globalStorageLock.Unlock()
	val, found := d.Data[key]
	if !found {
		return nil, false
	}
	if val.isExpired() {
		d.save()
		return nil, false
	}
	return val, found
}

func (d *dslStorage) set(key, val string) error {
	globalStorageLock.Lock()
	defer globalStorageLock.Unlock()
	d.Data[key] = &dslStorageValue{Value: val, SetAt: time.Now().Unix()}
	return d.save()
}

func (d *dslStorage) setTtl(key, val string, ttl int64) error {
	globalStorageLock.Lock()
	defer globalStorageLock.Unlock()
	d.Data[key] = &dslStorageValue{Value: val, SetAt: time.Now().Unix(), Ttl: ttl}
	return d.save()
}

func (c *Config) dslStorageGet(L *lua.LState) int {
	key := L.CheckString(2)
	storage, err := c.Storage(L.CheckString(1))
	if err != nil {
		log.Printf("[ERROR] load storage: %s\n", err.Error())
		L.RaiseError("load storage error: %s", err.Error())
		return -1
	}
	val, found := storage.get(key)
	if !found {
		L.Push(lua.LNil)
		return 1
	}
	L.Push(lua.LString(val.Value))
	return 1
}

func (c *Config) dslStorageSet(L *lua.LState) int {
	key := L.CheckString(2)
	val := L.CheckString(3)
	storage, err := c.Storage(L.CheckString(1))
	if err != nil {
		log.Printf("[ERROR] load storage: %s\n", err.Error())
		L.RaiseError("load storage error: %s", err.Error())
		return -1
	}
	if err := storage.set(key, val); err != nil {
		log.Printf("[ERROR] save storage: %s\n", err.Error())
		L.RaiseError("save storage error: %s", err.Error())
		return -1
	}
	return 0
}

func (c *Config) dslStorageExpire(L *lua.LState) int {
	key := L.CheckString(2)
	ttl := L.CheckNumber(3)
	storage, err := c.Storage(L.CheckString(1))
	if err != nil {
		log.Printf("[ERROR] load storage: %s\n", err.Error())
		L.RaiseError("load storage error: %s", err.Error())
		return -1
	}
	val, found := storage.get(key)
	if !found {
		return 0
	}
	if err := storage.setTtl(key, val.Value, int64(ttl)); err != nil {
		log.Printf("[ERROR] save storage: %s\n", err.Error())
		L.RaiseError("save storage error: %s", err.Error())
		return -1
	}
	return 0
}
