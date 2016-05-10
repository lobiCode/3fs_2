package main

import (
	"github.com/lobiCode/3fs_2/qman"
	"sync"
	"time"
)

type Queue struct {
	Q  []qman.Job
	Mu sync.Mutex
}

func (q *Queue) GetJob() *qman.Job {

	q.Mu.Lock()
	defer q.Mu.Unlock()

	if len(q.Q) == 0 {
		return nil
	}

	j := q.Q[0]

	q.Q = append(q.Q[:0], q.Q[1:]...)

	return &j
}

func (q *Queue) PushJob(j qman.Job) {

	q.Mu.Lock()
	q.Q = append(q.Q, j)
	q.Mu.Unlock()
}

type Worker struct {
	Jq *Queue
}

func CreateWorker(q *Queue) Worker {
	return Worker{q}
}

func (w *Worker) Start() {

	go func() {
		for {
			j := w.Jq.GetJob()
			if j != nil {
				j.Res.R()
			} else {
				time.Sleep(time.Millisecond)
			}
		}
	}()
}
