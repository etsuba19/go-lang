package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Task struct{
	id int
	data string
}

func worker(id int, taskChan <-chan Task, wg *sync.WaitGroup, failChan chan<-Task){
	defer wg.Done()

	for task := range taskChan{
		fmt.Printf("Worker %d started task %d: %s\n", id, task.id, task.data)

		if rand.Float32() < 0.3 {
			fmt.Printf("Worker %d failed on task %d\n", id, task.id)
			failChan <- task
			return
		}

		time.Sleep(time.Duration(rand.Intn(3)+1)*time.Second)
		fmt.Printf("Worker %d completed task %d\n", id, task.id)
	}
}

func main(){

	rand.Seed(time.Now().UnixNano())

	tasks := []Task{
		{id: 1, data: "Task 1"},
		{id: 2, data: "Task 2"},
		{id: 3, data: "Task 3"},
		{id: 4, data: "Task 4"},
		{id: 5, data: "Task 5"},
	}
	// taskChan := make(chan Task, len(tasks))
	// failChan := make(chan Task, len(tasks))

	// var wg sync.WaitGroup

	numWorkers := 3
	maxRetries := 2


	workerChannels := make([]chan Task, numWorkers)
	for i := range workerChannels {
		workerChannels[i] = make(chan Task, 5)
	}

	failChan := make(chan Task, len(tasks))

	var wg sync.WaitGroup
	retryCount := make(map[int]int)

	for i := 0; i <  numWorkers; i++ {
		wg.Add(1)
		go worker(i+1, workerChannels[i], &wg, failChan)
	}

	workerIndex := 0

	for _, task := range tasks{
		fmt.Printf("[Scheduler] Assigning Task %d → Worker %d\n", task.id, workerIndex+1)
		workerChannels[workerIndex] <- task

		workerIndex = (workerIndex + 1) % numWorkers
	}
	

	go func ()  {
		for task := range failChan{
			retryCount[task.id]++

			if retryCount[task.id] > maxRetries{
				fmt.Printf("[Scheduler] Task %d reached max retries. Marking as FAILED.\n", task.id)
				continue
			}

			w := rand.Intn(numWorkers)
			fmt.Printf("[Scheduler] Retrying Task %d to Worker %d (attempt %d)\n", task.id, w+1, retryCount[task.id])
			workerChannels[w] <- task
		}
	}()

	time.Sleep(8 * time.Second)

	for _, ch := range workerChannels {
		close(ch)
	}

	wg.Wait()
	close(failChan)

	fmt.Println("All workers finished. Program terminated.")
} 