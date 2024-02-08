package main

// Todo: write parser function to split and parse java log files
import time


cur = 0
while True:
    try:
        with open('myfile') as f:
            f.seek(0,2)
            if f.tell() < cur:
                f.seek(0,0)
            else:
                f.seek(cur,0)
            for line in f:
                print line.strip()
            cur = f.tell()
    except IOError, e:
        pass
    time.sleep(1)