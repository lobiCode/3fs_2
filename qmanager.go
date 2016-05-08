package main

import (
	_ "encoding/binary"
	"golang.org/x/crypto/bcrypt"
	"math"
	"net"
	"strconv"
)

const (
	R   = 2.236067977499789696409173668731276235440618359611525724270
	PHI = (1 + R) / 2
	Phi = (1 - R) / 2
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

type Fibonacci struct {
	B []byte
}

func (f Fibonacci) R(conn net.Conn) {

	n, err := strconv.Atoi(string(f.B))
	if err != nil {
		// TODO err msq, close, return
	}

	nfloat := float64(n)
	fib := (math.Pow(PHI, nfloat) - math.Pow(Phi, nfloat)) / R
	s := strconv.FormatFloat(fib, 'f', 0, 64)
	s = s + "\n"
	conn.Write([]byte(s))
	conn.Close()
}

type ReverseText struct {
	B []byte
}

func (rt ReverseText) R(conn net.Conn) {

	for i := 0; i < len(rt.B)/2; i++ {
		j := len(rt.B) - 1 - i
		rt.B[i], rt.B[j] = rt.B[j], rt.B[i]
	}

	rt.B = append(rt.B, '\n')
	conn.Write(rt.B)
	conn.Close()
}

type TextEncoder struct {
	B []byte
}

func (te TextEncoder) R(conn net.Conn) {

	b, err := bcrypt.GenerateFromPassword(te.B, bcrypt.DefaultCost)
	if err != nil {
		// TODO
	}

	conn.Write(b)
	conn.Close()

}
