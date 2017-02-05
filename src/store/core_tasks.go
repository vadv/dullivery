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

func (s *Storage) startTasksRoutine() {
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case <-ticker:
			for _, task := range s.Tasks.List {
				if task.StartAt < UnixNow() && task.State == StateAdded {
					log.Printf("[INFO] start task %d: it want to be started at: %d, now: %d.\n", task.Id, task.StartAt, UnixNow())
					// запускаем бэкграундом
					go func(task *Task) {
						if err := s.runTask(task); err != nil {
							task.State = StateFailed
							log.Printf("[ERROR] task %d failed: %s.\n", task.Id, err.Error())
						} else {
							task.State = StateCompleted
							log.Printf("[INFO] task %d: successfully ended.\n", task.Id)
						}
					}(task)
				}
			}
		}
	}
}

func (storage *Storage) runTask(t *Task) error {

	t.State = StateStarted

	// prepare log
	filename := storage.GetTaskLog(t.Id)
	fd, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640)
	if err != nil {
		log.Printf("[ERROR] Open task %d log: %s\n", t.Id, err.Error())
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
		fmt.Fprintf(fd, "[ERROR] Can't load config: %s\n", err.Error())
		return err
	}

	// run task
	if err := state.DoString(t.Content); err != nil {
		fmt.Fprintf(fd, "[ERROR] While exec task: %s\n", err.Error())
		return err
	}

	return nil
}
