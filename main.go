package main

import (
	"fmt"
	"os"
	"sync"
)

type Worker struct {
	id         int
	jobs       chan string
	workerPool chan chan string
	wg         *sync.WaitGroup
}

type WorkerPool struct {
	workers    []Worker
	workerPool chan chan string
	wg         sync.WaitGroup
}

func NewWorker(id int, workerPool chan chan string, wg *sync.WaitGroup) *Worker {
	worker := &Worker{
		id:         id,
		workerPool: workerPool,
		jobs:       make(chan string),
		wg:         wg,
	}
	go worker.start()
	return worker
}

func (w *Worker) start() {
	fmt.Println("Worker start: ", w.id)
	for {
		w.workerPool <- w.jobs
		select {
		case job := <-w.jobs:
			fmt.Println(w.id, " ", job)
			w.wg.Done()
		}
	}
}

func NewWorkerPool(numWorkers int) *WorkerPool {
	workerPool := &WorkerPool{
		workers:    make([]Worker, 0),
		workerPool: make(chan chan string),
	}
	for i := 0; i < numWorkers; i++ {
		worker := NewWorker(i, workerPool.workerPool, &workerPool.wg)
		workerPool.workers = append(workerPool.workers, *worker)
	}
	return workerPool
}

func (wp *WorkerPool) RemoveWorker() {
	if len(wp.workers) == 0 {
		fmt.Println("No workers to remove")
		return
	}
	lastWorker := wp.workers[len(wp.workers)-1]
	wp.workers = wp.workers[:len(wp.workers)-1]
	close(lastWorker.jobs)
	fmt.Println("Worker removed: ", lastWorker.id)
}

func (wp *WorkerPool) AddWorker() {
	worker := NewWorker(len(wp.workers), wp.workerPool, &wp.wg)
	wp.workers = append(wp.workers, *worker)
}

func (wp *WorkerPool) AddJob(job string) {
	wp.wg.Add(1)
	jobChannel := <-wp.workerPool
	jobChannel <- job
}

func (wp *WorkerPool) CountWorker() int {
	return len(wp.workers)
}

func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}

func main() {
	pool := NewWorkerPool(0)
	var choice int
	var job string
	finish := false
	for !finish {
		fmt.Println("Работников всего ", pool.CountWorker())
		fmt.Println("Для выхода из программы введите 0")
		fmt.Println("Для добавления работника введите 1")
		fmt.Println("Для удаления работника введите 2")
		fmt.Println("Для выдачи задачи введите 3")
		fmt.Scan(&choice)
		switch choice {
		case 0:
			finish = true
		case 1:
			pool.AddWorker()
		case 2:
			pool.RemoveWorker()
		case 3:
			fmt.Println("Введите сообщение: ")
			fmt.Fscan(os.Stdin, &job)
			pool.AddJob(job)
		}
	}
	pool.Wait()
}
