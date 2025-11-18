package gop

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/bytedance/gopkg/util/gopool"
)

// 状态枚举：Promise 的三种状态
type state int

const (
	pending   state = iota // 等待中
	fulfilled              // 完成
	rejected               // 失败
)

// Promise 接口：定义异步任务的核心能力
type Promise interface {
	// Then 任务成功后执行回调，返回新的 Promise（链式调用）
	Then(fn func(interface{}) (interface{}, error)) Promise
	// Catch 捕获任务执行过程中的错误
	Catch(fn func(error)) Promise
	// Await 阻塞等待任务完成，返回结果或错误
	Await() (interface{}, error)
	// WithTimeout 设置任务超时时间
	WithTimeout(timeout time.Duration) Promise
}

type pr struct {
	index  int
	result interface{}
}

// NewPromise 创建一个新的 Promise，提交任务到 gopool
func NewPromise(task func() (interface{}, error)) Promise {
	ctx, cancel := context.WithCancel(context.Background())
	p := &promise{
		state:  pending,
		ctx:    ctx,
		cancel: cancel,
	}

	// 提交任务到 gopool（协程池限制并发）
	gopool.Go(func() {
		defer func() {
			// 捕获 panic，转为错误
			if r := recover(); r != nil {
				p.reject(fmt.Errorf("task panicked: %v", r))
			}
		}()

		select {
		case <-p.ctx.Done():
			// 任务被取消（超时/主动取消）
			p.reject(p.ctx.Err())
			return
		default:
			// 执行任务
			result, err := task()
			if err != nil {
				p.reject(err)
			} else {
				p.fulfill(result)
			}
		}
	})

	return p
}

// All 等待所有 Promise 完成，返回结果切片（有一个失败则整体失败）
func All(ttl time.Duration, promises ...Promise) Promise {
	return NewPromise(func() (interface{}, error) {
		var wg sync.WaitGroup
		resultChan := make(chan *pr, len(promises))
		errChan := make(chan error, 1) // 只接收第一个错误

		wg.Add(len(promises))
		for index, p := range promises {
			go func(index int, p Promise) {
				defer wg.Done()
				res, err := p.Await()
				if err != nil {
					// 非阻塞发送错误（避免多个错误阻塞）
					errChan <- err
					return
				} else {
					resultChan <- &pr{index, res}
				}
			}(index, p)
		}

		// 等待所有任务完成或第一个错误出现
		go func() {
			wg.Wait()
			close(resultChan)
			close(errChan)
		}()

		// 优先处理错误
		select {
		case err := <-errChan:
			return nil, err
		case <-time.After(ttl):
			return nil, errors.New("all promises timeout")
		default:
			// 收集所有结果
			results := make([]interface{}, 0, len(promises))
			for res := range resultChan {
				results[res.index] = res.result
			}
			return results, nil
		}
	})
}

// Race 等待第一个完成的 Promise，返回其结果或错误
func Race(ttl time.Duration, promises ...Promise) Promise {
	return NewPromise(func() (interface{}, error) {
		resultChan := make(chan interface{}, 1)
		errChan := make(chan error, 1)

		for _, p := range promises {
			go func(p Promise) {
				res, err := p.Await()
				if err != nil {
					errChan <- err
				} else {
					resultChan <- res
				}
			}(p)
		}

		// 取第一个完成的结果
		select {
		case res := <-resultChan:
			return res, nil
		case err := <-errChan:
			return nil, err
		case <-time.After(ttl):
			return nil, errors.New("race promises timeout")
		}
	})
}
