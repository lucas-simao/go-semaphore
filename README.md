# go-semaphore
This is an example to test semaphore on Golang.

# Testing request with semaphore limit like 5.

Request: 
http://localhost:9000/with-semaphore?quantity=50

Response:
```json
{
    "TotalTime": 9.010918458,
    "Goroutines": 5,
    "Text": "total time execution 9s and created 5 goroutines"
}
```

# Testing request without semaphore limit.

Request:
http://localhost:9000/without-semaphore?quantity=50

Response:
```json
{
    "TotalTime": 0.000111833,
    "Goroutines": 50,
    "Text": "total time execution 0s and created 50 goroutines"
}
```

