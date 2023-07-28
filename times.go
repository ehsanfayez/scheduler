package scheduler

import "time"

func (t *task) EverySecond() *task {
	t.SetInterval(1 * time.Second)
	return t
}

func (t *task) EverySeconds(count int) *task {
	t.SetInterval(time.Duration(count) * time.Second)
	return t
}

func (t *task) EveryMinute() *task {
	t.SetInterval(1 * time.Minute)
	return t
}

func (t *task) EveryMinutes(count int) *task {
	t.SetInterval(time.Duration(count) * time.Minute)
	return t
}

func (t *task) EveryHour() *task {
	t.SetInterval(1 * time.Hour)
	return t
}

func (t *task) EveryHours(count int) *task {
	t.SetInterval(time.Duration(count) * time.Hour)
	return t
}

func (t *task) EveryDay() *task {
	t.SetInterval(24 * time.Hour)
	return t
}

func (t *task) EveryDays(count int) *task {
	t.SetInterval(time.Duration(count*24) * time.Hour)
	return t
}

func (t *task) EveryWeek() *task {
	t.SetInterval(24 * 7 * time.Hour)
	return t
}

func (t *task) EveryWeeks(count int) *task {
	t.SetInterval(time.Duration(count*24*7) * time.Hour)
	return t
}

func (t *task) EveryMonth() *task {
	t.SetInterval(30 * 24 * time.Hour)
	return t
}

func (t *task) EveryMonths(count int) *task {
	t.SetInterval(time.Duration(count*24*30) * time.Hour)
	return t
}

func (t *task) EveryYear() *task {
	t.SetInterval(365 * 24 * time.Hour)
	return t
}

func (t *task) EveryYears(count int) *task {
	t.SetInterval(time.Duration(count*365*24) * time.Hour)
	return t
}
