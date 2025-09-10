package caches

import "testing"

func TestNumber(t *testing.T) {
	t.Logf("%.2f, %d", 1e7, 1<<30)
}
