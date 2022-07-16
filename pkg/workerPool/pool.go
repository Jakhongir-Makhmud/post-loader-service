package workerPool

import (
	"context"
	"fmt"
)

type WorkerPoll interface {
	AddJob(func())
	Run(context.Context)
}

type pool struct {
	maxWokers int
	queue     chan func()
}

func NewWorkerPool(maxWokers int, queueSize int) WorkerPoll {
	return &pool{maxWokers: maxWokers, queue: make(chan func(), queueSize)}
}

func (p *pool) AddJob(f func()) {
	p.queue <- f
}

func (p *pool) Run(ctxFinish context.Context) {

	go func() {
		select {
		case <-ctxFinish.Done():
			close(p.queue) // if chanel closed workers will exit.
			fmt.Println("workers are done treir job.")
		}
	}()

	for i := 1; i <= p.maxWokers; i++ {

		go func() {
			for job := range p.queue {
				job()
			}
		}()

	}

}
