# Week 2


## Refactoring ITU-MiniTwit to another language and tecnhology
We decided to choose Go as our new programming language for our service. One of the reasons for choosing Go is that we have no experience with programming in Go and wish to gain experience with this. Moreover, since Go supports concurrency it makes it easy to write programs that get the most out of multicore and networked machines. It is good for running an app in the background as a single process because it is using goroutines instead of threads which requires much less RAM which reduces the risk of the app crashing due to lack of memory. It is also known for being good for writing web applications and working with API requests. Overall Go has good performance (somewhat similar to Java and it is in general a lot faster than python which we are refactoring from). 

One of the disadvantages of Go is that the runtime safety is not as good as in other programming languages. Since the compile time safety is good we don’t see this as a reason for not choosing this language.  

## Commands used

code --disable-gpu

export GOPATH=$HOME/go

Go run minitwit.go
