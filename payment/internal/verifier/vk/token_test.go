package vk

import (
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	t.Log(time.Now().UTC().Format("2006-01-02T15:04:05.000Z07:00"))
}
