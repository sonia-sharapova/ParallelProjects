package server

import (
	"encoding/json"
	"proj2/feed"
	"proj2/queue"
	"sync"
)

type Config struct {
	Encoder *json.Encoder // Represents the buffer to encode Responses
	Decoder *json.Decoder // Represents the buffer to decode Requests
	Mode    string        // Represents whether the server should execute
	// sequentially or in parallel
	// If Mode == "s"  then run the sequential version
	// If Mode == "p"  then run the parallel version
	// These are the only values for Version
	ConsumersCount int // Represents the number of consumers to spawn
}

type Response struct {
	ID      int         	// json:"id"
	Success bool        	// json:"success,omitempty"
	Feed    interface{} 	// json:"feed,omitempty"
}

// Run starts up the twitter server based on the configuration
// information provided and only returns when the server is fully
// shutdown.
func Run(config Config) {
	// Run sequential
	if config.Mode == "s" {
		runSequential(config)
	// Run parallel
	} else if config.Mode == "p" {
		runParallel(config)
	}
}

// Sequential execution of requests
func runSequential(config Config) {
	twitterFeed := feed.NewFeed()		// Initialize the feed

	for {
		var task queue.Request
		// Decode the next task
		if err := config.Decoder.Decode(&task); err != nil {
			break // Exit loop on EOF or error
		}

		// Check termination command
		if task.Command == "DONE" {
			break
		}

		handleTask(&task, twitterFeed, config.Encoder)
	}
}

// Parallel execution of requests
func runParallel(config Config) {
	twitterFeed := feed.NewFeed()			// Init the feed
	taskQueue := queue.NewLockFreeQueue()	// Init lock-free queue
	cond := sync.NewCond(&sync.Mutex{})
	done := false
	wg := sync.WaitGroup{}

	// Spawn consumers (goroutines)
	for i := 0; i < config.ConsumersCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()	// Mark this goroutine as done
			consumer(twitterFeed, taskQueue, cond, &done, config.Encoder)
		}()
	}

	// Producer logic
	producer(taskQueue, cond, &done, config.Decoder)

	// Notify consumers to exit
	cond.L.Lock()
	done = true
	cond.Broadcast()		// Wake up all waiting consumers
	cond.L.Unlock()

	// Wait for all consumers to finish
	wg.Wait()
}

// Producer reads tasks and adds them to the queue
func producer(taskQueue *queue.LockFreeQueue, cond *sync.Cond, done *bool, decoder *json.Decoder) {
	for {
		var task queue.Request
		if err := decoder.Decode(&task); err != nil {
			break // Exit loop on EOF or error
		}

		if task.Command == "DONE" {
			break
		}

		// Enqueue the task and signal a waiting consumer
		taskQueue.Enqueue(&task)
		cond.Signal()
	}
}

// Consumer processes tasks from the queue
func consumer(twitterFeed feed.Feed, taskQueue *queue.LockFreeQueue, cond *sync.Cond, done *bool, encoder *json.Encoder) {
	for {
		cond.L.Lock()
		for taskQueue.IsEmpty() && !*done {
			cond.Wait()
		}
		// Exit if no more tasks and done is true
		if *done && taskQueue.IsEmpty() {
			cond.L.Unlock()
			return 
		}
		// Dequeue a task
		task, ok := taskQueue.Dequeue()
		cond.L.Unlock()

		// Process task (if successfully dequeued)
		if ok {
			handleTask(task, twitterFeed, encoder)
		}
	}
}

// Handle a single task
func handleTask(task *queue.Request, twitterFeed feed.Feed, encoder *json.Encoder) {
	switch task.Command {
	case "ADD":
		twitterFeed.Add(task.Body, task.Timestamp)
		encoder.Encode(Response{Success: true, ID: task.ID})
	case "REMOVE":
		success := twitterFeed.Remove(task.Timestamp)
		encoder.Encode(Response{Success: success, ID: task.ID})
	case "CONTAINS":
		success := twitterFeed.Contains(task.Timestamp)
		encoder.Encode(Response{Success: success, ID: task.ID})
	case "FEED":
		posts := []map[string]interface{}{}

		for _, p := range twitterFeed.GetAllPosts() {
			posts = append(posts, map[string]interface{}{
				"body":      p.Body(),
				"timestamp": p.Timestamp(),
			})
		}
		encoder.Encode(Response{ID: task.ID, Feed: posts})
	}
}
