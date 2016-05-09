package main

import (
	"bytes"
	"github.com/lobiCode/3fs_2/qman"
	"net"
)

const (
	N      = 4
	BUFFER = 100
)

var (
	JobQueue chan qman.Job
)

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
		dispatcher(conn)
	}

}

func Init() {

	JobQueue = make(chan qman.Job, BUFFER)

	for i := 0; i < N; i++ {
		w := qman.CreateWorker(JobQueue)
		w.Start()
	}
}

func dispatcher(conn net.Conn) {

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
	JobQueue <- j
}
