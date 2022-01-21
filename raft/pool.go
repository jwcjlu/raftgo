package raft

import (
	"errors"
	"github.com/sirupsen/logrus"
	"io"
	"sync"
	"time"
)

// Pool 管理一组安全在多个goroutine间共享资源，被管理资源必须实现io.Closer接口
type Pool struct {
	m         sync.Mutex
	resources chan io.Closer            //通投类型为接口,可管理任意实现io.Closer接口的资源类型
	factory   func() (io.Closer, error) //创建新资源，由使用者提供
	closed    bool
}

var ErrPoolClosed = errors.New("Pool has been closed")

// New 函数工厂，指定有缓冲通道大小
func New(fn func() (io.Closer, error), size uint) (*Pool, error) {
	if size <= 0 {
		return nil, errors.New("Size value negative")
	}
	return &Pool{
		factory:   fn,
		resources: make(chan io.Closer, size),
	}, nil
}

// Acquire 池中获取资源
func (p *Pool) Acquire(ms time.Duration) (io.Closer, error) {
	timeout := time.After(ms)
	select {
	case r, ok := <-p.resources:
		logrus.Error("Acquire: Shared Resource")
		if !ok {
			return nil, ErrPoolClosed
		}
		return r, nil
	case <-timeout:
		logrus.Error("Acquire: New Resource")
		return p.factory()
	}
}

// Release 池中释放资源
func (p *Pool) Release(r io.Closer) {
	p.m.Lock()
	defer p.m.Unlock()

	if p.closed {
		r.Close()
		return
	}

	select {
	case p.resources <- r: //放入队列
		logrus.Println("Release: In Queue")
	default: //队列已满,则关闭
		logrus.Println("Release: Closing")
		r.Close()
	}
}

// Close 池关闭所有现有资源
func (p *Pool) Close() {
	p.m.Lock()
	defer p.m.Unlock()

	if p.closed {
		return
	}

	p.closed = true

	close(p.resources) //清空通道资源前将通道关闭，否则会产生死锁

	for r := range p.resources {
		r.Close()
	}
}
