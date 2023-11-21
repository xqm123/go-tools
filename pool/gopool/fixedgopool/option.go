package fixedgopool

type OpFunc func(p *GoPool)

func WithWorkNum(workerNum uint32) OpFunc {
	return func(p *GoPool) {
		p.workerNum = workerNum
	}
}

func WithJobQueueSize(jobQueueSize uint32) OpFunc {
	return func(p *GoPool) {
		p.jobQueueSize = jobQueueSize
	}
}
