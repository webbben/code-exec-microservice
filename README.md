## Code Execution Microservice
This is a code execution microservice that will execute code snippets in specific programming languages.
It's an HTTP server that exposes a POST API endpoint which allows you to send code, and it returns the output and whether an error occurred during its execution.

As of now, these are the supported programming languages:
* Python
* Go
* Bash

## How it works
The actual code execution is relatively simple: 
* this service is containerized with Docker and deployed on Fly.io. The docker image has runtimes installed for Python, Go, and Bash.
* the requested code is written to a file of the correct type (i.e. its programming language extension)
* the file is run, and its output is captured. If an error occurs, it also captures this information. This information is returned to the API consumer.
* Note: if a fatal error occurred, whether it be intentional (i.e. from a malicious code snippet) or not, the service will not have its availability affected. Fly.io handles deploying new containers as needed to keep this service operational.

This does mean that a code snippet needs to print something out to stdout (i.e. print() in Python) for our program to get any output data.
So, consuming applications should consider that and make sure to add code that prints to stdout whatever data they want from the execution.

## How to use
Make a POST request to the following URL: https://code-exec-microservice.fly.dev/
In the request body, include JSON data in the following format:

```json
{
  "code": "print(\"Hello world!\")",
  "lang": "python"
}
```

## Performance
(todo)
* add details here on average request speeds for different payload sizes/languages
* add details on average memory consumption per execution

## Future
Here are things I'll consider adding in the future:
* Support for more languages - it might be fun to add more popular languages, like Java, Javascript, C++, Rust, etc.
