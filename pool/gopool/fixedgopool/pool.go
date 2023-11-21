package fixedgopool

import (
	"context"
	"fmt"
)

type GoPool struct {
	workerNum    uint32
	jobQueueSize uint32
	jobs         chan JobI
	stop         context.CancelFunc
}

func NewPool(opts ...OpFunc) *GoPool {
	p := new(GoPool)
	for _, opt := range opts {
		opt(p)
	}
	ctx, cancel := context.WithCancel(context.Background())

	p.jobs = make(chan JobI, p.jobQueueSize)
	p.stop = cancel

	for i := uint32(0); i < p.workerNum; i++ {
		go p.worker(ctx)
	}
	return p
}

func (p *GoPool) worker(ctx context.Context) {
	for {
		select {
		case job := <-p.jobs:
			job.Run()
		case <-ctx.Done():
			// 停止
			fmt.Println("worker停止了...")
			return
		}
	}
}

func (p *GoPool) AddJob(job JobI) {
	p.jobs <- job
}

func (p *GoPool) Stop() {
	p.stop()
	close(p.jobs)
}
