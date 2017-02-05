package store

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type Scenario struct {
	Id      int64              `json:"id"`
	Name    string             `json:"name"`
	Content string             `json:"content"`
	Active  bool               `json:"active"`
	Run     bool               `json:"-"`
	History []*ScenarioHistory `json:"history"`
}

type ScenarioHistory struct {
	State     State    `json:"state"`
	StartedAt unixTime `json:"started_at"`
	EndedAt   unixTime `json:"ended_at"`
}

type Scenarios struct {
	List []*Scenario `json:"list"`
}

type FullScenarioHistory struct {
	Scenario *Scenario
	History  []*ScenarioHistory
	Limit    int
}

func (s *Scenario) Len() int {
	return len(s.History)
}

func (s *Scenario) Swap(i, j int) {
	s.History[i], s.History[j] = s.History[j], s.History[i]
}

func (s *Scenario) Less(i, j int) bool {
	return s.History[i].StartedAt > s.History[j].StartedAt
}

func (s *Scenarios) maxId() int64 {
	i := int64(0)
	for _, t := range s.List {
		if i < t.Id {
			i = t.Id
		}
	}
	return i + 1
}

func (s *Scenario) ActiveHuman() string {
	if s.Active {
		return "On"
	}
	return "Off"
}

func (s *Storage) AddScenario(scenario Scenario) {
	history := make([]*ScenarioHistory, 0)
	id := scenario.Id
	if id == 0 {
		id = s.Scenarios.maxId()
	}
	if prev, found := s.GetScenario(scenario.Id); found {
		log.Printf("[INFO] save backup of scenario (%d): %s\n", scenario.Id, scenario.Name)
		dir := filepath.Join(s.Dir, "backup", "scenario")
		os.MkdirAll(dir, 0750)
		ioutil.WriteFile(
			filepath.Join(dir, fmt.Sprintf("%s-%d.lua", prev.Name, time.Now().Unix())),
			[]byte(prev.Content), 0640)
		history = prev.History
		s.DelScenario(scenario.Id)
	}
	scenario.History = history
	scenario.Id = id
	s.Scenarios.List = append(s.Scenarios.List, &scenario)
	s.save()
}

func (s *Storage) DelScenario(id int64) {
	result := &Scenarios{List: make([]*Scenario, 0)}
	for _, scenario := range s.Scenarios.List {
		if scenario.Id != id {
			result.List = append(result.List, scenario)
		}
	}
	s.Scenarios = result
	s.save()
}

func (s *Storage) GetScenario(id int64) (Scenario, bool) {
	for _, scenario := range s.Scenarios.List {
		if scenario.Id == id {
			return Scenario{
				Id:      id,
				Name:    scenario.Name,
				Content: scenario.Content,
				History: scenario.History,
				Active:  scenario.Active,
				Run:     scenario.Run,
			}, true
		}
	}
	return Scenario{}, false
}

func (s *Scenario) LastHistory() *ScenarioHistory {
	sort.Sort(s)
	hs := s.History
	if len(hs) > 0 {
		h := hs[0]
		return h
	}
	return nil
}

func (s *Scenario) HistoryList(limit int) *FullScenarioHistory {
	sort.Sort(s)
	hs := s.History
	result := make([]*ScenarioHistory, 0)
	for i := 0; i < limit; i++ {
		if len(hs) > i {
			h := hs[i]
			result = append(result, h)
		}
	}
	return &FullScenarioHistory{History: result, Scenario: s, Limit: limit}
}

func (s *Storage) GetScenarioHistoryLogFile(id, log_id int64) string {
	dir := filepath.Join(s.Dir, "logs", "scenarios")
	os.MkdirAll(dir, 0750)
	filename := filepath.Join(dir, fmt.Sprintf("%d", id), fmt.Sprintf("%d.log", log_id))
	return filename
}

// поджимаем json, убираем инфу из админки, но логи остаются.
func (s *Storage) compactScenario() {
	for _, scenario := range s.Scenarios.List {
		if len(scenario.History) > 1000 {
			fh := scenario.HistoryList(500)
			scenario.History = fh.History
		}
	}
}
