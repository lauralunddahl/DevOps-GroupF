## Log containing arguments and reflections on choices of language, DBMS etc.

### Programming language

We decided to choose Go as the programming language for our service. Reasons for our choice are listed below:
- We have no experience with programming in Go and wish to gain experience with this
- Go supports concurrency it makes it easy to write programs that get the most out of multicore and networked machines. 
- Go is good for running an app in the background as a single process because it is using goroutines instead of threads. 
This requires much less RAM which reduces the risk of the app crashing due to lack of memory. 
- Go is also known for being good for writing web applications and working with API requests. 
- Overall Go has good performance (somewhat similar to Java and it is in general a lot faster than python which is currently used for the service).

### DB abstraction layer

When implementing a DB abstraction layer. such that we no longer communicate directly with the DB in the main application, 
we have chosen to move away from SQLite as our DBMS and instead chosen to use an online MySQL DB hosted on ITU's servers. 

Our main choices for this was:
- MySQL (compared to SQLite) is more flexible in regard to datatypes (it supports more datatypes)
- SQLite is more optimal for smaller databases, and is harder to optimize performance for, while MySQL is easilier scalable.
Although we do not have much data now, we wished to prepare for potentially having to deal with a lot more data once the simulation starts. 
- MySQL is more suitable for multiple user access.

(Source: https://www.hostinger.com/tutorials/sqlite-vs-mysql-whats-the-difference/)

We have chosen to use an Object-Relational Mapping (ORM) framework to decouple our application from the DBMS. 
This allows us to easily do basic CRUD operations and minimize repetitive code (e.g. mapping query results to fields). 
Moreover, for future use, in case we wish to switch to another DBMS it makes it much easier because of this abstration layer.
We have chosen to use the framework GORM as our ORM framework. GORM is a full featured ORM framework developed for Go.
Our main reason for choosing GORM over other ORM frameworks was that it has the features we need and at the same time 
has good documentation which makes it easy for us to learn to use.

### Virtualization techniques and deployment targets

We have used (?) as our local provider and DigitalOcean as a remote provider for deploying our VM as a droplet.  
We have used Docker for creating a Docker image for executing our program and a Docker compose file that specifies how it 
must be run on the droplet. 

We have chosen to use Docker (i.e. containers) instead of virtualizing hardware through a virtual machine as Docker enables us to
isolate our code into a single container which can be run virtually anywhere. It makes it easier for us to update our app and will
likely make it easier for us to introduce a CI/CD setup.

Deployment target: ? (Ubuntu?)

### CI/CD setup

We have chosen to use a continious integration system in which we can run automated builds and later on also automated tests.
This allows us to integrate our code often, into our master branch, from which we can automatically build our program.
We have chosen to use Travis CI as our CI tool. Travis was chosen as it is easy to integrate with our GitHub, it offers many
automated features/options and it is cloud based which means we do not need to run and maintain a server for it.

