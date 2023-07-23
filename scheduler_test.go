package scheduler

import (
	"github.com/stretchr/testify/suite"
)

type SchedulerTestSuite struct {
	suite.Suite
}

func (suite *SchedulerTestSuite) SetupSuite() {

}

func (suite *SchedulerTestSuite) Test_NewScheduler_Success() {
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
