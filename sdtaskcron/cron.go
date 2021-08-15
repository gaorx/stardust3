package sdtaskcron

import (
	"sync"
	"time"

	"github.com/gaorx/stardust3/sdcall"
	"github.com/gaorx/stardust3/sderr"
	"github.com/gaorx/stardust3/sdload"
	"github.com/gaorx/stardust3/sdlog"
	"github.com/gaorx/stardust3/sdtime"
)

type Cron struct {
	holder  holder
	actions map[string]Action
	mtx     sync.Mutex
	nexts   map[string]time.Time
	stopWg  sync.WaitGroup
}

func NewCron() *Cron {
	return &Cron{
		holder:  nil,
		actions: map[string]Action{},
		nexts:   map[string]time.Time{},
	}
}

// Load tasks

func (c *Cron) LoadTasks(tasks []Task) error {
	taskMap := map[string]*Task{}
	for _, t := range tasks {
		t1 := t
		if err := t1.check(); err != nil {
			return sderr.WithStack(err)
		}
		if err := t1.ensure(); err != nil {
			return sderr.WithStack(err)
		}
		taskMap[t1.Id] = &t1
	}
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.holder = staticHolder{m: taskMap}
	return nil
}

func (c *Cron) LoadConfig(loc string, pred TaskPred) error {
	data, err := sdload.Bytes(loc)
	if err != nil {
		return sderr.WithMessage(err, "load tasks json error")
	}
	tasks, err := parseTasks(data)
	if err != nil {
		return sderr.WithMessage(err, "parse json tasks error")
	}
	tasks = filterTasks(tasks, pred)
	err = c.LoadTasks(tasks)
	if err != nil {
		return sderr.WithMessage(err, "load tasks error")
	}
	return nil
}

// RegAction

func (c *Cron) RegAction(actionId string, action Action) error {
	if actionId == "" {
		return sderr.New("no action id")
	}
	if action == nil {
		return sderr.New("no action")
	}
	c.actions[actionId] = action
	return nil
}

// Run

func (c *Cron) runTask(task *Task, action Action) {
	c.stopWg.Add(1)
	defer c.stopWg.Done()

	var err error
	panicErr := sdcall.Safe(func() {
		err = action(task.Id, task.Action, task.Args)
	})
	if panicErr != nil {
		err = panicErr
	}
	if err != nil {
		sdlog.WithError(err).WithField("task_id", task.Id).Error("task error")
	} else {
		sdlog.WithField("task_id", task.Id).Debug("task done")
	}
}

// Start

func (c *Cron) tickTasks() {
	const (
		shouldRun  = 1
		notExpired = 2
	)

	taskMap := c.holder.TaskMap()
	now := sdtime.NowTruncateSecond()
	tickTask := func(task *Task) {
		task, ok := taskMap[task.Id]
		if !ok {
			sdlog.WithField("task_id", task.Id).Error("not found task")
			return
		}
		action, ok := c.actions[task.Action]
		if !ok || action == nil {
			sdlog.WithField("task_id", task.Id).WithField("action_id", task.Action).Error("not found action")
			return
		}

		sr := 0
		sdcall.Lock(&c.mtx, func() {
			next, ok := c.nexts[task.Id]
			if !ok {
				c.nexts[task.Id] = task.sched.Next(now)
				sr = notExpired
			} else {
				if now.After(next) {
					c.nexts[task.Id] = task.sched.Next(now)
					sr = shouldRun
				} else {
					sr = notExpired
				}
			}
		})
		if sr == shouldRun {
			go c.runTask(task, action)
		}
	}

	for _, task := range taskMap {
		tickTask(task)
	}
}

func (c *Cron) Start() {
	for {
		time.Sleep(100 * time.Millisecond)
		c.tickTasks()
	}
}

func (c *Cron) Stop() {
	c.stopWg.Wait()
}
