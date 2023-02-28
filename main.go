package main

import (
	"fmt"
	"time"
)

type Task struct {
	id         int64
	cT         string
	fT         string
	succeed    bool
	taskResult string
}

func taskGenerator(a chan Task) {
	for {
		cT := time.Now().Format(time.RFC3339)
		if time.Now().Nanosecond()%2 > 0 {
			cT = "Some error occurred"
		}
		a <- Task{cT: cT, id: time.Now().Unix()}
	}
}

func taskWorker(t Task) Task {
	cT, err := time.Parse(time.RFC3339, t.cT)
	now := time.Now()
	if cT.After(now.Add(-20*time.Second)) && err == nil {
		t.succeed = true
		t.taskResult = "task has been succeed"
	} else {
		t.taskResult = "something went wrong"
	}
	t.fT = now.Format(time.RFC3339Nano)

	return t
}

const WORKERS = 10

func main() {
	mainChan := make(chan Task, 10)
	doneChan := make(chan Task)
	errorChan := make(chan Task)

	go taskGenerator(mainChan)

	for i := 0; i < WORKERS; i++ {
		go func() {
			t := <-mainChan
			finishedTask := taskWorker(t)
			if finishedTask.succeed {
				doneChan <- t
			} else {
				errorChan <- t
			}
		}()
	}

	var doneTasksIds []int64
	var errors []error
	go func() {
		for {
			select {

			case t := <-doneChan:
				doneTasksIds = append(doneTasksIds, t.id)
			case t := <-errorChan:
				errors = append(errors, fmt.Errorf("task id %d time %s, error %s", t.id, t.cT, t.taskResult))
			}
		}
	}()

	time.Sleep(5 * time.Second)

	println("Done tasks:")
	for _, id := range doneTasksIds {
		println(id)
	}

	println("Errors:")
	for _, e := range errors {
		println(e)
	}
}
