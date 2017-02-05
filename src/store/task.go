package store

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
)

type Task struct {
	Id      int64    `json:"id"`
	Content string   `json:"content"`
	State   State    `json:"state"`
	StartAt unixTime `json:"start_at"`
}

type Tasks struct {
	List []*Task `json:"list"`
}

func (t *Tasks) Len() int {
	return len(t.List)
}

func (t *Tasks) Swap(i, j int) {
	t.List[i], t.List[j] = t.List[j], t.List[i]
}

func (t *Tasks) Less(i, j int) bool {
	return t.List[i].Id > t.List[j].Id
}

func (t *Tasks) Limit(count int) []*Task {
	result := make([]*Task, 0)
	sort.Sort(t)
	for i := 0; i < count; i++ {
		if t.Len() > i {
			result = append(result, t.List[i])
		}
	}
	return result
}

func (t *Tasks) maxId() int64 {
	i := int64(0)
	for _, task := range t.List {
		if i < task.Id {
			i = task.Id
		}
	}
	return i + 1
}

func (s *Storage) AddTask(task *Task) {
	if task.Id > 0 {
		newList := make([]*Task, 0)
		for _, t := range s.Tasks.List {
			if t.Id != task.Id {
				newList = append(newList, t)
			}
		}
		s.Tasks.List = newList
	} else {
		task.Id = s.Tasks.maxId()
	}
	task.State = StateAdded
	log.Printf("[INFO] task: %v saved\n", task)
	s.Tasks.List = append(s.Tasks.List, task)
	s.save()
}

func (s *Storage) GetTask(id int64) *Task {
	for _, task := range s.Tasks.List {
		if task.Id == id {
			return task
		}
	}
	return nil
}

func (s *Storage) GetTaskLog(id int64) string {
	dir := filepath.Join(s.Dir, "logs", "tasks")
	os.MkdirAll(dir, 0750)
	return filepath.Join(dir, fmt.Sprintf("%d.log", id))
}

func (t *Task) IsExecuted() bool {
	return t.State != StateAdded
}
