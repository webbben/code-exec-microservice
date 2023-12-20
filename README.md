## Code Execution Microservice
This is a code execution microservice that will execute code snippets in specific programming languges; as of now, the list includes Python, Go, and Bash.
It's essentially an HTTP server that exposes a POST API endpoint which allows you to send code, and it returns the output and whether an error occurred during its execution.

## How it works
The actual code execution is relatively simple: 
* the requested code is written to a file of the correct type (i.e. its programming language extension)
* a docker container is created using an image that supports the requested code's language (Python, Go, etc).
* the file we created for the code is imported into the docker container as a volume, and the file is run.
* any output from the code is read from the container's stdout and stderr to our Go program. container is closed.

This does mean that a code snippet needs to print something out to stdout (i.e. print() in Python) for our program to get any output data.
So, consuming applications should consider that and make sure to add code that prints to stdout whatever data they want from the execution.

## How to use
Run this server by navigating to the project root folder and enter:
```go
go run main.go
```

Then, you can make POST requests to the API on localhost:8081/ (it's just the root endpoint).

The request body should be like this:
```json
{
  "lang": "python",
  "code": "print(\"Hello world!\")"
}
```
Note that double quotation marks should be back-slashed, and \n should be used for line breaks.

## Performance
(todo)
* add details here on average request speeds for different payload sizes/languages
* add details on average memory consumption per execution

## Todo
This isn't deployed or in its final state yet, so here are future things I want to implement:
* code execution time limits - abort a container if execution takes longer than an arbitrary amount of time, like 10 seconds
* memory and cpu consumption limit per execution - abort a container (or restrict execution some way) if code execution exceeds an arbitrary resource limit.
  or simply make limit containers to a certain cpu and memory limit, and this will be handled automatically.
* load balancing for containers - implement a mechanism to delay code execution/container creation if the server is handling too many requests simultaneously (in other words, the server exceeds a resource limit)
* add support for more languages - currently only supports: python, go, bash.  consider adding support for: java, javascript, rust, c++
