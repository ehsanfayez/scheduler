# Scheduler Package
The Scheduler package provides a simple and flexible way to schedule and manage tasks in Go. It allows you to schedule tasks with specific intervals, start and finish times, execution counts, and options to prevent overlapping executions.

## Features
Task Scheduling: Schedule tasks to run at specified intervals.
Start and Finish Times: Define when tasks should start and stop executing.
Execution Counts: Limit the number of times a task should run.
Overlapping Control: Prevent tasks from overlapping executions.
Concurrency Control: Configure the number of worker goroutines.
Failure Tracking: Keep track of task failures with detailed error messages and timestamps.
Installation
To use the Scheduler package, simply import it into your Go project:

```go
import "path/to/your/package/scheduler"

```
Replace "path/to/your/package/scheduler" with the actual import path where the Scheduler package resides.

### Usage
Creating a Scheduler
Start by creating a new instance of the Scheduler:

```go
s := scheduler.NewScheduler()
```
### Configuring Worker Count
By default, the Scheduler uses 3 worker goroutines to execute tasks. You can change this number using SetWorkerCount:

```go
s.SetWorkerCount(5) // Sets the worker count to 5
```
### Adding Tasks
Add a task by providing a function that returns an error:

```go
task := s.AddTask(func() error {
	// Your task logic here
	return nil
})
```
### Configuring Tasks
After adding a task, you can configure it using method chaining:

```go
task.
	SetInterval(10 * time.Second).
	StartAt(time.Now().Add(1 * time.Minute)).
	FinishAt(time.Now().Add(1 * time.Hour)).
	SetCount(5).
	WithoutOverLapping()
```
### Task Configuration Methods
```go
SetInterval(interval time.Duration) *task: Sets the interval 
```
between task executions.
```go
StartAt(start time.Time) *task: Sets the time when the task 
```
should start executing.
```go
FinishAt(finish time.Time) *task: Sets the time when the task 
```
should stop executing.
```go
SetCount(count int) *task: Limits the number of times the 
```
task should execute.
```go
WithoutOverLapping() *task: Ensures the task doesn't start a 
```
new execution until the previous one finishes.
```go
GetTaskId() int: Retrieves the unique ID of the task.
```
### Starting the Scheduler
Start the Scheduler to begin executing tasks:

```go
stopChan := s.Start()
```
### Stopping the Scheduler
To gracefully stop the Scheduler and all running tasks:

```go
stopChan <- true
```
### Managing Tasks
Stopping a Specific Task
To stop a specific task from being scheduled:

```go
s.StopTaskById(task.GetTaskId())
```
Forcing a Task to Stop Immediately
To forcefully stop a task and remove any pending executions:

```go
s.ForceStopTaskById(task.GetTaskId())
```
Resuming a Force-Stopped Task
To remove a task from the force-stop list and allow it to be scheduled again:

```go
s.RemoveTaskIdFromForce(task.GetTaskId())
```
### Monitoring Task Failures
The Scheduler keeps track of task failures in the failed_tasks slice. Each failure includes the task ID, failure count, and detailed error messages with timestamps.

You can access and inspect s.failed_tasks to monitor task failures.

Example
Here's a complete example demonstrating how to use the Scheduler:

```go
package main

import (
	"fmt"
	"time"

	"path/to/your/package/scheduler"
)

func main() {
	// Create a new Scheduler
	s := scheduler.NewScheduler()

	// Set the number of worker goroutines
	s.SetWorkerCount(5)

	// Add a task
	task := s.AddTask(func() error {
		fmt.Println("Task executed at", time.Now())
		return nil
	})

	// Configure the task
	task.
		SetInterval(10 * time.Second).
		StartAt(time.Now().Add(1 * time.Minute)).
		FinishAt(time.Now().Add(1 * time.Hour)).
		SetCount(5).
		WithoutOverLapping()

	// Start the Scheduler
	stopChan := s.Start()

	// Run the Scheduler for a certain duration
	time.Sleep(2 * time.Hour)

	// Stop the Scheduler
	stopChan <- true
}
```
### Notes
Thread Safety: Ensure that your task functions are thread-safe if they access shared resources.
Execution Timing: The Scheduler checks for tasks to run every second. Tasks may not execute at the exact specified time but will run as close as possible.
Error Handling: If a task returns an error, the failure is recorded but the task continues to be scheduled unless the execution count is reached or the finish time has passed.
Overlapping Executions: Use WithoutOverLapping() to prevent a task from starting a new execution before the previous one finishes.
Contributing
Contributions are welcome! If you encounter any issues or have suggestions for improvements, please open an issue or submit a pull request on the project's repository.

License
[Specify the license under which the package is distributed, e.g., MIT License]