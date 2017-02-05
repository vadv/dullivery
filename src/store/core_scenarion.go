package store

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/yuin/gopher-lua"

	"dsl"
)

func (s *Storage) startScenariosRoutine() {
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case <-ticker:
			for _, scenario := range s.Scenarios.List {
				if scenario.Active && !scenario.Run {
					log.Printf("[INFO] start scenario %d.\n", scenario.Id)
					// запускаем бэкграундом
					go func(scenario *Scenario) {
						if err := s.runScenario(scenario); err != nil {
							log.Printf("[ERROR] scenario %d failed: %s.\n", scenario.Id, err.Error())
						} else {
							log.Printf("[INFO] scenario %d: successfully ended.\n", scenario.Id)
						}
					}(scenario)
				}
			}
		}
	}
}

func (storage *Storage) runScenario(s *Scenario) error {

	s.Run = true
	defer func() {
		s.Run = false
	}()

	// create history
	if len(s.History) == 0 {
		s.History = make([]*ScenarioHistory, 0)
	}
	now := UnixNow()
	history := &ScenarioHistory{State: StateStarted, StartedAt: now}
	s.History = append(s.History, history)

	// prepare log
	logDir := filepath.Join(storage.Dir, "logs", "scenarios", fmt.Sprintf("%d", s.Id))
	os.MkdirAll(logDir, 0750)
	filename := filepath.Join(logDir, fmt.Sprintf("%d.log", history.StartedAt))
	fd, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640)
	if err != nil {
		log.Printf("[ERROR] Open scenario %d log: %s\n", s.Id, err.Error())
		return err
	}
	defer fd.Close()

	// prepare storage
	storageDir := filepath.Join(storage.Dir, "storage")
	os.MkdirAll(storageDir, 0750)

	// start dsl
	state := lua.NewState()
	defer state.Close()
	dsl.Register(state, &dsl.Config{LogFd: fd, StorageDir: storageDir})

	// load config
	if err := state.DoString(storage.Config); err != nil {
		history.State = StateFailed
		fmt.Fprintf(fd, "[ERROR] Can't load config: %s\n", err.Error())
		return err
	}

	// run scenario
	if err := state.DoString(s.Content); err != nil {
		history.State = StateFailed
		fmt.Fprintf(fd, "[ERROR] While exec scenario: %s\n", err.Error())
		return err
	}

	history.EndedAt = UnixNow()
	history.State = StateCompleted
	return nil
}
