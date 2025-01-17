package ioc

import (
	"time"
)

func NewTimeDuration() time.Duration {
	return time.Minute * 5
}
