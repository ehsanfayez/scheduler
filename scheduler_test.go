package scheduler

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type SchedulerTestSuite struct {
	suite.Suite
}

func (suite *SchedulerTestSuite) SetupSuite() {
}

func (suite *SchedulerTestSuite) Test_NewScheduler() {
	require := suite.Require()

	// Create a new scheduler instance using NewScheduler
	sch := NewScheduler()

	// Check that the scheduler is not nil
	require.NotNil(sch)

	// Check the jobs_count field
	expectedJobsCount := 3
	require.Equal(expectedJobsCount, sch.jobs_count)

	// Since the id field is a function, we can call it to get the value and check its correctness.
	// Assuming generateId() is defined in the same package as NewScheduler.
	expectedID := generateId()()
	require.Equal(expectedID, sch.id())

	// Check that the tasks and pending_tasks slices are initialized and empty
	require.Nil(sch.tasks)
	require.Nil(sch.pending_tasks)
}

func (suite *SchedulerTestSuite) Test_GenerateId() {
	require := suite.Require()

	// Get the generated ID function from generateId()
	generateIDFunc := generateId()

	// Generate the first ID and check if it is 1
	expectedID := 1
	id := generateIDFunc()
	require.Equal(expectedID, id)

	// Generate a few more IDs to verify incrementing behavior
	for i := 2; i <= 10; i++ {
		expectedID = i
		id = generateIDFunc()
		require.Equal(expectedID, id)
	}
}

func (suite *SchedulerTestSuite) Test_SetJobsCount() {
	require := suite.Require()

	// Create a new scheduler instance using NewScheduler
	sch := NewScheduler()

	// Test case 1: Set valid jobs count
	first_count := 5
	updatedScheduler, err := sch.SetJobsCount(first_count)
	require.NoError(err)

	// Check that the jobs_count field is updated correctly
	require.Equal(first_count, updatedScheduler.jobs_count)

	// Test case 2: Set jobs count less than 1
	second_count := -1
	_, err = sch.SetJobsCount(second_count)

	// Check that it returns the expected error
	expectedError := "count should be more than 0"
	require.EqualError(errors.New(expectedError), err.Error())

	// Check that the jobs_count field remains unchanged
	require.Equal(first_count, updatedScheduler.jobs_count)
}

func (suite *SchedulerTestSuite) Test_AddTask() {
	require := suite.Require()

	// Create a new scheduler instance using NewScheduler
	sch := NewScheduler()

	// Add the first task and check if it's assigned the correct ID
	taskInstruction1 := func() error { return nil }
	sch.AddTask(taskInstruction1).SetInterval(1 * time.Second)
	require.Len(sch.tasks, 1)

	// Assuming the generateId() function is defined in the same package.
	id := generateId()
	expectedID1 := id()
	require.Equal(expectedID1, sch.tasks[0].id)

	// Add the second task and check if it's assigned the correct ID
	taskInstruction2 := func() error { return nil }
	sch.AddTask(taskInstruction2).SetInterval(1 * time.Second)
	require.Len(sch.tasks, 2)

	// Assuming the generateId() function is defined in the same package.
	expectedID2 := id()
	require.Equal(expectedID2, sch.tasks[1].id)
}

func (suite *SchedulerTestSuite) Test_SetInterval() {
	require := suite.Require()

	// Create a new scheduler instance using NewScheduler
	sch := NewScheduler()

	// Add a task to the scheduler
	taskInstruction := func() error { return nil }
	sch.AddTask(taskInstruction)

	// Set an interval for the last added task
	interval := time.Second * 5
	sch.SetInterval(interval)

	// Check if the interval is set correctly for the last task
	lastTaskIndex := len(sch.tasks) - 1
	require.Equal(interval, sch.tasks[lastTaskIndex].interval)
}

func (suite *SchedulerTestSuite) Test_PushToPendingTasks() {
	require := suite.Require()

	// Create a new scheduler instance using NewScheduler
	sch := NewScheduler()

	// Create a task
	taskInstruction := func() error { return nil }
	task := task{
		id:          1, // You can assign a specific ID here.
		instruction: taskInstruction,
	}

	// Push the task to the pending tasks
	sch.PushToPendingTasks(task)

	// Check if the task is added correctly to the pending tasks
	require.Len(sch.pending_tasks, 1)

	// Check if the task content is correctly pushed to the pending tasks
	require.Equal(task.id, sch.pending_tasks[0].id)
}

func (suite *SchedulerTestSuite) Test_RemoveFromPendingTasks() {
	require := suite.Require()

	// Create a new scheduler instance using NewScheduler
	sch := NewScheduler()

	// Create a task
	taskInstruction := func() error { return nil }
	task := task{
		id:          1, // You can assign a specific ID here.
		instruction: taskInstruction,
	}

	// Push the task to the pending tasks
	sch.PushToPendingTasks(task)

	//	Remove from pending tasks
	sch.RemoveFromPendingTasks(task.id)

	// Check if the task is added correctly to the pending tasks
	require.Len(sch.pending_tasks, 0)
}

func (suite *SchedulerTestSuite) Test_AddPendingTasks() {
	require := suite.Require()

	// Create a new scheduler instance using NewScheduler
	sch := NewScheduler()

	// Create a task with interval of 1 second (1000 milliseconds)
	task1 := task{
		id:                  1,
		instruction:         func() error { return nil },
		interval:            time.Second,
		last_time_performed: time.Now().Add(-2 * time.Second), // Last execution 2 seconds ago
	}
	// Add the task to the scheduler's tasks
	sch.tasks = append(sch.tasks, task1)

	// Create a task with interval of 2 seconds (2000 milliseconds)
	task2 := task{
		id:                  2,
		instruction:         func() error { return nil },
		interval:            2 * time.Second,
		last_time_performed: time.Now().Add(-1 * time.Second), // Last execution 1 second ago
	}
	// Add the task to the scheduler's tasks
	sch.tasks = append(sch.tasks, task2)

	// Create a task with interval of 3 seconds (3000 milliseconds)
	task3 := task{
		id:                  3,
		instruction:         func() error { return nil },
		interval:            3 * time.Second,
		last_time_performed: time.Now().Add(-4 * time.Second), // Last execution 4 seconds ago
	}
	// Add the task to the scheduler's tasks
	sch.tasks = append(sch.tasks, task3)

	// Call AddPendingTasks to add the pending tasks to the scheduler's pending_tasks
	sch.AddPendingTasks()

	// Check if the pending_tasks contains the expected tasks
	expectedPendingTasks := []task{task1, task3}
	require.Len(sch.pending_tasks, 2)

	// Check if the pending tasks are the expected ones
	for i := 0; i < len(expectedPendingTasks); i++ {
		require.Equal(expectedPendingTasks[i].id, sch.pending_tasks[i].id)
	}
}

func (suite *SchedulerTestSuite) Test_FindFailureTask() {
	require := suite.Require()

	// Create a new scheduler instance using NewScheduler
	sch := NewScheduler()

	// Add a failure task with ID 1 and error message "error1"
	sch.AddFailureTask(1, "error1")

	// check index of exist failure task return correctly
	expectedIndex := 0
	index := sch.FindFailureTask(1)
	require.Equal(expectedIndex, index)

	// check index of does not exist failure task return correctly
	expectedIndex = -1
	index = sch.FindFailureTask(2)
	require.Equal(expectedIndex, index)
}

func (suite *SchedulerTestSuite) Test_AddFailureTask() {
	require := suite.Require()

	// Create a new scheduler instance using NewScheduler
	sch := NewScheduler()

	// Add a failure task with ID 1 and error message "error1"
	sch.AddFailureTask(1, "error1")

	// Check if the failed_tasks slice contains the expected failure task
	expectedFailure := failure{
		task_id: 1,
		count:   1,
		fails: []fail{
			{
				time:          time.Now(), // The time when the failure was added
				error_message: "error1",
			},
		},
	}

	require.Len(sch.failed_tasks, 1)

	// Compare the content of the failure tasks
	require.Equal(expectedFailure, sch.failed_tasks[0])

	// Add another failure task with ID 1 and a different error message "error2"
	sch.AddFailureTask(1, "error2")

	// Check if the count is incremented and the error message is added to the existing failure task
	expectedFailure.count = 2
	expectedFailure.fails = append(expectedFailure.fails, fail{
		time:          time.Now(),
		error_message: "error2",
	})

	require.Len(sch.failed_tasks, 1)

	// Compare the content of the failure tasks
	require.Equal(expectedFailure, sch.failed_tasks[0])

	// Add a failure task with a new ID 2 and error message "error3"
	sch.AddFailureTask(2, "error3")

	// Check if a new failure task is added to the failed_tasks slice
	expectedFailure2 := failure{
		task_id: 2,
		count:   1,
		fails: []fail{
			{
				time:          time.Now(),
				error_message: "error3",
			},
		},
	}

	require.Len(sch.failed_tasks, 2)

	// Compare the content of the failure tasks
	require.Equal(expectedFailure2, sch.failed_tasks[1])
}

func (suite *SchedulerTestSuite) Test_RunPendingTasks() {
	require := suite.Require()

	// Create a new scheduler instance using NewScheduler
	sch := NewScheduler()

	// Create a task with interval of 1 second (1000 milliseconds)
	task1 := task{
		id:                  1,
		instruction:         func() error { return errors.New("error task 1") },
		interval:            time.Second,
		last_time_performed: time.Time{},
	}
	// Add the task to the scheduler's tasks
	sch.tasks = append(sch.tasks, task1)

	// Run the pending tasks
	sch.AddPendingTasks().RunPendingTasks()

	// Wait for goroutines to finish
	time.Sleep(time.Second)

	// Get the current time
	now := time.Now()

	// Check if the failure task was added
	expectedFailureTask := failure{
		task_id: 1,
		count:   1,
		fails: []fail{
			{
				time:          now.Add(-1 * time.Second),
				error_message: "Some error occurred",
			},
		},
	}

	require.Len(sch.failed_tasks, 1)

	// Compare the content of the failure tasks
	require.Equal(expectedFailureTask.task_id, sch.failed_tasks[0].task_id)
	require.Equal(expectedFailureTask.fails[0].time.Unix(), sch.failed_tasks[0].fails[0].time.Unix())
	require.Len(sch.failed_tasks[0].fails, 1)

	// Create a mock task with ID 2 that returns no error
	noErrTask := task{
		id:                  2,
		instruction:         func() error { return nil },
		interval:            time.Second,
		last_time_performed: time.Time{},
	}

	// Add the mock task to the scheduler's pending_tasks
	sch.tasks = append(sch.tasks, noErrTask)

	// Run the pending tasks
	sch.AddPendingTasks()

	// Check if the last_time_performed was updated and the task was removed from pending_tasks
	require.Len(sch.pending_tasks, 2)
	sch.RunPendingTasks()

	// Wait for goroutines to finish
	time.Sleep(time.Second)
	require.Len(sch.failed_tasks, 1)
	require.Len(sch.failed_tasks[0].fails, 2)
}

func (suite *SchedulerTestSuite) Test_Start() {
	require := suite.Require()

	// Create a new scheduler instance using NewScheduler
	s := NewScheduler()

	// Start the scheduler and get the stopped channel
	stopped := s.Start()

	// Wait for 3 seconds to let the scheduler run
	time.Sleep(3 * time.Second)

	// Stop the scheduler by sending a signal on the stopped channel
	stopped <- true

	// Wait for the scheduler to stop
	time.Sleep(1 * time.Second)

	// Check if the scheduler stopped by verifying that no pending tasks are left
	require.Empty(s.pending_tasks)
}

func TestSchedulerTestSuite(t *testing.T) {
	suite.Run(t, new(SchedulerTestSuite))
}
