package main

import (
	"bytes"
	"github.com/lobiCode/3fs_2/qman"
	"net"
	"time"
)

const (
	N = 10
)

var (
	FibQ Queue
	AriQ Queue
	EncQ Queue
	RevQ Queue

	QM Queue
)

func main() {

	ln, err := net.Listen("tcp", ":1234")

	if err != nil {
		// TODO
	}

	start()

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue

		}
		createJob(conn)
	}

}

func start() {

	FibQ = Queue{Q: make([]qman.Job, 0, 0)}
	AriQ = Queue{Q: make([]qman.Job, 0, 0)}
	EncQ = Queue{Q: make([]qman.Job, 0, 0)}
	RevQ = Queue{Q: make([]qman.Job, 0, 0)}

	for i := 0; i < N; i++ {
		fibw := CreateWorker(&FibQ)
		fibw.Start()
		ariw := CreateWorker(&AriQ)
		ariw.Start()
		encw := CreateWorker(&EncQ)
		encw.Start()
		revw := CreateWorker(&RevQ)
		revw.Start()
	}

	QM = Queue{Q: make([]qman.Job, 0, 0)}
	go jobDispatcher(&QM)
}

func jobDispatcher(q *Queue) {

	for {
		j := q.GetJob()
		if j != nil {
			switch j.Res.(type) {
			case qman.Fibonacci:
				sendJob(j, &FibQ)
			case qman.TextEncoder:
				sendJob(j, &EncQ)
			case qman.BasicArithmetic:
				sendJob(j, &AriQ)
			case qman.ReverseText:
				sendJob(j, &RevQ)
			}
		} else {
			time.Sleep(time.Millisecond)
		}
	}
}

func sendJob(j *qman.Job, q *Queue) {

	q.PushJob(*j)
}

func createJob(conn net.Conn) {

	var j qman.Job
	b := make([]byte, 512)
	conn.Read(b) // TODO len n,

	req := bytes.Split(b, []byte{'\n'})
	req = bytes.Split(req[0], []byte{' '})
	if len(req) < 2 {
		conn.Write([]byte("Wrong usage\n"))
		conn.Close()
		return
	}

	switch {
	case bytes.Compare(req[0], qman.FIBONACCI) == 0:
		j = qman.Job{qman.Fibonacci{req[1], conn}}
	case bytes.Compare(req[0], qman.REVERSE) == 0:
		j = qman.Job{qman.ReverseText{req[1], conn}}
	case bytes.Compare(req[0], qman.ENCODER) == 0:
		j = qman.Job{qman.TextEncoder{req[1], conn}}
	case bytes.Compare(req[0], qman.ARITHMETIC) == 0:
		j = qman.Job{qman.BasicArithmetic{req[1], conn}}
	default:
		b = []byte("is unknown resolver")
		req = append(req, b)
		conn.Write(bytes.Join(req, []byte{' '}))
		conn.Close()
		return
	}

	QM.PushJob(j)
}
