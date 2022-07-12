package workerPool

type WorkerPoll interface {
	AddJob(func()) 
	Run()
}

type pool struct {
	maxWokers int
	queue chan func()
}

func NewWorkerPool(maxWokers int, queueSize int)  WorkerPoll {
	return &pool{maxWokers: maxWokers, queue: make(chan func(),queueSize )}
}

func (p *pool) AddJob(f func()) {
	p.queue <- f
}

func (p *pool) Run() {

	for i:= 1; i <= p.maxWokers; i++ {
		
		go func ()  {
			for job := range p.queue {
				job()
			}
		}()

	}

}