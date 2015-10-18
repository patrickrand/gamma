package agent

import (
	"time"
)

type Interval struct {
	Days         time.Duration `json:"days"`
	Hours        time.Duration `json:"hours"`
	Minutes      time.Duration `json:"minutes"`
	Seconds      time.Duration `json:"seconds"`
	Milliseconds time.Duration `json:"milliseconds"`
}

func NewInterval(d, h, m, s, ms time.Duration) *Interval {
	return &Interval{
		Days:         d,
		Hours:        h,
		Minutes:      m,
		Seconds:      s,
		Milliseconds: ms,
	}
}

func (i *Interval) Duration() (t time.Duration) {
	t += i.Days / time.Nanosecond
	t += i.Hours / time.Nanosecond
	t += i.Minutes / time.Nanosecond
	t += i.Seconds / time.Nanosecond
	t += i.Milliseconds / time.Nanosecond
	return
}
