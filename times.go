package scheduler

import "time"

func (s *scheduler) EverySecond() *scheduler {
	s.SetInterval(1 * time.Second)
	return s
}

func (s *scheduler) EverySeconds(t int) *scheduler {
	s.SetInterval(time.Duration(t) * time.Second)
	return s
}

func (s *scheduler) EveryMinute() *scheduler {
	s.SetInterval(1 * time.Minute)
	return s
}

func (s *scheduler) EveryMinutes(t int) *scheduler {
	s.SetInterval(time.Duration(t) * time.Minute)
	return s
}

func (s *scheduler) EveryHour() *scheduler {
	s.SetInterval(1 * time.Hour)
	return s
}

func (s *scheduler) EveryHours(t int) *scheduler {
	s.SetInterval(time.Duration(t) * time.Hour)
	return s
}

func (s *scheduler) EveryDay() *scheduler {
	s.SetInterval(24 * time.Hour)
	return s
}

func (s *scheduler) EveryDays(t int) *scheduler {
	s.SetInterval(time.Duration(t*24) * time.Hour)
	return s
}

func (s *scheduler) EveryWeek() *scheduler {
	s.SetInterval(24 * 7 * time.Hour)
	return s
}

func (s *scheduler) EveryWeeks(t int) *scheduler {
	s.SetInterval(time.Duration(t*24*7) * time.Hour)
	return s
}

func (s *scheduler) EveryMonth() *scheduler {
	s.SetInterval(30 * 24 * time.Hour)
	return s
}

func (s *scheduler) EveryMonths(t int) *scheduler {
	s.SetInterval(time.Duration(t*24*30) * time.Hour)
	return s
}

func (s *scheduler) EveryYear() *scheduler {
	s.SetInterval(365 * 24 * time.Hour)
	return s
}

func (s *scheduler) EveryYears(t int) *scheduler {
	s.SetInterval(time.Duration(t*365*24) * time.Hour)
	return s
}
