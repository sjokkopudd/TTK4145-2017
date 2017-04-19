# Exeercice 9

__Properties__
__Task 1__

- Why do we assign priorities to tasks?

Some taks in a process may be time critical, or necessary to complete for the system to function correctly. If two tasks are ready to be 
executed, and one is time critical, it is natural to assign higher priority to this one, so that it will be completed first in the case that the two tasks demand the use of the same recources.

- What features must a scheduler have for it to be usable for real-time systems?

To be usable in a real-time system, a scheduler must be able to finish tasks within deadlines. In real time systems, the output from a time critical task is useless if it is not computed within a relevant time line. Nor can computing recources be assinged to a single task for too long, stopping the system from progressing. Therefore, the ability to meet deadlines is necessary for a scheduler in a real time system. 

__Inversion and inheritance__

__Task 2__

- Without priority inheritance

| Task\Time | 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10| 11| 12| 13| 14| 
|-----------|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| a | | | | |`E`| | | | | | |`Q`|`V`|`E`| |
| b | | |`E`|`V`| |`V`|`E`|`E`|`E`| | | | | | |
| c |`E`|`Q`| | | | | | | |`Q`|`Q`| | | |`E`|

- With priority inheritance

| Task\Time | 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10| 11| 12| 13| 14|
|-----------|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
|  a        |   |   |   |   |`E`|   |   |`Q`|   |`V`|`E`|   |   |   |   |
|  b        |   |   |`E`|`V`|   |   |   |   |`V`|   |   |`E`|`E`|`E`|   |
|  c        |`E`|`Q`|   |   |   |`Q`|`Q`|   |   |   |   |   |   |   |`E`|

__Task 3__

- What is priority inversion? What is unbounded priority inversion?

  Priority inversion is when a task of lower priority is assigned recources, unintentionally preempting a higher priority task and unabling it to execute as intended by the static priotiry assignment. This inverts the priority assignment, causng the lower priority task to exctue as a high prority task and vise versa. 


  Unbounded priority inversion occurs when lower priority tasks take and give access to receources to each outher, preempting higher priority tasks indefinetly. There is no guaranteee that the higher priority tasks will ever finish their execution, as the access to the recources are blocked by the lower priority tasks.

- Does priority inheritance avoid deadlocks?

  No

__Utilization and response time__

__Task 4__

- There are a number of assumptions/conditions that must be true for the utilization and response time tests to be usable (The "simple task model"). What are these assumptions? Comment on how realistic they are.

   - Fixed set of tasks (No sporadic tasks. Not optimal, but can be worked around)
   - Periodic tasks with known periods (Realistic in many systems)
   - The tasks are independent (Completely realistic in an embedded system)
   - Overheads, switching times can be ignored (Depends)
   - Deadline == Period (Inflexible, but fair enough)
   - Fixed Worst Case Execution Time. (Not realistic to know a tight (not overly conservative) estimate here)
   - And in addition: Rate-Monotonic Priority ordering (Our choice, so ok)
  
 - Perform the utilization test for task set 2. Is the task set schedulable?
  
 U = 15/50 + 10/30 + 5/20 = 0.8833, 3*(2^(1/3)-1) = 0.7798
 
 Since 0.8833 > 0.7798, the utilization test fails, and we might not be able to schedule the task set. 
 
- Perform response-time analysis for task set 2. Is the task set schedulable? If you got different results than in 2), explain why.

  - Task `c`:  
     `w0 = 5`
     => `Rc = 5 <= 20`, ok
   - Task `b`:  
     `w0 = 10`  
     `w1 = 10 + ceil(10/20)*5 = 15`  
     `w2 = 10 + ceil(15/20)*5 = 15`  
     => `Rb = 15 <= 30`, ok
   - Task `a`:  
     `w0 = 15`  
     `w1 = 15 + ceil(15/30)*10 + ceil(15/20)*5 = 15 + 10 + 5 = 30`  
     `w2 = 15 + ceil(30/30)*10 + ceil(30/20)*5 = 15 + 10 + 10 = 35`  
     `w3 = 15 + ceil(35/30)*10 + ceil(35/20)*5 = 15 + 20 + 10 = 45`  
     `w4 = 15 + ceil(45/30)*10 + ceil(45/20)*5 = 15 + 20 + 15 = 50`  
     `w5 = 15 + ceil(50/30)*10 + ceil(50/20)*5 = 15 + 20 + 15 = 50`  
     => `Ra = 50 <= 50`, ok
     
    Conclusion: Task set is schedulable  

    The utilization test is sufficient, but not necessary. The response-time analysis is both sufficient and necessary.

