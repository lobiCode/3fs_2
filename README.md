### About

Build a simple queue manager which receives requests of four different tasks 
and then distributes them among four types of processes:

Markup : *Fibonacci resolver: Takes an integer and returns the result of Fibonacci function.
		*Basic arithmetic resolver: Takes basic arithmetic problem and returns the result.
		*Reverse text resolver: Takes mirrored text and returns it reversed.
		*Text encoder: Takes string of text and returns BCrypt encrypted hash.

Each type of process/server runs in several instances. Queue manager has to decide
which instance to use based on availability (meaning that the instance is not processing any task at the time).

###### First solution:

Queue manager is implemented as JobQueue of a buffered channels of Job.
At the beginning  N processes are started (workers). When a request is recieved dispatcher 
creates a Job and sends it on the JobQueue. JobQueue then sends Job  to an available process (worker).
Despite the fact that this solution is not exactly what we want (because  any process (worker) can handle any task (job))
I used it because it is very idiomatic and elegant. Also it can handle a lot requests in very short amount of time.

###### Second solution:

Queue manager is implemented as main queue and four taskQueue.
At the beginning  N processes are started (workers) for four different tasks. 
For each request  a Job is created and pushed to main queue. Main
queue then dispatch a Job (regarding to the type) to the taskQueue.
Each available instance then pull job from a taskQueue.

### INSTALLATION


```go
go get  golang.org/x/crypto/bcrypt
```

**First solution:**
```go
go run main.go
```
**Second solution:**
	
```go
go run sol2.go qman2.go
```

### Example
Request should be send through a TCP socket.

```shell
$ echo Fibonacci 8 | nc localhost 1234
-> 21
```

```shell
$ echo ReverseText foo | nc localhost 1234
-> oof
```

```shell
$ echo TextEncoder foo | nc localhost 1234
-> $2a$10$zUBd22nVobnU76wt8gN32uZBBwX08IjSU.1IrvuBy3HTrRXv36w8.
```

```shell
$ echo BasicArithmetic 4+4*8-2/7 | nc localhost 1234
-> 35.71428571
```

Or you can use simple client in client folder

```go
go run client BasicArithmetic 4+4*8-2/7
-> 35.71428571
```





