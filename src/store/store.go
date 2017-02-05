package store

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

var Box *Storage

type Storage struct {
	Dir       string     `json:"-"`
	Users     *Users     `json:"users"`
	Scenarios *Scenarios `json:"scenarios"`
	Tasks     *Tasks     `json:"tasks"`
	Config    string     `json:"config.lua"`
}

func NewStorage(dir string) (*Storage, error) {

	if err := os.MkdirAll(dir, 0750); err != nil {
		return nil, err
	}

	result := &Storage{
		Dir:       dir,
		Users:     &Users{List: make([]*User, 0)},
		Scenarios: &Scenarios{List: make([]*Scenario, 0)},
		Tasks:     &Tasks{List: make([]*Task, 0)},
	}

	needBackup := false
	if data, err := ioutil.ReadFile(result.scenariosFilename()); err == nil {
		o := result.Scenarios.List
		if err := json.Unmarshal(data, &o); err == nil {
			result.Scenarios.List = o
		} else {
			log.Printf("[ERROR] unmarshal: %s\n", err.Error())
			needBackup = true
		}
	}
	if data, err := ioutil.ReadFile(result.tasksFilename()); err == nil {
		o := result.Tasks
		if err := json.Unmarshal(data, o); err == nil {
			result.Tasks = o
		} else {
			log.Printf("[ERROR] unmarshal: %s\n", err.Error())
			needBackup = true
		}
	}
	if data, err := ioutil.ReadFile(result.usersFilename()); err == nil {
		o := result.Users
		if err := json.Unmarshal(data, o); err == nil {
			result.Users = o
		} else {
			log.Printf("[ERROR] unmarshal: %s\n", err.Error())
			needBackup = true
		}
	}
	if data, err := ioutil.ReadFile(result.configFilename()); err == nil {
		result.Config = string(data)
	}

	if needBackup {
		result.backup()
	}

	go result.saveRoutine()
	go result.startScenariosRoutine()
	go result.startTasksRoutine()
	return result, result.save()
}

func (s *Storage) SaveConfig(lua string) {
	dir := filepath.Join(s.Dir, "backup", "config")
	os.MkdirAll(dir, 0750)
	ioutil.WriteFile(filepath.Join(dir, fmt.Sprintf("%d.lua", time.Now().Unix())), []byte(s.Config), 0640)
	s.Config = lua
	s.save()
}

func (s *Storage) saveRoutine() {
	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-ticker.C:
			s.compactScenario()
			if err := s.save(); err != nil {
				log.Printf("[ERROR] store: %s\n", err.Error())
			} else {
				log.Printf("[INFO] saved")
			}
		}
	}
}

func (s *Storage) configFilename() string {
	return filepath.Join(s.Dir, "config.lua")
}

func (s *Storage) usersFilename() string {
	return filepath.Join(s.Dir, "users.json")
}

func (s *Storage) scenariosFilename() string {
	return filepath.Join(s.Dir, "scenarios.json")
}

func (s *Storage) tasksFilename() string {
	return filepath.Join(s.Dir, "tasks.json")
}

func (s *Storage) backup() {

	dir := filepath.Join(s.Dir, "backup")
	t := fmt.Sprintf("%d", time.Now().Unix())

	if data, err := ioutil.ReadFile(s.scenariosFilename()); err == nil {
		ioutil.WriteFile(filepath.Join(dir, "scenarios-"+t+".json"), data, 0640)
	}
	if data, err := ioutil.ReadFile(s.usersFilename()); err == nil {
		ioutil.WriteFile(filepath.Join(dir, "users-"+t+".json"), data, 0640)
	}
	if data, err := ioutil.ReadFile(s.tasksFilename()); err == nil {
		ioutil.WriteFile(filepath.Join(dir, "tasks-"+t+".json"), data, 0640)
	}
	if data, err := ioutil.ReadFile(s.configFilename()); err == nil {
		ioutil.WriteFile(filepath.Join(dir, "config-"+t+".lua"), data, 0640)
	}

}

func (s *Storage) save() error {
	// scenarios
	o1 := s.Scenarios.List
	if data, err := json.Marshal(o1); err != nil {
		return err
	} else {
		if err := ioutil.WriteFile(s.scenariosFilename(), data, 0640); err != nil {
			return err
		}
	}
	// users
	o2 := s.Users
	if data, err := json.Marshal(o2); err != nil {
		return err
	} else {
		if err := ioutil.WriteFile(s.usersFilename(), data, 0640); err != nil {
			return err
		}
	}
	// tasks
	o3 := s.Tasks
	if data, err := json.Marshal(o3); err != nil {
		return err
	} else {
		if err := ioutil.WriteFile(s.tasksFilename(), data, 0640); err != nil {
			return err
		}
	}
	// config
	if err := ioutil.WriteFile(s.configFilename(), []byte(s.Config), 0640); err != nil {
		return nil
	}
	return nil
}
