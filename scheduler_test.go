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
	sch.AddTask(taskInstruction1)
	require.Len(sch.tasks, 1)

	// Assuming the generateId() function is defined in the same package.
	id := generateId()
	expectedID1 := id()
	require.Equal(expectedID1, sch.tasks[0].id)

	// Add the second task and check if it's assigned the correct ID
	taskInstruction2 := func() error { return nil }
	sch.AddTask(taskInstruction2)
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

func TestSchedulerTestSuite(t *testing.T) {
	suite.Run(t, new(SchedulerTestSuite))
}
