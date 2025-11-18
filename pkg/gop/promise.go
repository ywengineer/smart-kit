package gop

import (
	"context"
	"errors"
	"sync"
	"time"
)

// promise 实现类：封装状态、结果、错误和回调
type promise struct {
	mu        sync.Mutex                               // 并发安全锁
	state     state                                    // 当前状态
	result    interface{}                              // 任务成功结果
	err       error                                    // 任务失败错误
	onFulfill []func(interface{}) (interface{}, error) // 成功回调链
	onReject  []func(error)                            // 失败回调链
	ctx       context.Context                          // 上下文（用于超时/取消）
	cancel    context.CancelFunc                       // 取消函数
}

// fulfill：任务成功，触发回调链
func (p *promise) fulfill(result interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.state != pending {
		return // 状态已变更，忽略
	}

	p.state = fulfilled
	p.result = result

	// 执行所有成功回调（链式调用）
	for _, fn := range p.onFulfill {
		nextResult, err := fn(p.result)
		if err != nil {
			// 回调执行失败，转为 rejected 状态
			p.state = rejected
			p.err = err
			p.triggerReject()
			return
		}
		p.result = nextResult // 更新结果，传递给下一个回调
	}

	// 回调执行完毕，清空回调链
	p.onFulfill = nil
}

// reject：任务失败，触发错误回调
func (p *promise) reject(err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.state != pending {
		return // 状态已变更，忽略
	}

	p.state = rejected
	p.err = err
	p.triggerReject()
}

// triggerReject：执行所有错误回调
func (p *promise) triggerReject() {
	for _, fn := range p.onReject {
		fn(p.err)
	}
	p.onReject = nil // 清空回调链
}

// Then 添加成功回调，支持链式调用
func (p *promise) Then(fn func(interface{}) (interface{}, error)) Promise {
	p.mu.Lock()
	defer p.mu.Unlock()

	switch p.state {
	case pending:
		// 任务未完成，添加回调到链
		p.onFulfill = append(p.onFulfill, fn)
	case fulfilled:
		// 任务已完成，立即执行回调（同步）
		nextResult, err := fn(p.result)
		if err != nil {
			p.state = rejected
			p.err = err
			p.triggerReject()
		} else {
			p.result = nextResult
		}
	case rejected:
		// 任务已失败，忽略 Then 回调
	}

	return p // 返回自身，支持链式调用
}

// Catch 添加错误回调
func (p *promise) Catch(fn func(error)) Promise {
	p.mu.Lock()
	defer p.mu.Unlock()

	switch p.state {
	case pending:
		// 任务未完成，添加错误回调
		p.onReject = append(p.onReject, fn)
	case rejected:
		// 任务已失败，立即执行回调（同步）
		fn(p.err)
	case fulfilled:
		// 任务已成功，忽略 Catch 回调
	}

	return p // 返回自身，支持链式调用
}

// Await 阻塞等待任务完成，获取结果或错误
func (p *promise) Await() (interface{}, error) {
	// 已完成，直接返回
	if p.state != pending {
		return p.result, p.err
	}

	// 未完成，等待信号
	done := make(chan struct{})
	defer close(done)

	// 添加临时回调，触发信号
	p.Then(func(result interface{}) (interface{}, error) {
		done <- struct{}{}
		return result, nil
	}).Catch(func(err error) {
		done <- struct{}{}
	})

	// 等待任务完成或上下文取消
	select {
	case <-done:
		return p.result, p.err
	case <-p.ctx.Done():
		return nil, p.ctx.Err()
	}
}

// WithTimeout 设置任务超时时间
func (p *promise) WithTimeout(timeout time.Duration) Promise {
	// 创建带超时的上下文，替换原有上下文
	timeoutCtx, timeoutCancel := context.WithTimeout(p.ctx, timeout)
	defer timeoutCancel()

	// 监听超时信号
	go func() {
		select {
		case <-timeoutCtx.Done():
			// 超时，取消任务
			p.cancel()
			p.reject(errors.New("promise timeout"))
		case <-p.ctx.Done():
			// 任务已完成/取消，退出
			return
		}
	}()

	return p
}
