package waitgroup

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

func counter(gId int, sharedCount *int, amount int, delayAmount time.Duration,  group WaitGroup, mutex *sync.Mutex) {

	//Critical section use synchronization before proceeding
	mutex.Lock()
	for i := 0; i < amount; i++ {
		*sharedCount = *sharedCount  + 1
		if delayAmount > 0 {
			time.Sleep(delayAmount)
		}
	}
	mutex.Unlock()
	group.Done()
}

func Test1(t *testing.T) {

	spawnAmount := 10
	iterationAmount := 10
	var sharedCount int

	// Create the WaitGroup
	group := NewWaitGroup()
	var mutex sync.Mutex
	for i:=0; i < spawnAmount; i++ {
		group.Add(1)
		go counter(i,&sharedCount,iterationAmount,0,group,&mutex)
	}
	group.Wait()

	if sharedCount != iterationAmount*spawnAmount  {
		t.Errorf("Shared counter value incorrect.\ngot:%v\nexpected:%v", sharedCount, 100)
	}
}
func Test2(t *testing.T) {

	spawnAmount := 10
	iterationAmount := 10
	var sharedCount int

	// Create the WaitGroup
	group := NewWaitGroup()

	group.Add(10)

	var mutex sync.Mutex
	for i:=0; i < spawnAmount; i++ {
		go counter(i,&sharedCount,iterationAmount,0,group,&mutex)
	}
	group.Wait()

	if sharedCount != iterationAmount*spawnAmount  {
		t.Errorf("Shared counter value incorrect.\ngot:%v\nexpected:%v", sharedCount, 100)
	}
}

func Test3(t *testing.T) {

	spawnAmount := 10
	iterationAmount := 10
	var sharedCount int

	// Create the WaitGroup
	group := NewWaitGroup()

	group.Add(10)
	var mutex sync.Mutex

	for i:=0; i < spawnAmount; i++ {
		delayAmount := time.Duration(rand.Int31n(250))
		go counter(i,&sharedCount,iterationAmount,delayAmount,group,&mutex)
	}
	group.Wait()

	if sharedCount != iterationAmount*spawnAmount   {
		t.Errorf("Shared counter value incorrect.\ngot:%v\nexpected:%v", sharedCount, 100)
	}
}

func Test4(t *testing.T) {

	spawnAmount := 10000
	iterationAmount := 10
	var sharedCount int

	// Create the WaitGroup
	group := NewWaitGroup()

	var mutex sync.Mutex

	for i:=0; i < spawnAmount; i++ {
		group.Add(1)
		delayAmount := time.Duration(rand.Int31n(250))
		go counter(i,&sharedCount,iterationAmount,delayAmount,group,&mutex)
	}
	group.Wait()

	if sharedCount != iterationAmount*spawnAmount  {
		t.Errorf("Shared counter value incorrect.\ngot:%v\nexpected:%v", sharedCount, 100)
	}
}
func Test5(t *testing.T) {

	spawnAmount := 100
	iterationAmount := 10
	var sharedCount int

	// Create the WaitGroup
	group := NewWaitGroup()

	var mutex sync.Mutex

	for i:=0; i < spawnAmount; i++ {
		group.Add(1)
		delayAmount := time.Duration(rand.Int31n(250))
		go func() {
			group.Add(1)
			go counter(i,&sharedCount,iterationAmount,delayAmount,group,&mutex)
			group.Done()
		}()
	}
	group.Wait()

	if sharedCount != iterationAmount*spawnAmount  {
		t.Errorf("Shared counter value incorrect.\ngot:%v\nexpected:%v", sharedCount, 100)
	}
}
func Test6(t *testing.T) {

	spawnAmount := 100
	iterationAmount := 10
	var sharedCount int

	// Create the WaitGroup
	group := NewWaitGroup()

	var mutex sync.Mutex

	for i:=0; i < spawnAmount; i++ {
		group.Add(1)
		delayAmount := time.Duration(rand.Int31n(250))
		go func() {
			group.Add(1)
			go counter(i,&sharedCount,iterationAmount,delayAmount,group,&mutex)
			group.Done()
		}()
	}
	group.Wait()

	for i:=0; i < spawnAmount; i++ {
		group.Add(1)
		delayAmount := time.Duration(rand.Int31n(250))
		go func() {
			group.Add(1)
			go counter(i,&sharedCount,iterationAmount,delayAmount,group,&mutex)
			group.Done()
		}()
	}
	group.Wait()

	if sharedCount != (iterationAmount*spawnAmount) * 2  {
		t.Errorf("Shared counter value incorrect.\ngot:%v\nexpected:%v", sharedCount, 100)
	}
}


func TestWait(t *testing.T) {

	// Create the WaitGroup
	group := NewWaitGroup()

	// The main goroutine should exit immediately from Wait since no
	// goroutines are part of the group.
	group.Wait()
}