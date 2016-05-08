package main

import (
	"net"
)

const (
	N      = 4
	BUFFER = 100
)

var JobQueue chan Job

func main() {

	ln, err := net.Listen("tcp", ":1234")

	if err != nil {
		// TODO
	}

	Init()

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue

		}
		go dispatchers(conn)
	}

}

func Init() {
	JobQueue = make(chan Job, BUFFER)

	for i := 0; i < N; i++ {
		w := CreateWorker(JobQueue)
		w.Start()
	}
}

func dispatchers(conn net.Conn) {
	j := Job{conn, new(Fibonacci)}
	JobQueue <- j
}
