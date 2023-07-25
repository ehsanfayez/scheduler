package scheduler

import (
	"errors"
	"time"
)

// Define the possible states of the scheduler
type schedulerState int

const (
	stateNew schedulerState = iota
	stateTaskAdded
)

type task struct {
	id                  int
	instruction         func() error
	interval            time.Duration
	last_time_performed time.Time
}

type fail struct {
	time          time.Time
	error_message string
}

type failure struct {
	task_id int
	count   int
	fails   []fail
}

type scheduler struct {
	tasks         []task     // tasks to perform with their period
	pending_tasks []task     // tasks that are their time to perform
	failed_tasks  []failure  // times of a task failed with details
	jobs_count    int        // count of worker to run pending_jobs
	id            func() int // func for get unique and sort id for us in task

	last_added_task_index int            // index of the most recently added task
	state                 schedulerState // current state of the scheduler
}

func NewScheduler() *scheduler {
	return &scheduler{
		id:         generateId(),
		jobs_count: 3,
		state:      stateNew, // Initialize the state to "New"
	}
}

func generateId() func() int {
	id := 0
	return func() int {
		id++
		return id
	}
}

func (s *scheduler) SetJobsCount(count int) (*scheduler, error) {
	if count < 1 {
		return s, errors.New("count should be more than 0")
	}

	s.jobs_count = count
	return s, nil
}

func (s *scheduler) AddTask(ins func() error) *scheduler {
	// Check if AddTask() has been called before SetInterval()
	if s.state != stateNew {
		panic(errors.New("you can call AddTask() twice before call SetInterval for last AddTask()"))
	}

	t := task{
		id:          s.id(),
		instruction: ins,
	}

	s.tasks = append(s.tasks, t)
	// Set the index of the most recently added task
	s.last_added_task_index = len(s.tasks) - 1

	// Update the state to "TaskAdded"
	s.state = stateTaskAdded
	return s
}

func (s *scheduler) SetInterval(interval time.Duration) *scheduler {
	// Check if AddTask() has been called before SetInterval()
	if s.state != stateTaskAdded {
		panic(errors.New("SetInterval() should be called after AddTask()"))
	}

	if s.last_added_task_index == -1 {
		// No task has been added yet, return without setting the interval
		return s
	}

	// Set the interval for the most recently added task
	s.tasks[s.last_added_task_index].interval = interval
	s.last_added_task_index = -1
	s.state = stateNew
	return s
}

func (s *scheduler) PushToPendingTasks(task task) {
	s.pending_tasks = append(s.pending_tasks, task)
}

func (s *scheduler) RemoveFromPendingTasks(id int) {
	for i, item := range s.pending_tasks {
		if item.id == id {
			s.pending_tasks = append(s.pending_tasks[:i], s.pending_tasks[i+1:]...)
			break
		}
	}
}

func (s *scheduler) AddPendingTasks() *scheduler {
	for _, task := range s.tasks {
		var zeroTime time.Time
		if task.last_time_performed == zeroTime {
			s.PushToPendingTasks(task)
		} else if time.Now().Add(-1*task.interval).Unix() >= task.last_time_performed.Unix() {
			s.PushToPendingTasks(task)
		}
	}

	return s
}

func (s *scheduler) FindFailureTask(id int) int {
	for index, failed_task := range s.failed_tasks {
		if failed_task.task_id == id {
			return index
		}
	}
	return -1
}

func (s *scheduler) AddFailureTask(id int, err string) *scheduler {
	index := s.FindFailureTask(id)
	if index >= 0 {
		new_failure := fail{
			time:          time.Now(),
			error_message: err,
		}
		s.failed_tasks[index].fails = append(s.failed_tasks[index].fails, new_failure)
		s.failed_tasks[index].count++
	} else {
		s.failed_tasks = append(s.failed_tasks, failure{
			task_id: id,
			count:   1,
			fails: []fail{
				{
					time:          time.Now(),
					error_message: err,
				},
			},
		})
	}

	return s
}

func (s *scheduler) RunPendingTasks() {
	for _, pending_task := range s.pending_tasks {
		go func(t task) {
			err := t.instruction()
			if err != nil {
				s.AddFailureTask(t.id, err.Error())
				s.RemoveFromPendingTasks(t.id)
			} else {
				for i := range s.tasks {
					if s.tasks[i].id == t.id {
						s.tasks[i].last_time_performed = time.Now()
						break
					}
				}
				s.RemoveFromPendingTasks(t.id)
			}
		}(pending_task)
	}
}

func (s *scheduler) Start() chan bool {
	stopped := make(chan bool, 1)
	ticker := time.NewTicker(1 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				s.AddPendingTasks().RunPendingTasks()
			case <-stopped:
				ticker.Stop()
				return
			}
		}
	}()

	return stopped
}
