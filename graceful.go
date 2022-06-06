package gdown

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Priority int

const (
	PriorityLevel0  Priority = 0
	PriorityLevel10 Priority = 10
	PriorityLevel20 Priority = 20
	PriorityLevel30 Priority = 30
	PriorityLevel40 Priority = 40
	PriorityLevel50 Priority = 50
	PriorityLevel60 Priority = 60
	PriorityLevel70 Priority = 70
	PriorityLevel80 Priority = 80
	PriorityLevel90 Priority = 90
)

type Graceful interface {
	Priorities() []Priority
	Register(p Priority, f func())
	Listen()
}

// Graceful
// 平滑关闭
type graceful struct {
	ctx     context.Context
	cancel  context.CancelFunc
	wg      *sync.WaitGroup
	lock    sync.Mutex
	jobs    map[Priority][]func()
	signals []os.Signal
}

func (g *graceful) Priorities() []Priority {
	return []Priority{PriorityLevel0, PriorityLevel10, PriorityLevel20, PriorityLevel30, PriorityLevel40, PriorityLevel50, PriorityLevel60, PriorityLevel70, PriorityLevel80, PriorityLevel90}
}

// Register
// 回调注册
func (g *graceful) Register(p Priority, f func()) {
	g.lock.Lock()
	defer g.lock.Unlock()
	g.wg.Add(1)
	if _, ok := g.jobs[p]; !ok {
		g.jobs[p] = []func(){}
	}
	g.jobs[p] = append(g.jobs[p], f)
}

// Listen
// 信号监听
func (g *graceful) Listen() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, g.signals...)

	go func() {
		<-sig
		g.cancel()
		for _, p := range g.Priorities() {
			g.lock.Lock()
			if _, ok := g.jobs[p]; !ok {
				g.lock.Unlock()
				continue
			}
			for _, job := range g.jobs[p] {
				job()
				g.wg.Done()
			}
			g.lock.Unlock()
			time.Sleep(300 * time.Millisecond)
		}
		// 确保所有的回调都被执行
		g.wg.Done()
	}()
	g.wg.Wait()
}

func NewGraceful(ctx context.Context, cancel context.CancelFunc, signals ...os.Signal) Graceful {
	if len(signals) == 0 || signals == nil {
		signals = []os.Signal{syscall.SIGTERM, syscall.SIGINT}
	}
	shutdown := &graceful{
		ctx:     ctx,
		cancel:  cancel,
		wg:      &sync.WaitGroup{},
		lock:    sync.Mutex{},
		jobs:    make(map[Priority][]func()),
		signals: signals,
	}
	// 确保所有的回调都被执行
	shutdown.wg.Add(1)
	return shutdown
}
