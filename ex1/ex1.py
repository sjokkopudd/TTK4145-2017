#!/usr/bin/python3

import threading

arg = 0

# Define a function for the thread
def add():
    for i in range(1000000):
        global arg
        arg += 1

def sub():
    for i in range(1000000):
        global arg
        arg -= 1

# Create two threads as follows

try:
    t1 = threading.Thread(target=add)
    t2 = threading.Thread(target=sub)
    
    t1.start()
    t2.start()

    t1.join()
    t2.join()
except:
    print ("Error: unable to start thread")


print(arg)

