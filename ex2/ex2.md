#
-What is an atomic operation? A Semaphore? A Mutex? A critical section?

An __atomic operation__ is a data operation acting on shared memory that can be finished in a single bus operation. Mostly this means reading and writing data in one cycle relative to the other threads, so that other threads cannot observe the modification. This hinders the operation on the shared variable to be interrupted and potentially corrupted by other threads. 

A __semaphore__ is a variable that describes the availability of a shared resource. A basic semaphore is a counter that describes how much of the shared memory or system resources that is being used by the multiple processes. 

A __critical section__ is the time of execution when a program is accessing a shared resource. This section of the operation is protected so that concurrent operation on the resource is not allowed within this time interval.

A __mutex__ (mutual exclusion) is the requirement that only one thread may enter its critical section at a time. This hinders race conditions. 