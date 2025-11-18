package gop

import (
	"testing"
	"time"
)

func TestPromise(t *testing.T) {
	// 测试 NewPromise
	p := NewPromise(func() (interface{}, error) {
		return "success", nil
	})

	// 测试 Then 回调
	result, err := p.Then(func(val interface{}) (interface{}, error) {
		return val.(string) + " appended", nil
	}).Await()
	if err != nil {
		t.Errorf("Then callback failed: %v", err)
	}
	if result.(string) != "success appended" {
		t.Errorf("Then callback result mismatch: got %v, want %v", result, "success appended")
	}
}

func TestAll(t *testing.T) {
	// 测试 All 方法
	p1 := NewPromise(func() (interface{}, error) {
		return "success1", nil
	})
	p2 := NewPromise(func() (interface{}, error) {
		return "success2", nil
	})

	results, err := All(time.Second*2, p1, p2).Await()
	if err != nil {
		t.Errorf("All failed: %v", err)
	}
	rs := results.([]interface{})
	if len(rs) != 2 {
		t.Errorf("All result length mismatch: got %v, want %v", len(rs), 2)
	}
	if rs[0].(string) != "success1" || rs[1].(string) != "success2" {
		t.Errorf("All result values mismatch: got %v, want %v", rs, []string{"success1", "success2"})
	}
}

func TestRace(t *testing.T) {
	// 测试 Race 方法
	p1 := NewPromise(func() (interface{}, error) {
		time.Sleep(200 * time.Millisecond)
		return "success1", nil
	})
	p2 := NewPromise(func() (interface{}, error) {
		time.Sleep(100 * time.Millisecond)
		return "success2", nil
	})

	result, err := Race(time.Second*2, p1, p2).Await()
	if err != nil {
		t.Errorf("Race failed: %v", err)
	}
	if result.(string) != "success2" {
		t.Errorf("Race result mismatch: got %v, want %v", result, "success2")
	}
}

func TestSelect(t *testing.T) {
	errChan := make(chan error, 1)
	go func() {
		time.Sleep(5 * time.Second)
		close(errChan)
	}()
	// 优先处理错误
	select {
	case err := <-errChan:
		t.Logf("Select error: %v", err)
	case <-time.After(time.Second * 10): // 防止无限等待（可自定义超时）
		t.Errorf("Select timeout: %v", time.Second*10)
	}
	t.Logf("Select succeed")
}
