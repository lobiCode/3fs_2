package qman

import (
	"github.com/soniah/evaler"
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

var (
	FIBONACCI  = []byte("Fibonacci")
	REVERSE    = []byte("ReverseText")
	ENCODER    = []byte("TextEncoder")
	ARITHMETIC = []byte("BasicArithmetic")
)

type Resolver interface {
	R()
}

type Job struct {
	Res Resolver
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
				job.Res.R()
			}
		}
	}()
}

type Fibonacci struct {
	B    []byte
	Conn net.Conn
}

func (f Fibonacci) R() {

	n, err := strconv.Atoi(string(f.B))
	if err != nil {
		// TODO err msq, close, return
	}

	nfloat := float64(n)
	fib := (math.Pow(PHI, nfloat) - math.Pow(Phi, nfloat)) / R
	s := strconv.FormatFloat(fib, 'f', 0, 64)
	s = s + "\n"
	f.Conn.Write([]byte(s))
	f.Conn.Close()
}

type ReverseText struct {
	B    []byte
	Conn net.Conn
}

func (rt ReverseText) R() {

	for i := 0; i < len(rt.B)/2; i++ {
		j := len(rt.B) - 1 - i
		rt.B[i], rt.B[j] = rt.B[j], rt.B[i]
	}

	rt.B = append(rt.B, '\n')
	rt.Conn.Write(rt.B)
	rt.Conn.Close()
}

type TextEncoder struct {
	B    []byte
	Conn net.Conn
}

func (te TextEncoder) R() {

	b, err := bcrypt.GenerateFromPassword(te.B, bcrypt.DefaultCost)
	if err != nil {
		// TODO
	}

	te.Conn.Write(b)
	te.Conn.Close()

}

type BasicArithmetic struct {
	B    []byte
	Conn net.Conn
}

func (ba BasicArithmetic) R() {

	// TODO najdi boljsi lib
	x, err := evaler.Eval(string(ba.B))

	if err != nil {
		ba.Conn.Write([]byte("wrong arg"))
		ba.Conn.Close()
		return
	}
	ba.Conn.Write([]byte(x.String()))
	ba.Conn.Close()
}
