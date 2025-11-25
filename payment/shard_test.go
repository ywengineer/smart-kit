package main

import (
	"testing"

	"gitee.com/ywengineer/smart-kit/pkg/utilk"
)

func TestShard(t *testing.T) {
	transactionID := "10000015114"
	shards := 64
	t.Logf("shard = %02d", utilk.Hash(transactionID)%uint64(shards))
}
