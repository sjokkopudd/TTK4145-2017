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

-Without priority inheritance

| Task\Time | 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10| 11| 12| 13| 14| 
|-----------|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| a | | | | |`E`| | | | | | |`Q`|`V`|`E`| |
| b | | |`E`|`V`| |`V`|`E`|`E`|`E`| | | | | | |
| c |`E`|`Q`| | | | | | | |`Q`|`Q`| | | |`E`|

-With priority inheritance

| Task\Time | 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10| 11| 12| 13| 14|
|-----------|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
|  a        |   |   |   |   |`E`|   |   |`Q`|   |`V`|`E`|   |   |   |   |
|  b        |   |   |`E`|`V`|   |   |   |   |`V`|   |   |`E`|`E`|`E`|   |
|  c        |`E`|`Q`|   |   |   |`Q`|`Q`|   |   |   |   |   |   |   |`E`|



