package main

import (
	"net"
)

type Resolver interface {
	R(conn net.Conn)
}

type Job struct {
	Conn net.Conn
	Res  Resolver
}

type Worker struct {
	JobQueue chan Job
}

func CreateWorker(jobQueue chan Job) *Worker {
	return &Worker{jobQueue}
}

func (w *Worker) Start() {

	go func() {
		for {
			select {
			case job := <-w.JobQueue:
				job.Res.R(job.Conn)
			}
		}
	}()
}

type Fibonacci int

func (f Fibonacci) R(conn net.Conn) {

	conn.Write([]byte("Closed\n"))
	conn.Close()
}
