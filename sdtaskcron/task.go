package sdtaskcron

import (
	"github.com/gaorx/stardust3/sderr"
	"github.com/gaorx/stardust3/sdjson"
	"github.com/robfig/cron/v3"
)

type Task struct {
	Id     string        `json:"id"`
	At     string        `json:"at"`
	Action string        `json:"action"`
	Args   sdjson.Object `json:"args"`
	sched  cron.Schedule `json:"-"`
}

type TaskPred func(*Task) bool

func (t *Task) check() error {
	if t.Action == "" {
		return sderr.New("no action")
	}
	if t.At == "" {
		return sderr.New("no 'at'")
	}
	return nil
}

func (t *Task) ensure() error {
	if t.Id == "" {
		t.Id = t.Action
	}
	if t.sched != nil {
		return nil
	}
	sched, err := parseSched(t.At)
	if err != nil {
		return sderr.WithStack(err)
	}
	t.sched = sched
	return nil
}

func (t *Task) filter(pred TaskPred) bool {
	if t == nil {
		return false
	}
	if pred == nil {
		return true
	}
	return pred(t)
}

func parseTasks(data []byte) ([]Task, error) {
	var r []Task
	err := sdjson.Unmarshal(data, &r)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return r, nil
}

func filterTasks(tasks []Task, pred TaskPred) []Task {
	r := make([]Task, 0)
	for _, t := range tasks {
		if t.filter(pred) {
			r = append(r, t)
		}
	}
	return r
}

type holder interface {
	TaskMap() map[string]*Task
}

type staticHolder struct {
	m map[string]*Task
}

func (h staticHolder) TaskMap() map[string]*Task {
	return h.m
}
