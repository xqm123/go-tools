package fixedgopool

import (
	"context"
	"fmt"
	"sync"
)

type GoPool struct {
	workerNum    uint32
	jobQueueSize uint32

	jobs   chan Func
	jobsWg sync.WaitGroup

	wockerWg  sync.WaitGroup
	closed    bool
	closeFunc context.CancelFunc
}

func NewPool(opts ...OpFunc) *GoPool {
	p := new(GoPool)
	for _, opt := range opts {
		opt(p)
	}
	ctx, cancel := context.WithCancel(context.Background())

	p.jobs = make(chan Func, p.jobQueueSize)
	p.closeFunc = cancel

	for i := uint32(0); i < p.workerNum; i++ {
		p.wockerWg.Add(1)
		go p.worker(ctx)
	}
	return p
}

func (p *GoPool) AddJob(job Func) {
	p.addJob(job)
}

func (p *GoPool) WaitJobs() {
	p.waitJobs()
}

func (p *GoPool) Close() {
	p.close()
}

func (p *GoPool) worker(ctx context.Context) {
	defer p.wockerWg.Done()
	for {
		select {
		case job := <-p.jobs:
			func() {
				defer func() {
					p.jobsWg.Done()
					if r := recover(); r != nil {
						fmt.Println(r)
					}
				}()
				job()
			}()
		case <-ctx.Done():
			// 停止
			//fmt.Println("worker stop...")
			return
		}
	}
}

func (p *GoPool) addJob(job Func) {
	if p.closed {
		return
	}
	defer p.jobsWg.Add(1)
	p.jobs <- job

}

func (p *GoPool) waitJobs() {
	p.jobsWg.Wait()
}

func (p *GoPool) close() {
	p.closed = true
	p.closeFunc()
	p.wockerWg.Wait()
	close(p.jobs)
}
