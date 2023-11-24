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

	// Check the worker_count field
	expectedWorkerCount := 3
	require.Equal(expectedWorkerCount, sch.worker_count)

	// Since the id field is a function, we can call it to get the value and check its correctness.
	// Assuming generateId() is defined in the same package as NewScheduler.
	expectedID := generateId()()
	require.Equal(expectedID, sch.id())

	// Check that the tasks and pending_tasks slices are initialized and empty
	require.Nil(sch.tasks)
	require.Len(sch.pending_tasks, 0)
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

func (suite *SchedulerTestSuite) Test_GetTaskId() {
	require := suite.Require()

	// Create a new scheduler instance using NewScheduler
	sch := NewScheduler()

	expectedID := generateId()()
	taskId := sch.AddTask(func() error { return nil }).GetTaskId()
	require.Equal(expectedID, taskId)
}

func (suite *SchedulerTestSuite) Test_StopTaskById() {
	require := suite.Require()

	// Create a new scheduler instance using NewScheduler
	sch := NewScheduler()

	taskId := sch.AddTask(func() error { return nil }).GetTaskId()
	require.Len(sch.tasks, 1)

	sch.StopTaskById(taskId)
	require.Len(sch.tasks, 0)
}

func (suite *SchedulerTestSuite) Test_ForceStopTaskById() {
	require := suite.Require()

	// Create a new scheduler instance using NewScheduler
	sch := NewScheduler()

	taskId := sch.AddTask(func() error { return nil }).GetTaskId()
	require.Len(sch.tasks, 1)

	sch.ForceStopTaskById(taskId)
	require.Len(sch.tasks, 0)
	require.Len(sch.force_pending_tasks, 1)
}

func (suite *SchedulerTestSuite) Test_RemoveTaskIdFromForce() {
	require := suite.Require()

	// Create a new scheduler instance using NewScheduler
	sch := NewScheduler()
	taskId := sch.AddTask(func() error { return nil }).GetTaskId()

	require.Len(sch.tasks, 1)
	sch.ForceStopTaskById(taskId)
	require.Len(sch.force_pending_tasks, 1)

	sch.RemoveTaskIdFromForce(taskId)
	require.Len(sch.force_pending_tasks, 0)
}

func (suite *SchedulerTestSuite) Test_ForceRemoveTaskExist() {
	require := suite.Require()

	// Create a new scheduler instance using NewScheduler
	sch := NewScheduler()
	taskId := sch.AddTask(func() error { return nil }).GetTaskId()

	require.Len(sch.tasks, 1)
	sch.ForceStopTaskById(taskId)
	check := sch.ForceRemoveTaskExist(taskId)
	require.Equal(true, check)
}

func (suite *SchedulerTestSuite) Test_SetWorkerCount() {
	require := suite.Require()

	// Create a new scheduler instance using NewScheduler
	sch := NewScheduler()

	// Test case 1: Set valid worker count
	first_count := 5
	updatedScheduler, err := sch.SetWorkerCount(first_count)
	require.NoError(err)

	// Check that the worker_count field is updated correctly
	require.Equal(first_count, updatedScheduler.worker_count)

	// Test case 2: Set worker count less than 1
	second_count := -1
	_, err = sch.SetWorkerCount(second_count)

	// Check that it returns the expected error
	expectedError := "count should be more than 0"
	require.EqualError(errors.New(expectedError), err.Error())

	// Check that the worker_count field remains unchanged
	require.Equal(first_count, updatedScheduler.worker_count)
}

func (suite *SchedulerTestSuite) Test_CheckAndRunTask() {
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
	interval := time.Second * 5
	sch.AddTask(taskInstruction).SetInterval(interval)

	// Check if the interval is set correctly for the last task
	lastTaskIndex := len(sch.tasks) - 1
	require.Equal(interval, sch.tasks[lastTaskIndex].interval)
}

func (suite *SchedulerTestSuite) Test_StartAt() {
	require := suite.Require()

	// Create a new scheduler instance using NewScheduler
	sch := NewScheduler()

	// Add a task to the scheduler
	taskInstruction := func() error { return nil }
	time := time.Now().Add(1 * time.Hour)
	sch.AddTask(taskInstruction).StartAt(time)

	// Check if the interval is set correctly for the last task
	lastTaskIndex := len(sch.tasks) - 1
	require.Equal(time, sch.tasks[lastTaskIndex].start_time)
}

func (suite *SchedulerTestSuite) Test_FinishAt() {
	require := suite.Require()

	// Create a new scheduler instance using NewScheduler
	sch := NewScheduler()

	// Add a task to the scheduler
	taskInstruction := func() error { return nil }
	time := time.Now().Add(2 * time.Hour)
	sch.AddTask(taskInstruction).FinishAt(time)

	// Check if the interval is set correctly for the last task
	lastTaskIndex := len(sch.tasks) - 1
	require.Equal(time, sch.tasks[lastTaskIndex].finish_time)
}

func (suite *SchedulerTestSuite) Test_SetCount() {
	require := suite.Require()

	// Create a new scheduler instance using NewScheduler
	sch := NewScheduler()

	// Add a task to the scheduler
	taskInstruction := func() error { return nil }
	sch.AddTask(taskInstruction).SetCount(5)

	// Check if the interval is set correctly for the last task
	lastTaskIndex := len(sch.tasks) - 1
	require.Equal(5, sch.tasks[lastTaskIndex].count)
}

func (suite *SchedulerTestSuite) Test_IsStartTime() {
	require := suite.Require()

	// Create a new scheduler instance using NewScheduler
	sch := NewScheduler()

	// Add a task to the scheduler with start at 2 hours later
	taskInstruction := func() error { return nil }
	time1 := time.Now().Add(2 * time.Hour)
	sch.AddTask(taskInstruction).StartAt(time1)
	lastTaskIndex := len(sch.tasks) - 1
	require.Equal(false, sch.IsStartTime(sch.tasks[lastTaskIndex]))

	// Add a task to the scheduler with start at 2 hours before
	time2 := time.Now().Add(-2 * time.Hour)
	sch.AddTask(taskInstruction).StartAt(time2)
	lastTaskIndex = len(sch.tasks) - 1
	require.Equal(true, sch.IsStartTime(sch.tasks[lastTaskIndex]))
}

func (suite *SchedulerTestSuite) Test_IsFinishTime() {
	require := suite.Require()

	// Create a new scheduler instance using NewScheduler
	sch := NewScheduler()

	// Add a task to the scheduler with start at 2 hours later
	taskInstruction := func() error { return nil }
	time1 := time.Now().Add(2 * time.Hour)
	sch.AddTask(taskInstruction).FinishAt(time1)
	lastTaskIndex := len(sch.tasks) - 1
	require.Equal(false, sch.IsFinishTime(sch.tasks[lastTaskIndex].id))

	// Add a task to the scheduler with start at 2 hours before
	time2 := time.Now().Add(-2 * time.Hour)
	sch.AddTask(taskInstruction).FinishAt(time2)
	lastTaskIndex = len(sch.tasks) - 1
	removed := sch.IsFinishTime(sch.tasks[lastTaskIndex].id)
	require.Equal(true, removed)
	require.Len(sch.tasks, 1)
}

func (suite *SchedulerTestSuite) Test_CheckCount() {
	require := suite.Require()

	// Create a new scheduler instance using NewScheduler
	sch := NewScheduler()

	// Add a task to the scheduler with count 2
	taskInstruction := func() error { return nil }
	sch.AddTask(taskInstruction).EverySecond().SetCount(2)
	lastTaskIndex := len(sch.tasks) - 1

	// check count when count is 2
	countState := sch.CheckCount(sch.tasks[lastTaskIndex].id)
	require.Equal(true, countState)
	require.Equal(1, sch.tasks[lastTaskIndex].count)

	// check count when count is 1
	countState = sch.CheckCount(sch.tasks[lastTaskIndex].id)
	require.Equal(true, countState)
	require.Equal(0, sch.tasks[lastTaskIndex].count)

	// check count when count is 0
	countState = sch.CheckCount(sch.tasks[lastTaskIndex].id)
	require.Equal(false, countState)
	require.Len(sch.tasks, 0)
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
	require.Equal(expectedFailure.task_id, sch.failed_tasks[0].task_id)
	require.Equal(expectedFailure.fails[0].error_message, sch.failed_tasks[0].fails[0].error_message)

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
	require.Len(sch.failed_tasks[0].fails, 2)

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
	require.Equal(expectedFailure2.task_id, sch.failed_tasks[1].task_id)
}

func (suite *SchedulerTestSuite) Test_Start() {
	require := suite.Require()

	// Create a new scheduler instance using NewScheduler
	s := NewScheduler()

	// Start the scheduler and get the stopped channel
	_, stopped := s.Start()

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
