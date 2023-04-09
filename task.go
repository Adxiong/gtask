/*
 * @Description:
 * @version:
 * @Author: Adxiong
 * @Date: 2023-04-09 00:21:39
 * @LastEditors: Adxiong
 * @LastEditTime: 2023-04-09 17:41:56
 */
package gtask

import (
	"context"
	"fmt"
	"sync"
)

type task struct {
	ctx             context.Context
	cancel          context.CancelFunc
	g               sync.WaitGroup
	err             error
	allowPartFailed bool // 允许部分失败
	once            sync.Once
}

func NewTask(ctx context.Context, cancel context.CancelFunc, allowPartFailed bool) *task {
	return &task{
		ctx:             ctx,
		cancel:          cancel,
		g:               sync.WaitGroup{},
		allowPartFailed: allowPartFailed,
		once:            sync.Once{},
	}
}

// Do 执行任务
func (t *task) Do(f func() error) {
	t.g.Add(1)
	go func() {
		defer func() {
			fmt.Println("关闭")
			t.g.Done()
		}()
		execFuncErr := f()

		if execFuncErr != nil {
			t.once.Do(func() {
				t.err = execFuncErr // 记录下错误
				// 不允许部分错误，携程任务都会取消
				if !t.allowPartFailed && t.cancel != nil {
					t.cancel()
				}
			})
		}
	}()
}

// Wait 等待任务执行完成
func (t *task) Wait() error {
	go func() {
		t.g.Wait()
		t.cancel()
	}()

	select {
	case <-t.ctx.Done():
		fmt.Println("receive cancel")
		return t.err
	}
}
