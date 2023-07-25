package scheduler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type TimesTestSuite struct {
	suite.Suite
}

func (suite *TimesTestSuite) SetupSuite() {
}

func (suite *TimesTestSuite) Test_EverySecond() {
	require := suite.Require()

	sch := NewScheduler()
	taskInstruction := func() error { return nil }
	sch.AddTask(taskInstruction).EverySecond()

	// Check if the interval is set correctly to 1 second
	expectedInterval := 1 * time.Second
	require.Equal(expectedInterval, sch.tasks[0].interval)
}

func (suite *TimesTestSuite) Test_EverySeconds() {
	require := suite.Require()

	sch := NewScheduler()
	taskInstruction := func() error { return nil }
	sch.AddTask(taskInstruction).EverySeconds(5)

	// Check if the interval is set correctly to 1 second
	expectedInterval := 5 * time.Second
	require.Equal(expectedInterval, sch.tasks[0].interval)
}

func (suite *TimesTestSuite) Test_EveryMinute() {
	require := suite.Require()

	sch := NewScheduler()
	taskInstruction := func() error { return nil }
	sch.AddTask(taskInstruction).EveryMinute()

	// Check if the interval is set correctly to 1 second
	expectedInterval := 1 * time.Minute
	require.Equal(expectedInterval, sch.tasks[0].interval)
}

func (suite *TimesTestSuite) Test_EveryMinutes() {
	require := suite.Require()

	sch := NewScheduler()
	taskInstruction := func() error { return nil }
	sch.AddTask(taskInstruction).EveryMinutes(5)

	// Check if the interval is set correctly to 1 second
	expectedInterval := 5 * time.Minute
	require.Equal(expectedInterval, sch.tasks[0].interval)
}

func (suite *TimesTestSuite) Test_EveryHour() {
	require := suite.Require()

	sch := NewScheduler()
	taskInstruction := func() error { return nil }
	sch.AddTask(taskInstruction).EveryHour()

	// Check if the interval is set correctly to 1 second
	expectedInterval := 1 * time.Hour
	require.Equal(expectedInterval, sch.tasks[0].interval)
}

func (suite *TimesTestSuite) Test_EveryHours() {
	require := suite.Require()

	sch := NewScheduler()
	taskInstruction := func() error { return nil }
	sch.AddTask(taskInstruction).EveryHours(5)

	// Check if the interval is set correctly to 1 second
	expectedInterval := 5 * time.Hour
	require.Equal(expectedInterval, sch.tasks[0].interval)
}

func (suite *TimesTestSuite) Test_EveryDay() {
	require := suite.Require()

	sch := NewScheduler()
	taskInstruction := func() error { return nil }
	sch.AddTask(taskInstruction).EveryDay()

	// Check if the interval is set correctly to 1 second
	expectedInterval := 24 * time.Hour
	require.Equal(expectedInterval, sch.tasks[0].interval)
}

func (suite *TimesTestSuite) Test_EveryDays() {
	require := suite.Require()

	sch := NewScheduler()
	taskInstruction := func() error { return nil }
	sch.AddTask(taskInstruction).EveryDays(5)

	// Check if the interval is set correctly to 1 second
	expectedInterval := 5 * 24 * time.Hour
	require.Equal(expectedInterval, sch.tasks[0].interval)
}

func (suite *TimesTestSuite) Test_EveryWeek() {
	require := suite.Require()

	sch := NewScheduler()
	taskInstruction := func() error { return nil }
	sch.AddTask(taskInstruction).EveryWeek()

	// Check if the interval is set correctly to 1 second
	expectedInterval := 7 * 24 * time.Hour
	require.Equal(expectedInterval, sch.tasks[0].interval)
}

func (suite *TimesTestSuite) Test_EveryWeeks() {
	require := suite.Require()

	sch := NewScheduler()
	taskInstruction := func() error { return nil }
	sch.AddTask(taskInstruction).EveryWeeks(5)

	// Check if the interval is set correctly to 1 second
	expectedInterval := 5 * 7 * 24 * time.Hour
	require.Equal(expectedInterval, sch.tasks[0].interval)
}

func (suite *TimesTestSuite) Test_EveryMonth() {
	require := suite.Require()

	sch := NewScheduler()
	taskInstruction := func() error { return nil }
	sch.AddTask(taskInstruction).EveryMonth()

	// Check if the interval is set correctly to 1 second
	expectedInterval := 30 * 24 * time.Hour
	require.Equal(expectedInterval, sch.tasks[0].interval)
}

func (suite *TimesTestSuite) Test_EveryMonths() {
	require := suite.Require()

	sch := NewScheduler()
	taskInstruction := func() error { return nil }
	sch.AddTask(taskInstruction).EveryMonths(5)

	// Check if the interval is set correctly to 1 second
	expectedInterval := 5 * 30 * 24 * time.Hour
	require.Equal(expectedInterval, sch.tasks[0].interval)
}

func (suite *TimesTestSuite) Test_EveryYear() {
	require := suite.Require()

	sch := NewScheduler()
	taskInstruction := func() error { return nil }
	sch.AddTask(taskInstruction).EveryYear()

	// Check if the interval is set correctly to 1 second
	expectedInterval := 365 * 24 * time.Hour
	require.Equal(expectedInterval, sch.tasks[0].interval)
}

func (suite *TimesTestSuite) Test_EveryYears() {
	require := suite.Require()

	sch := NewScheduler()
	taskInstruction := func() error { return nil }
	sch.AddTask(taskInstruction).EveryYears(5)

	// Check if the interval is set correctly to 1 second
	expectedInterval := 5 * 365 * 24 * time.Hour
	require.Equal(expectedInterval, sch.tasks[0].interval)
}

func TestTimesTestSuite(t *testing.T) {
	suite.Run(t, new(TimesTestSuite))
}
