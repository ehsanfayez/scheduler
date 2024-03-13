package scheduler

import (
	"errors"
	"time"
)

type task struct {
	id                   int
	instruction          func() error
	interval             time.Duration
	start_time           time.Time
	finish_time          time.Time
	last_time_performed  time.Time
	without_over_lapping bool
	count                int
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

type Scheduler struct {
	tasks               []task     // tasks to perform with their period
	pending_tasks       chan task  // tasks that are their time to perform
	running_tasks_id    []int      // tasks that are running
	failed_tasks        []failure  // times of a task failed with details
	worker_count        int        // count of worker to run pending_jobs
	id                  func() int // func for get unique and sort id for us in task
	force_pending_tasks []int
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		id:            generateId(),
		worker_count:  3,
		pending_tasks: make(chan task),
	}
}

func generateId() func() int {
	id := 0
	return func() int {
		id++
		return id
	}
}

func (s *Scheduler) SetWorkerCount(count int) (*Scheduler, error) {
	if count < 1 {
		return s, errors.New("count should be more than 0")
	}

	s.worker_count = count
	return s, nil
}

func (s *Scheduler) AddTask(ins func() error) *task {
	t := task{
		id:                   s.id(),
		instruction:          func() error { return ins() },
		count:                -1,
		without_over_lapping: false,
	}

	s.tasks = append(s.tasks, t)
	last_index := len(s.tasks) - 1
	return &s.tasks[last_index]
}

func (t *task) SetInterval(interval time.Duration) *task {
	// Set the interval for the most recently added task
	t.interval = interval
	return t
}

func (t *task) StartAt(start time.Time) *task {
	t.start_time = start
	return t
}

func (t *task) FinishAt(finish time.Time) *task {
	t.finish_time = finish
	return t
}

func (t *task) SetCount(count int) *task {
	if count < 0 {
		panic("Invalid Count")
	}

	t.count = count
	return t
}

func (t *task) WithoutOverLapping() *task {
	t.without_over_lapping = true
	return t
}

func (t *task) GetTaskId() int {
	return t.id
}

func (s *Scheduler) StopTaskById(id int) {
	for index, task := range s.tasks {
		if task.id == id {
			s.tasks = append(s.tasks[:index], s.tasks[index+1:]...)
		}
	}
}

func (s *Scheduler) ForceStopTaskById(id int) {
	for index, task := range s.tasks {
		if task.id == id {
			s.tasks = append(s.tasks[:index], s.tasks[index+1:]...)
			s.force_pending_tasks = append(s.force_pending_tasks, id)
		}
	}
}

func (s *Scheduler) RemoveTaskIdFromForce(id int) {
	for index, value := range s.force_pending_tasks {
		if value == id {
			s.force_pending_tasks = append(s.force_pending_tasks[:index], s.force_pending_tasks[index+1:]...)
		}
	}
}

func (s *Scheduler) ForceRemoveTaskExist(id int) bool {
	for _, value := range s.force_pending_tasks {
		if value == id {
			return true
		}
	}

	return false
}

func (s *Scheduler) CheckAndRunTask() *Scheduler {
	for _, task := range s.tasks {
		var zeroTime time.Time
		if s.IsFinishTime(task.id) || !s.IsStartTime(task) || !s.CheckCount(task.id) || s.IsRunning(task) {
			continue
		}

		if s.ForceRemoveTaskExist(task.id) {
			s.RemoveTaskIdFromForce(task.id)
			continue
		}

		if task.last_time_performed == zeroTime || time.Now().Add(-1*task.interval).Unix() >= task.last_time_performed.Unix() {
			s.pending_tasks <- task
		}
	}

	return s
}

func (s *Scheduler) IsStartTime(t task) bool {
	var zeroTime time.Time
	if zeroTime.Unix() != t.start_time.Unix() && time.Now().Unix() <= t.start_time.Unix() {
		return false
	}

	return true
}

func (s *Scheduler) IsFinishTime(id int) bool {
	var zeroTime time.Time
	for index, task := range s.tasks {
		if task.id == id {
			if zeroTime.Unix() != task.finish_time.Unix() && time.Now().Unix() >= task.finish_time.Unix() {
				// Task has reached its finish time, remove it from the list
				s.tasks = append(s.tasks[:index], s.tasks[index+1:]...)
				return true
			}
		}
	}

	return false
}

func (s *Scheduler) IsRunning(t task) bool {
	if t.without_over_lapping {
		// if task id is in running_tasks_id, return true
		for _, id := range s.running_tasks_id {
			if id == t.id {
				return true
			}
		}
	}

	return false
}

func (s *Scheduler) CheckCount(id int) bool {
	for index, task := range s.tasks {
		if task.id == id {
			if task.count == 0 {
				// Task has reached its count, remove it from the list
				s.tasks = append(s.tasks[:index], s.tasks[index+1:]...)
				return false
			} else {
				s.tasks[index].count = s.tasks[index].count - 1
				return true
			}
		}
	}

	return true
}

func (s *Scheduler) FindFailureTask(id int) int {
	for index, failed_task := range s.failed_tasks {
		if failed_task.task_id == id {
			return index
		}
	}

	return -1
}

func (s *Scheduler) AddFailureTask(id int, err string) *Scheduler {
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

// Define the worker function that processes tasks
func (s *Scheduler) worker(taskQueue <-chan task, workerID int) {
	for task := range taskQueue {
		s.running_tasks_id = append(s.running_tasks_id, task.id)
		err := task.instruction() // Execute the task
		s.running_tasks_id = deleteFromArray(s.running_tasks_id, task.id)
		if err != nil {
			s.AddFailureTask(task.id, err.Error())
		} else {
			for i := range s.tasks {
				if s.tasks[i].id == task.id {
					s.tasks[i].last_time_performed = time.Now()
					break
				}
			}
		}
	}
}

func deleteFromArray(array []int, value int) []int {
	// Find index of value to delete
	index := -1
	for i, v := range array {
		if v == value {
			index = i
			break
		}
	}

	// If value is not found, return original array
	if index == -1 {
		return array
	}

	// Remove value from array
	newArray := append(array[:index], array[index+1:]...)

	return newArray
}

func (s *Scheduler) Start() chan bool {
	stopped := make(chan bool, 1)

	// Create worker pool
	for i := 0; i < s.worker_count; i++ {
		go s.worker(s.pending_tasks, i)
	}

	ticker := time.NewTicker(1 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				s.CheckAndRunTask()
			case <-stopped:
				ticker.Stop()
				// Close the pending_tasks channel to signal workers to stop
				close(s.pending_tasks)
				return
			}
		}
	}()

	return stopped
}
