package main

import (
	"fmt"
	"sync"
	"time"
)

type Worker struct {
	id      int
	jobChan chan string
	wg      *sync.WaitGroup
}

func NewWorker(id int, wg *sync.WaitGroup) *Worker {
	w := &Worker{
		id:      id,
		jobChan: make(chan string),
		wg:      wg,
	}
	go w.start()
	return w
}



func (w *Worker) start() {
	defer w.wg.Done()
	for job := range w.jobChan {

		fmt.Printf("Worker %d processing job: %s\n", w.id, job)
		time.Sleep(time.Second) //Время обработки
	}
}


type WorkerPool struct {
	workers []*Worker
	jobChan chan string
	wg      sync.WaitGroup
}


func NewWorkerPool(size int) *WorkerPool {
	pool := &WorkerPool{
		workers: make([]*Worker, 0, size),
		jobChan: make(chan string),
	}
	for i := 0; i < size; i++ {
		pool.addWorker(i + 1)
	}
	return pool
}

func (p *WorkerPool) addWorker(id int) {
	p.wg.Add(1)
	worker := NewWorker(id, &p.wg)
	p.workers = append(p.workers, worker)
}


func (p *WorkerPool) removeWorker() {
	if len(p.workers) == 0 {
		return // если нет ничего не делаем
	}
	worker := p.workers[len(p.workers)-1]
	close(worker.jobChan)
	p.workers = p.workers[:len(p.workers)-1]
}


func (p *WorkerPool) Start() {
	go func() {
		for job := range p.jobChan {
			if len(p.workers) > 0 {
				p.workers[0].jobChan <- job
			}
		}
	}()
}


func (p *WorkerPool) AddJob(job string) {
	p.jobChan <- job
}


func (p *WorkerPool) Stop() {
	close(p.jobChan)
	p.wg.Wait()
}

func main() {
	pool := NewWorkerPool(3)
	pool.Start()

	for i := 1; i <= 10; i++ {
		pool.AddJob(fmt.Sprintf("Job %d", i))
	}
	pool.addWorker(4)

	for i := 11; i <= 15; i++ {
		pool.AddJob(fmt.Sprintf("Job %d", i))
	}
	pool.removeWorker()
	pool.Stop()
}
