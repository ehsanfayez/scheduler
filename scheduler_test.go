package scheduler

import (
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
	require.Equal(expectedID, sch.id)

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
