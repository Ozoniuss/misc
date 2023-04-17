# Context Madness

The following program creates a binary tree of goroutines, and makes use of context to progpagate cancellation signals to the children. Initially you have:

```
    0 (main goroutine)
    |
    1 (child goroutine 1)
   / \ 
  2   3 (child goroutines 2 and 3)
 / \ / \
4  5 6  7 (child goroutines 4-7)
   ...    
   ...    (child goroutines 8-15 and so on) 
```

Every goroutine, starting from the main goroutine, creates a new context and passes it to the direct descendants. When a goroutine is killed, it generates a chain reaction that stops both the goroutine and the children. The program stops once all goroutines are killed (thus it always required stopping goroutine 1, and stopping goroutine 1 stops the program at once).

Killing goroutine 2 does the following:

```
    0
    |
    1 
   / \ 
  X   3 (child goroutine 2 killed) 
 / \ / \ 
X  X 6  7 (child goroutines 4 and 5 canceled via propagation) 
```

Running
-------

Run with 

```
go run main.go
```

Input a valid depth (between 1 and 10 should be reasonable, at some point the OS will get overwhelmed with all those goroutines and stdout printing) which will generate a goroutine tree and start killing goroutines!

You can input 0 to find out which goroutines are still alive.

Sample run
----------

https://user-images.githubusercontent.com/56697671/232488220-532cd79a-8ce2-4ea7-b3b0-9cc1bb33711e.mp4

Limitations
-----------

- Since the goroutine getting user input uses `stdin` and the goroutine getting user output uses `stdout`, they share the same windows, so it's difficult to synchronize writing a suggestive message just before the user input, because it might get interlaced with the goroutines printing the output. A solution for that would be, from the main thread, use a channel that notifies when all goroutines that are supposed to stop at one operation have actually stopped (taking into account those who were already stopped previously) and notify the reader goroutine, which is blocked until then. Since the main goroutine receives the signal another goroutine died _at the same time_ with the goroutine stopping through an unbuffered channel, it is guaranteed that at that point the goroutine outputted the message. This should not be that hard to implement but I kindof got bored with it.
- If sending two kill signals to the same goroutine fast enough, it might be possible to get a panic if the map with the alive goroutines has not been updated before the reader thread sends the signal again. This is highly unlikely, but for completeness can be fixed with the solution mentioned above. Now that I think about it I might just as well do it...
