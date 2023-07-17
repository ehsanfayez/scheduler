package scheduler

import "time"

type Task struct {
	id                  int
	instruction         func()
	interval            time.Duration
	last_time_performed time.Time
}

type Scheduler struct {
	tasks         []Task //
	pending_tasks []Task
	jobs_count    int
	id            func() int
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		id:         generateId(),
		jobs_count: 3,
	}
}

func generateId() func() int {
	id := 0
	return func() int {
		id++
		return id
	}
}

func (s *Scheduler) SetJobsCount(count int) {
	s.jobs_count = count
}

func (s *Scheduler) Add(task func()) *Scheduler {
	t := Task{
		id:          s.id(),
		instruction: task,
	}
	s.tasks = append(s.tasks, t)
	return s
}

func (s *Scheduler) SetInterval(interval time.Duration) *Scheduler {
	index := len(s.tasks) - 1
	s.tasks[index].interval = interval
	return s
}

func (s *Scheduler) PushToPending(task Task) {
	s.pending_tasks = append(s.pending_tasks, task)
}

func (s *Scheduler) RemoveFromPending(id int) {
	for i, item := range s.pending_tasks {
		if item.id == id {
			s.pending_tasks = append(s.pending_tasks[:i], s.pending_tasks[i+1:]...)
			break
		}
	}
}

func (s *Scheduler) AddPendingJobs() *Scheduler {
	for _, task := range s.tasks {
		var zeroTime time.Time
		if task.last_time_performed == zeroTime {
			s.PushToPending(task)
		} else if time.Now().Add(-1*task.interval).Unix() >= task.last_time_performed.Unix() {
			s.PushToPending(task)
		}
	}
	return s
}

func (s *Scheduler) RunPendingJobs() {
	for _, task := range s.pending_tasks {
		go task.instruction()
		for i := range s.tasks {
			if s.tasks[i].id == task.id {
				s.tasks[i].last_time_performed = time.Now()
				break
			}
		}
		s.RemoveFromPending(task.id)
	}
}

func (s *Scheduler) Start() chan bool {
	stopped := make(chan bool, 1)
	ticker := time.NewTicker(1 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				s.AddPendingJobs().RunPendingJobs()
			case <-stopped:
				ticker.Stop()
				return
			}
		}
	}()

	return stopped
}
