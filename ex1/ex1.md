# Reasons for concurrency and parallelism

- What is concurrency? What is parallelism? What's the difference?
 
__Concurrency__ is when you can have multiple tasks start, run and stop in arbitrary order and in overlapping time periods while producingg the same result. The tasks does not necessarily have to both be running at the same instant. Eg. multitasking on a single-core machine.

__Parallelism__ is when tasks run at the same time, eg. on a multicore processor.

#
 
- Why have machines become increasingly multicore in the past decade?
 
Stagnation in the development of higher clock speed procecors, limited by heat, means better preformance is made possible by multicore processors.
 
#
 
- What kinds of problems motivates the need for concurrent execution? (Or phrased differently: What problems do concurrency help in solving?)

Problems that can be separated into subproblems. The subproblems can be solved in overlapping time intervals.

#
 
- Does creating concurrent programs make the programmer's life easier? Harder? Maybe both? (Come back to this after you have worked on part 4 of this exercise)
 
When creating a concurrent program correct momdularisation of the program is necessary. This may be very difficult, and if done wrongly be hard to debug. On the other hand, a correct concurrent program can be expanded more easily that senquential programs and be easier to maintain. 

#

<!-- -->
- What are the differences between processes, threads, green threads, and coroutines?
 
__Threads__ are an execution of a sequence of programmed instructions. Thay are not independent, as processes, but exists as subsets of processes. Threads within a process may share process state, memory and other recources. 

__Processes__ are instances of a computer program, and are the actual execution of the instructions of the program. 
 
Where normal threads are typically scheduled by the operating system, __green threads__ are scheduled by a runtime library or a virtual machine. Green threads can be used when lacking a native multithread support; we may emulate a multithreating environment by using a VM.
 
__Coroutines__ can be thought of as threads that can pause their operations in favour of the program executing another coroutine. The coroutines cooperate, there is no return statements, but a well defined relationship between them. There is no concurrency as functions do not execute independently, nor in parallel, but are given the controll when needing it, and yielding the controll back to the other coroutine when finished or able. 

#

- Which one of these do `pthread_create()` (C/POSIX), `threading.Thread()` (Python), `go` (Go) create?
 
`pthread_create()` creates a thread.

`threading.Thread()` creates a thread.

`go` creates a _goroutine_ (?)

#

- How does pythons Global Interpreter Lock (GIL) influence the way a python Thread behaves?
 
Some languages, such as CPython, are not _thread-safe_. This is to say that the language cannot distribute resources between multiple threads running at the same time guaranteeing that the result will be correct. For this reason, there exists a __Global Interpreter Lock__ (GIL), which ensures that only one thread will be able to run at a time, even when running on a muticore processor. This will lower performance, but is a safe way to ensure correct program execution. 

#

- With this in mind: What is the workaround for the GIL (Hint: it's another module)?
 
There is a hack to spawn an interperter per thread. This makes it thread safe, but is not true parrallelssimsms.

#

- What does `func GOMAXPROCS(n int) int` change? 

GOMAXPROCS sets the maximum number of CPUs that can be executing simultaneously and returns the previous setting. If n < 1, it does not change the current setting. The number of logical CPUs on the local machine can be queried with NumCPU. This call will go away when the scheduler improves.

