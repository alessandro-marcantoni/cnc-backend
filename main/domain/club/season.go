package club

import "time"

type Season struct {
	ID       int64
	Code     string
	Name     string
	StartsAt time.Time
	EndsAt   time.Time
}

func (s Season) GetCode() string {
	return s.Code
}
