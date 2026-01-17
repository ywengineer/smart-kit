package types

import "time"

func MustParseDuration(d string) time.Duration {
	dur, err := time.ParseDuration(d)
	if err != nil {
		panic(err)
	}
	return dur
}
