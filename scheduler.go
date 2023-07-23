package scheduler

import "time"

type task struct {
	id                  int
	instruction         func()
	interval            time.Duration
	last_time_performed time.Time
}

type scheduler struct {
	tasks         []task //
	pending_tasks []task
	jobs_count    int
	id            func() int
}

func NewScheduler() *scheduler {
	return &scheduler{
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

func (s *scheduler) SetJobsCount(count int) {
	s.jobs_count = count
}

func (s *scheduler) Add(ins func()) *scheduler {
	t := task{
		id:          s.id(),
		instruction: ins,
	}
	s.tasks = append(s.tasks, t)
	return s
}

func (s *scheduler) SetInterval(interval time.Duration) *scheduler {
	index := len(s.tasks) - 1
	s.tasks[index].interval = interval
	return s
}

func (s *scheduler) PushToPending(task task) {
	s.pending_tasks = append(s.pending_tasks, task)
}

func (s *scheduler) RemoveFromPending(id int) {
	for i, item := range s.pending_tasks {
		if item.id == id {
			s.pending_tasks = append(s.pending_tasks[:i], s.pending_tasks[i+1:]...)
			break
		}
	}
}

func (s *scheduler) AddPendingJobs() *scheduler {
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

func (s *scheduler) RunPendingJobs() {
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

func (s *scheduler) Start() chan bool {
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
